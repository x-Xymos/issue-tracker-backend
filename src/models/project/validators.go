package project

import (
	strMaxLength "issue-tracker-backend/src/models/validators/string/maxLength"
	strMinLength "issue-tracker-backend/src/models/validators/string/minLength"
	"issue-tracker-backend/src/models/validators/unique"
	v "issue-tracker-backend/src/models/validators/validator"
	u "issue-tracker-backend/src/utils"
	"net/http"
)

var titleValidators v.Validators

//InitValidators :
func InitValidators(DBConnection interface{}) {
	titleValidators = v.Assign(
		v.Create(strMaxLength.Validator, strMaxLength.Options(32)),
		v.Create(strMinLength.Validator, strMinLength.Options(3)),
		v.Create(unique.Validator, unique.Options(DBConnection, collectionName, "title", false)))
}

func (project *Project) validateTitle() error {
	for _, validator := range *titleValidators {
		err := validator.Validate(project.Title, validator.Options)
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
