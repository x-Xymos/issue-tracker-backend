package signupcont

import (
	//"context"

	"encoding/json"
	"net/http"
	AccountModel "todo-backend/src/models/account"
	Service "todo-backend/src/servicetemplates"
	u "todo-backend/src/utils"
)

func home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user").(string)
	json.NewEncoder(w).Encode(map[string]string{"Name": "This is the Singup-Api", "UserID": userID})

}

func signup(w http.ResponseWriter, r *http.Request) {

	account := &AccountModel.Account{}

	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create(Service.DBConn) //Create account
	u.Respond(w, resp)
}

//Routes : an array of route bindings
var Routes = []Service.RouteBinding{
	Service.RouteBinding{"/api/", home, "GET"},
	Service.RouteBinding{"/api/signup", signup, "POST"}}

//ServiceName : service name
var ServiceName = "Signup-api"

//Port : service port
var Port = "8081"
