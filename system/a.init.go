package system

import (
	"log"
	"ibsi/dbase"
)

func init() {
	connection, ok := dbase.Connections["DBSecure"]
	if !ok {
		log.Println("DBSecure connection was not found")
		return
	}
	
	//***************************************************************************************************
	//	SYSTEM 
	//***************************************************************************************************
	connection.NewCommand("AddVisit", "System_AddVisit", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("visit_id", "int", "inout", 0, 0)
		cmd.NewParameter("application_id", "int", "in", 0, 0)
		cmd.NewParameter("session_id", "string", "in", 48, "")
		cmd.NewParameter("method", "string", "in", 10, "")
		cmd.NewParameter("local_ip", "string", "in", 20, "")
		cmd.NewParameter("remote_ip", "string", "in", 20, "")
		cmd.NewParameter("remote_host", "string", "in", 100, "")
		cmd.NewParameter("user_agent", "string", "in", 100, "")
		cmd.NewParameter("referrer_url", "string", "in", 200, "")
		cmd.NewParameter("request_url", "string", "in", 200, "")
	})

	connection.NewCommand("Login", "System_Login", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("user_id", "int", "inout", 0, -20)
		cmd.NewParameter("user_name", "string", "in", 200, "")
		cmd.NewParameter("password", "string", "in", 200, "")
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	})

	connection.NewCommand("Logout", "System_Logout", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	})

	//====================================================================================================
	//	PERMISSIONS
	//====================================================================================================
	connection.NewCommand("AddPermission", "AddPermission", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("role_id", "int", "in", 0, 0)
		cmd.NewParameter("action_id", "int", "in", 0, 0)
		cmd.NewParameter("permissions", "string", "in", 100, "")

		cmd.NewParameter("action", "int", "in", 0, 10)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
		cmd.NewParameter("action_status_id", "int", "inout", 0, 0)
		cmd.NewParameter("action_msg", "string", "inout", 200, "")
	}) 

	connection.NewCommand("GetManagePermissions", "GetManagePermissions", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("role_id", "int", "in", 0, 0)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	}) 

	connection.NewCommand("GetPermissions", "GetPermissions", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("role_id", "int", "in", 0, 0)
		cmd.NewParameter("action_id", "int", "in", 0, 0)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	}) 

	connection.NewCommand("GetAllowAction", "GetAllowAction", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("action_code", "string", "in", 20, "")
		cmd.NewParameter("allow", "bool", "inout", 1, 0)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	}) 

	connection.NewCommand("GetMyPermission", "GetMyPermission", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("action_code", "string", "in", 20, "")
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	}) 

	//====================================================================================================
	// 
	//====================================================================================================
	connection.NewCommand("GetText", "GetText", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "in", 0, 0)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	})

	connection.NewCommand("AddText", "AddText", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "inout", 0, 0)
		cmd.NewParameter("label", "string", "in", 100, "")
		cmd.NewParameter("text", "string", "in", -1, "")
		cmd.NewParameter("action", "int", "in", 0, 10)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
		cmd.NewParameter("action_status_id", "int", "inout", 0, 0)
		cmd.NewParameter("action_msg", "string", "inout", 200, "")
	})

	//====================================================================================================
	// 
	//====================================================================================================
	dbase.Connections["DBApp"].NewCommand("GetOwner", "GetOwner", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	})
}
