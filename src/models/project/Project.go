package project

import (
	"context"
	u "issue-tracker-backend/src/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Project : Project struct
type Project struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	OwnerID   primitive.ObjectID `bson:"ownerID, omitempty"`
	Title     string             `json:"title"`
	CreatedAt string             `json:"createdAt"`
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
func (project *Project) Create(authenticatedUserID string, DBConn *mongo.Client) map[string]interface{} {

	if resp := project._titleValidator(DBConn); resp["status"] == false {
		return resp
	}

	collection := DBConn.Database("issue-tracker").Collection("projects")
	project.ID = primitive.NewObjectID()
	project.OwnerID, _ = primitive.ObjectIDFromHex(authenticatedUserID)

	_, err := collection.InsertOne(context.TODO(), project)
	if err != nil {
		return u.Message(false, "Failed to create project: "+err.Error())
	}

	project.CreatedAt = project.ID.Timestamp().String()

	response := u.Message(true, "")
	response["project"] = project
	return response
}

//Get :
func (project *Project) Get(authenticatedUserID string, DBConn *mongo.Client) map[string]interface{} {

	collection := DBConn.Database("issue-tracker").Collection("projects")

	currID, _ := primitive.ObjectIDFromHex(authenticatedUserID)

	projectFilter := bson.D{{"_id", project.ID}}

	err := collection.FindOne(context.TODO(), projectFilter).Decode(project)

	resp := map[string]interface{}{"project": project}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return u.Message(false, "Failed to retrieve project: Project not found")
		}
		return u.Message(false, "Failed to retrieve project: "+err.Error())
	}

	if project.OwnerID == currID {
		resp["projectOwner"] = true
		project.CreatedAt = project.ID.Timestamp().String()
	} else {
		resp["projectOwner"] = false
	}

	resp["status"] = true
	return resp
}

//GetAll :
func (project *Project) GetAll(lastID string, DBConn *mongo.Client) map[string]interface{} {

	collection := DBConn.Database("issue-tracker").Collection("projects")

	projectFilter := bson.D{{}}

	if lastID != "" {
		_lastID, _ := primitive.ObjectIDFromHex(lastID)
		projectFilter = bson.D{
			{"_id", bson.D{
				{"$gt", _lastID},
			},
			}}
	}

	findOptions := options.Find()
	findOptions.SetLimit(10)
	var results []*Project

	cur, err := collection.Find(context.TODO(), projectFilter, findOptions)
	if err != nil {
		return u.Message(false, "Failed to retrieve projects: "+err.Error())
	}

	for cur.Next(context.TODO()) {

		var elem Project
		err := cur.Decode(&elem)
		if err != nil {
			return u.Message(false, "Failed to decode projects: "+err.Error())
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return u.Message(false, "Cursor error: "+err.Error())
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	resp := u.Message(true, "Success")
	resp["results"] = results
	return resp
	///https://arpitbhayani.me/blogs/fast-and-efficient-pagination-in-mongodb
}

func (project *Project) validateUpdate(updatedProject map[string]string, DBConn *mongo.Client) map[string]interface{} {

	if project.Title != "" {
		if resp := project._titleValidator(DBConn); resp["status"] == false {
			return resp
		}
		updatedProject["title"] = project.Title
	}

	return u.Message(false, "No fields specified for the update or fields were identical")
}

//Update :
func (project *Project) Update(DBConn *mongo.Client) map[string]interface{} {

	updatedProject := map[string]string{}

	if resp := project.validateUpdate(updatedProject, DBConn); resp["status"] == false {
		return resp
	}

	collection := DBConn.Database("issue-tracker").Collection("accounts")

	userIDFilter := bson.D{{"_id", project.ID}}
	update := bson.D{
		{"$set", updatedProject},
	}
	_, err := collection.UpdateOne(context.TODO(), userIDFilter, update)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return u.Message(false, "Failed to update project: project not found")
		}
		return u.Message(false, "Failed to update project: "+err.Error())
	}

	return u.Message(true, "")
}
