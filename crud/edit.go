package crud

import (
	"fmt"
	"strings"
	"net/http"
	"ibsi/utils"
	"ibsi/dbase"
	"ibsi/session"
)

func HandleEdit(w http.ResponseWriter, r *http.Request, dsh CrudHandler) {
	
	sn := session.GetSession(r)
	
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	r.ParseForm() // Parses the request body
	
	key := r.Form.Get(dsh.KeyName)
	mode := utils.Ifs(key == "new", "new", "edit")

	var params dbase.TParameters = make(dbase.TParameters)
	params["action"] = 10  // action: 0:list, 1:lookup, 10:for editing, 20:for new record, 50:fetch updated data
	
	for k, v := range r.Form {
		if mode == "edit" || (mode == "new" && k != dsh.KeyName) {
			params[k] = strings.Join(v, "")
		}
    }
	
	params["visit_id"] = sn.VisitorId
	
	if dsh.OnInitCrudParams != nil {
		dsh.OnInitCrudParams(mode, params, w, r)
	}

	cname, command := dbase.ParseNames(utils.Ifs(dsh.EditDataSource == "", dsh.ListDataSource, dsh.EditDataSource), "DBApp")
	if qry, err := dbase.Connections[cname].Open(command, params); err == nil {
		data, kv := qry.GetDataTable(0), utils.NewKeyValue()
		kv.Set("status", 0)
		kv.Set("message", "")
		kv.Set("mode", mode)
		
		crud := make(map[string]bool)
		crud["view"] = true
		crud["add"] = true
		crud["edit"] = true
		crud["delete"] = true
		if dsh.OnInitCrud != nil {
			dsh.OnInitCrud(crud)
		}
		
		kv.Set("crud", crud)
		
		if mode == "new" {		
			row := make(dbase.TDataTableRow)
			un, uc := dbase.ParseNames(dsh.UpdateDataSource, "DBApp")
			cmd, ok := dbase.Connections[un].GetCommand(uc)
			if ok {
				for _, c := range data.GetColumns() {
					if val, ok := cmd.GetParameter(c); ok {
						row[c] = val.GetDefaultValue()
					} else {
						row[c] = nil
					}
				}
				
				if dsh.OnNewRecord != nil {
					dsh.OnNewRecord(mode, row, w, r)
				}
			}
		
			data.Add(row)
		}

		kv.Set("edit", data.GetRows())
		w.Write([]byte(kv.Json()))
	} else {
		http.Error(w, fmt.Sprintf("Error opening %s.%s", cname, command), 500)
	}
}
