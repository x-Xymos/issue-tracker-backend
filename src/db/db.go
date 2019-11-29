package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func mapToBSON(data map[string]interface{}, bsonDoc *bson.D) error {
	for key, value := range data {

		if key == "$regex" {
			pattern, ok := value.(map[string]string)["pattern"]
			if !ok {
				return errors.New("Error in mapToBSON, invalid map type")
			}

			options, ok := value.(map[string]string)["options"]
			if !ok {
				return errors.New("Error in mapToBSON, invalid map type")
			}
			*bsonDoc = append(*bsonDoc, bson.E{key, primitive.Regex{Pattern: pattern, Options: options}})
			continue
		}

		if reflect.TypeOf(value).String() == "map[string]interface {}" {
			var bsonD bson.D

			valueCast, ok := value.(map[string]interface{})
			if !ok {
				return errors.New("Error in mapToBSON, invalid map type")
			}

			mapToBSON(valueCast, &bsonD)
			*bsonDoc = append(*bsonDoc, bson.E{key, bsonD})
			continue
		}
		*bsonDoc = append(*bsonDoc, bson.E{key, value})
	}
	return nil
}

func getCollection(DBConnection interface{}, collectionName string) (*mongo.Collection, error) {
	db, ok := DBConnection.(*mongo.Database)
	if !ok {
		return nil, errors.New("Error in getCollection, invalid database, expecting mongodb")
	}
	return db.Collection(collectionName), nil
}

//FindMany : Returns many results based on supplied paramaters
// resultStruct : interface{}
// an empty struct used to dynamically create a struct of that type to decode the results into
// has to be passed in as an array in order for reflect.TypeOf to work, i.e FindMany([]Struct{}...)
//
func FindMany(DBConnection interface{}, collectionName string, filter map[string]interface{}, resultStruct interface{}, resultLimit int64) ([]*interface{}, error) {

	var searchFilter bson.D
	err := mapToBSON(filter, &searchFilter)
	if err != nil {
		return nil, err
	}

	findOptions := options.Find()
	findOptions.SetLimit(resultLimit)

	collection, err := getCollection(DBConnection, collectionName)
	if err != nil {
		return nil, err
	}

	cur, err := collection.Find(context.TODO(), searchFilter, findOptions)
	if err != nil {
		return nil, err
	}

	aType := reflect.TypeOf(resultStruct).Elem()

	var results []*interface{}

	for cur.Next(context.TODO()) {
		elem := reflect.New(aType).Interface()
		err := cur.Decode(elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	return results, nil
}

//FindOne : Returns one result based on search parameters
func FindOne(DBConnection interface{}, collectionName string, filter map[string]interface{}, resultStruct interface{}, caseSensitive bool) (interface{}, error) {

	var searchFilter bson.D
	err := mapToBSON(filter, &searchFilter)
	if err != nil {
		return nil, err
	}

	findOptions := options.FindOne()
	if !caseSensitive {
		findOptions = findOptions.SetCollation(&options.Collation{Strength: 2, Locale: "en"})
	}

	collection, err := getCollection(DBConnection, collectionName)
	if err != nil {
		return nil, err
	}

	aType := reflect.TypeOf(resultStruct).Elem()
	result := reflect.New(aType).Interface()

	err = collection.FindOne(context.TODO(), searchFilter, findOptions).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//InsertOne : Inserts one result
func InsertOne(DBConnection interface{}, collectionName string, insertDocument interface{}) error {

	collection, err := getCollection(DBConnection, collectionName)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(context.TODO(), insertDocument)
	if err != nil {
		return err
	}

	return nil
}

//UpdateOne : Update one result
func UpdateOne(DBConnection interface{}, collectionName string, filter map[string]interface{}, update map[string]interface{}) error {

	var searchFilter bson.D
	err := mapToBSON(filter, &searchFilter)
	if err != nil {
		return err
	}

	var updateDoc bson.D
	err = mapToBSON(update, &updateDoc)
	if err != nil {
		return err
	}

	collection, err := getCollection(DBConnection, collectionName)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.TODO(), searchFilter, updateDoc) //this doesn't produce an error if the object that needs to be updated isn't found

	if err != nil {
		return err
	}

	return nil
}

// Connect to MongoDB
func Connect(DBName *string) *mongo.Database {
	//Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		fmt.Println("Failed to establish a connection with the database")
		log.Fatal(err)
	}

	return client.Database(*DBName)
}
