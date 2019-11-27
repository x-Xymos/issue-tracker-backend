package servicetemplates

import (
	"context"
	"fmt"
	"issue-tracker-backend/src/auth"
	"issue-tracker-backend/src/servicetemplates/db"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//RouteBinding : route binding definition
type RouteBinding struct {
	Path     string
	Function func(http.ResponseWriter, *http.Request)
	Method   []string
}

//DB : creates a DB connection that is used by the service
var DB *mongo.Database
var Router *mux.Router
var ctx context.Context
var cancel context.CancelFunc

//Stop : stop the listener
func Stop() {
	cancel()
}

//Start : start listener
//requires a routes array containing all
//the RouteBindings that this service will accept
func Start(routes *[]RouteBinding, port *string, serviceName *string, DBName *string) {
	ctx, cancel = context.WithCancel(context.Background())

	DB = db.Connect().Database(*DBName)

	Router = mux.NewRouter().StrictSlash(true)
	Router.Use(auth.JwtAuthentication) //attach JWT auth middleware

	if len(*routes) == 0 {
		fmt.Println("Error: no bind routes specified for service")
		os.Exit(1)
	}
	for _, route := range *routes {
		Router.HandleFunc(route.Path, route.Function).Methods(route.Method...)

	}
	fmt.Printf("Started %s on 0.0.0.0:%s\n", *serviceName, *port)

	srv := &http.Server{Addr: ":" + *port, Handler: handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "params"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(Router)}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil && err != context.Canceled {
		log.Println(err)
	}
	log.Println("done.")
}
