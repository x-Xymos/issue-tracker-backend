package signupcont

import (
	//"context"
	"context"
	"encoding/json"
	"log"
	"net/http"
	Service "todo-backend/src/servicetemplates"

	"go.mongodb.org/mongo-driver/bson"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func home(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")

	collection := Service.DB.Database("test").Collection("trainers")

	var result Trainer
	filter := bson.D{{"name", "Ash"}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Found a single document: %+v\n", result)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

var mainRoute = Service.RouteBinding{"/", home, "GET"}

var Routes = []Service.RouteBinding{mainRoute}

var ServiceName = "Signup-api"

var Port = "8081"

//fmt.Println(template.MyAsshole)
// router := mux.NewRouter().StrictSlash(true)

// router.HandleFunc("/", home)
// //router.HandleFunc("/events", getAllEvents).Methods("GET")
// //router.HandleFunc("/events/{id}", getEvent).Methods("GET")

// fmt.Println("Started Service on 0.0.0.0:8080")

// // Set client options
// // clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

// // // Connect to MongoDB
// // client, err := mongo.Connect(context.TODO(), clientOptions)
// // if err != nil {
// // 	log.Fatal(err)
// // }
// // // Check the connection
// // err = client.Ping(context.TODO(), nil)
// // if err != nil {
// // 	log.Fatal(err)
// // }

// // fmt.Println("Connected to MongoDB!")
// // collection := client.Database("test").Collection("trainers")

// // ash := Trainer{"Ash", 10, "Pallet Town"}
// // misty := Trainer{"Misty", 10, "Cerulean City"}
// // brock := Trainer{"Brock", 15, "Pewter City"}

// // insertResult, err := collection.InsertOne(context.TODO(), ash)
// // if err != nil {
// // 	log.Fatal(err)
// // }
// // fmt.Println("Inserted a single document: ", insertResult.InsertedID)

// // trainers := []interface{}{misty, brock}

// // insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
// // if err != nil {
// // 	log.Fatal(err)
// // }

// // fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

// // filter := bson.D{{"name", "Ash"}}

// // update := bson.D{
// // 	{"$inc", bson.D{
// // 		{"age", 1},
// // 	}},
// // }

// // updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
// // if err != nil {
// // 	log.Fatal(err)
// // }

// // fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

// log.Fatal(http.ListenAndServe(":8080", router))
