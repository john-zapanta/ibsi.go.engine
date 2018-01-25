package dbase

import (
	"log"
	"strings"
)

type TParameters map[string]interface{}
type TDataTableRow map[string]interface{}
// type TDataTableRows []map[string]interface{}
type TDataTableRows []TDataTableRow
type TDataTableOutput map[string]interface{}

var Connections map[string]IConnection

func CloseConnections() {
	log.Println("Closing database connections...")
	for _, v := range Connections {
		v.CloseConnection()
    }
}

func ParseNames(name string, defaultConnection string) (string, string) {
	parts := strings.Split(name, ".")
	
	var cname, command string = defaultConnection, ""
	if len(parts) == 2 {
		cname = parts[0]
		command = parts[1]
	} else {
		command = parts[0]
	}
	
	return cname, command
}

func init() {
	Connections = make(map[string]IConnection)
	// Connections = config.Connections
}
