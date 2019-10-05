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
	UserID   primitive.ObjectID `bson:"_id, omitempty"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

func insert(db *mongo.Client) {
	collection := db.Database("test").Collection("accounts")

	newAcc := Account{UserID: primitive.NewObjectID(), Username: "jbpratt", Email: "jbpratt@gmail.com", Password: "rootroot4"}
	//newAcc := Account{}
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

	filter := bson.D{{"email", "jbpratt@gmail.com"}}

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

	DBConn := Connect()

	objID, _ := primitive.ObjectIDFromHex("5d97999a4f55a15083dffd8f")
	searchByObjectId(DBConn, &objID)

	//insert(DBConn)
	//newAcc := Account{Username: "hehexd", Email: "email@email.com", Password: "root", Token: "asd"}
	//fmt.Println(newAcc)
	// pass := "hehexd"
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	// Password := string(hashedPassword)

	// fmt.Println(Password)

	// pass2 := "hehexd2"
	// err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(pass2))
	// fmt.Println(err)

}
