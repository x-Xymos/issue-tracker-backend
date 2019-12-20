package validator

type Validators *[]*Function

//Function :
type Function struct {
	Validate func(interface{}, *[]*Option) error
	Options  *[]*Option
}

//Option :
type Option struct {
	Name  string
	Value interface{}
}

//Create :
func Create(function func(interface{}, *[]*Option) error, options *[]*Option) *Function {
	return &Function{Validate: function, Options: options}
}

//Assign :
//Returns an array of validation functions that can be used on a variable
func Assign(functions ...*Function) Validators {
	funcArr := &[]*Function{}
	for _, v := range functions {
		*funcArr = append(*funcArr, v)
	}
	return funcArr
}

//Options :
func Options(opts ...*Option) *[]*Option {
	optArr := &[]*Option{}
	for _, v := range opts {
		*optArr = append(*optArr, &Option{Name: v.Name, Value: v.Value})
	}
	return optArr
}
