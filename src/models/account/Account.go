package account

import (
	"context"
	"log"
	"strings"
	"todo-backend/env"
	u "todo-backend/src/utils"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//Account : user account struct
type Account struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

//Token : JWT token struct
type Token struct {
	UserID string
	jwt.StandardClaims
}

//Validate incoming user details when they register for a new account
func (account *Account) Validate(DBConn *mongo.Client) (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	tempAcc := &Account{}

	collection := DBConn.Database("todo").Collection("accounts")

	//Email must be unique
	emailFilter := bson.D{{"email", account.Email}}

	err := collection.FindOne(context.TODO(), emailFilter).Decode(&tempAcc)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Fatal(err)
		return u.Message(false, "Connection error, please try again"), false
	}

	if tempAcc.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	//Username must be unique
	usernameFilter := bson.D{{"username", account.Username}}

	err = collection.FindOne(context.TODO(), usernameFilter).Decode(&tempAcc)

	if err != mongo.ErrNoDocuments {
		log.Fatal(err)
		return u.Message(false, "Connection error, please try again"), false
	}

	if tempAcc.Username != "" {
		return u.Message(false, "Username already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

//Create : account creation
func (account *Account) Create(DBConn *mongo.Client) map[string]interface{} {

	if resp, ok := account.Validate(DBConn); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	collection := DBConn.Database("todo").Collection("accounts")

	insertResult, err := collection.InsertOne(context.TODO(), account)
	if err != nil {
		log.Fatal(err)
		return u.Message(false, "Failed to create account, connection error.")
	}

	UserID := insertResult.InsertedID.(primitive.ObjectID).Hex()

	//Create new JWT token for the newly registered account
	tk := &Token{UserID: UserID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(env.TokenPassword))
	account.Token = tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

//Login : login
// func Login(email string, password string, DBConn *mongo.Client) map[string]interface{} {

// 	account := &Account{}

// 	collection := DBConn.Database("todo").Collection("accounts")
// 	filter := bson.D{{"email", email}}

// 	err := collection.FindOne(context.TODO(), filter).Decode(&account)

// 	if err != nil {
// 		if err.Error() == "mongo: no documents in result" {
// 			return u.Message(false, "Email address not found.")
// 		}
// 		log.Fatal(err)
// 		return u.Message(false, "Connection error.")
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
// 	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
// 		return u.Message(false, "Invalid login credentials. Please try again")
// 	}
// 	//Worked! Logged In
// 	account.Password = ""

// 	//Create JWT token
// 	tk := &Token{UserId: account.ID}
// 	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
// 	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
// 	account.Token = tokenString //Store the token in the response

// 	resp := u.Message(true, "Logged In")
// 	resp["account"] = account
// 	return resp
// }

// func GetUser(u uint) *Account {

// 	acc := &Account{}
// 	GetDB().Table("accounts").Where("id = ?", u).First(acc)
// 	if acc.Email == "" { //User not found!
// 		return nil
// 	}

// 	acc.Password = ""
// 	return acc
// }
