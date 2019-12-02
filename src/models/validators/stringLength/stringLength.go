package stringlength

import (
	"errors"
	v "issue-tracker-backend/src/models/validators/validator"
	"strconv"
	"unicode/utf8"
)

//Validator :
func Validator(input interface{}, options *[]*v.Option) error {
	for _, v := range *options {
		switch v.Name {
		case "minLength":
			input, ok := input.(string)
			if !ok {
				return errors.New("Error casting value in validator")
			}

			value, ok := v.Value.(int)
			if !ok {
				return errors.New("Error casting value in validator")
			}

			if utf8.RuneCountInString(input) < value {
				return errors.New("has to be at least " + strconv.Itoa(value) + " characters long")
			}

		case "maxLength":
			input, ok := input.(string)
			if !ok {
				return errors.New("Error casting value in validator")
			}

			value, ok := v.Value.(int)
			if !ok {
				return errors.New("Error casting value in validator")
			}

			if utf8.RuneCountInString(input) > value {
				return errors.New("has to be less than " + strconv.Itoa(value) + " characters long")
			}

		default:
			return errors.New("Error, option unsupported by this function or no options provided")
		}
	}
	return nil
}

//Min :
func Min(value interface{}) *v.Option {
	return &v.Option{Name: "minLength", Value: value}
}

//Max :
func Max(value interface{}) *v.Option {
	return &v.Option{Name: "maxLength", Value: value}
}
