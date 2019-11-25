package project

import (
	"context"
	u "issue-tracker-backend/src/utils"
	"net/http"
	"strings"
	"unicode/utf8"

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
	createdAt string             `json:"createdAt"`
}

func newProjectCollection(DBConn *mongo.Client) *mongo.Collection {
	return DBConn.Database("issue-tracker").Collection("projects")
}

func (project *Project) _titleValidator(DBConn *mongo.Client) (map[string]interface{}, int) {

	project.Title = strings.TrimSpace(project.Title)
	titleLen := utf8.RuneCountInString(project.Title)
	if titleLen < 1 || titleLen > 64 {
		return u.Message(false, "Title has to be between 1-64 characters long"), http.StatusBadRequest
	}

	tempProj := &Project{}
	collection := newProjectCollection(DBConn)

	//Project Title must be unique
	projFilter := bson.D{{"title", project.Title}}

	findOptions := options.FindOne().SetCollation(&options.Collation{Strength: 2, Locale: "en"})
	err := collection.FindOne(context.TODO(), projFilter, findOptions).Decode(&tempProj)

	if err != nil && err != mongo.ErrNoDocuments {
		return u.Message(false, "Connection error, please try again"), http.StatusInternalServerError
	}

	if tempProj.Title != "" {
		return u.Message(false, "Title has to be unique."), http.StatusBadRequest
	}

	return u.Message(true, ""), 0
}

//Create :
func (project *Project) Create(DBConn *mongo.Client) (map[string]interface{}, int) {

	if resp, statusCode := project._titleValidator(DBConn); resp["status"] == false {
		return resp, statusCode
	}

	collection := newProjectCollection(DBConn)
	project.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.TODO(), project)
	if err != nil {
		return u.Message(false, "Failed to create project: "+err.Error()), http.StatusInternalServerError
	}

	project.createdAt = project.ID.Timestamp().String()

	response := u.Message(true, "")
	response["project"] = project
	return response, http.StatusOK
}

//Get :
func (project *Project) Get(authenticatedUserID string, DBConn *mongo.Client) (map[string]interface{}, int) {

	collection := newProjectCollection(DBConn)

	currID, _ := primitive.ObjectIDFromHex(authenticatedUserID)

	var projFilter bson.D
	if project.Title != "" {
		projFilter = bson.D{{"title", project.Title}}
	} else {
		return u.Message(false, "Missing project title"), http.StatusBadRequest
	}

	findOptions := options.FindOne().SetCollation(&options.Collation{Strength: 2, Locale: "en"})

	err := collection.FindOne(context.TODO(), projFilter, findOptions).Decode(project)

	resp := map[string]interface{}{"project": project}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return u.Message(false, "Failed to retrieve project: Project not found"), http.StatusNotFound
		}
		return u.Message(false, "Failed to retrieve project: "+err.Error()), http.StatusInternalServerError
	}

	if project.OwnerID == currID {
		resp["projectOwner"] = true
		project.createdAt = project.ID.Timestamp().String()
	} else {
		resp["projectOwner"] = false
	}

	resp["status"] = true
	return resp, http.StatusOK
}

//GetAll :
func (project *Project) GetAll(lastID string, DBConn *mongo.Client) (map[string]interface{}, int) {

	collection := newProjectCollection(DBConn)

	var projectFilter bson.D

	if lastID != "" {
		_lastID, _ := primitive.ObjectIDFromHex(lastID)
		projectFilter = append(projectFilter,
			bson.E{"_id", bson.D{
				{"$gt", _lastID},
			},
			})
	}

	projectFilter = append(projectFilter, bson.E{"title", bson.D{
		{"$regex", primitive.Regex{Pattern: project.Title, Options: "i"}},
	}},
	)

	findOptions := options.Find()
	findOptions.SetLimit(10)
	var results []*Project

	cur, err := collection.Find(context.TODO(), projectFilter, findOptions)
	if err != nil {
		return u.Message(false, "Failed to retrieve projects: "+err.Error()), http.StatusInternalServerError
	}

	for cur.Next(context.TODO()) {

		var elem Project
		err := cur.Decode(&elem)
		if err != nil {
			return u.Message(false, "Failed to decode projects: "+err.Error()), http.StatusInternalServerError
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return u.Message(false, "Cursor error: "+err.Error()), http.StatusInternalServerError
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	resp := u.Message(true, "Success")
	resp["results"] = results
	return resp, http.StatusOK
	///https://arpitbhayani.me/blogs/fast-and-efficient-pagination-in-mongodb
}

func (project *Project) validateUpdate(updatedProject map[string]string, DBConn *mongo.Client) (map[string]interface{}, int) {

	if project.Title != "" {
		if resp, statusCode := project._titleValidator(DBConn); resp["status"] == false {
			return resp, statusCode
		}
		updatedProject["title"] = project.Title
	}

	for _, v := range updatedProject {
		if v != "" {
			return u.Message(true, ""), 0
		}
	}
	return u.Message(false, "No fields specified for the update or fields were identical"), http.StatusBadRequest
}

//Update :
func (project *Project) Update(newProject *Project, DBConn *mongo.Client) (map[string]interface{}, int) {

	updatedProject := map[string]string{}

	if resp, statusCode := newProject.validateUpdate(updatedProject, DBConn); resp["status"] == false {
		return resp, statusCode
	}

	collection := newProjectCollection(DBConn)

	projectFilter := bson.D{{"title", project.Title}, {"ownerID", project.OwnerID}}
	update := bson.D{
		{"$set", updatedProject},
	}
	_, err := collection.UpdateOne(context.TODO(), projectFilter, update) //this doesn't produce an error if the object that needs to be updated isn't found

	if err != nil {
		return u.Message(false, "Failed to update project: "+err.Error()), http.StatusInternalServerError
	}

	return u.Message(true, ""), http.StatusOK
}
