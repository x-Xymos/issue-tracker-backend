package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
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
		log.Fatal(err)
		fmt.Println("Failed to establish a connection with the database")
	}

	return client

	// return client, err
}

type Trainer struct {
	Name string
	Age  int
	City string
}

func search(db *mongo.Client) {
	collection := db.Database("test").Collection("trainers")

	//var result Trainer
	filter := bson.D{{"name", "Ashs"}}

	err := collection.FindOne(context.TODO(), filter).Err()

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			fmt.Println(err)
			return
		} else {
			log.Fatal(err)
		}
	}

}

func main() {

	DBConn := Connect()
	search(DBConn)
	//collection := DBConn.Database("test").Collection("trainers")

	// var result Trainer
	// filter := bson.D{{"name", "Ashs"}}

	// err := collection.FindOne(context.TODO(), filter).Err()

	// if err != nil {
	// 	if err.Error() == "mongo: no documents in result" {
	// 		fmt.Println(err)
	// 		return
	// 	} else {
	// 		log.Fatal(err)
	// 	}
	// }

	//fmt.Printf("%+v\n", result)

}
