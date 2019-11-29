package projectcontroller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	ProjectModel "issue-tracker-backend/src/models/project"
	Server "issue-tracker-backend/src/server"
	u "issue-tracker-backend/src/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func project(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//GET:
		//Returns a project
		//The request payload should contain the title of the project to retrieve
		// {title: required}
		//Example payload:
		//{"title": "projectTitle"}

		project := &ProjectModel.Project{}
		err := json.NewDecoder(r.Body).Decode(project) //decode the request body into struct and failed if any error occur
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
		resp, statusCode := project.Get(userID.(string), Server.DBConnection)
		u.Respond(w, resp, statusCode)

	case http.MethodPost:
		//POST:
		//Creates a new project
		//The request payload should cointain information required to create a new project
		// {title: required}
		//Example payload:
		//{"title": "projectTitle"}
		userID := r.Context().Value("userID")
		if userID != nil {
			objID, _ := primitive.ObjectIDFromHex(userID.(string))
			project := &ProjectModel.Project{OwnerID: objID}
			err := json.NewDecoder(r.Body).Decode(project) //decode the request body into struct and failed if any error occur
			if err != nil {
				u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
				return
			}
			resp, statusCode := project.Create(Server.DBConnection)
			u.Respond(w, resp, statusCode)
		} else {
			u.Respond(w, u.Message(false, "Error retrieving userID"), http.StatusUnauthorized)
		}

	case http.MethodPut:
		//PUT:
		//Updates a project
		//The request payload should be an array containing 2 elements:
		// elem1: {title: required}
		// elem2: {title: required}
		//Example payload:
		// [{"title": "projectTitle"},{"title": "changedProjectTitle"}]

		userID := r.Context().Value("userID")
		if userID != nil {
			objID, _ := primitive.ObjectIDFromHex(userID.(string))

			var projects []*ProjectModel.Project
			data, err := ioutil.ReadAll(r.Body)
			r.Body.Close()

			if err != nil {
				u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
				return
			}
			if err := json.Unmarshal([]byte(data), &projects); err != nil {
				u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
				return
			}
			projects[0].OwnerID = objID
			resp, statusCode := projects[0].Update(projects[1], Server.DBConnection)
			u.Respond(w, resp, statusCode)
		} else {
			u.Respond(w, u.Message(false, "Error retrieving userID"), http.StatusUnauthorized)
		}
	default:
		u.Respond(w, u.Message(false, "Error: Method unsupported"), http.StatusMethodNotAllowed)
	}
}

func projects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//GET:
		//Returns 10 projects based on an optional search query,
		//can be paginated by providing the id of the last project that was retrieved
		//{title: optional, lastID: optional}
		//Payload examples:
		//{"title": "optionalSearchQuery", "lastID","5dd1c87c662abc93b4ddbaa9"}
		//{"title": "", "lastID","5dd1c87c662abc93b4ddbaa9"}
		//{"title": "optionalSearchQuery", "lastID",""}
		defer r.Body.Close()
		var data map[string]interface{}

		reqBody, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		if err := json.Unmarshal(reqBody, &data); err != nil {
			u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

		project := &ProjectModel.Project{}
		err := json.NewDecoder(r.Body).Decode(project) //decode the request body into struct and failed if any error occur
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
			return
		}
		resp, statusCode := project.GetAll(data["lastID"].(string), Server.DBConnection) //data["lastID"].(string) check  this first
		u.Respond(w, resp, statusCode)
	default:
		u.Respond(w, u.Message(false, "Error: Method unsupported"), http.StatusMethodNotAllowed)
	}
}

//Routes : an array of route bindings
var Routes = []Server.RouteBinding{
	Server.RouteBinding{"/api/project", project, []string{"GET", "POST", "PUT", "DELETE"}},
	Server.RouteBinding{"/api/projects", projects, []string{"GET"}},
}

var DBName = "issue-tracker"

//ServerName : Server name
var ServerName = "Project-api"

//Port : Server port
var Port = "8882"
