package accountapi

import (
	"encoding/json"
	AccountModel "issue-tracker-backend/src/models/account"
	Server "issue-tracker-backend/src/server"
	u "issue-tracker-backend/src/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func profile(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		//GET:
		//Returns a user profile
		//The request payload should contain the username of the user to retrieve
		// {username: required}
		//Example payload:
		//{"username": "Tom"}
		account := &AccountModel.Account{}
		account.Username = r.URL.Query().Get("username")
		if account.Username == "" {
			u.Respond(w, u.Message(false, "Invalid request: Missing username"), http.StatusBadRequest)
			return
		}
		userID := r.Context().Value("userID")

		if userID != nil {
			var ok bool
			userID, ok = userID.(string)
			if !ok {
				u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
			}
		} else {
			userID = ""
		}
		resp, statusCode := account.Get(userID.(string), Server.DBConnection) // remove the  userID.(string) part and change to userID
		u.Respond(w, resp, statusCode)

	case http.MethodPut:
		//PUT:
		//Updates a user profile
		//The request payload should contain the fields you want to update
		//Example payload:
		//{"username": "Tom", "email": "newEmail@gmail.com"}

		userID := r.Context().Value("userID")
		if userID != nil {
			objID, _ := primitive.ObjectIDFromHex(userID.(string)) //converting to object id from hex needs to be moved into a function for abstraction
			account := &AccountModel.Account{ID: objID}
			err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
			if err != nil {
				u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
				return
			}
			resp, statusCode := account.Update(Server.DBConnection)
			u.Respond(w, resp, statusCode)
		} else {
			u.Respond(w, u.Message(false, "Error retrieving userID"), http.StatusNotFound)
		}
	default:
		u.Respond(w, u.Message(false, "Invalid request: Method unsupported"), http.StatusMethodNotAllowed)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

		account := &AccountModel.Account{}
		err := json.NewDecoder(r.Body).Decode(&account) //decode the request body into struct and failed if any error occur
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
			return
		}
		resp, statusCode := account.Create(Server.DBConnection) //Create account
		u.Respond(w, resp, statusCode)
	default:
		u.Respond(w, u.Message(false, "Invalid request: Method unsupported"), http.StatusMethodNotAllowed)
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	account := &AccountModel.Account{}

	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	resp, statusCode := account.Login(Server.DBConnection)

	u.Respond(w, resp, statusCode)

}
