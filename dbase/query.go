package dbase

import (
	"fmt"
	"errors"
	// "log"
	"strings"
	"strconv"
	"database/sql"
	"ibsi/utils"
)

type IQuery interface {
	Open() error
	Execute() error
	GetDataTable(index int) IDataTable
	// GetDataTables() []IDataTable
	GetOutput() (IDataTable, error)
	TablesCount() int
}

type Query struct {
	Command ICommand
	Parameters TParameters
	DataTables []IDataTable // array DataTables
	// exec bool
}

func NewQuery(cname string, cmd string, parameters TParameters) (IQuery, error) {
	command, ok := Connections[cname].GetCommand(cmd)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Command %s.%s not found!", cname, cmd))
	}
	
    qry := &Query{
		Command: command,
		Parameters: parameters,
		DataTables: make([]IDataTable, 0),
		// exec: false,
	}
	
	// fmt.Println("here at NewQuery")
	return qry, nil
}

func (s *Query) Execute() error {
	sqlValues, sql := s.PrepareSql()
	rows, err := s.Command.RunQuery(sql, sqlValues...)
    if err != nil {
		return err
    }
	defer rows.Close()
	
	// s.ProcessRows(rows)
	// return nil
	return s.ProcessRows(rows)
	// return s.ProcessOutputRow(rows)
}

func (s Query) GetDataTable(i int) IDataTable {
	return s.DataTables[i]
}

// func (s Query) GetDataTables() []IDataTable {
	// return s.DataTables
// }

func (s Query) GetOutput() (IDataTable, error) {
	if i := len(s.DataTables) - 1; i < 0 {
		return nil, errors.New("No outut table found")
	} else {
		return s.DataTables[i], nil
	}
}

func (s *Query) Open() error {
	sqlValues, sql := s.PrepareSql()
	rows, err := s.Command.RunQuery(sql, sqlValues...)
    if err != nil {
		return err
    }
	defer rows.Close()
	
	// s.ProcessRows(rows)
	// return nil
	return s.ProcessRows(rows)
	// return errors.New("Open(), not yet implemented")
}

// PrepareQuery() is not an interface method, internal only
func (s *Query) PrepareSql() ([]interface{}, string) {

	sqlParams := []string{}
	sqlOutput := []string{}
	sqlValues := []interface{}{}
	
	i, cmd := 1, s.Command.GetCommand()
	
	for k := range cmd.Parameters {
		p := cmd.Parameters[k].GetParameter()
		if p.Direction == "in" { 
			sqlParams = append(sqlParams, fmt.Sprintf("@%s=?", p.Name))
		} else {
			sqlParams = append(sqlParams, fmt.Sprintf("@%s=? OUTPUT", p.Name))
			sqlOutput = append(sqlOutput, fmt.Sprintf("%s=@p%d", p.Name, i))
		}
		
		val, ok := s.Parameters[p.Name]
		if ok {
			sqlValues = append(sqlValues, val)
		} else {
			sqlValues = append(sqlValues, p.Default)
		}		
		
		i++
	}
	
	var sql string
	if cmd.CommandType == "text" {
		sql = fmt.Sprintf("%s %s", cmd.CommandText, strings.Join(sqlParams, ","))
	} else {
		if len(sqlOutput) > 0 {
			sql = fmt.Sprintf("EXEC %s %s; SELECT _outrow_=1, %s", cmd.CommandText, strings.Join(sqlParams, ","), strings.Join(sqlOutput, ","))
		} else {
			sql = fmt.Sprintf("EXEC %s %s", cmd.CommandText, strings.Join(sqlParams, ","))
		}
	}
	
	return sqlValues, sql
}

// ProcessRows() is not an interface method, internal only
func (s *Query) ProcessRows(rows *sql.Rows) error {
    columns, err := rows.Columns()
    if err != nil {
		return err
    }
	
	count := len(columns)
    values := make([]interface{}, count)
    scanArgs := make([]interface{}, count)
    for i := range values {
        scanArgs[i] = &values[i]
    }

	data := NewDataTable()
	data.SetColumns(columns)

	for rows.Next() {
		err := rows.Scan(scanArgs...) // this is how you get the values with unknown columns, pass scanArgs array
		if err != nil {
			return err
		}
		
		row := make(TDataTableRow)
		for i, v := range values {
			switch v.(type) {
				case []byte:
					row[columns[i]], _ = strconv.ParseFloat(string(v.([]byte)), 64)
				default:
					row[columns[i]] = v
			}
		}
		
		data.Add(row)
	}
	
	s.DataTables = append(s.DataTables, data)
	
	if rows.NextResultSet() {
		columns, _ := rows.Columns()
		if utils.Contains(columns, "_outrow_") {
			s.ProcessOutputRow(rows)
		} else {
			s.ProcessRows(rows)
		}
	}
	
	return nil
}

func (s *Query) ProcessOutputRow(rows *sql.Rows) error {
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	
	// fmt.Println(s.Command.GetCommand().Name)		
	for rows.Next() {
		err := rows.Scan(scanArgs...) // this is how you get the values with unknown columns, pass scanArgs array
		if err != nil {
			return err
		}
		
		for i, v := range values {
			if columns[i] != "_outrow_" {			
				s.GetDataTable(0).SetOutput(columns[i], v)
			}
		}
	}
	
	return nil
}

func (s Query) TablesCount() int {
	return len(s.DataTables)
}
