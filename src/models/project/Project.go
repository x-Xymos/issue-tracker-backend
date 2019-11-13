package project

import (
	"context"
	u "issue-tracker-backend/src/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Project : Project struct
type Project struct {
	_id   primitive.ObjectID `bson:"_id, omitempty"`
	Title string             `json:"title"`
}

func (project *Project) _titleValidator(DBConn *mongo.Client) map[string]interface{} {

	if len(project.Title) < 1 || len(project.Title) > 64 {
		return u.Message(false, "Title has to be between 1-64 characters long")
	}

	tempProj := &Project{}
	collection := DBConn.Database("issue-tracker").Collection("projects")

	//Project Title must be unique
	projFilter := bson.D{{"title", project.Title}}

	err := collection.FindOne(context.TODO(), projFilter).Decode(&tempProj)
	if err != nil && err != mongo.ErrNoDocuments {
		return u.Message(false, "Connection error, please try again")
	}

	if tempProj.Title != "" {
		return u.Message(false, "Title has to be unique.")
	}

	return u.Message(true, "")
}

//Create :
func (project *Project) Create(DBConn *mongo.Client) map[string]interface{} {

	if resp := project._titleValidator(DBConn); resp["status"] == false {
		return resp
	}

	collection := DBConn.Database("issue-tracker").Collection("projects")

	_, err := collection.InsertOne(context.TODO(), project)
	if err != nil {
		return u.Message(false, "Failed to create project, connection error.")
	}

	response := u.Message(true, "Project has been created")
	response["project"] = project
	return response
}
