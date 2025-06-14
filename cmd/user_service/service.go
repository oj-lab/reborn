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
		return nil, errors.New("name and email are required")
	}

	user := &userpb.CreateUserRequest{
		Name:  req.GetName(),
		Email: req.GetEmail(),
		Role:  req.GetRole(),
	}

	// If password is provided, hash it
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		user.Password = &hashedPassword
	}

	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &userpb.GetUserResponse{User: user}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		user.Name = req.GetName()
	}
	if req.Email != nil {
		user.Email = req.GetEmail()
	}
	updateReq := &userpb.UpdateUserRequest{
		Id:    user.Id,
		Name:  &user.Name,
		Email: &user.Email,
	}
	err = s.repo.UpdateUser(ctx, updateReq)
	if err != nil {
		return nil, err
	}
	return &userpb.UpdateUserResponse{User: user}, nil
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
	if userModel.Password == nil || *userModel.Password == "" {
		return nil, errors.New("user has not set a password")
	}

	// Verify password
	valid, err := VerifyPassword(req.GetPassword(), *userModel.Password)
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
