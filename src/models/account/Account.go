package account

import (
	"context"
	"fmt"
	"strings"
	"todo-backend/env"
	u "todo-backend/src/utils"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//Account : user account struct
type Account struct {
	UserID   primitive.ObjectID `bson:"_id, omitempty"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
}

//Token : JWT token struct
type Token struct {
	UserID string `json:"UserID"`
	jwt.StandardClaims
}

func emailValidator(email *string, DBConn *mongo.Client) bool {
	return true
}

//Validate incoming user details when they register for a new account
func (account *Account) Validate(DBConn *mongo.Client) (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Invalid email address"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password has to be at least 6 characters long"), false
	}

	tempAcc := &Account{}

	collection := DBConn.Database("todo").Collection("accounts")

	//Email must be unique
	emailFilter := bson.D{{"email", account.Email}}

	err := collection.FindOne(context.TODO(), emailFilter).Decode(&tempAcc)
	if err != nil && err != mongo.ErrNoDocuments {
		return u.Message(false, "Connection error, please try again"), false
	}

	if tempAcc.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	//Username must be unique
	usernameFilter := bson.D{{"username", account.Username}}

	err = collection.FindOne(context.TODO(), usernameFilter).Decode(&tempAcc)

	if err != mongo.ErrNoDocuments && err != nil {
		return u.Message(false, "Connection error, please try again"), false
	}

	if tempAcc.Username != "" {
		return u.Message(false, "Username already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

//Create : account creation
func (account *Account) Create(DBConn *mongo.Client) map[string]interface{} {

	if resp, ok := account.Validate(DBConn); !ok {
		return resp
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	account.UserID = primitive.NewObjectID()

	collection := DBConn.Database("todo").Collection("accounts")

	_, err := collection.InsertOne(context.TODO(), account)
	if err != nil {
		return u.Message(false, "Failed to create account, connection error.")
	}

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

//Login : login
func (account *Account) Login(DBConn *mongo.Client) map[string]interface{} {

	storedAccount := &Account{}

	collection := DBConn.Database("todo").Collection("accounts")

	emailFilter := bson.D{{"email", account.Email}}

	err := collection.FindOne(context.TODO(), emailFilter).Decode(&storedAccount)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return u.Message(false, "Email address doesn't match any accounts in our records, please try again")
		}
		return u.Message(false, "Connection error, please try again")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedAccount.Password), []byte(account.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &Token{UserID: storedAccount.UserID.Hex()}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(env.TokenPassword))

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	resp["token"] = tokenString
	return resp
}

//GetProfile : Retrieve user information
func (account *Account) GetProfile(DBConn *mongo.Client) map[string]interface{} {

	collection := DBConn.Database("todo").Collection("accounts")

	userIDFilter := bson.D{{"_id", account.UserID}}

	err := collection.FindOne(context.TODO(), userIDFilter).Decode(account)

	account.Password = ""
	resp := map[string]interface{}{"account": account}

	if err != nil {
		resp["status"] = false
		if err == mongo.ErrNoDocuments {
			return resp
		}
		return resp
	}

	resp["status"] = true
	return resp
}

//UpdateProfile : Retrieve user information
func (account *Account) UpdateProfile(DBConn *mongo.Client) map[string]interface{} {

	if resp, ok := account.Validate(DBConn); !ok {
		return resp
	}

	updatedAccount := map[string]string{
		"username": account.Username,
		"email":    account.Email,
	}

	if account.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
		updatedAccount["password"] = string(hashedPassword)
	}

	collection := DBConn.Database("todo").Collection("accounts")

	userIDFilter := bson.D{{"_id", account.UserID}}
	update := bson.D{
		{"$set", updatedAccount},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), userIDFilter, update)

	fmt.Println(updateResult)

	resp := map[string]interface{}{}

	if err != nil {
		resp["status"] = false
		if err == mongo.ErrNoDocuments {
			return resp
		}
		return resp
	}

	resp["status"] = true
	return resp
}
