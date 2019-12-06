package main

import (
	ProjectModel "issue-tracker-backend/src/models/project"
	Server "issue-tracker-backend/src/server"
)

//Routes : an array of route bindings
var Routes = []Server.RouteBinding{
	Server.RouteBinding{"/api/project", project, []string{"GET", "POST", "PUT", "DELETE"}},
	Server.RouteBinding{"/api/projects", projects, []string{"GET"}},
}

//InitValidators : pointer to the function that initializes the validators for the account models, this is ran before the server starts
var InitValidators = ProjectModel.InitValidators

//DBName : name of the database used by the service
var DBName = "issue-tracker"

//ServerName : Server name
var ServerName = "Project-api"

//Port : Server port
var Port = "8882"

func main() {
	Server.Start(&Routes, &Port, &ServerName, &DBName, InitValidators)
}
