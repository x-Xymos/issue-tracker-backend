package account

import (
	strLen "issue-tracker-backend/src/models/validators/stringLength"
	"issue-tracker-backend/src/models/validators/unique"
	v "issue-tracker-backend/src/models/validators/validator"
	u "issue-tracker-backend/src/utils"
	"net/http"
	"strings"
)

type Validators struct {
	Username *[]*v.Function
	Email    *[]*v.Function //todo add an email validator
	Password *[]*v.Function
}

var validators Validators

//InitValidators :
func InitValidators(DBConnection interface{}) {
	validators.Username = v.Assign(
		v.Create(strLen.Validator, v.Options(strLen.Max(16), strLen.Min(3))),
		v.Create(unique.Validator, unique.Options(DBConnection, collectionName, "username", false)))

	validators.Email = v.Assign(
		v.Create(strLen.Validator, v.Options(strLen.Max(32), strLen.Min(5))),
		v.Create(unique.Validator, unique.Options(DBConnection, collectionName, "email", false)))

	validators.Password = v.Assign(
		v.Create(strLen.Validator, v.Options(strLen.Max(32), strLen.Min(6))))

}

//todo can we make a function for the loop part that just takes in those params???
func (account *Account) validateUsername() error {
	for _, vFunc := range *validators.Username {
		err := vFunc.Function(account.Username, vFunc.Options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (account *Account) validateEmail() error {
	for _, vFunc := range *validators.Email {
		err := vFunc.Function(account.Email, vFunc.Options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (account *Account) validatePassword() error {
	for _, vFunc := range *validators.Password {
		err := vFunc.Function(account.Password, vFunc.Options)
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
