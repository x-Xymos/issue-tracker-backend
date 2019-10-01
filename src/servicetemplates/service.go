package servicetemplates

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo-backend/src/servicetemplates/db"

	"github.com/gorilla/mux"
)

//RouteBinding : route binding definition
type RouteBinding struct {
	Path     string
	Function func(http.ResponseWriter, *http.Request)
	Method   string
}

//DB : creates a DB connection that is used by a service
var DB = db.Connect()

//Start : start listener
//requires a routes array containing all
//the RouteBindings that this service will accept
func Start(routes *[]RouteBinding, port *string, serviceName *string) {

	router := mux.NewRouter().StrictSlash(true)
	if len(*routes) == 0 {
		fmt.Println("Error: no bind routes specified for service")
		os.Exit(1)
	}
	for _, route := range *routes {
		router.HandleFunc(route.Path, route.Function).Methods(route.Method)
	}
	fmt.Printf("Started %s on 0.0.0.0:%s\n", *serviceName, *port)

	log.Fatal(http.ListenAndServe(":"+*port, router))
}
