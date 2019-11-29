package main

import (
	Server "issue-tracker-backend/src/server"
	Controller "issue-tracker-backend/src/services/project-api/controllers"
)

func main() {
	Server.Start(&Controller.Routes, &Controller.Port, &Controller.ServerName, &Controller.DBName)
}
