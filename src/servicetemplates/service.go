package servicetemplates

import (
	"fmt"
	"issue-tracker-backend/src/auth"
	"issue-tracker-backend/src/servicetemplates/db"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//RouteBinding : route binding definition
type RouteBinding struct {
	Path     string
	Function func(http.ResponseWriter, *http.Request)
	Method   []string
}

//DBConn : creates a DB connection that is used by the service
var DBConn = db.Connect()

//Start : start listener
//requires a routes array containing all
//the RouteBindings that this service will accept
func Start(routes *[]RouteBinding, port *string, serviceName *string) {

	router := mux.NewRouter().StrictSlash(true)
	router.Use(auth.JwtAuthentication) //attach JWT auth middleware

	if len(*routes) == 0 {
		fmt.Println("Error: no bind routes specified for service")
		os.Exit(1)
	}
	for _, route := range *routes {
		router.HandleFunc(route.Path, route.Function).Methods(route.Method...)

	}
	fmt.Printf("Started %s on 0.0.0.0:%s\n", *serviceName, *port)

	log.Fatal(http.ListenAndServe(":"+*port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
