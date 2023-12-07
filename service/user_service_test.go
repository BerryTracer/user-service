package service_test

import (
	"context"
	"errors"
	"testing"

	mock_crypto "github.com/BerryTracer/common-service/crypto/mock"
	"github.com/BerryTracer/user-service/model"
	mock_repository "github.com/BerryTracer/user-service/repository/mock"
	service "github.com/BerryTracer/user-service/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestUserServiceImpl_CreateUser tests the CreateUser method of the UserServiceImpl
func TestUserServiceImpl_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_crypto.NewMockPasswordHasher(ctrl)
	userService := service.NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	username := "testuser"
	email := "testuser@example.com"
	password := "password"

	// Mock successful creation
	mockRepo.EXPECT().
		CreateUser(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, u *model.User) error {
			assert.Equal(t, username, u.Username)
			assert.Equal(t, email, u.Email)
			assert.NotEmpty(t, u.HashedPassword)
			return nil
		}).
		Times(1)

	// Mock successful hashing
	mockHasher.EXPECT().
		HashPassword(gomock.Any()).
		DoAndReturn(func(password string) (string, error) {
			assert.Equal(t, password, password)
			return "hashedPassword", nil
		}).
		Times(1)

	// Call CreateUser
	result, err := userService.CreateUser(ctx, username, email, password)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, username, result.Username)
	assert.Equal(t, email, result.Email)
	assert.NotEmpty(t, result.HashedPassword)
}

// TestUserServiceImpl_CreateUser_Validation_Error tests the CreateUser method of the UserServiceImpl
func TestUserServiceImpl_CreateUser_Repository_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_crypto.NewMockPasswordHasher(ctrl)
	userService := service.NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	username := "testuser"
	email := "testuser@example.com"
	password := "password"

	// Mock error
	mockRepo.EXPECT().
		CreateUser(ctx, gomock.Any()).
		Return(assert.AnError).
		Times(1)

	// Mock successful hashing
	mockHasher.EXPECT().
		HashPassword(gomock.Any()).
		DoAndReturn(func(password string) (string, error) {
			assert.Equal(t, password, password)
			return "hashedPassword", nil
		}).
		Times(1)

	// Call CreateUser
	result, err := userService.CreateUser(ctx, username, email, password)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserServiceImpl_CreateUser_BcryptError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_crypto.NewMockPasswordHasher(ctrl)
	userService := service.NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	username := "testuser"
	email := "testuser@example.com"
	password := "password"

	// Mock error
	mockHasher.EXPECT().
		HashPassword(gomock.Any()).
		Return("", assert.AnError).
		Times(1)

	// Call CreateUser expecting a bcrypt error
	_, err := userService.CreateUser(ctx, username, email, password)

	// Assertions
	assert.Error(t, err)
}

// TestUserServiceImpl_CreateUser_ValidationFail tests the CreateUser method of the UserServiceImpl
func TestUserServiceImpl_CreateUser_ValidationFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_crypto.NewMockPasswordHasher(ctrl)
	userService := service.NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	// Providing invalid data for validation to fail
	username := ""
	email := "invalidemail"
	password := "pwd"

	// Mock successful hashing
	mockHasher.EXPECT().
		HashPassword(gomock.Any()).
		DoAndReturn(func(password string) (string, error) {
			assert.Equal(t, password, password)
			return "hashedPassword", nil
		}).
		Times(1)

	// Call CreateUser with invalid data
	_, err := userService.CreateUser(ctx, username, email, password)

	// Assertions
	assert.Error(t, err)
}

func TestUserServiceImpl_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_crypto.NewMockPasswordHasher(ctrl)
	userService := service.NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	testID := "12345"
	expectedUser := &model.User{ID: testID, Username: "testuser", Email: "test@example.com"}

	// Mock successful retrieval
	mockRepo.EXPECT().
		GetUserById(ctx, testID).
		Return(expectedUser, nil).
		Times(1)

	// Call GetUserById
	user, err := userService.GetUserById(ctx, testID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserServiceImpl_GetUserById_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_crypto.NewMockPasswordHasher(ctrl)
	userService := service.NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	testID := "12345"

	// Mock failure in retrieval
	mockRepo.EXPECT().
		GetUserById(ctx, testID).
		Return(nil, errors.New("user not found")).
		Times(1)

	// Call GetUserById
	user, err := userService.GetUserById(ctx, testID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserServiceImpl_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_crypto.NewMockPasswordHasher(ctrl)
	userService := service.NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	testEmail := "test@example.com"
	expectedUser := &model.User{ID: "12345", Username: "testuser", Email: testEmail}

	// Mock successful retrieval
	mockRepo.EXPECT().
		GetUserByEmail(ctx, testEmail).
		Return(expectedUser, nil).
		Times(1)

	// Call GetUserByEmail
	user, err := userService.GetUserByEmail(ctx, testEmail)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserServiceImpl_GetUserByEmail_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_crypto.NewMockPasswordHasher(ctrl)
	userService := service.NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	testEmail := "test@example.com"

	// Mock failure in retrieval
	mockRepo.EXPECT().
		GetUserByEmail(ctx, testEmail).
		Return(nil, errors.New("user not found")).
		Times(1)

	// Call GetUserByEmail
	user, err := userService.GetUserByEmail(ctx, testEmail)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserServiceImpl_GetUserByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_crypto.NewMockPasswordHasher(ctrl)
	userService := service.NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	testUsername := "testuser"
	expectedUser := &model.User{ID: "12345", Username: testUsername, Email: "test@example.com"}

	// Mock successful retrieval
	mockRepo.EXPECT().
		GetUserByUsername(ctx, testUsername).
		Return(expectedUser, nil).
		Times(1)

	// Call GetUserByUsername
	user, err := userService.GetUserByUsername(ctx, testUsername)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserServiceImpl_GetUserByUsername_Fail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	mockHasher := mock_crypto.NewMockPasswordHasher(ctrl)
	userService := service.NewUserService(mockRepo, mockHasher)

	ctx := context.Background()
	testUsername := "testuser"

	// Mock failure in retrieval
	mockRepo.EXPECT().
		GetUserByUsername(ctx, testUsername).
		Return(nil, errors.New("user not found")).
		Times(1)

	// Call GetUserByUsername
	user, err := userService.GetUserByUsername(ctx, testUsername)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)
}
