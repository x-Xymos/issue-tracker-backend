package stringlength

import (
	"errors"
	"issue-tracker-backend/src/models/validators/validator"
	"strconv"
)

//Validator :
//Options : minLength int, maxLength int
func Validator(input interface{}, options []validator.Option) error {
	for _, v := range options {
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

			if !(len(input) < value) {
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

			if !(len(input) > value) {
				return errors.New("has to be less than " + strconv.Itoa(value) + " characters long")
			}

		}
	}
	return nil
}
