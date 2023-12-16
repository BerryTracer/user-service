package repository_test

import (
	"context"
	"testing"

	mock "github.com/BerryTracer/common-service/adapter/database/mongodb/mock"
	"github.com/BerryTracer/user-service/model"
	"github.com/BerryTracer/user-service/repository"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TestUserMongoRepository_CreateUser tests the CreateUser method of the UserMongoRepository
func TestUserMongoRepository_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mock.NewMockMongoAdapter(ctrl)
	repo := repository.NewUserMongoRepository(mockAdapter)

	ctx := context.Background()
	user := model.NewUser("test", "test@mail.com", "test")

	userDB, _ := user.ToUserDB()

	mockAdapter.EXPECT().
		InsertOne(ctx, userDB, gomock.Any()).
		Return(&mongo.InsertOneResult{}, nil).
		Times(1)

	err := repo.CreateUser(ctx, user)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestUserMongoRepository_CreateUser_InvalidID tests the CreateUser method of the UserMongoRepository
func TestUserMongoRepository_CreateUser_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mock.NewMockMongoAdapter(ctrl)
	repo := repository.NewUserMongoRepository(mockAdapter)

	ctx := context.Background()
	user := model.NewUser("test", "test@mail.com", "test")
	user.ID = "invalid"

	err := repo.CreateUser(ctx, user)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUserMongoRepository_CreateUser_FindOne_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mock.NewMockMongoAdapter(ctrl)
	repo := repository.NewUserMongoRepository(mockAdapter)

	ctx := context.Background()
	user := model.NewUser("test", "test@mail.com", "test")

	userDB, _ := user.ToUserDB()

	mockAdapter.EXPECT().
		InsertOne(ctx, userDB, gomock.Any()).
		Return(&mongo.InsertOneResult{}, mongo.ErrNoDocuments).
		Times(1)

	err := repo.CreateUser(ctx, user)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

// TestUserMongoRepository_GetUserById tests the GetUserById method of the UserMongoRepository
func TestUserMongoRepository_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMongoAdapter := mock.NewMockMongoAdapter(ctrl)
	mockSingleResult := mock.NewMockSingleResult(ctrl)
	userRepo := repository.NewUserMongoRepository(mockMongoAdapter)

	ctx := context.Background()
	testID := primitive.NewObjectID().Hex()
	objectID, _ := primitive.ObjectIDFromHex(testID)

	var userDB model.UserDB
	userDB.ID = objectID

	// Setup mock expectations
	mockMongoAdapter.EXPECT().
		FindOne(ctx, primitive.M{"_id": objectID}).
		Return(mockSingleResult).
		Times(1)

	mockSingleResult.EXPECT().
		Decode(gomock.Any()).
		SetArg(0, userDB).
		Return(nil).
		Times(1)

	// Call the method
	resultUser, err := userRepo.GetUserById(ctx, testID)

	// Assertions
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if resultUser.ID != testID {
		t.Errorf("expected user id %s, got %s", testID, resultUser.ID)
	}

	if resultUser.Username != userDB.Username {
		t.Errorf("expected username %s, got %s", userDB.Username, resultUser.Username)
	}

	if resultUser.Email != userDB.Email {
		t.Errorf("expected email %s, got %s", userDB.Email, resultUser.Email)
	}
}

// TestUserMongoRepository_GetUserById_InvalidID tests the GetUserById method of the UserMongoRepository
func TestUserMongoRepository_GetUserById_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMongoAdapter := mock.NewMockMongoAdapter(ctrl)
	userRepo := repository.NewUserMongoRepository(mockMongoAdapter)

	ctx := context.Background()
	testID := "invalid"

	// Call the method
	resultUser, err := userRepo.GetUserById(ctx, testID)

	// Assertions
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if resultUser != nil {
		t.Errorf("expected nil user, got %v", resultUser)
	}
}

func TestUserMongoRepository_GetUserById_FindOne_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMongoAdapter := mock.NewMockMongoAdapter(ctrl)
	mockSingleResult := mock.NewMockSingleResult(ctrl)
	userRepo := repository.NewUserMongoRepository(mockMongoAdapter)

	ctx := context.Background()
	testID := primitive.NewObjectID().Hex()
	objectID, _ := primitive.ObjectIDFromHex(testID)

	// Setup mock expectations
	mockMongoAdapter.EXPECT().
		FindOne(ctx, primitive.M{"_id": objectID}).
		Return(mockSingleResult).
		Times(1)

	mockSingleResult.EXPECT().
		Decode(gomock.Any()).
		Return(mongo.ErrNoDocuments).
		Times(1)

	// Call the method
	resultUser, err := userRepo.GetUserById(ctx, testID)

	// Assertions
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if resultUser != nil {
		t.Errorf("expected nil user, got %v", resultUser)
	}
}

