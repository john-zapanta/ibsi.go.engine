package crud

import (
	"fmt"
	"net/http"
	"ibsi/system"
	"ibsi/dbase"
	"ibsi/utils"
	"ibsi/session"
)

type TInitCommand func(CommandHandler, dbase.TParameters, http.ResponseWriter, *http.Request)
// type TInitCrudParams func(string, dbase.TParameters, http.ResponseWriter, *http.Request)

type CommandHandler struct {
	Name string `json:"name"`
	Action string `json:"action"`
	Path string `json:"path"`
	DataSource string `json:"datasource"`
	OnInitCommand TInitCommand
	// OnInitCrudParams TInitCrudParams
	// OnNewRecord TNewRecord
}

func CommandHandle(dsh CommandHandler) {
	system.Router.HandleFunc(fmt.Sprintf("/%s/command/%s", utils.Ifs(dsh.Path == "", "app", dsh.Path), dsh.Name), func(w http.ResponseWriter, r *http.Request) {
		HandleCommand(w, r, dsh)
	})
}

func HandleCommand(w http.ResponseWriter, r *http.Request, dsh CommandHandler) {

	r.ParseForm() // Parses the request body
	sn := session.GetSession(r)
	
	// record := make([]dbase.TParameters, 0)[0]
	record := make(dbase.TParameters, 0)
	
	// sn := session.GetSession(r)
	record["visit_id"] = sn.VisitorId
	record["action_status_id"] = 0
	record["action_msg"] = ""

	var status int64 = 0
	message := ""
	
	kv := utils.NewKeyValue()
	un, uc := dbase.ParseNames(dsh.DataSource, "DBApp")	
	if cmd, ok := dbase.Connections[un].GetCommand(uc); ok {
		params := make(dbase.TParameters)
		for _, c := range cmd.GetParameters() {
			if v, ok := record[c.GetName()]; ok {
				params[c.GetName()] = v
			}
		}
		
		if dsh.OnInitCommand != nil {
			dsh.OnInitCommand(dsh, params, w, r)
		}
		
		if qry, err := cmd.Execute(params); err != nil {
			status = -500
			message = err.Error()	
			http.Error(w, message, 500)			
		} else {
			if output, err := qry.GetOutput(); err == nil {
				status = output.GetIf("action_status_id", 0).(int64)
				message = output.GetIf("action_msg", "").(string)
				
				// result := make(dbase.TDataTableOutput)
				output.EachRowItem(func(name string, value interface{}) {
					if name != "action_msg" && name != "action_status_id" {
						kv.Set(name, value)
					}
				})
				
				// message = utils.Ifs(status == 0, "", "Database error: ") + message
			} else {
				status = 2
				message = err.Error()
			}
		}
	} else {
		status = 1
		message = fmt.Sprintf("Command %s.%s was not found", un, uc)
	}
	
	kv.Set("status", status)
	kv.Set("message", message)
	
	if status != 0 {
		
	}
	
	w.Write([]byte(kv.Json()))
}
