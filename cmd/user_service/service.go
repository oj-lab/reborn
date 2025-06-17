package main

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	userpb "github.com/oj-lab/reborn/protobuf/user"
	"golang.org/x/crypto/argon2"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type UserService struct {
	repo UserRepository
	userpb.UnimplementedUserServiceServer
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// UserServiceServer implements userpb.UserServiceServer
var _ userpb.UserServiceServer = (*UserService)(nil)

func (s *UserService) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*emptypb.Empty, error) {
	if req.GetName() == "" || req.GetEmail() == "" {
		return nil, errors.New("name and email are required for new user registration")
	}

	if (req.Password == nil || *req.Password == "") && (req.GithubId == nil || *req.GithubId == "") {
		return nil, errors.New("github_id or password is required for new user registration")
	}

	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, req.GetEmail())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check for existing user: %w", err)
	}

	// User exists
	if existingUser != nil {
		// User already has a password, so they are fully registered
		if existingUser.HashedPassword != nil && *existingUser.HashedPassword != "" {
			return nil, errors.New("user with this email already exists")
		}

		// User exists from a Github login and is trying to set a password
		if req.Password != nil && *req.Password != "" {
			hashedPassword, err := HashPassword(*req.Password)
			if err != nil {
				return nil, fmt.Errorf("failed to hash password: %w", err)
			}
			err = s.repo.SetPassword(ctx, uint64(existingUser.ID), hashedPassword)
			if err != nil {
				return nil, fmt.Errorf("failed to set password for existing user: %w", err)
			}
			return &emptypb.Empty{}, nil
		}
	}

	// User does not exist, create a new one
	var hashedPassword *string
	if req.Password != nil && *req.Password != "" {
		h, err := HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		hashedPassword = &h
	}

	createUserReq := &userpb.CreateUserRequest{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Role:     req.GetRole(),
		Password: hashedPassword,
		GithubId: req.GithubId,
	}

	err = s.repo.CreateUser(ctx, createUserReq)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) GithubLogin(ctx context.Context, req *userpb.GithubLoginRequest) (*userpb.LoginResponse, error) {
	if req.GetGithubId() == "" || req.GetEmail() == "" {
		return nil, errors.New("github id and email are required")
	}

	// Case 1: User exists with this github_id
	user, err := s.repo.GetUserByGithubID(ctx, req.GetGithubId())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get user by github id: %w", err)
	}

	if user != nil {
		// User found, generate token and log them in
		token := fmt.Sprintf("user_%d_token", user.ID)
		return &userpb.LoginResponse{
			User:  user.ToPb(),
			Token: token,
		}, nil
	}

	// Case 2: No user with github_id, check if email exists
	user, err = s.repo.GetUserByEmail(ctx, req.GetEmail())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user != nil {
		// User with this email exists, link the github_id
		updateReq := &userpb.UpdateUserRequest{
			Id:       uint64(user.ID),
			GithubId: &req.GithubId,
		}
		err := s.repo.UpdateUser(ctx, updateReq)
		if err != nil {
			return nil, fmt.Errorf("failed to link github id: %w", err)
		}
		// Generate token and log them in
		token := fmt.Sprintf("user_%d_token", user.ID)
		user.GithubID = &req.GithubId // manually update for response
		return &userpb.LoginResponse{
			User:  user.ToPb(),
			Token: token,
		}, nil
	}

	// Case 3: New user, create them
	createUserReq := &userpb.CreateUserRequest{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Role:     userpb.UserRole_USER, // Default role
		GithubId: &req.GithubId,
	}
	err = s.repo.CreateUser(ctx, createUserReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Fetch the newly created user to get their ID
	newUser, err := s.repo.GetUserByGithubID(ctx, req.GetGithubId())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch newly created user: %w", err)
	}

	token := fmt.Sprintf("user_%d_token", newUser.ID)
	return &userpb.LoginResponse{
		User:  newUser.ToPb(),
		Token: token,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &userpb.GetUserResponse{User: user}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	userModel, err := s.repo.GetUserModelByID(ctx, req.GetId())
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Determine the final state of GithubID after the potential update.
	githubIdWillBeEmpty := (req.GithubId != nil && *req.GithubId == "") ||
		(req.GithubId == nil && (userModel.GithubID == nil || *userModel.GithubID == ""))

	// Determine the final state of the password after the potential update.
	passwordWillBeEmpty := (req.Password != nil && *req.Password == "") ||
		(req.Password == nil && (userModel.HashedPassword == nil || *userModel.HashedPassword == ""))

	// Validation: Ensure the user retains at least one login method.
	if githubIdWillBeEmpty && passwordWillBeEmpty {
		return nil, errors.New("user cannot have both github_id and password empty")
	}

	var hashedPassword *string
	if req.Password != nil {
		if *req.Password != "" {
			h, err := HashPassword(*req.Password)
			if err != nil {
				return nil, err
			}
			hashedPassword = &h
		}
	}

	// Build the request for the repository
	repoReq := &userpb.UpdateUserRequest{
		Id:       req.GetId(),
		Name:     req.Name,
		Email:    req.Email,
		GithubId: req.GithubId,
		Role:     req.Role,
	}

	if req.Password != nil {
		if *req.Password == "" {
			emptyString := ""
			repoReq.Password = &emptyString
		} else {
			repoReq.Password = hashedPassword
		}
	}

	err = s.repo.UpdateUser(ctx, repoReq)
	if err != nil {
		return nil, err
	}

	updatedUser, err := s.repo.GetUserByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userpb.UpdateUserResponse{User: updatedUser}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	err := s.repo.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &userpb.DeleteUserResponse{Id: req.GetId()}, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	users, err := s.repo.ListUsers(ctx, req.GetPage(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return &userpb.ListUsersResponse{Users: users}, nil
}

// Password hashing configuration
type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

var defaultPasswordConfig = PasswordConfig{
	time:    1,
	memory:  64 * 1024,
	threads: 4,
	keyLen:  32,
}

// HashPassword hashes password using Argon2ID
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, defaultPasswordConfig.time,
		defaultPasswordConfig.memory, defaultPasswordConfig.threads, defaultPasswordConfig.keyLen)

	// Format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		defaultPasswordConfig.memory, defaultPasswordConfig.time,
		defaultPasswordConfig.threads, encodedSalt, encodedHash), nil
}

// VerifyPassword verifies password
func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(decodedHash)))

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}

func (s *UserService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, errors.New("email and password are required")
	}

	// Get user by email
	userModel, err := s.repo.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user has set a password
	if userModel.HashedPassword == nil || *userModel.HashedPassword == "" {
		return nil, errors.New("user has not set a password, please use github login or set a password")
	}

	// Verify password
	valid, err := VerifyPassword(req.GetPassword(), *userModel.HashedPassword)
	if err != nil || !valid {
		return nil, errors.New("invalid email or password")
	}

	// Generate simple token (should use JWT in production)
	token := fmt.Sprintf("user_%d_token", userModel.ID)

	return &userpb.LoginResponse{
		User:  userModel.ToPb(),
		Token: token,
	}, nil
}

func (s *UserService) SetPassword(ctx context.Context, req *userpb.SetPasswordRequest) (*userpb.SetPasswordResponse, error) {
	if req.GetUserId() == 0 || req.GetPassword() == "" {
		return nil, errors.New("user_id and password are required")
	}

	// Check if user exists
	_, err := s.repo.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Hash password
	hashedPassword, err := HashPassword(req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Update password
	err = s.repo.SetPassword(ctx, req.GetUserId(), hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to set password: %v", err)
	}

	return &userpb.SetPasswordResponse{}, nil
}
