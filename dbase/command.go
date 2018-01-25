package dbase

import (
	"database/sql"
)

type ICommand interface {
	Execute(params TParameters) (IQuery, error)
	GetCommand() Command
	GetCommandText() string
	GetParameter(name string) (IParameter, bool)
	GetParameters() map[string]IParameter
	NewParameter(pname string, ptype string, pdirection string, psize int, pdefault interface{}) IParameter
	RunQuery(query string, args ...interface{}) (*sql.Rows, error)
}

type Command struct {
	Connection IConnection
	Name string
	CommandType string
	CommandText string
	Parameters map[string]IParameter
}

func NewCommand(connection IConnection, cname string, ctext string, ctype string) ICommand {
    return &Command{
		Connection: connection,
		Name: cname,
		CommandType: ctype,
		CommandText: ctext,
		Parameters: make(map[string]IParameter),
	}
}

// func (s Command) AddParameter(name string, param IParameter) {
	// s.Parameters[name] = param
// }

func (s Command) Execute(params TParameters) (IQuery, error) {
	// return s.Connection.Execute(s.CommandText, params)
	return s.Connection.Execute(s.Name, params)
}

func (s Command) GetCommand() Command {
	return s
}

func (s Command) GetCommandText() string {
	return s.CommandText
}

func (s Command) GetParameter(name string) (IParameter, bool) {
	p, ok := s.Parameters[name]
	return p, ok
}

func (s Command) GetParameters() map[string]IParameter {
	return s.Parameters
}

func (s *Command) NewParameter(pname string, ptype string, pdirection string, psize int, pdefault interface{}) IParameter {
	p := NewParameter(s, pname, ptype, pdirection, psize, pdefault)
	s.Parameters[pname] = p
	return p
}

func (s Command) RunQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return s.Connection.RunQuery(query, args...)
}

