package main

import (
	"context"

	userpb "github.com/oj-lab/reborn/protobuf/user"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *userpb.CreateUserRequest) error
	GetUserByID(ctx context.Context, id uint64) (*userpb.User, error)
	UpdateUser(ctx context.Context, user *userpb.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id uint64) error
	ListUsers(ctx context.Context) ([]*userpb.User, error)
}

type UserModel struct {
	gorm.Model
	Name  string
	Email string
}

func (UserModel) TableName() string {
	return "users"
}

func (m *UserModel) ToUser() *userpb.User {
	return &userpb.User{
		Id:    uint64(m.ID),
		Name:  m.Name,
		Email: m.Email,
	}
}

func (m *UserModel) FromUser(user *userpb.User) {
	m.ID = uint(user.Id)
	m.Name = user.Name
	m.Email = user.Email
}

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) CreateUser(ctx context.Context, user *userpb.CreateUserRequest) error {
	model := &UserModel{
		Name:  user.Name,
		Email: user.Email,
	}

	return r.db.WithContext(ctx).Create(model).Error
}

func (r *GormUserRepository) GetUserByID(ctx context.Context, id uint64) (*userpb.User, error) {
	model := &UserModel{}

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(model).Error; err != nil {
		return nil, err
	}

	return model.ToUser(), nil
}

func (r *GormUserRepository) UpdateUser(ctx context.Context, user *userpb.UpdateUserRequest) error {
	model := &UserModel{}
	model.ID = uint(user.Id)
	if user.Name != nil {
		model.Name = *user.Name
	}
	if user.Email != nil {
		model.Email = *user.Email
	}

	return r.db.WithContext(ctx).Save(model).Error
}

func (r *GormUserRepository) DeleteUser(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&UserModel{}).Error
}

func (r *GormUserRepository) ListUsers(ctx context.Context) ([]*userpb.User, error) {
	var models []*UserModel

	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]*userpb.User, len(models))
	for i, model := range models {
		users[i] = model.ToUser()
	}

	return users, nil
}
