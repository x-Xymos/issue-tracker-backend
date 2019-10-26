package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to MongoDB
func Connect() *mongo.Client {
	//Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		fmt.Println("Failed to establish a connection with the database")
	}

	return client

	// return client, err
}
