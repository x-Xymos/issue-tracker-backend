package issuecontroller

import (
	"encoding/json"
	"fmt"
	IssueModel "issue-tracker-backend/src/models/issue"
	Server "issue-tracker-backend/src/server"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {

	issue := &IssueModel.Issue{}

	fmt.Println(issue)

	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value("user").(string)
	json.NewEncoder(w).Encode(map[string]string{"Name": "This is the Issue Tracking Api", "UserID": userID})

}

//Routes : an array of route bindings
var Routes = []Server.RouteBinding{
	Server.RouteBinding{"/api/", home, []string{"GET"}},
}

var DBName = "issue-tracker"

//ServerName : Server name
var ServerName = "Issue-api"

//Port : Server port
var Port = "8881"
