package template

import (
	"net/http"
	// "github.com/gorilla/mux"
)

func HandleDynamicTemplate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// vars := mux.Vars(r)
	// pid := vars["pid"]
	
}
