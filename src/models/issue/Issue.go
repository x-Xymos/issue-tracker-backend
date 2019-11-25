package issue

import (
	u "issue-tracker-backend/src/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Issue : Issue struct
type Issue struct {
	UserID      primitive.ObjectID `bson:"_id, omitempty"`
	IssueID     uint64             `json:"issueID"`
	ProjectID   uint64             `json:"projectID"`
	Title       string             `json:"title"`
	Body        string             `json:"body"`
	Tags        string             `json:"tags"`
	DateCreated time.Time
}

func newIssueCollection(DBConn *mongo.Client) *mongo.Collection {
	return DBConn.Database("issue-tracker").Collection("issues")
}

func (issue *Issue) _titleValidator(DBConn *mongo.Client) map[string]interface{} {

	if len(issue.Title) < 1 {
		return u.Message(false, "Title has to be at least 1 character long")
	}
	return u.Message(true, "")
}

// //ValidateIssueCreation :
// func (issue *Issue) ValidateIssueCreation(DBConn *mongo.Client) map[string]interface{} {

// 	if resp := account._emailValidator(DBConn); resp["status"] == false {
// 		return resp
// 	}

// 	if resp := account._usernameValidator(DBConn); resp["status"] == false {
// 		return resp
// 	}

// 	return account._passwordValidator()
// }

// //Create : Issue creation
// func (issue *Issue) Create(DBConn *mongo.Client) map[string]interface{} {

// 	if resp := account.ValidateAccountCreation(DBConn); resp["status"] == false {
// 		return resp
// 	}

// 	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
// 	account.Password = string(hashedPassword)
// 	account.UserID = primitive.NewObjectID()

// 	collection := DBConn.Database("issue-tracker").Collection("accounts")

// 	_, err := collection.InsertOne(context.TODO(), account)
// 	if err != nil {
// 		return u.Message(false, "Failed to create account, connection error.")
// 	}

// 	account.Password = "" //delete password

// 	response := u.Message(true, "Account has been created")
// 	response["account"] = account
// 	return response
// }
