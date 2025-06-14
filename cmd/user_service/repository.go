package main

import (
	"context"
	"time"

	userpb "github.com/oj-lab/reborn/protobuf/user"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *userpb.CreateUserRequest) error
	GetUserByID(ctx context.Context, id uint64) (*userpb.User, error)
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	UpdateUser(ctx context.Context, user *userpb.UpdateUserRequest) error
	SetPassword(ctx context.Context, userID uint64, passwordHash string) error
	DeleteUser(ctx context.Context, id uint64) error
	ListUsers(ctx context.Context, page, pageSize uint64) ([]*userpb.User, error)
}

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

func (role UserRole) ToPb() userpb.UserRole {
	switch role {
	case UserRoleUser:
		return userpb.UserRole_USER
	case UserRoleAdmin:
		return userpb.UserRole_ADMIN
	default:
		return userpb.UserRole_USER // Default to user if unknown
	}
}

func (role *UserRole) FromPb(pbRole userpb.UserRole) {
	switch pbRole {
	case userpb.UserRole_USER:
		*role = UserRoleUser
	case userpb.UserRole_ADMIN:
		*role = UserRoleAdmin
	default:
		*role = UserRoleUser // Default to user if unknown
	}
}

type UserModel struct {
	gorm.Model
	Name        string
	Email       string
	Password    *string // Password hash, optional
	LastLoginAt *time.Time
	Role        UserRole
}

func (UserModel) TableName() string {
	return "users"
}

func (m UserModel) ToPb() *userpb.User {
	return &userpb.User{
		Id:    uint64(m.ID),
		Name:  m.Name,
		Email: m.Email,
		Role:  m.Role.ToPb(),
	}
}

func (m *UserModel) FromPb(user *userpb.User) {
	m.ID = uint(user.Id)
	m.Name = user.Name
	m.Email = user.Email
	m.Role.FromPb(user.Role)
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
	model.Role.FromPb(user.Role)

	// If password hash is provided, set password
	if user.Password != nil && *user.Password != "" {
		model.Password = user.Password
	}

	return r.db.WithContext(ctx).Create(model).Error
}

func (r *GormUserRepository) GetUserByID(ctx context.Context, id uint64) (*userpb.User, error) {
	model := &UserModel{}
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(model).Error; err != nil {
		return nil, err
	}
	return model.ToPb(), nil
}

func (r *GormUserRepository) GetUserByEmail(ctx context.Context, email string) (*UserModel, error) {
	model := &UserModel{}
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(model).Error; err != nil {
		return nil, err
	}
	return model, nil
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
	if user.Role != nil {
		model.Role.FromPb(*user.Role)
	}
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *GormUserRepository) SetPassword(ctx context.Context, userID uint64, passwordHash string) error {
	return r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ?", userID).
		Update("password", passwordHash).Error
}

func (r *GormUserRepository) DeleteUser(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&UserModel{}).Error
}

func (r *GormUserRepository) ListUsers(
	ctx context.Context, page, pageSize uint64,
) ([]*userpb.User, error) {
	var models []*UserModel
	err := r.db.WithContext(ctx).
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	users := make([]*userpb.User, len(models))
	for i, model := range models {
		users[i] = model.ToPb()
	}
	return users, nil
}
