package dbase

import (
	"fmt"
	"log"
	"database/sql"
)

// type TNewConnection func(IConnection)
type TNewCommand func(ICommand)
type TGetCommandText func() string

type IConnection interface {
	CloseConnection()
	Execute(cmd string, params TParameters) (IQuery, error)
	GetCommand(name string) (ICommand, bool)
	HasCommand(name string) (bool)
	NewCommand(cname string, ctext string, ctype string, callback TNewCommand) ICommand
	NewQuery(cname string, init TNewCommand, initText TGetCommandText) ICommand
	Open(cmd string, params TParameters) (IQuery, error)
	OpenDataTable(cmd string, params TParameters) (IDataTable)
	RunQuery(query string, args ...interface{}) (*sql.Rows, error)
	SetDB(*sql.DB)
}

type Connection struct {
	db *sql.DB
    Name string
    Server string
    Port int
    Database string
    User string
    Password string
	Commands map[string]ICommand
}

func NewConnection(name string, server string, port int, database string, user string, password string) IConnection {
    c := &Connection{
		Name: name,
		Server: server, 
		Port: port, 
		Database: database, 
		User: user, 
		Password: password, 
		Commands: make(map[string]ICommand),
	}

	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=disable", c.Server, c.Port, c.Database, c.User, c.Password)

	// var err error
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("NewConnection()... sql.Open():", err.Error())
	} else {
		c.SetDB(db)
	}
	
	Connections[name] = c
	
	return c
}

func (c Connection) CloseConnection() {
	c.db.Close()
}

func (c Connection) Execute(cmd string, params TParameters) (IQuery, error) {
	if qry, err := NewQuery(c.Name, cmd, params); err == nil {	
		if err := qry.Execute(); err != nil {
			return nil, err
		} else {
			return qry, nil
		}
	} else {
		return nil, err
	}
}

func (s Connection) GetCommand(name string) (ICommand, bool) {
	cmd, ok := s.Commands[name]
	return cmd, ok
}

func (s Connection) HasCommand(name string) (bool) {
	_, ok := s.Commands[name]
	return ok
}

// func (s Connection) GetDB() *sql.DB {
	// return s.db
// }

func (s *Connection) NewCommand(cname string, ctext string, ctype string, callback TNewCommand) ICommand {
	cmd := NewCommand(s, cname, ctext, ctype)

	// will initialize the parameters in the callback
	if callback != nil { 
		callback(cmd)
	}

	s.Commands[cname] = cmd
	
	return cmd
}

func (s *Connection) NewQuery(cname string, init TNewCommand, initText TGetCommandText) ICommand {
	return s.NewCommand(cname, initText(), "text", init)
}

func (c Connection) Open(cmd string, params TParameters) (IQuery, error) {
	if qry, err := NewQuery(c.Name, cmd, params); err == nil {	
		if err := qry.Open(); err != nil {
			return nil, err
		} else {
			return qry, nil
		}
	} else {
		return nil, err
	}
}

func (c Connection) OpenDataTable(cmd string, params TParameters) (IDataTable) {
	if qry, err := NewQuery(c.Name, cmd, params); err == nil {	
		if err := qry.Open(); err != nil {
			return nil
		} else {
			return qry.GetDataTable(0)
		}
	} else {
		return nil
	}
}

func (s *Connection) RunQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

func (s *Connection) SetDB(db *sql.DB) {
	s.db = db
}
