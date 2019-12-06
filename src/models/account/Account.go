package account

import (
	"issue-tracker-backend/env"
	"issue-tracker-backend/src/db"
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
