package easy

import (
	"net/http"
	"ibsi/dbase"
	"ibsi/crud"
	"ibsi/system"
)

func init() {

	crud.Handler(crud.CrudHandler {
		Name: "sys-rights",
		Action: "sys-rights",
		Path: "engine",
		KeyName: "id",
		ListDataSource: "DBSecure.GetRights",
		UpdateDataSource: "DBSecure.AddRights",
		OnInitCrudParams: func(mode string, params dbase.TParameters, w http.ResponseWriter, r *http.Request) {
			// params["visit_id"] = 13389  // to-do: remove later, for testing only
			
			if mode == "list" {
				params["application_id"] = system.Settings.AppID
				if r.URL.Query().Get("lookup") != "" {  // to-do: there is a better way of of doing this,
					params["mode"] = 50
				}
			// } else if mode == "edit" || mode == "new" { // to-do: this is not necessary in the updated engine
				// params["mode"] = 10
			}
		},
		OnNewRecord: func(mode string, row map[string]interface{}, w http.ResponseWriter, r *http.Request) {
			row["application_id"] = system.Settings.AppID
			// row["action_type_id"] = 10
			row["status_code_id"] = 10
		},
	})
	
	dbase.Connections["DBSecure"].NewCommand("AddRights", "AddRights", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "inout", 0, 0)
		cmd.NewParameter("code", "string", "in", 20, "")
		cmd.NewParameter("rights", "string", "in", 100, "")
		cmd.NewParameter("description", "string", "in", 200, "")
		cmd.NewParameter("application_id", "int", "in", 0, 0)
		cmd.NewParameter("position", "int", "in", 0, 0)
		cmd.NewParameter("status_code_id", "int", "in", 0, 0)
		cmd.NewParameter("action", "int", "in", 0, 10)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
		cmd.NewParameter("action_status_id", "int", "inout", 0, 0)
		cmd.NewParameter("action_msg", "string", "inout", 200, "")
	}) 

	dbase.Connections["DBSecure"].NewCommand("GetRights", "GetRights", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "in", 0, 0)
		cmd.NewParameter("ids", "string", "in", 200, "")
		cmd.NewParameter("filter", "string", "in", 200, "")
		cmd.NewParameter("action", "int", "in", 0, 0)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
		cmd.NewParameter("application_id", "int", "in", 0, 0)
		cmd.NewParameter("page", "int", "in", 0, 1)
		cmd.NewParameter("pagesize", "int", "in", 0, 25)
		cmd.NewParameter("row_count", "int", "inout", 0, 0)
		cmd.NewParameter("page_count", "int", "inout", 0, 0)
		cmd.NewParameter("sort", "string", "in", 200, "position")
		cmd.NewParameter("order", "string", "in", 10, "asc")
	}) 

	dbase.Connections["DBSecure"].NewCommand("GetActionRights", "GetActionRights", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("action_id", "int", "in", 0, 0)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	}) 
}
