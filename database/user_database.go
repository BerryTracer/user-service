package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	Disconnect() error
}

type UserMongoDatabase struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

// NewUserMongoDatabaseConnection returns a new UserMongoDatabase.
func NewUserMongoDatabaseConnection(connStr, databaseStr, collectionStr string) (*UserMongoDatabase, error) {
	clientOptions := options.Client().ApplyURI(connStr).SetMaxPoolSize(50)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	db := client.Database(databaseStr)
	collection := db.Collection(collectionStr)

	emailIndexModel := mongo.IndexModel{
		Keys:    map[string]int{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	usernameIndexModel := mongo.IndexModel{
		Keys:    map[string]int{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{emailIndexModel, usernameIndexModel})
	if err != nil {
		return nil, err
	}

	_, err = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    map[string]int{"email": 1, "username": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	return &UserMongoDatabase{Client: client, Collection: collection}, nil
}

// Disconnect implements Database.
func (d *UserMongoDatabase) Disconnect() error {
	return d.Client.Disconnect(context.Background())
}

// Ensure UserMongoDatabase implements Database interface.
var _ Database = &UserMongoDatabase{}
