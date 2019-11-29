package project

import (
	"issue-tracker-backend/src/db"
	u "issue-tracker-backend/src/utils"
	"net/http"
	"strings"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

	project.ID = primitive.NewObjectID()

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

	currID, _ := primitive.ObjectIDFromHex(authenticatedUserID)
	filterMap := map[string]interface{}{
		"title": project.Title,
	}

	results, err := db.FindOne(DBConnection, "projects", filterMap, []Project{}, false)

	if err != nil {
		return u.Message(false, "Failed to retrieve project: "+err.Error()), http.StatusInternalServerError
	}

	_, ok := results.(*Project)
	if !ok {
		return u.Message(false, "Failed to cast interface to struct: "), http.StatusInternalServerError
	}

	resp := map[string]interface{}{}

	if results.(*Project).OwnerID == currID {
		resp["projectOwner"] = true
		results.(*Project).CreatedAt = results.(*Project).ID.Timestamp().String()
	} else {
		resp["projectOwner"] = false
	}
	resp["project"] = results
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

	results, err := db.FindMany(DBConnection, "projects", filterMap, []Project{}, 10)

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

	if resp, statusCode := newProject.ValidateUpdate(updatedProject, DBConnection); resp["status"] == false {
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

func TitleValidator(project *Project, DBConnection interface{}) (map[string]interface{}, int) {

	project.Title = strings.TrimSpace(project.Title)
	titleLen := utf8.RuneCountInString(project.Title)
	if titleLen < 1 || titleLen > 64 {
		return u.Message(false, "Title has to be between 1-64 characters long"), http.StatusBadRequest
	}

	filterMap := map[string]interface{}{
		"title": project.Title,
	}

	_, err := db.FindOne(DBConnection, "projects", filterMap, []Project{}, false)

	if err == mongo.ErrNoDocuments {
		return u.Message(true, ""), 0
	} else if err != nil {
		return u.Message(false, "Server error, please try again"), http.StatusInternalServerError
	}

	return u.Message(false, "Title has to be unique."), http.StatusBadRequest

}

func (project *Project) ValidateUpdate(updatedProject map[string]string, DBConnection interface{}) (map[string]interface{}, int) {

	if project.Title != "" {
		if resp, statusCode := TitleValidator(project, DBConnection); resp["status"] == false {
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
