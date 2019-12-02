package unique

import (
	"errors"
	v "issue-tracker-backend/src/models/validators/validator"
	"strconv"
	"unicode/utf8"
)

//Validator :
func Validator(input interface{}, options *[]*v.Option) error {
	for _, v := range *options {
		
		// default:
		// 	return errors.New("Error, option unsupported by this function or no options provided")
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
