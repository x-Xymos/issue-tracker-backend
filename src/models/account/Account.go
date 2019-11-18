package account

import (
	"context"
	"issue-tracker-backend/env"
	u "issue-tracker-backend/src/utils"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//Account : user account struct
type Account struct {
	UserID    primitive.ObjectID `bson:"_id, omitempty"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	CreatedAt string             `json:"createdAt"`
}

//Token : JWT token struct
type Token struct {
	UserID string `json:"UserID"`
	jwt.StandardClaims
}

func (account *Account) _emailValidator(DBConn *mongo.Client) (map[string]interface{}, int) {

	if !strings.Contains(account.Email, "@") || len(account.Email) < 5 {
		return u.Message(false, "Invalid email address"), http.StatusBadRequest
	}

	tempAcc := &Account{}
	collection := DBConn.Database("issue-tracker").Collection("accounts")

	//Email must be unique
	emailFilter := bson.D{{"email", account.Email}}

	err := collection.FindOne(context.TODO(), emailFilter).Decode(&tempAcc)
	if err != nil && err != mongo.ErrNoDocuments {
		return u.Message(false, "Connection error, please try again"), http.StatusInternalServerError
	}

	if tempAcc.Email != "" {
		return u.Message(false, "Email address already in use by another user."), http.StatusBadRequest
	}
	return u.Message(true, ""), 0
}

func (account *Account) _usernameValidator(DBConn *mongo.Client) (map[string]interface{}, int) {

	if len(account.Username) < 3 {
		return u.Message(false, "Username has to be at least 3 characters long"), http.StatusBadRequest
	}

	tempAcc := &Account{}
	collection := DBConn.Database("issue-tracker").Collection("accounts")
	//Username must be unique
	usernameFilter := bson.D{{"username", account.Username}}

	err := collection.FindOne(context.TODO(), usernameFilter).Decode(&tempAcc)

	if err != mongo.ErrNoDocuments && err != nil {
		return u.Message(false, "Connection error, please try again"), http.StatusInternalServerError
	}

	if tempAcc.Username != "" {
		return u.Message(false, "Username already in use by another user."), http.StatusBadRequest
	}
	return u.Message(true, ""), 0
}

func (account *Account) _passwordValidator() (map[string]interface{}, int) {

	if len(account.Password) < 6 {
		return u.Message(false, "Password has to be at least 6 characters long"), http.StatusBadRequest
	}
	return u.Message(true, ""), 0

}

func (account *Account) validateAccountCreation(DBConn *mongo.Client) (map[string]interface{}, int) {

	if resp, statusCode := account._emailValidator(DBConn); resp["status"] == false {
		return resp, statusCode
	}

	if resp, statusCode := account._usernameValidator(DBConn); resp["status"] == false {
		return resp, statusCode
	}

	return account._passwordValidator()
}

//Create : account creation
func (account *Account) Create(DBConn *mongo.Client) (map[string]interface{}, int) {

	if resp, statusCode := account.validateAccountCreation(DBConn); resp["status"] == false {
		return resp, statusCode
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	account.UserID = primitive.NewObjectID()
	account.Email = strings.ToLower(account.Email)
	collection := DBConn.Database("issue-tracker").Collection("accounts")

	_, err := collection.InsertOne(context.TODO(), account)
	if err != nil {
		return u.Message(false, "Failed to create account: "+err.Error()), http.StatusInternalServerError
	}

	account.Password = "" //delete password

	resp := u.Message(true, "Account has been created")
	resp["account"] = account
	return resp, http.StatusCreated
}

//Login : login
func (account *Account) Login(DBConn *mongo.Client) map[string]interface{} {

	storedAccount := &Account{}

	collection := DBConn.Database("issue-tracker").Collection("accounts")

	account.Email = strings.ToLower(account.Email)
	emailFilter := bson.D{{"email", account.Email}}

	err := collection.FindOne(context.TODO(), emailFilter).Decode(&storedAccount)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return u.Message(false, "Failed to log in: Email address doesn't match any accounts in our records, please try again")
		}
		return u.Message(false, "Failed to log in: "+err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedAccount.Password), []byte(account.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""
	storedAccount.Password = ""

	//Create JWT token
	tk := &Token{UserID: storedAccount.UserID.Hex()}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(env.TokenPassword))

	storedAccount.CreatedAt = storedAccount.UserID.Timestamp().String()

	resp := u.Message(true, "Logged In")
	resp["account"] = storedAccount
	resp["token"] = tokenString
	return resp
}

//Get : Retrieve account information
func (account *Account) Get(authenticatedUserID string, DBConn *mongo.Client) (map[string]interface{}, int) {

	collection := DBConn.Database("issue-tracker").Collection("accounts")

	currID, _ := primitive.ObjectIDFromHex(authenticatedUserID)

	userFilter := bson.D{{}}
	if account.Username != "" {
		userFilter = bson.D{{"username", account.Username}}
	} else {
		userFilter = bson.D{{"_id", currID}}
	}

	err := collection.FindOne(context.TODO(), userFilter).Decode(account)
	account.Password = ""

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return u.Message(false, "Failed to retrieve account: user not found"), http.StatusNotFound
		}
		return u.Message(false, "Failed to retrieve account: "+err.Error()), http.StatusInternalServerError
	}

	resp := u.Message(true, "")

	if account.UserID == currID {
		resp["accountOwner"] = true
		account.CreatedAt = account.UserID.Timestamp().String()
	} else {
		resp["accountOwner"] = false
		account.Email = ""
	}
	resp["account"] = account
	return resp, http.StatusOK
}

func (account *Account) validateUpdate(updatedAccount map[string]string, DBConn *mongo.Client) (map[string]interface{}, int) {

	if account.Email != "" {
		if resp, statusCode := account._emailValidator(DBConn); resp["status"] == false {
			return resp, statusCode
		}
		updatedAccount["email"] = strings.ToLower(account.Email)
	}

	if account.Username != "" {
		if resp, statusCode := account._usernameValidator(DBConn); resp["status"] == false {
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

//Update : Update account information
func (account *Account) Update(DBConn *mongo.Client) (map[string]interface{}, int) {

	updatedAccount := map[string]string{}

	if resp, statusCode := account.validateUpdate(updatedAccount, DBConn); resp["status"] == false {
		return resp, statusCode
	}

	collection := DBConn.Database("issue-tracker").Collection("accounts")

	userIDFilter := bson.D{{"_id", account.UserID}}
	update := bson.D{
		{"$set", updatedAccount},
	}
	_, err := collection.UpdateOne(context.TODO(), userIDFilter, update)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return u.Message(false, "Failed to update account: account not found"), http.StatusNotFound
		}
		return u.Message(false, "Failed to update account: "+err.Error()), http.StatusInternalServerError
	}
	return u.Message(true, ""), http.StatusOK
}
