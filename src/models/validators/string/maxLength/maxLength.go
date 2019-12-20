package strmaxlength

import (
	"errors"
	v "issue-tracker-backend/src/models/validators/validator"
	"strconv"
	"unicode/utf8"
)

//Validator : Checks if the length of the input doesn't exceed the provided length
func Validator(input interface{}, options *[]*v.Option) error {

	_options := make(map[string]interface{})
	for _, v := range *options {
		_options[v.Name] = v.Value
	}

	_input, ok := input.(string)
	if !ok {
		return errors.New("Error casting input value in maxLength validator")
	}

	_maxLength, _ := _options["maxLength"].(int)

	if _maxLength < 1 {
		return errors.New("Error, passed length parameter should be more than 0 in maxLength validator")
	}

	if utf8.RuneCountInString(_input) > _maxLength {
		return errors.New("has to be less than " + strconv.Itoa(_maxLength) + " characters long")
	}
	return nil
}

//Options : assign options
func Options(length int) *[]*v.Option {
	return &[]*v.Option{
		&v.Option{Name: "maxLength", Value: length},
	}
}
