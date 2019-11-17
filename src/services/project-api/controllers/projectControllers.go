package projectcontroller

import (
	"encoding/json"
	ProjectModel "issue-tracker-backend/src/models/project"
	Service "issue-tracker-backend/src/servicetemplates"
	u "issue-tracker-backend/src/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func project(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		paramHeader := r.Header.Get("params")
		if paramHeader == "" {
			u.Respond(w, u.Message(false, "Invalid request: Missing query parameters"), http.StatusBadRequest)
			return
		}
		project := &ProjectModel.Project{}

		err := json.Unmarshal([]byte(paramHeader), project)
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
		resp, statusCode := project.Get(userID.(string), Service.DBConn)
		u.Respond(w, resp, statusCode)

	case http.MethodPost:
		userID := r.Context().Value("userID")
		if userID != nil {
			objID, _ := primitive.ObjectIDFromHex(userID.(string))
			project := &ProjectModel.Project{OwnerID: objID}
			err := json.NewDecoder(r.Body).Decode(project) //decode the request body into struct and failed if any error occur
			if err != nil {
				u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
				return
			}
			resp, statusCode := project.Create(Service.DBConn)
			u.Respond(w, resp, statusCode)
		} else {
			u.Respond(w, u.Message(false, "Error retrieving userID"), http.StatusUnauthorized)
		}

	case http.MethodPut:
		userID := r.Context().Value("userID")
		if userID != nil {
			objID, _ := primitive.ObjectIDFromHex(userID.(string))
			project := &ProjectModel.Project{OwnerID: objID}
			err := json.NewDecoder(r.Body).Decode(project) //decode the request body into struct and failed if any error occur
			if err != nil {
				u.Respond(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
				return
			}
			resp, statusCode := project.Update(Service.DBConn)
			u.Respond(w, resp, statusCode)
		} else {
			u.Respond(w, u.Message(false, "Error retrieving userID"), http.StatusUnauthorized)
		}
	default:
		u.Respond(w, u.Message(false, "Error: Method unsupported"), http.StatusMethodNotAllowed)
	}
}

//Routes : an array of route bindings
var Routes = []Service.RouteBinding{
	Service.RouteBinding{"/api/project", project, []string{"GET", "POST", "PUT", "DELETE"}},
}

//ServiceName : service name
var ServiceName = "Project-api"

//Port : service port
var Port = "8882"
