package account

import (
	"context"
	"issue-tracker-backend/env"
	"issue-tracker-backend/src/db"
	u "issue-tracker-backend/src/utils"
	"net/http"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var collectionName = "accounts"

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

func NewAccountCollection(DBConnection interface{}) *mongo.Collection {
	return DBConnection.(*mongo.Database).Collection("accounts")
}

func (account *Account) _emailValidator(DBConnection interface{}) (map[string]interface{}, int) {

	account.Email = strings.ReplaceAll(account.Email, " ", "")
	if !strings.Contains(account.Email, "@") || utf8.RuneCountInString(account.Email) < 5 {
		return u.Message(false, "Invalid email address"), http.StatusBadRequest
	}

	tempAcc := &Account{}
	collection := NewAccountCollection(DBConnection)

	//Email must be unique
	emailFilter := bson.D{{"email", account.Email}}

	findOptions := options.FindOne().SetCollation(&options.Collation{Strength: 2, Locale: "en"})
	err := collection.FindOne(context.TODO(), emailFilter, findOptions).Decode(&tempAcc)
	if err != nil && err != mongo.ErrNoDocuments {
		return u.Message(false, "Connection error, please try again"), http.StatusInternalServerError
	}

	if tempAcc.Email != "" {
		return u.Message(false, "Email address already in use by another user"), http.StatusBadRequest
	}
	return u.Message(true, ""), 0
}

func (account *Account) _usernameValidator(DBConnection interface{}) (map[string]interface{}, int) {

	account.Username = strings.ReplaceAll(account.Username, " ", "")

	if utf8.RuneCountInString(account.Username) < 3 {
		return u.Message(false, "Username has to be at least 3 characters long"), http.StatusBadRequest
	}

	tempAcc := &Account{}
	collection := NewAccountCollection(DBConnection)
	//Username must be unique
	usernameFilter := bson.D{{"username", account.Username}}

	findOptions := options.FindOne().SetCollation(&options.Collation{Strength: 2, Locale: "en"})
	err := collection.FindOne(context.TODO(), usernameFilter, findOptions).Decode(&tempAcc)

	if err != mongo.ErrNoDocuments && err != nil {
		return u.Message(false, "Connection error, please try again"), http.StatusInternalServerError
	}

	if tempAcc.Username != "" {
		return u.Message(false, "Username already in use by another user"), http.StatusBadRequest
	}
	return u.Message(true, ""), 0
}

func (account *Account) _passwordValidator() (map[string]interface{}, int) {

	account.Password = strings.ReplaceAll(account.Password, " ", "")

	if utf8.RuneCountInString(account.Password) < 6 {
		return u.Message(false, "Password has to be at least 6 characters long"), http.StatusBadRequest
	}
	return u.Message(true, ""), 0

}

func (account *Account) validateAccountCreation(DBConnection interface{}) (map[string]interface{}, int) {

	if resp, statusCode := account._emailValidator(DBConnection); resp["status"] == false {
		return resp, statusCode
	}

	if resp, statusCode := account._usernameValidator(DBConnection); resp["status"] == false {
		return resp, statusCode
	}

	return account._passwordValidator()
}

func (account *Account) validateUpdate(updatedAccount map[string]string, DBConnection interface{}) (map[string]interface{}, int) {

	if account.Email != "" {
		if resp, statusCode := account._emailValidator(DBConnection); resp["status"] == false {
			return resp, statusCode
		}
		updatedAccount["email"] = strings.ToLower(account.Email)
	}

	if account.Username != "" {
		if resp, statusCode := account._usernameValidator(DBConnection); resp["status"] == false {
			return resp, statusCode
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

	if resp, statusCode := account.validateAccountCreation(DBConnection); resp["status"] == false {
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

	if resp, statusCode := account.validateUpdate(updatedAccount, DBConnection); resp["status"] == false {
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
