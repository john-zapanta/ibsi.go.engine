package crud

import (
	"fmt"
	"net/http"
	// "github.com/gorilla/mux"
	"ibsi/dbase"
	"ibsi/utils"
	"ibsi/system"
)

// var Router *mux.Router

// type TRequestHandler func(http.ResponseWriter, *http.Request)
type TInitCrud func(map[string]bool)
type TInitCrudParams func(string, dbase.TParameters, http.ResponseWriter, *http.Request)
type TNewRecord func(string, map[string]interface{}, http.ResponseWriter, *http.Request)

type CrudHandler struct {
	Name string `json:"name"`
	Action string `json:"action"`
	KeyName string `json:"key_name"`
	Path string `json:"path"`
	ListDataSource string `json:"list_datasource"`
	EditDataSource string `json:"edit_datasource"`
	UpdateDataSource string `json:"update_datasource"`
	Rights map[string]bool `json:"rights"`
	OnInitCrud TInitCrud
	OnInitCrudParams TInitCrudParams
	OnNewRecord TNewRecord
}

func Handler(dsh CrudHandler) {
	system.Router.HandleFunc(fmt.Sprintf("/%s/get/list/%s", utils.Ifs(dsh.Path == "", "app", dsh.Path), dsh.Name), func(w http.ResponseWriter, r *http.Request) {
		// list handler
		HandleList(w, r, dsh)
	})
	
	system.Router.HandleFunc(fmt.Sprintf("/get/edit/%s", dsh.Name), func(w http.ResponseWriter, r *http.Request) {
		// edit handler
		HandleEdit(w, r, dsh)
	})
	
	system.Router.HandleFunc(fmt.Sprintf("/%s/get/edit/%s", utils.Ifs(dsh.Path == "", "app", dsh.Path), dsh.Name), func(w http.ResponseWriter, r *http.Request) {
		// edit handler
		HandleEdit(w, r, dsh)
	})
		
	system.Router.HandleFunc(fmt.Sprintf("/%s/get/update/%s", utils.Ifs(dsh.Path == "", "app", dsh.Path), dsh.Name), func(w http.ResponseWriter, r *http.Request) {
		// update handler
		HandleUpdate(w, r, dsh)
	})
		
	system.Router.HandleFunc(fmt.Sprintf("/%s/get/delete/%s", utils.Ifs(dsh.Path == "", "app", dsh.Path), dsh.Name), func(w http.ResponseWriter, r *http.Request) {
		// delete handler
		DeleteUpdate(w, r, dsh)
	})
}

// iniialize crud handlers defined in json file
func init() {
	// system.Router.HandleFunc("/app/get/list/{crud}", HandleListGeneric)
}