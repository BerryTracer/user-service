package service

import (
	"context"

	"github.com/BerryTracer/user-service/model"
	"github.com/BerryTracer/user-service/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, username, email, password string) (*model.User, error)
	GetUserById(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
}

// CreateUser implements UserService.
func (s *UserServiceImpl) CreateUser(ctx context.Context, username string, email string, password string) (*model.User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := model.NewUser(username, email, string(hashedPassword))

	err = s.UserRepository.CreateUser(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser implements UserService.
func (s *UserServiceImpl) GetUserById(ctx context.Context, id string) (*model.User, error) {
	return s.UserRepository.GetUserById(ctx, id)
}

// GetUserByEmail implements UserService.
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.UserRepository.GetUserByEmail(ctx, email)
}

// GetUserByUsername implements UserService.
func (s *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.UserRepository.GetUserByUsername(ctx, username)
}

// Ensure UserServiceImpl implements UserService.
var _ UserService = &UserServiceImpl{}
