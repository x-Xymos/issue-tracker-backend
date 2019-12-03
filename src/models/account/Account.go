package account

import (
	"issue-tracker-backend/env"
	"issue-tracker-backend/src/db"
	strLen "issue-tracker-backend/src/models/validators/stringLength"
	"issue-tracker-backend/src/models/validators/unique"
	v "issue-tracker-backend/src/models/validators/validator"
	u "issue-tracker-backend/src/utils"
	"net/http"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const collectionName = "accounts"

//Account : user account struct
type Account struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Username  string             `json:"username,omitempty"`
	Email     string             `json:"email,omitempty"`
	Password  string             `json:"password,omitempty"`
	CreatedAt string             `json:"createdAt,omitempty" bson:"-"`
	//https://willnorris.com/2014/05/go-rest-apis-and-pointers/
}

//Token : JWT token struct
type Token struct {
	UserID string `json:"UserID"`
	jwt.StandardClaims
}

type Validators struct {
	Username *[]*v.Function
	Email    *[]*v.Function //todo add an email validator
	Password *[]*v.Function
}

var validators Validators

func InitValidators(DBConnection interface{}) {
	validators.Username = v.Assign(
		v.Create(strLen.Validator, v.Options(strLen.Max(16), strLen.Min(3))),
		v.Create(unique.Validator, v.Options(unique.CaseSensitive(false), unique.Collection(collectionName), unique.Database(DBConnection), unique.SearchField("username"))))

	validators.Email = v.Assign(
		v.Create(strLen.Validator, v.Options(strLen.Max(32), strLen.Min(5))),
		v.Create(unique.Validator, v.Options(unique.CaseSensitive(false), unique.Collection(collectionName), unique.Database(DBConnection), unique.SearchField("email"))))

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

//Create : account creation
func (account *Account) Create(DBConnection interface{}) (map[string]interface{}, int) {

	if resp, statusCode := account.validateAccountCreation(); resp["status"] == false {
		return resp, statusCode
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	account.ID = db.NewID()
	account.Email = strings.ToLower(account.Email)

	err := db.InsertOne(DBConnection, collectionName, account)
	if err != nil {
		return u.Message(false, "Failed to create account: "+err.Error()), http.StatusInternalServerError
	}

	resp := u.Message(true, "Account has been created")
	return resp, http.StatusCreated
}

//Login : login
func (account *Account) Login(DBConnection interface{}) (map[string]interface{}, int) {

	account.Email = strings.ToLower(account.Email)

	filterMap := map[string]interface{}{
		"email": account.Email,
	}

	results, err := db.FindOne(DBConnection, collectionName, filterMap, nil, reflect.TypeOf(account).Elem(), false)
	if err != nil {
		return u.Message(false, "Failed to retrieve account: "+err.Error()), http.StatusInternalServerError
	}

	storedAccount := &Account{}

	storedAccount, ok := results.(*Account)
	if !ok {
		return u.Message(false, "Failed to cast interface to struct: "), http.StatusInternalServerError
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedAccount.Password), []byte(account.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again"), http.StatusUnauthorized
	}
	//Worked! Logged In
	account.Password = ""
	storedAccount.Password = ""

	//Create JWT token
	tk := &Token{UserID: storedAccount.ID.Hex()}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(env.TokenPassword))

	createdAt := storedAccount.ID.Timestamp().String()
	storedAccount.CreatedAt = createdAt

	resp := u.Message(true, "Logged In")
	resp["account"] = storedAccount
	resp["token"] = tokenString
	return resp, http.StatusOK
}

//Get : Retrieve account information
func (account *Account) Get(authenticatedUserID string, DBConnection interface{}) (map[string]interface{}, int) {

	currID, _ := primitive.ObjectIDFromHex(authenticatedUserID)

	filterMap := map[string]interface{}{
		"username": account.Username,
	}

	projectionMap := map[string]interface{}{
		"password": 0,
	}

	results, err := db.FindOne(DBConnection, collectionName, filterMap, projectionMap, reflect.TypeOf(account).Elem(), false)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return u.Message(false, "Account not found: "+err.Error()), http.StatusNotFound
		}
		return u.Message(false, "Failed to retrieve account: "+err.Error()), http.StatusInternalServerError
	}

	retrievedAccount := &Account{}
	retrievedAccount, ok := results.(*Account)
	if !ok {
		return u.Message(false, "Failed to cast interface to struct: "), http.StatusInternalServerError
	}

	resp := u.Message(true, "")

	if retrievedAccount.ID == currID {
		resp["accountOwner"] = true
		retrievedAccount.CreatedAt = retrievedAccount.ID.Timestamp().String()
	} else {
		resp["accountOwner"] = false
		retrievedAccount.Email = ""
	}
	resp["account"] = retrievedAccount
	return resp, http.StatusOK
}

//Update : Update account information
func (account *Account) Update(DBConnection interface{}) (map[string]interface{}, int) {

	updatedAccount := map[string]string{}

	if resp, statusCode := account.validateUpdate(updatedAccount); resp["status"] == false {
		return resp, statusCode
	}

	//this should take an array of validators
	filterMap := map[string]interface{}{
		"_id": account.ID,
	}

	updateMap := map[string]interface{}{
		"$set": updatedAccount,
	}

	err := db.UpdateOne(DBConnection, "accounts", filterMap, updateMap)
	if err != nil {
		return u.Message(false, "Failed to update: "+err.Error()), http.StatusInternalServerError
	}

	return u.Message(true, ""), http.StatusOK
}
