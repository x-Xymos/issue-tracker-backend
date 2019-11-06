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

func (account *Account) _emailValidator(DBConn *mongo.Client) map[string]interface{} {

	if !strings.Contains(account.Email, "@") || len(account.Email) < 5 {
		return u.Message(false, "Invalid email address")
	}

	tempAcc := &Account{}
	collection := DBConn.Database("todo").Collection("accounts")

	//Email must be unique
	emailFilter := bson.D{{"email", account.Email}}

	err := collection.FindOne(context.TODO(), emailFilter).Decode(&tempAcc)
	if err != nil && err != mongo.ErrNoDocuments {
		return u.Message(false, "Connection error, please try again")
	}

	if tempAcc.Email != "" {
		return u.Message(false, "Email address already in use by another user.")
	}
	return u.Message(true, "")
}

func (account *Account) _usernameValidator(DBConn *mongo.Client) map[string]interface{} {

	if len(account.Username) < 3 {
		return u.Message(false, "Username has to be at least 3 characters long")
	}

	tempAcc := &Account{}
	collection := DBConn.Database("todo").Collection("accounts")
	//Username must be unique
	usernameFilter := bson.D{{"username", account.Username}}

	err := collection.FindOne(context.TODO(), usernameFilter).Decode(&tempAcc)

	if err != mongo.ErrNoDocuments && err != nil {
		return u.Message(false, "Connection error, please try again")
	}

	if tempAcc.Username != "" {
		return u.Message(false, "Username already in use by another user.")
	}
	return u.Message(true, "")
}

func (account *Account) _passwordValidator() map[string]interface{} {

	if len(account.Password) < 6 {
		return u.Message(false, "Password has to be at least 6 characters long")
	}
	return u.Message(true, "")

}

//ValidateAccountCreation : validate incoming user details in the account.Create method
func (account *Account) ValidateAccountCreation(DBConn *mongo.Client) map[string]interface{} {

	if resp := account._emailValidator(DBConn); resp["status"] == false {
		return resp
	}

	if resp := account._usernameValidator(DBConn); resp["status"] == false {
		return resp
	}

	return account._passwordValidator()
}

//Create : account creation
func (account *Account) Create(DBConn *mongo.Client) map[string]interface{} {

	if resp := account.ValidateAccountCreation(DBConn); resp["status"] == false {
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

//ValidateProfileUpdate : validate incoming user details in the account.UpdateProfile method
func (account *Account) ValidateProfileUpdate(updatedAccount map[string]string, DBConn *mongo.Client) map[string]interface{} {

	if account.Email != "" {
		if resp := account._emailValidator(DBConn); resp["status"] == false {
			return resp
		}
		updatedAccount["email"] = account.Email
	}

	if account.Username != "" {
		if resp := account._usernameValidator(DBConn); resp["status"] == false {
			return resp
		}
		updatedAccount["username"] = account.Username
	}

	for _, v := range updatedAccount {
		if v != "" {
			return u.Message(true, "")
		}
	}
	return u.Message(false, "No fields specified for the update or fields were identical")
}

//UpdateProfile : Retrieve user information
func (account *Account) UpdateProfile(DBConn *mongo.Client) map[string]interface{} {

	updatedAccount := map[string]string{}

	if resp := account.ValidateProfileUpdate(updatedAccount, DBConn); resp["status"] == false {
		return resp
	}

	collection := DBConn.Database("todo").Collection("accounts")

	userIDFilter := bson.D{{"_id", account.UserID}}
	update := bson.D{
		{"$set", updatedAccount},
	}
	_, err := collection.UpdateOne(context.TODO(), userIDFilter, update)
	fmt.Println(err)
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
