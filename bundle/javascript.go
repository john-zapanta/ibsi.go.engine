package bundle

import (
    "fmt"
	// "encoding/json"
	"path/filepath"
	"io/ioutil"
	"bufio"
    "net/http"
    // "html/template"
    "strings"
    "log"
	"os"
	"github.com/gorilla/mux"
	// "github.com/gorilla/handlers"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
	
	// database 
	// "database/sql"
	// _ "github.com/denisenkom/go-mssqldb"
)

// func init() {
	// Router.HandleFunc("/loadscript/{path}/scripts/{script}", LoadScript).Host(config.Domain)
// }

func getFile(line string, ext string) (string) {
	var value string = strings.Split(line, "=")[1]
	var parts []string = strings.Split(value, "\\")
	var path, file, fullPath string = "", "", ""

	if len(parts) == 1 {
		path = "scripts"
		file = parts[0]
	} else {
		path = parts[0]
		file = parts[1]
		
		if path == "engine" || path == "app" {
			path = filepath.Join(path, "scripts")
		} else {
			path = filepath.Join("scripts", value)
			file = ""
		}
	}

	if file == "" {
		fullPath = filepath.Join(ResourcePath, path)
	} else {
		fullPath = filepath.Join(ResourcePath, path, file)
	}

	var extension = filepath.Ext(fullPath)
	// var name = fullPath[0:len(fullPath)-len(extension)]
	fullPath = fullPath[0:len(fullPath)-len(extension)] + ext
	
	// if filepath.Ext(fullPath) == "" {
		// fullPath += ext
	// }

	return fullPath
}

func parseTemplate(fileName string, files []string) ([]string) {
    file, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
	
    for scanner.Scan() {
		var line string = strings.TrimSpace(scanner.Text())
		if line != "" && strings.Index(line, "//") != 0 {
			if strings.Index(line, "##include=") == 0 {
				files = append(files, getFile(line, ".js"))
			} else if strings.Index(line, "##includetemplate=") == 0 {
				files = parseTemplate(getFile(line, ".jst"), files)
			} else if strings.Index(line, "##run=") == 0 {
				var value string = strings.Split(line, "=")[1]
				// var parts []string = strings.Split(value, "\\")
				files = append(files, fmt.Sprintf("##code=desktop.subModule.load(%s)", value))
				// files = append(files, fmt.Sprintf("##code=%s", value))
				// files = parseTemplate(getFile(line, ".jst"), files)
			}
		}
    }
	
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
	
	return files
}

func outputFile(w http.ResponseWriter, fileName string) {
	// fmt.Fprintf(w, fileName + "\n")
	
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	
	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)
	
	// b, err = m.Bytes(content, b)
	content, err = m.Bytes("text/javascript", content)
	if err != nil {
		panic(err)
	}
	
	w.Write(content)
}

func LoadScript(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var path, script string = vars["path"], vars["script"]
	var file = "=" + path + "\\" + script
	var files []string = parseTemplate(getFile(file, ".jst"), []string{})
	
	w.Header().Set("Content-Type", "application/x-javascript; charset=utf-8")
	
	for i := 0; i < len(files); i++ {
		var line string = files[i]
		if strings.Index(line, "##code=") == 0 {
			var value string = strings.Split(line, "=")[1]
			w.Write([]byte(value))
		} else {
			outputFile(w, line)
			}
	}
}
