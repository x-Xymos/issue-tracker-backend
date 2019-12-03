package unique

import (
	"errors"
	"issue-tracker-backend/src/db"
	v "issue-tracker-backend/src/models/validators/validator"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	searchField   = "field"
	caseSensitive = "caseSensitive"
	database      = "database"
	collection    = "collection"
	filter        = "filter"
)

//Validator :
func Validator(input interface{}, options *[]*v.Option) error {

	_options := make(map[string]interface{})
	for _, v := range *options {
		_options[v.Name] = v.Value
	}

	_collection, ok := _options[collection].(string)
	if !ok {
		return errors.New("Error, collection requires a string paramater in unique validator")
	}

	_searchField, ok := _options[searchField].(string)
	if !ok {
		return errors.New("Error, collection requires a string paramater in unique validator")
	}

	_caseSensitive, ok := _options[caseSensitive].(bool)
	if !ok {
		return errors.New("Error, collection requires a bool paramater in unique validator")
	}

	_, err := db.FindOne(_options[database], _collection, map[string]interface{}{_searchField: input}, nil, nil, _caseSensitive)

	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		return err
	}
	return errors.New("has already been taken")
}

//SearchField : the name of the field in the database to perform the search against: accepted values: string
func SearchField(value interface{}) *v.Option {
	return &v.Option{Name: searchField, Value: value}
}

//CaseSensitive : whether to perform a case sensitive or insensitive search, accepted values : bool
func CaseSensitive(value interface{}) *v.Option {
	return &v.Option{Name: caseSensitive, Value: value}
}

//Database : a pointer to a database connection
func Database(value interface{}) *v.Option {
	return &v.Option{Name: database, Value: value}
}

//Collection : name of the collection to perform the search on, accepted values: string
func Collection(value interface{}) *v.Option {
	return &v.Option{Name: collection, Value: value}
}
