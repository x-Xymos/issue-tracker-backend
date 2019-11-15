package logincont

import (
	"encoding/json"
	AccountModel "issue-tracker-backend/src/models/account"
	Service "issue-tracker-backend/src/servicetemplates"
	u "issue-tracker-backend/src/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func profile(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		paramHeader := r.Header.Get("params")
		if paramHeader == "" {
			response := u.Message(false, "Missing query parameters")
			w.WriteHeader(http.StatusBadRequest)
			u.Respond(w, response)
			return
		}
		account := &AccountModel.Account{}

		err := json.Unmarshal([]byte(paramHeader), account)
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"))
			return
		}
		userID := r.Context().Value("userID")

		if userID != nil {
			userID = userID.(string)
		} else {
			userID = ""
		}
		resp := account.Get(userID.(string), Service.DBConn)
		u.Respond(w, resp)

	case http.MethodPost:
		userID := r.Context().Value("userID")
		if userID != nil {
			objID, _ := primitive.ObjectIDFromHex(userID.(string))
			account := &AccountModel.Account{UserID: objID}
			err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
			if err != nil {
				u.Respond(w, u.Message(false, "Invalid request"))
				return
			}
			resp := account.Update(Service.DBConn)
			u.Respond(w, resp)
		} else {
			u.Respond(w, u.Message(false, "Error retrieving userID"))
		}
	default:
		u.Respond(w, u.Message(false, "Error: Method unsupported"))
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

//Routes : an array of route bindings
var Routes = []Service.RouteBinding{
	Service.RouteBinding{"/api/account/login", login, []string{"POST"}},
	Service.RouteBinding{"/api/account/signup", signup, []string{"POST"}},
	Service.RouteBinding{"/api/account/profile", profile, []string{"GET", "POST"}},
}

//ServiceName : service name
var ServiceName = "Account-api"

//Port : service port
var Port = "8880"
