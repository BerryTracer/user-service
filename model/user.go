package model

import (
	"errors"
	"regexp"

	userservice "github.com/BerryTracer/user-service/grpc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             string
	Username       string
	Email          string
	HashedPassword string
}

type UserDB struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username       string             `bson:"username" json:"username"`
	Email          string             `bson:"email" json:"email"`
	HashedPassword string             `bson:"hashed_password" json:"hashed_password"`
}

// NewUser creates a new User instance
func NewUser(username, email, hashedPassword string) *User {
	return &User{
		ID:             primitive.NewObjectID().Hex(),
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
	}
}

// ToUserDB converts a User domain model to a UserDB database model
func (u *User) ToUserDB() (*UserDB, error) {
	var id primitive.ObjectID
	var err error

	if u.ID != "" {
		id, err = primitive.ObjectIDFromHex(u.ID)
		if err != nil {
			return nil, err
		}
	}

	return &UserDB{
		ID:             id,
		Username:       u.Username,
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
	}, nil
}

// ToUser converts a UserDB database model to a User domain model
func (udb *UserDB) ToUser() *User {
	return &User{
		ID:             udb.ID.Hex(),
		Username:       udb.Username,
		Email:          udb.Email,
		HashedPassword: udb.HashedPassword,
	}
}

// ConvertToProto converts a User domain model to a User proto model
func (u *User) ConvertToProto() *userservice.User {
	return &userservice.User{
		Id:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
	}
}

// Validate checks if the user's fields meet basic requirements
func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	if u.HashedPassword == "" {
		return errors.New("hashed password is required")
	}
	// Further validation logic goes here...
	return nil
}

// isValidEmail validates the email format
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
