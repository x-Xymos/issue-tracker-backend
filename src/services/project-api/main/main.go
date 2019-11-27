package main

import (
	Controller "issue-tracker-backend/src/services/project-api/controllers"
	Service "issue-tracker-backend/src/servicetemplates"
)

func main() {
	Service.Start(&Controller.Routes, &Controller.Port, &Controller.ServiceName, &Controller.DBName)
}
