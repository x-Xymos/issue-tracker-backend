package logincontroller

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

		account := &AccountModel.Account{}
		err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
			return
		}

		userID := r.Context().Value("userID")

		if userID != nil {
			userID = userID.(string)
		} else {
			userID = ""
		}
		resp, statusCode := account.Get(userID.(string), Service.DBConn)
		u.Respond(w, resp, statusCode)

	case http.MethodPut:
		userID := r.Context().Value("userID")
		if userID != nil {
			objID, _ := primitive.ObjectIDFromHex(userID.(string))
			account := &AccountModel.Account{UserID: objID}
			err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
			if err != nil {
				u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
				return
			}
			resp, statusCode := account.Update(Service.DBConn)
			u.Respond(w, resp, statusCode)
		} else {
			u.Respond(w, u.Message(false, "Error retrieving userID"), http.StatusNotFound)
		}
	default:
		u.Respond(w, u.Message(false, "Invalid request: Method unsupported"), http.StatusMethodNotAllowed)
	}

}

func signup(w http.ResponseWriter, r *http.Request) {

	account := &AccountModel.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}
	resp, statusCode := account.Create(Service.DBConn) //Create account
	u.Respond(w, resp, statusCode)
}

func login(w http.ResponseWriter, r *http.Request) {

	account := &AccountModel.Account{}

	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	//resp := account.Login(&account.Email, &account.Password, Service.DBConn)
	resp := account.Login(Service.DBConn)

	u.Respond(w, resp, http.StatusOK)

}

//Routes : an array of route bindings
var Routes = []Service.RouteBinding{
	Service.RouteBinding{"/api/account/login", login, []string{"POST"}},
	Service.RouteBinding{"/api/account/signup", signup, []string{"POST"}},
	Service.RouteBinding{"/api/account/profile", profile, []string{"GET", "PUT"}},
}

//ServiceName : service name
var ServiceName = "Account-api"

//Port : service port
var Port = "8880"
