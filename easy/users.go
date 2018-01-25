package easy

import (
	"net/http"
	"ibsi/dbase"
	"ibsi/crud"
)

func init() {

	crud.Handler(crud.CrudHandler {
		Name: "sys-users",
		Action: "sys-users",
		Path: "engine",
		KeyName: "id",
		ListDataSource: "DBSecure.GetUsers",
		UpdateDataSource: "DBSecure.AddUser",
		OnInitCrudParams: func(mode string, params dbase.TParameters, w http.ResponseWriter, r *http.Request) {
			if mode == "edit" || mode == "new" {
				params["mode"] = 10
			}
		},
		OnNewRecord: func(mode string, row map[string]interface{}, w http.ResponseWriter, r *http.Request) {
		},
	})
	
	dbase.Connections["DBSecure"].NewCommand("AddUser", "AddUser", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "inout", 0, 0)
		cmd.NewParameter("organisation_id", "int", "in", 0, 0)
		cmd.NewParameter("user_name", "string", "in", 200, "")
		cmd.NewParameter("last_name", "string", "in", 60, "")
		cmd.NewParameter("middle_name", "string", "in", 60, "")
		cmd.NewParameter("first_name", "string", "in", 60, "")
		cmd.NewParameter("gender", "string", "in", 1, "")
		cmd.NewParameter("dob", "datetime", "in", 0, nil)
		cmd.NewParameter("email", "string", "in", 200, "")
		cmd.NewParameter("status_code_id", "int", "in", 0, 0)
		cmd.NewParameter("roles", "string", "in", 100, "")
		cmd.NewParameter("action", "int", "in", 0, 10)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
		cmd.NewParameter("action_status_id", "int", "inout", 0, 0)
		cmd.NewParameter("action_msg", "string", "inout", 200, "")
	}) 

	dbase.Connections["DBSecure"].NewCommand("lookup_users", "GetUsersLookup", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("ids", "string", "in", 100, "")
		cmd.NewParameter("filter", "string", "in", 200, "")
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
		cmd.NewParameter("sort", "string", "in", 200, "full_name")
		cmd.NewParameter("order", "string", "in", 10, "asc")
	}) 

	dbase.Connections["DBSecure"].NewCommand("GetUsers", "GetUsers", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "in", 0, 0)
		cmd.NewParameter("ids", "string", "in", 200, "")
		cmd.NewParameter("filter", "string", "in", 200, "")
		cmd.NewParameter("mode", "int", "in", 0, 0)
		cmd.NewParameter("page", "int", "in", 0, 1)
		cmd.NewParameter("pagesize", "int", "in", 0, 25)
		cmd.NewParameter("row_count", "int", "inout", 0, 0)
		cmd.NewParameter("page_count", "int", "inout", 0, 0)
		cmd.NewParameter("sort", "string", "in", 200, "user_name")
		cmd.NewParameter("order", "string", "in", 10, "asc")
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	}) 

	dbase.Connections["DBSecure"].NewCommand("GetUserSessionInfo", "System_GetUserSessionInfo", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("visit_id", "int", "in", 0, nil)
	})

	dbase.Connections["DBSecure"].NewCommand("GetMyRights", "System_GetMyRights", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("action_id", "int", "in", 0, nil)
		cmd.NewParameter("visit_id", "int", "in", 0, nil)
		cmd.NewParameter("error_log_id", "int", "in", 0, 0)
		cmd.NewParameter("verbose", "int", "in", 0, 1)
	})

	dbase.Connections["DBSecure"].NewCommand("GetUserDetails", "System_GetUserDetails", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "in", 0, 0)
	})

	dbase.Connections["DBSecure"].NewCommand("GetUserActions", "System_GetUserActions", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("user_id", "int", "in", 0, 0)
		cmd.NewParameter("role_id", "int", "in", 0, 0)
		cmd.NewParameter("visit_id", "int", "in", 0, nil)
		cmd.NewParameter("error_log_id", "int", "inout", 0, 0)
		cmd.NewParameter("verbose", "int", "in", 0, 1)
	})
}
