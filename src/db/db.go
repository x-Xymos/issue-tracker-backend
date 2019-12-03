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

//NewID : Returns a new unique id
func NewID() primitive.ObjectID {
	return primitive.NewObjectID()
}

//FindMany : Returns many results based on supplied paramaters
func FindMany(DBConnection interface{}, collectionName string, filter map[string]interface{}, resultType reflect.Type, resultLimit int64) ([]*interface{}, error) {

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

	var results []*interface{}

	for cur.Next(context.TODO()) {
		elem := reflect.New(resultType).Interface()
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
func FindOne(DBConnection interface{}, collectionName string, filter map[string]interface{}, projection map[string]interface{}, resultType reflect.Type, caseSensitive bool) (interface{}, error) {

	var searchFilter bson.D
	err := mapToBSON(filter, &searchFilter)
	if err != nil {
		return nil, err
	}

	findOptions := options.FindOne()
	if !caseSensitive {
		findOptions = findOptions.SetCollation(&options.Collation{Strength: 2, Locale: "en"})
	}
	if projection != nil {
		var projectionDoc bson.D
		err = mapToBSON(projection, &projectionDoc)
		if err != nil {
			return nil, err
		}
		findOptions = findOptions.SetProjection(projectionDoc)
	}

	collection, err := getCollection(DBConnection, collectionName)
	if err != nil {
		return nil, err
	}

	singleResult := collection.FindOne(context.TODO(), searchFilter, findOptions)

	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}

	if resultType == nil {
		return nil, nil
	}
	result := reflect.New(resultType).Interface()
	singleResult.Decode(result)

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
