package bundle

import (
    // "fmt"
	// "path/filepath"
	// "io/ioutil"
	// "bufio"
    // "net/http"
    // "strings"
    // "log"
	"ibsi/utils"
	// "github.com/tdewolff/minify"
	// "github.com/tdewolff/minify/js"
)

type FileData struct {
	Keys utils.IKeyValue
}

func (fd *FileData) AddParameter(name string, value interface{}, overwrite bool) {
	if overwrite {
		fd.Keys.Set(name, value)
	} else if !fd.Keys.ContainsKey(name) {
		fd.Keys.Set(name, value)
	}			
}

func (fd *FileData) AddParameterEx(values []string, overwrite bool) {
	fd.AddParameter(values[0], values[1], overwrite)
}

func (fd *FileData) Append(variables utils.IKeyValue) {
}

func NewFileData(initValues string) FileData {
	fd := FileData{
		Keys: utils.NewKeyValueEx(initValues),
	}
	
	return fd
}

func NewFileDataEx(fileName string, initValues string) FileData {
	fd := NewFileData(initValues)
	fd.Keys.Set("file", fileName)
	
	return fd
}
