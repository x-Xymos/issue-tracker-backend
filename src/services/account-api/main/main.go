package main

import (
	Controller "todo-backend/src/services/account-api/controllers"
	Service "todo-backend/src/servicetemplates"
)

func main() {
	Service.Start(&Controller.Routes, &Controller.Port, &Controller.ServiceName)
}
