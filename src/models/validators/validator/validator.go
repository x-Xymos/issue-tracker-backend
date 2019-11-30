package validator

import "reflect"

type Option struct {
	Name  string
	Value interface{}
}

func AddValidators(field string, structType reflect.Type, function func(interface{}, []Option) error) interface{} {

	p := reflect.New(structType).Interface()
	f := []func(interface{}, []Option) error{function}
	reflect.ValueOf(p).Elem().FieldByName(field).Set(reflect.ValueOf(&f).Elem())
	return p
	//todo make this safe https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field
}
