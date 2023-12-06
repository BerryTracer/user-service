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

	return &UserMongoDatabase{Client: client, Collection: collection}, nil
}

// Disconnect implements Database.
func (d *UserMongoDatabase) Disconnect() error {
	return d.Client.Disconnect(context.Background())
}

// Ensure UserMongoDatabase implements Database interface
var _ Database = &UserMongoDatabase{}
