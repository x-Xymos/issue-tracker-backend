package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type Account struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func insert(db *mongo.Client) {
	collection := db.Database("test").Collection("accounts")

	info := Account{"Xymos4", "xymos4@gmail.com", "rootroot4", "testtoken4"}

	newAcc := Account{info.Username, info.Email, info.Password, info.Token}

	insertResult, err := collection.InsertOne(context.TODO(), newAcc)
	if err != nil {
		log.Fatal(err)
	}
	insertResultId := insertResult.InsertedID.(primitive.ObjectID).Hex()

	fmt.Println("Inserted a single document: ", insertResultId)

	//searchByObjectId(db, &insertResultId)

}

func searchByObjectId(db *mongo.Client, id *primitive.ObjectID) {

	collection := db.Database("test").Collection("accounts")

	//searching by object id
	result := &Account{}
	//objID, _ := primitive.ObjectIDFromHex("5d97999a4f55a15083dffd8f")
	filter := bson.D{{"_id", id}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println(err)
		} else {
			log.Fatal(err)
		}
	}
	fmt.Println(result)
}

func main() {

	//DBConn := Connect()

	//search(DBConn)
	//insert(DBConn)
	newAcc := Account{Username: "hehexd", Email: "email@email.com", Password: "root", Token: "asd"}
	fmt.Println(newAcc)
	// pass := "hehexd"
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	// Password := string(hashedPassword)

	// fmt.Println(Password)

	// pass2 := "hehexd2"
	// err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(pass2))
	// fmt.Println(err)

}
