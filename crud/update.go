package crud

import (
	"fmt"
	// "strings"
	"net/http"
	"encoding/json"
	"ibsi/dbase"
	"ibsi/utils"
	"ibsi/session"
)

func HandleUpdate(w http.ResponseWriter, r *http.Request, dsh CrudHandler) {
	
	sn := session.GetSession(r)
	
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	r.ParseForm() // Parses the request body
	
	data := make([]dbase.TParameters, 0)
    if err := json.Unmarshal([]byte(r.Form.Get("data")), &data); err != nil {
		http.Error(w, "Expecting input 'data'", 500)
		return
    }
	
	record := data[0] // we only expect an array with one item
	
	// sn := session.GetSession(r)
	record["visit_id"] = sn.VisitorId
	record["action_status_id"] = 0
	record["action_msg"] = ""
	
	// mode := params["mode"]
	mode := r.Form.Get("mode")
	if mode == "edit" {
		record["action"] = 10
	} else {
		record["action"] = 20
	}
	
	kv := utils.NewKeyValue()
	kv.Set("status", 0)
	kv.Set("message", "")
	kv.Set("mode", mode)
	
	un, uc := dbase.ParseNames(dsh.UpdateDataSource, "DBApp")	
	if cmd, ok := dbase.Connections[un].GetCommand(uc); ok {
		params := make(dbase.TParameters)
		for _, c := range cmd.GetParameters() {
			if v, ok := record[c.GetName()]; ok {
				params[c.GetName()] = v
			}
		}
		
		// params["visit_id"] = sn.VisitorId
		
		if dsh.OnInitCrudParams != nil {
			dsh.OnInitCrudParams("update", params, w, r)
		}
		
		if qry, err := cmd.Execute(params); err != nil {
			kv.Set("status", -2)
			kv.Set("message", err.Error())
			http.Error(w, err.Error(), 500)
			// return
		} else {
			if output, err := qry.GetOutput(); err == nil {
				kv.Set("status", output.GetIf("action_status_id", 0))
				kv.Set("message", output.GetIf("action_msg", ""))
				
				result := make(dbase.TDataTableOutput)
				output.EachRowItem(func(name string, value interface{}) {
					if name != "action_msg" && name != "action_status_id" {
						result[name] = value
					}
				})
				kv.Set("result", result)
			} else {
				kv.Set("output_error", err.Error())
			}
		}
	} else {
		kv.Set("status", -1)
		kv.Set("message", fmt.Sprintf("Command %s.%s was not found", un, uc))
	}
	
	w.Write([]byte(kv.Json()))
}