// TestUserMongoRepository_GetUserByEmail tests the GetUserByEmail method of the UserMongoRepository
func TestUserMongoRepository_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMongoAdapter := mock.NewMockMongoAdapter(ctrl)
	mockSingleResult := mock.NewMockSingleResult(ctrl)

	userRepo := repository.NewUserMongoRepository(mockMongoAdapter)

	ctx := context.Background()
	testEmail := "test@mail.com"

	var userDB model.UserDB
	userDB.Email = testEmail

	// Setup mock expectations
	mockMongoAdapter.EXPECT().
		FindOne(ctx, primitive.M{"email": testEmail}).
		Return(mockSingleResult).
		Times(1)

	mockSingleResult.EXPECT().
		Decode(gomock.Any()).
		SetArg(0, userDB).
		Return(nil).
		Times(1)

	// Call the method
	resultUser, err := userRepo.GetUserByEmail(ctx, testEmail)

	// Assertions
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if resultUser.Email != testEmail {
		t.Errorf("expected email %s, got %s", testEmail, resultUser.Email)
	}

	if resultUser.Username != userDB.Username {
		t.Errorf("expected username %s, got %s", userDB.Username, resultUser.Username)
	}

	if resultUser.ID != userDB.ID.Hex() {
		t.Errorf("expected id %s, got %s", userDB.ID.Hex(), resultUser.ID)
	}
}

func TestUserMongoRepository_GetUserByEmail_FindOne_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMongoAdapter := mock.NewMockMongoAdapter(ctrl)
	mockSingleResult := mock.NewMockSingleResult(ctrl)
	userRepo := repository.NewUserMongoRepository(mockMongoAdapter)

	ctx := context.Background()
	testEmail := "test@mail.com"

	// Setup mock expectations
	mockMongoAdapter.EXPECT().
		FindOne(ctx, primitive.M{"email": testEmail}).
		Return(mockSingleResult).
		Times(1)

	mockSingleResult.EXPECT().
		Decode(gomock.Any()).
		Return(mongo.ErrNoDocuments).
		Times(1)

	// Call the method
	resultUser, err := userRepo.GetUserByEmail(ctx, testEmail)

	// Assertions
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if resultUser != nil {
		t.Errorf("expected nil user, got %v", resultUser)
	}
}

// TestUserMongoRepository_GetUserByUsername tests the GetUserByUsername method of the UserMongoRepository
func TestUserMongoRepository_GetUserByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMongoAdapter := mock.NewMockMongoAdapter(ctrl)
	mockSingleResult := mock.NewMockSingleResult(ctrl)

	userRepo := repository.NewUserMongoRepository(mockMongoAdapter)

	ctx := context.Background()
	testUsername := "test"

	var userDB model.UserDB
	userDB.Username = testUsername

	// Setup mock expectations
	mockMongoAdapter.EXPECT().
		FindOne(ctx, primitive.M{"username": testUsername}).
		Return(mockSingleResult).
		Times(1)

	mockSingleResult.EXPECT().
		Decode(gomock.Any()).
		SetArg(0, userDB).
		Return(nil).
		Times(1)

	// Call the method
	resultUser, err := userRepo.GetUserByUsername(ctx, testUsername)

	// Assertions
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if resultUser.Username != testUsername {
		t.Errorf("expected username %s, got %s", testUsername, resultUser.Username)
	}

	if resultUser.Email != userDB.Email {
		t.Errorf("expected email %s, got %s", userDB.Email, resultUser.Email)
	}

	if resultUser.ID != userDB.ID.Hex() {
		t.Errorf("expected id %s, got %s", userDB.ID.Hex(), resultUser.ID)
	}
}

func TestUserMongoRepository_GetUserByUsername_FindOne_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMongoAdapter := mock.NewMockMongoAdapter(ctrl)
	mockSingleResult := mock.NewMockSingleResult(ctrl)
	userRepo := repository.NewUserMongoRepository(mockMongoAdapter)

	ctx := context.Background()
	testUsername := "test"

	// Setup mock expectations
	mockMongoAdapter.EXPECT().
		FindOne(ctx, primitive.M{"username": testUsername}).
		Return(mockSingleResult).
		Times(1)

	mockSingleResult.EXPECT().
		Decode(gomock.Any()).
		Return(mongo.ErrNoDocuments).
		Times(1)

	// Call the method
	resultUser, err := userRepo.GetUserByUsername(ctx, testUsername)

	// Assertions
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if resultUser != nil {
		t.Errorf("expected nil user, got %v", resultUser)
	}
}
