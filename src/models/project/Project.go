package project

import (
	"issue-tracker-backend/src/db"
	u "issue-tracker-backend/src/utils"
	"net/http"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "projects"

//Project : Project struct
type Project struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OwnerID   primitive.ObjectID `json:"ownerID,omitempty" bson:"ownerID,omitempty"`
	Title     string             `json:"title"`
	CreatedAt string             `json:"createdAt,omitempty" bson:"-"`
}

//Create :
func (project *Project) Create(DBConnection interface{}) (map[string]interface{}, int) {

	// if resp, statusCode := v._titleValidator(DB); resp["status"] == false {
	// 	return resp, statusCode
	// }
	////this should take an array of validators

	project.ID = db.NewID()

	err := db.InsertOne(DBConnection, "projects", project)
	if err != nil {
		return u.Message(false, err.Error()), http.StatusInternalServerError
	}
	project.CreatedAt = project.ID.Timestamp().String()

	response := u.Message(true, "")
	response["project"] = project
	return response, http.StatusOK
}

//Get :
func (project *Project) Get(authenticatedUserID string, DBConnection interface{}) (map[string]interface{}, int) {

	if project.Title == "" {
		return u.Message(false, "Missing project title"), http.StatusBadRequest
	}

	filterMap := map[string]interface{}{
		"title": project.Title,
	}

	results, err := db.FindOne(DBConnection, collectionName, filterMap, nil, reflect.TypeOf(project).Elem(), false)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return u.Message(false, "Project not found: "+err.Error()), http.StatusNotFound
		}
		return u.Message(false, "Failed to retrieve project: "+err.Error()), http.StatusInternalServerError
	}

	retrievedProject := &Project{}

	retrievedProject, ok := results.(*Project)
	if !ok {
		return u.Message(false, "Failed to cast interface to struct: "), http.StatusInternalServerError
	}

	resp := map[string]interface{}{}

	currID, _ := primitive.ObjectIDFromHex(authenticatedUserID)
	if retrievedProject.OwnerID == currID {
		resp["projectOwner"] = true
		retrievedProject.CreatedAt = retrievedProject.ID.Timestamp().String()
	} else {
		resp["projectOwner"] = false
	}
	resp["project"] = retrievedProject
	resp["status"] = true

	return resp, http.StatusOK
}

//GetAll :
func (project *Project) GetAll(lastID string, DBConnection interface{}) (map[string]interface{}, int) {

	_lastID, _ := primitive.ObjectIDFromHex(lastID)
	filterMap := map[string]interface{}{
		"_id":   map[string]interface{}{"$gt": _lastID},
		"title": map[string]interface{}{"$regex": map[string]string{"pattern": project.Title, "options": "i"}},
	}

	results, err := db.FindMany(DBConnection, "projects", filterMap, reflect.TypeOf(project).Elem(), 10)

	if err != nil {
		return u.Message(false, err.Error()), http.StatusInternalServerError
	}

	resp := u.Message(true, "Success")
	resp["results"] = results
	return resp, http.StatusOK
	///https://arpitbhayani.me/blogs/fast-and-efficient-pagination-in-mongodb
}

//Update :
func (project *Project) Update(newProject *Project, DBConnection interface{}) (map[string]interface{}, int) {

	updatedProject := map[string]string{}

	if resp, statusCode := newProject.validateUpdate(updatedProject); resp["status"] == false {
		return resp, statusCode
	}
	//this should take an array of validators
	filterMap := map[string]interface{}{
		"title":   project.Title,
		"ownerID": project.OwnerID,
	}

	updateMap := map[string]interface{}{
		"$set": updatedProject,
	}

	err := db.UpdateOne(DBConnection, "projects", filterMap, updateMap)
	if err != nil {
		return u.Message(false, "Failed to update project: "+err.Error()), http.StatusInternalServerError
	}

	return u.Message(true, ""), http.StatusOK
}
