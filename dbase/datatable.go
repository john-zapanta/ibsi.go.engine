package dbase

import (
	"encoding/json"
)

// type TEachOutput func(name string, value interface{})
type TEachRowItem func(string, interface{})

type IDataTable interface {
	Add(row TDataTableRow)
	EachRowItem(TEachRowItem)
	Get(name string) interface{}
	GetColumns() []string
	GetIf(name string, def interface{}) interface{}
	GetOutput(name string) (interface{}, bool)
	// GetOutputs() TDataTableOutput
	GetRows() TDataTableRows
	Json(min bool) []byte
	RowCount() int
	SetColumns(columns []string)
	SetOutput(name string, value interface{})
}

type DataTable struct {
	rowNo int
	columns []string
	data TDataTableRows
	output TDataTableOutput
}

func NewDataTable() IDataTable {
    return &DataTable{
		rowNo: 0,
		columns: make([]string, 0),
		data: make(TDataTableRows, 0),
		output: make(TDataTableOutput),
	}
}

func (s *DataTable) Add(row TDataTableRow) {
	s.data = append(s.data, row)
}

func (s *DataTable) EachRowItem(callback TEachRowItem) {
	for n := range s.data[s.rowNo] {
		if n != "_outrow_" {
			callback(n, s.data[s.rowNo][n])
		}
	}
}

func (s DataTable) Get(name string) interface{} {
	return s.data[s.rowNo][name]
}

func (s DataTable) GetIf(name string, def interface{}) interface{} {
	v, ok := s.data[s.rowNo][name]
	if ok {
		return v
	} else {
		return def
	}
}

func (s DataTable) GetColumns() []string {
	return s.columns
}

func (s DataTable) GetOutput(name string) (interface{}, bool) {
	v, ok := s.output[name]
	return v, ok
}

func (s DataTable) GetRows() TDataTableRows {
	return s.data
}

func (s *DataTable) Json(min bool) []byte {
	tab := "\t"
	if min {
		tab = ""
	}
	
	json, err := json.MarshalIndent(s.data, "", tab) // \t is formatted, remove it to minify
    if err != nil {
        return []byte("{}")
    }
	
	return json
}

func (s DataTable) RowCount() int {
	return len(s.data)
}

func (s *DataTable) SetColumns(columns []string) {
	s.columns = columns
}

func (s *DataTable) SetOutput(name string, value interface{}) {
	s.output[name] = value
}
