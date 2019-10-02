package account

import (
	"os"
	"strings"
	u "todo-backend/src/utils"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//Account : user account struct
type Account struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Token    string `json:"token"`
}

//Token : JWT token struct
type Token struct {
	UserID uint
	jwt.StandardClaims
}

//Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

//Create : account creation
func (account *Account) Create(DBConn *mongo.Client) map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
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
