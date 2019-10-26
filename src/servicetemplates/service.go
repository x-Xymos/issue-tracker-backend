package servicetemplates

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo-backend/src/auth"
	"todo-backend/src/servicetemplates/db"

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

	//handler := cors.Default().Handler(router)

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowCredentials: true,
	// 	//AllowedHeaders:   []string{"Authorization", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Origin", "Content-Length", "Content-Type", "transformRequest", "transformResponse", "xsrfCookieName", "xsrfHeaderName"},
	// 	// Enable Debugging for testing, consider disabling in production
	// 	Debug: true,
	// })

	// Insert the middleware
	//handler = c.Handler(handler)

	//log.Fatal(http.ListenAndServe(":"+*port, handler))
	log.Fatal(http.ListenAndServe(":"+*port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
