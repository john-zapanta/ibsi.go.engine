package bundle

import (
    "fmt"
	"path/filepath"
	"io/ioutil"
	"bufio"
    "net/http"
    "strings"
    "log"
	"os"
	// "github.com/gorilla/mux"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"	
)

var ResourcePath string

type Bundler struct {
	// Actions map[string]
}

// type Bundler struct {
	// Actions map[string]
// }

func GetFile(line string, ext string) (string) {
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
	fullPath = fullPath[0:len(fullPath)-len(extension)] + ext

	return fullPath
}

func ParseTemplate(fileName string, files []string) ([]string) {
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
				files = append(files, fmt.Sprintf("##code=desktop.subModule.load(%s)", value))
			}
		}
    }
	
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
	
	return files
}

func OutputFile(w http.ResponseWriter, fileName string) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	
	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)
	
	content, err = m.Bytes("text/javascript", content)
	if err != nil {
		panic(err)
	}
	
	w.Write(content)
}
