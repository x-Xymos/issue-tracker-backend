package main

import (
	AccountModel "issue-tracker-backend/src/models/account"
	Server "issue-tracker-backend/src/server"
)

//Routes : an array of route bindings
var Routes = []Server.RouteBinding{
	Server.RouteBinding{"/api/account/login", login, []string{"POST"}},
	Server.RouteBinding{"/api/account/signup", signup, []string{"POST"}},
	Server.RouteBinding{"/api/account/profile", profile, []string{"GET", "PUT"}},
}

//InitValidators : pointer to the function that initializes the validators for the account models, this is ran before the server starts
var InitValidators = AccountModel.InitValidators

//DBName : name of the database used by the service
var DBName = "issue-tracker"

//ServerName : Server name
var ServerName = "Account-api"

//Port : Server port
var Port = "8880"

func main() {
	Server.Start(&Routes, &Port, &ServerName, &DBName, InitValidators)
}
