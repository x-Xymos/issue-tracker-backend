package logincont

import (
	"encoding/json"
	"net/http"
	AccountModel "todo-backend/src/models/account"
	Service "todo-backend/src/servicetemplates"
	u "todo-backend/src/utils"
)

func getUser(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID")

	if userID != nil {
		resp := AccountModel.GetUser(userID.(string), Service.DBConn)
		u.Respond(w, resp)
	} else {
		u.Respond(w, map[string]interface{}{"status": false, "message": "Error retrieving user information"})
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	account := &AccountModel.Account{}

	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := AccountModel.Login(&account.Email, &account.Password, Service.DBConn)
	u.Respond(w, resp)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("userID").(string)
	//json.NewEncoder(w).Encode(map[string]string{"Name": "This is the Login-Api", "UserID": userID})
	//resp := u.Message(true, "Logged In")
	u.Respond(w, map[string]interface{}{"status": true, "message": "This is the Login-Api", "UserID": userID})
}

//Routes : an array of route bindings
var Routes = []Service.RouteBinding{
	Service.RouteBinding{"/api/", home, []string{"GET"}},
	Service.RouteBinding{"/api/login", login, []string{"POST"}},
	Service.RouteBinding{"/api/user", getUser, []string{"GET"}},
}

//ServiceName : service name
var ServiceName = "Login-api"

//Port : service port
var Port = "8880"
