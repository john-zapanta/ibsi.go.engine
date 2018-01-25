package crud

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"ibsi/utils"
	"ibsi/dbase"
	"ibsi/session"
)

func HandleLookup(w http.ResponseWriter, r *http.Request) {
	
	// system.Router.HandleFunc("/app/get/list/{list}", crud.HandleListGeneric)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	// vars := mux.Vars(r)
	// lookup := vars["list"]
	lookup := r.URL.Query().Get("name")
	// fmt.Println("lookup", lookup)
	
	r.ParseForm() // Parses the request body
	sn := session.GetSession(r)
	
	data := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(r.Form.Get("qry")), &data)
	if err != nil {
		fmt.Println(err.Error(), ", expecting input 'qry'")
	}

	params := make(dbase.TParameters)
	if err == nil {
		for k, v := range data[0] {
			params[k] = v
		}
	}

	params["action"] = 1
	params["visit_id"] = sn.VisitorId
	// if dsh.OnInitCrudParams != nil {
		// dsh.OnInitCrudParams("list", params, w, r)
	// }

	cname, command := dbase.ParseNames(lookup, "DBApp")
	if qry, err := dbase.Connections[cname].Open(command, params); err == nil {
		data, kv := qry.GetDataTable(0), utils.NewKeyValue()
		kv.Set("status", 0)
		kv.Set("message", "")
		
		if v, ok := data.GetOutput("page"); ok {
			kv.Set("page", v)
		} else {
			kv.Set("page", 1)
		}
		
		if v, ok := data.GetOutput("row_count"); ok {
			kv.Set("row_count", v)
		}
		
		if v, ok := data.GetOutput("page_count"); ok {
			kv.Set("page_count", v)
		}
		
		crud := make(map[string]bool)
		crud["view"] = true
		crud["add"] = false
		crud["edit"] = false
		crud["delete"] = false
		
		kv.Set("crud", crud)
		
		kv.Set("table_count", qry.TablesCount())
		for i := 0; i < qry.TablesCount(); i++ {
			kv.Set(fmt.Sprintf("data_%d", i), qry.GetDataTable(i).GetRows())
		}
		
		w.Write([]byte(kv.Json()))

	} else {
		http.Error(w, err.Error(), 500)
	}
}

func HandleListGeneric(w http.ResponseWriter, r *http.Request) {
	
	// system.Router.HandleFunc("/app/get/list/{list}", crud.HandleListGeneric)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	vars := mux.Vars(r)
	raw, err := ioutil.ReadFile(fmt.Sprintf("crud/%s.json", vars["list"]))
	if err != nil {
		fmt.Println(err.Error())
		// os.Exit(1)
	}
	
	var dsh CrudHandler
	err = json.Unmarshal(raw, &dsh)
	if err != nil {
		fmt.Println(err.Error())
	}
	
	r.ParseForm() // Parses the request body
	sn := session.GetSession(r)
	
	data := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(r.Form.Get("qry")), &data)
	if err != nil {
		fmt.Println(err.Error(), ", expecting input 'qry'")
	}

	params := make(dbase.TParameters)
	if err == nil {
		for k, v := range data[0] {
			params[k] = v
		}
	}

	params["visit_id"] = sn.VisitorId
	// if dsh.OnInitCrudParams != nil {
		// dsh.OnInitCrudParams("list", params, w, r)
	// }

	cname, command := dbase.ParseNames(dsh.ListDataSource, "DBApp")
	if qry, err := dbase.Connections[cname].Open(command, params); err == nil {
		data, kv := qry.GetDataTable(0), utils.NewKeyValue()
		kv.Set("status", 0)
		kv.Set("message", "")
		
		if v, ok := data.GetOutput("page"); ok {
			kv.Set("page", v)
		} else {
			kv.Set("page", 1)
		}
		
		if v, ok := data.GetOutput("row_count"); ok {
			kv.Set("row_count", v)
		}
		
		if v, ok := data.GetOutput("page_count"); ok {
			kv.Set("page_count", v)
		}
		
		// crud := make(map[string]bool)
		// crud["view"] = true
		// crud["add"] = true
		// crud["edit"] = true
		// crud["delete"] = true
		// if dsh.OnInitCrud != nil {
			// dsh.OnInitCrud(crud)
		// }
		
		kv.Set("crud", dsh.Rights)
		
		kv.Set("table_count", qry.TablesCount())
		for i := 0; i < qry.TablesCount(); i++ {
			kv.Set(fmt.Sprintf("data_%d", i), qry.GetDataTable(i).GetRows())
		}
		
		w.Write([]byte(kv.Json()))

	} else {
		http.Error(w, err.Error(), 500)
	}
	
}

func HandleList(w http.ResponseWriter, r *http.Request, dsh CrudHandler) {
	
	sn := session.GetSession(r)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	r.ParseForm() // Parses the request body

	data := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(r.Form.Get("qry")), &data)
    if err != nil {
        fmt.Println(err.Error(), ", expecting input 'qry'")
    }

	params := make(dbase.TParameters)
	
	if err == nil {
		for k, v := range data[0] {
			params[k] = v
			// if k == "visit_id" {
				// params[k] = sn.VisitorId
			// }
		}
	}

	params["visit_id"] = sn.VisitorId
	if dsh.OnInitCrudParams != nil {
		dsh.OnInitCrudParams("list", params, w, r)
	}

	cname, command := dbase.ParseNames(dsh.ListDataSource, "DBApp")
	if qry, err := dbase.Connections[cname].Open(command, params); err == nil {
		data, kv := qry.GetDataTable(0), utils.NewKeyValue()
		kv.Set("status", 0)
		kv.Set("message", "")
		
		if v, ok := data.GetOutput("page"); ok {
			kv.Set("page", v)
		} else {
			kv.Set("page", 1)
		}
		
		if v, ok := data.GetOutput("row_count"); ok {
			kv.Set("row_count", v)
		}
		
		if v, ok := data.GetOutput("page_count"); ok {
			kv.Set("page_count", v)
		}
		
		crud := make(map[string]bool)
		crud["view"] = true
		crud["add"] = true
		crud["edit"] = true
		crud["delete"] = true
		if dsh.OnInitCrud != nil {
			dsh.OnInitCrud(crud)
		}
		
		kv.Set("crud", crud)
		
		kv.Set("table_count", qry.TablesCount())
		for i := 0; i < qry.TablesCount(); i++ {
			kv.Set(fmt.Sprintf("data_%d", i), qry.GetDataTable(i).GetRows())
		}
		//kv.Set("data_0", data.GetRows())
		
		w.Write([]byte(kv.Json()))

	} else {
		http.Error(w, err.Error(), 500)
	}
}
