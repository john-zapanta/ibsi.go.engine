package bundle

import (
    // "fmt"
    "os"
    "net/http"
	"path/filepath"
	"github.com/gorilla/mux"
)

type CssLoader struct {
	Module string
	Extension string
	TemplateExtension string
	SaveFilePath string
	ContentType string
	CompileFileName string
	
	w http.ResponseWriter
	r *http.Request
}

func (lo *CssLoader) Init() {

	// vars := mux.Vars(lo.r)
	// var path, css string = vars["path"], vars["css"]
	
	// names := make([]string, 0)

	// lo.w.Write([]byte(
	// lo.CompileFileName = fmt.Fprintf("%s/%s/%s", path, lo.Module, css)
	
	// lo.w.Write([]byte(lo.CompileFileName))
	// names := append.
}

func (lo *CssLoader) Output(line string) {
	lo.w.Write([]byte(line))
}

func (lo *CssLoader) Load() {
	vars := mux.Vars(lo.r)
	path, css := vars["path"], vars["css"]

	fileName := filepath.Join(ResourcePath, path, lo.Module, css)
	
	var ext string
	if ext = filepath.Ext(fileName); ext == "" {
		fileName += lo.TemplateExtension
	} else {
		fileName = fileName[0:len(fileName)-len(ext)] + lo.TemplateExtension
	}
	
	fullPath := fileName
	
	bundler := NewBundlerCss(lo)	
	bundler.Variables.Set("pid", "")
	bundler.Variables.Set("keyid", "")
	
	if pid := lo.r.URL.Query().Get("pid"); pid == "" {
		bundler.Variables.Set("pid", pid)
	}
	
	if keyid := lo.r.URL.Query().Get("keyid"); keyid == "" {
		bundler.Variables.Set("keyid", keyid)
	}
	
	ext = filepath.Ext(fullPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// fmt.Println("here 1") // debug
		// log := NewFileDataEx("file=" + fileName[0:len(fullPath)-len(ext)] + lo.Extension, "")
		log := NewFileDataEx(fileName[0:len(fullPath)-len(ext)] + lo.Extension, "")
		log.Keys.Append(bundler.Variables);
		bundler.Log(log);
		
		// bundler.ParseFile(fileName[0:len(fullPath)-len(ext)] + lo.Extension, nil) // debug
		// bundler.ParseFile(fileName[0:len(fullPath)-len(ext)] + lo.Extension, nil) // debug
	} else {
		// fmt.Println("here 2") // debug
		bundler.ParseFile(fullPath, nil)
	}
	
	
	bundler.BundleFiles(lo.r)
	
	// lo.w.Write([]byte(fullPath))  // debug
}

func NewCssLoader(w http.ResponseWriter, r *http.Request) CssLoader {
	css := CssLoader{
		w: w,
		r: r,
		Module: "css",
		Extension: ".css",
		TemplateExtension: ".csst",
		SaveFilePath: "",
		ContentType: "text/css",
		CompileFileName: "",
	}
	
	css.Init()
	
	return css
}