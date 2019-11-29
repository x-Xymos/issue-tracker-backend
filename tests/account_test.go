package main

import (
	"context"
	"fmt"
	"issue-tracker-backend/src/auth"
	AccountModel "issue-tracker-backend/src/models/account"
	Controller "issue-tracker-backend/src/services/account-api/controllers"
	Service "issue-tracker-backend/src/servicetemplates"
	"issue-tracker-backend/src/servicetemplates/db"
	u "issue-tracker-backend/src/utils"
	tu "issue-tracker-backend/tests/testUtils"
	"log"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var Router *mux.Router

type AddedAccount struct {
	AccountModel.Account
	token string
}

var AddedAccounts []AddedAccount

func (account *AddedAccount) create(email string, username string, password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	account.Account = AccountModel.Account{UserID: primitive.NewObjectID(), Email: email, Username: username, Password: string(hashedPassword)}
}

func clearDatabase(DB *mongo.Database) {
	var ctx context.Context
	DB.Collection("accounts").Drop(ctx)
}

func TestMain(m *testing.M) {
	Service.DB = db.Connect().Database("test")

	Router = mux.NewRouter().StrictSlash(true)
	Router.Use(auth.JwtAuthentication) //attach JWT auth middleware

	if len(Controller.Routes) == 0 {
		fmt.Println("Error: no bind routes specified for service")
		os.Exit(1)
	}
	for _, route := range Controller.Routes {
		Router.HandleFunc(route.Path, route.Function).Methods(route.Method...)
	}

	//inserting some accounts into the database before running tests
	newAddedAccount := AddedAccount{}
	newAddedAccount.create("randomemail@gmail.com", "someusername", "testpassword123")
	AddedAccounts = append(AddedAccounts, newAddedAccount)

	newAddedAccount = AddedAccount{}
	newAddedAccount.create("email@email.com", "aUser", "rootroot")
	AddedAccounts = append(AddedAccounts, newAddedAccount)

	newAddedAccount = AddedAccount{}
	newAddedAccount.create("myEmail@gmail.com", "testing", "verysecurepassword1234")
	AddedAccounts = append(AddedAccounts, newAddedAccount)

	accounts := []interface{}{AddedAccounts[0].Account, AddedAccounts[1].Account, AddedAccounts[2].Account}

	collection := AccountModel.NewAccountCollection(Service.DB)

	_, err := collection.InsertMany(context.TODO(), accounts)
	if err != nil {
		log.Panic("Test initialization failed - Error: " + err.Error())
	}

	code := m.Run()
	clearDatabase(Service.DB)
	os.Exit(code)
}

func TestSignup(t *testing.T) {
	requests := []tu.Request{}
	requests = append(requests, tu.Request{
		URL:              "/api/account/signup",
		Method:           "POST",
		Payload:          map[string]interface{}{"username": "testUser", "email": "testEmail@gmail.com", "password": "testPassword"},
		ExpectedCode:     201,
		ExpectedResponse: u.Message(true, "Account has been created")})

	requests = append(requests, tu.Request{
		URL:              "/api/account/signup",
		Method:           "POST",
		Payload:          map[string]interface{}{"username": "testUser", "email": "testEmail@gmail.com", "password": "testPassword"},
		ExpectedCode:     400,
		ExpectedResponse: u.Message(false, "Email address already in use by another user")})

	requests = append(requests, tu.Request{
		URL:              "/api/account/signup",
		Method:           "POST",
		Payload:          map[string]interface{}{"username": "testUser", "email": "test1Email@gmail.com", "password": "testPassword"},
		ExpectedCode:     400,
		ExpectedResponse: u.Message(false, "Username already in use by another user")})

	for _, request := range requests {
		request.ExecuteRequest(t, Router)
	}
}

func TestProfile(t *testing.T) {

	requests := []tu.Request{}
	requests = append(requests, tu.Request{
		URL:              "/api/account/profile?username=" + AddedAccounts[0].Username,
		Method:           "GET",
		Payload:          map[string]interface{}{},
		ExpectedCode:     200,
		ExpectedResponse: map[string]interface{}{"account": map[string]interface{}{"UserID": AddedAccounts[0].UserID.Hex(), "createdAt": "", "email": "", "password": "", "username": AddedAccounts[0].Username}, "accountOwner": false, "message": "", "status": true}})

	requests = append(requests, tu.Request{
		URL:              "/api/account/profile?username=" + AddedAccounts[1].Username,
		Method:           "GET",
		Payload:          map[string]interface{}{},
		ExpectedCode:     200,
		ExpectedResponse: map[string]interface{}{"account": map[string]interface{}{"UserID": AddedAccounts[1].UserID.Hex(), "createdAt": AddedAccounts[1].UserID.Timestamp().String(), "email": AddedAccounts[1].Email, "password": "", "username": AddedAccounts[1].Username}, "accountOwner": true, "message": "", "status": true}})

	for _, request := range requests {
		request.ExecuteRequest(t, Router)
	}

}
