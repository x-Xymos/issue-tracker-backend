package signupcont

import (
	//"context"

	"encoding/json"
	"net/http"
	AccountModel "todo-backend/src/models/account"
	Service "todo-backend/src/servicetemplates"
	u "todo-backend/src/utils"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func home(w http.ResponseWriter, r *http.Request) {

	//w.Header().Add("Content-Type", "application/json")
	//json.NewEncoder(w).Encode("{Title: "This is the Signup-Api"}")
	u.Respond(w, u.Message(false, "This is the Signup-Api"))
}

func signup(w http.ResponseWriter, r *http.Request) {

	account := &AccountModel.Account{}

	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create(Service.DBConn) //Create account
	u.Respond(w, resp)

	// w.Header().Add("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(result)
}

//Routes : an array of route bindings
var Routes = []Service.RouteBinding{
	Service.RouteBinding{"/api/", home, "GET"},
	Service.RouteBinding{"/api/signup", signup, "POST"}}

//ServiceName : service name
var ServiceName = "Signup-api"

//Port : service port
var Port = "8081"

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
