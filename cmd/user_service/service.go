package main

import (
	"context"
	"errors"

	userpb "github.com/oj-lab/reborn/protobuf/user"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	repo UserRepository
	userpb.UnimplementedUserServiceServer
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// UserServiceServer 实现 userpb.UserServiceServer
var _ userpb.UserServiceServer = (*UserService)(nil)

func (s *UserService) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*emptypb.Empty, error) {
	if req.GetName() == "" || req.GetEmail() == "" {
		return nil, errors.New("name and email are required")
	}
	user := &userpb.CreateUserRequest{
		Name:  req.GetName(),
		Email: req.GetEmail(),
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
