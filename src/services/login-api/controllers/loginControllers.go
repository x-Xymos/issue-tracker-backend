package logincont

import (
	//"context"
	"encoding/json"
	"net/http"
	Service "todo-backend/src/servicetemplates"
	//"strconv"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

// type event struct {
// 	ID          string `json:"ID"`
// 	Title       string `json:"Title"`
// 	Description string `json:"Description"`
// }

// type Trainer struct {
// 	Name string
// 	Age  int
// 	City string
// }

// type allEvents []event

// var events = allEvents{
// 	{
// 		ID:          "1",
// 		Title:       "Introduction to Golang",
// 		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
// 	},
// }

// func getEvent(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	eventID := mux.Vars(r)["id"]
// 	id, _ := strconv.Atoi(eventID)
// 	json.NewEncoder(w).Encode(events[id])

// }

// func getAllEvents(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(events)
// }

type msg struct {
	Title string `json:"Title"`
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(msg{Title: "This is the Login-Api"})
}

var mainRoute = Service.RouteBinding{"/api/", home, "GET"}

var Routes = []Service.RouteBinding{mainRoute}

var ServiceName = "Login-api"

var Port = "8080"

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
