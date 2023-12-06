package repository

import (
	"context"

	"github.com/BerryTracer/common-service/adapter"
	"github.com/BerryTracer/user-service/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserById(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, name string) (*model.User, error)
}

type UserMongoRepository struct {
	Collection adapter.MongoAdapter
}

// CreateUser implements UserRepository.
func (r *UserMongoRepository) CreateUser(ctx context.Context, user *model.User) error {
	userDB, err := user.ToUserDB()

	if err != nil {
		return err
	}

	_, err = r.Collection.InsertOne(ctx, userDB)

	if err != nil {
		return err
	}

	return nil
}

// GetUser implements UserRepository.
func (r *UserMongoRepository) GetUserById(ctx context.Context, id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var userDB model.UserDB
	err = r.Collection.FindOne(ctx, primitive.M{"_id": objectID}).Decode(&userDB)
	if err != nil {
		return nil, err
	}

	return userDB.ToUser(), nil
}

// GetUserByEmail implements UserRepository.
func (r *UserMongoRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var userDB model.UserDB
	err := r.Collection.FindOne(ctx, primitive.M{"email": email}).Decode(&userDB)
	if err != nil {
		return nil, err
	}

	return userDB.ToUser(), nil
}

// GetUserByUsername implements UserRepository.
func (r *UserMongoRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var userDB model.UserDB
	err := r.Collection.FindOne(ctx, primitive.M{"username": username}).Decode(&userDB)
	if err != nil {
		return nil, err
	}

	return userDB.ToUser(), nil
}

// Ensure UserMongoRepository implements the UserRepository interface
var _ UserRepository = &UserMongoRepository{}
