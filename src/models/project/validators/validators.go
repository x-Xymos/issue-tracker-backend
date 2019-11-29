package validators

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

func TitleValidator(project *Project, DBConnection *mongo.Database) (map[string]interface{}, int) {

	project.Title = strings.TrimSpace(project.Title)
	titleLen := utf8.RuneCountInString(project.Title)
	if titleLen < 1 || titleLen > 64 {
		return u.Message(false, "Title has to be between 1-64 characters long"), http.StatusBadRequest
	}

	filterMap := map[string]interface{}{
		"title": project.Title,
	}

	result, err := db.FindOne(DBConnection, "projects", filterMap, []Project{}, false)

	if err != nil && err != mongo.ErrNoDocuments {
		return u.Message(false, "Server error, please try again"), http.StatusInternalServerError
	}

	_, ok := result.(*Project)
	if !ok {
		return u.Message(false, "Failed to cast interface to struct: "), http.StatusInternalServerError
	}

	if result.(*Project).Title != "" {
		return u.Message(false, "Title has to be unique."), http.StatusBadRequest
	}

	return u.Message(true, ""), 0
}

func ValidateUpdate(project *Project, updatedProject map[string]string, DBConnection *mongo.Database) (map[string]interface{}, int) {

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
