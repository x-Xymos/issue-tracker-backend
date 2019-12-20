package strminlength

import (
	"errors"
	v "issue-tracker-backend/src/models/validators/validator"
	"strconv"
	"unicode/utf8"
)

//Validator : Checks if the length of the input is at least the provided length
func Validator(input interface{}, options *[]*v.Option) error {

	_options := make(map[string]interface{})
	for _, v := range *options {
		_options[v.Name] = v.Value
	}

	_input, ok := input.(string)
	if !ok {
		return errors.New("Error casting input value in minLength validator")
	}

	_minLength, ok := _options["minLength"].(int)

	if _minLength < 1 {
		return errors.New("Error, passed length parameter should be more than 0 in minLength validator")
	}

	if utf8.RuneCountInString(_input) < _minLength {
		return errors.New("has to be more than " + strconv.Itoa(_minLength) + " characters long")
	}
	return nil
}

//Options : assign options
func Options(length int) *[]*v.Option {
	return &[]*v.Option{
		&v.Option{Name: "minLength", Value: length},
	}
}
