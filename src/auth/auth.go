package auth

import (
	"context"
	"net/http"
	"strings"
	"todo-backend/env"
	AccountModel "todo-backend/src/models/account"
	u "todo-backend/src/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

//JwtAuthentication :  function used for all requests that require authentication to check for validity
//of the users token
var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//"/api/",
		notAuth := []string{"/api/account/login", "/api/account/signup"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                                        //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
		if tokenHeader == "" {                       //Token is missing, returns with error code 403 Unauthorized
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		splitToken := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitToken) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenStr := splitToken[1] //Grab the token part
		decodedToken := &AccountModel.Token{}

		token, err := jwt.ParseWithClaims(tokenStr, decodedToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(env.TokenPassword), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", decodedToken.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
