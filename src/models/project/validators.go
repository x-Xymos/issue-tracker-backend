package project

import (
	strLen "issue-tracker-backend/src/models/validators/stringLength"
	"issue-tracker-backend/src/models/validators/unique"
	v "issue-tracker-backend/src/models/validators/validator"
	u "issue-tracker-backend/src/utils"
	"net/http"
)

//Validators :
type Validators struct {
	Title *[]*v.Function
}

var validators Validators

//InitValidators :
func InitValidators(DBConnection interface{}) {
	validators.Title = v.Assign(
		v.Create(strLen.Validator, v.Options(strLen.Max(32))), v.Create(strLen.Validator, v.Options(strLen.Min(3))),
		v.Create(unique.Validator, unique.Options(DBConnection, collectionName, "title", false)))
}

func (project *Project) validateTitle() error {
	for _, vFunc := range *validators.Title {
		err := vFunc.Function(project.Title, vFunc.Options)
		if err != nil {
			return err
		}
	}
	return nil
}

//ValidateUpdate :
func (project *Project) validateUpdate(updatedProject map[string]string) (map[string]interface{}, int) {

	if project.Title != "" {
		if err := project.validateTitle(); err != nil {
			return map[string]interface{}{"status": false, "message": err.Error(), "badField": "title"}, http.StatusBadRequest
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
