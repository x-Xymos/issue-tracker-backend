package strminlength

import (
	"errors"
	"math"
	"strconv"
	"testing"
)

func TestMinLengthValidator(t *testing.T) {

	inputTests := []Test{
		Test{map[string]interface{}{"String": "TestString", "Length": 3}, nil},
		Test{map[string]interface{}{"String": "TestString", "Length": 32}, errors.New("has to be more than " + strconv.Itoa(32) + " characters long")},
		Test{map[string]interface{}{"String": "TestString", "Length": 0}, errors.New("Error, passed length parameter should be more than 0 in minLength validator")},
		Test{map[string]interface{}{"String": "TestString", "Length": -1}, errors.New("Error, passed length parameter should be more than 0 in minLength validator")},
		Test{map[string]interface{}{"String": 1, "Length": 16}, errors.New("Error casting input value in minLength validator")},
		Test{map[string]interface{}{"String": "TestString", "Length": math.MaxInt32}, errors.New("has to be more than " + strconv.Itoa(math.MaxInt32) + " characters long")},
	}

	for _, test := range inputTests {

		input, _ := test.Input.(map[string]interface{})["String"]
		optionLength, _ := test.Input.(map[string]interface{})["Length"].(int)
		opts := Options(optionLength)

		err := Validator(input, opts)

		if test.Expected == nil && err != nil || test.Expected != nil && err == nil {
			t.Errorf("TestMaxLengthValidator failed. expected Validator to return %v, got %v", test.Expected, err)
		} else if test.Expected != nil && err != nil {
			if test.Expected.(error).Error() != err.Error() {
				t.Errorf("TestMaxLengthValidator failed. expected Validator to return %v, got %v", test.Expected, err)
			}
		}
	}
}

func TestMinLengthOptions(t *testing.T) {

	lengthTests := []Test{
		Test{0, 1},
		Test{2, 1},
		Test{3, 1},
	}

	for _, test := range lengthTests {
		input, _ := test.Input.(int)
		expected, _ := test.Expected.(int)
		opts := Options(input)
		if len(*opts) != expected {
			t.Errorf("TestMinLengthOptions failed. expected options array to be of length %v, got %v", expected, len(*opts))
		}
	}

	valueTests := []Test{
		Test{1, 1},
		Test{2, 2},
		Test{math.MaxInt32, math.MaxInt32},
	}

	for _, test := range valueTests {
		input, _ := test.Input.(int)
		expected, _ := test.Expected.(int)
		opts := Options(input)
		if (*opts)[0].Value != expected {
			t.Errorf("TestMinLengthOptions failed. expected options array Value at index 0 to be %v, got %v", expected, (*opts)[0].Value)
		}
	}

	nameTests := []Test{
		Test{0, "minLength"},
		Test{1, "minLength"},
		Test{math.MaxInt32, "minLength"},
	}

	for _, test := range nameTests {
		input, _ := test.Input.(int)
		expected, _ := test.Expected.(string)
		opts := Options(input)
		if (*opts)[0].Name != expected {
			t.Errorf("TestMinLengthOptions failed. expected options array Name at index 0 to be %v, got %v", expected, (*opts)[0].Name)
		}
	}

}

type Test struct {
	Input    interface{}
	Expected interface{}
}
