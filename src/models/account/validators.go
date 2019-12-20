package account

import (
	strMaxLength "issue-tracker-backend/src/models/validators/string/maxLength"
	strMinLength "issue-tracker-backend/src/models/validators/string/minLength"
	"issue-tracker-backend/src/models/validators/unique"
	v "issue-tracker-backend/src/models/validators/validator"
	u "issue-tracker-backend/src/utils"
	"net/http"
	"strings"
)

var usernameValidators v.Validators
var emailValidators v.Validators
var passwordValidators v.Validators

//InitValidators :
func InitValidators(DBConnection interface{}) {
	usernameValidators = v.Assign(
		v.Create(strMaxLength.Validator, strMaxLength.Options(16)),
		v.Create(strMinLength.Validator, strMinLength.Options(3)),
		v.Create(unique.Validator, unique.Options(DBConnection, collectionName, "username", false)))

	emailValidators = v.Assign(
		v.Create(strMaxLength.Validator, strMaxLength.Options(32)),
		v.Create(strMinLength.Validator, strMinLength.Options(5)),
		v.Create(unique.Validator, unique.Options(DBConnection, collectionName, "email", false)))

	passwordValidators = v.Assign(
		v.Create(strMaxLength.Validator, strMaxLength.Options(32)),
		v.Create(strMinLength.Validator, strMinLength.Options(6)))
}

func (account *Account) validateUsername() error {
	for _, validator := range *usernameValidators {
		err := validator.Validate(account.Username, validator.Options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (account *Account) validateEmail() error {
	for _, validator := range *emailValidators {
		err := validator.Validate(account.Email, validator.Options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (account *Account) validatePassword() error {
	for _, validator := range *passwordValidators {
		err := validator.Validate(account.Password, validator.Options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (account *Account) validateAccountCreation() (map[string]interface{}, int) {

	if err := account.validateEmail(); err != nil {
		return map[string]interface{}{"status": false, "message": err.Error(), "badField": "email"}, http.StatusBadRequest
	}

	if err := account.validateUsername(); err != nil {
		return map[string]interface{}{"status": false, "message": err.Error(), "badField": "username"}, http.StatusBadRequest
	}

	if err := account.validatePassword(); err != nil {
		return map[string]interface{}{"status": false, "message": err.Error(), "badField": "password"}, http.StatusBadRequest
	}

	return u.Message(true, ""), 0
}

func (account *Account) validateUpdate(updatedAccount map[string]string) (map[string]interface{}, int) {

	if account.Email != "" {
		if err := account.validateEmail(); err != nil {
			return map[string]interface{}{"status": false, "message": err.Error(), "badField": "email"}, http.StatusBadRequest
		}
		updatedAccount["email"] = strings.ToLower(account.Email)
	}

	if account.Username != "" {
		if err := account.validateUsername(); err != nil {
			return map[string]interface{}{"status": false, "message": err.Error(), "badField": "username"}, http.StatusBadRequest
		}
		updatedAccount["username"] = account.Username
	}

	for _, v := range updatedAccount {
		if v != "" {
			return u.Message(true, ""), 0
		}
	}
	return u.Message(false, "No fields specified for the update or fields were identical"), http.StatusBadRequest
}
