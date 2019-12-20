package unique

import (
	"errors"
	"issue-tracker-backend/src/db"
	v "issue-tracker-backend/src/models/validators/validator"

	"go.mongodb.org/mongo-driver/mongo"
)

//Validator : Tests if the field entry is unique to the database collection
func Validator(input interface{}, options *[]*v.Option) error {

	_options := make(map[string]interface{})
	for _, v := range *options {
		_options[v.Name] = v.Value
	}

	_collection, ok := _options["collection"].(string)
	if !ok {
		return errors.New("Error, collection requires a string paramater in unique validator")
	}

	_searchField, ok := _options["searchField"].(string)
	if !ok {
		return errors.New("Error, searchField requires a string paramater in unique validator")
	}

	_caseSensitive, ok := _options["caseSensitive"].(bool)
	if !ok {
		return errors.New("Error, caseSensitive requires a bool paramater in unique validator")
	}

	_, err := db.FindOne(_options["database"], _collection, map[string]interface{}{_searchField: input}, nil, nil, _caseSensitive)

	if err == mongo.ErrNoDocuments {
		return nil
	} else if err != nil {
		return err
	}
	return errors.New("has already been taken")
}

//Options : assign options
func Options(DBConnection interface{}, collectionName string, searchField string, caseSensitive bool) *[]*v.Option {
	return &[]*v.Option{
		&v.Option{Name: "database", Value: DBConnection},
		&v.Option{Name: "collection", Value: collectionName},
		&v.Option{Name: "searchField", Value: searchField},
		&v.Option{Name: "caseSensitive", Value: caseSensitive},
	}
}
