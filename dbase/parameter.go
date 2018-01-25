package dbase

type IParameter interface {
	GetDefaultValue() interface{}
	GetName() string
	GetParameter() Parameter
}

type Parameter struct { 
	Command ICommand
	Name string
	DataType string
	Direction string
	Size int
	Default interface{}
}

func NewParameter(cmd ICommand, pname string, ptype string, pdirection string, psize int, pdefault interface{}) IParameter {
	param := &Parameter{
		Command: cmd,
		Name: pname,
		DataType: ptype,
		Direction: pdirection,
		Size: psize,
		Default: pdefault,
	}
	
	return param
}

func (s Parameter) GetDefaultValue() interface{} {
	return s.Default
}

func (s Parameter) GetName() string {
	return s.Name
}

func (s Parameter) GetParameter() Parameter {
	return s
}

