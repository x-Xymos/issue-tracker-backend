package logincont

import (
	"encoding/json"
	"net/http"
	AccountModel "todo-backend/src/models/account"
	Service "todo-backend/src/servicetemplates"
	u "todo-backend/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getProfile(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID")

	if userID != nil {
		objID, _ := primitive.ObjectIDFromHex(userID.(string))
		account := &AccountModel.Account{UserID: objID}
		resp := account.GetProfile(Service.DBConn)
		u.Respond(w, resp)
	} else {
		u.Respond(w, u.Message(false, "Error retrieving userID"))
	}
}

func updateProfile(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID")
	if userID != nil {
		objID, _ := primitive.ObjectIDFromHex(userID.(string))
		account := &AccountModel.Account{UserID: objID}
		err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"))
			return
		}
		resp := account.UpdateProfile(Service.DBConn)
		u.Respond(w, resp)
	} else {
		u.Respond(w, u.Message(false, "Error retrieving userID"))
	}
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

func login(w http.ResponseWriter, r *http.Request) {

	account := &AccountModel.Account{}

	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	//resp := account.Login(&account.Email, &account.Password, Service.DBConn)
	resp := account.Login(Service.DBConn)

	u.Respond(w, resp)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("userID").(string)
	u.Respond(w, map[string]interface{}{"status": true, "message": "This is the Account-Api", "UserID": userID})
}

//Routes : an array of route bindings
var Routes = []Service.RouteBinding{
	Service.RouteBinding{"/api/account", home, []string{"GET"}},
	Service.RouteBinding{"/api/account/login", login, []string{"POST"}},
	Service.RouteBinding{"/api/account/signup", signup, []string{"POST"}},
	Service.RouteBinding{"/api/account/profile", getProfile, []string{"GET"}},
	Service.RouteBinding{"/api/account/profile/update", updateProfile, []string{"POST"}},
}

//ServiceName : service name
var ServiceName = "Account-api"

//Port : service port
var Port = "8880"
