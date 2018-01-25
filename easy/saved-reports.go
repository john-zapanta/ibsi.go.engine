package easy

import (
	"net/http"
	"ibsi/dbase"
	"ibsi/crud"
	"ibsi/utils"
	// "ibsi/session"
)

func init() {

	crud.Handler(crud.CrudHandler {
		Name: "sys-saved-reports",
		Action: "sys-saved-reports",
		Path: "engine",
		KeyName: "id",
		ListDataSource: "DBReporting.GetSavedReports",
		UpdateDataSource: "DBReporting.AddSavedReportQuery",
		OnInitCrudParams: func(mode string, params dbase.TParameters, w http.ResponseWriter, r *http.Request) {
			// if mode == "edit" || mode == "new" {
				// params["mode"] = 10
			// }
		},
		OnNewRecord: func(mode string, row map[string]interface{}, w http.ResponseWriter, r *http.Request) {
			// row["user_id"] = session.GetUserId(r)
			// row["user_id"] = 1
			row["report_type_id"] = utils.StrToInt(r.Form.Get("report_type_id"))
		},
	})
	
	dbase.Connections["DBReporting"].NewCommand("GetSavedReports", "SysGetSavedReports", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "in", 0, 0)
		cmd.NewParameter("report_type_id", "int", "in", 0, 0)
		cmd.NewParameter("name", "string", "in", 100, "")
		cmd.NewParameter("action", "int", "in", 0, 0)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
		cmd.NewParameter("sort", "string", "in", 200, "")
		cmd.NewParameter("order", "string", "in", 10, "")
	})

	dbase.Connections["DBReporting"].NewCommand("AddSavedReportQuery", "SysAddSavedReportQuery", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "inout", 0, 0)
		cmd.NewParameter("name", "string", "in", 500, "")
		cmd.NewParameter("report_type_id", "int", "in", 0, 0)
		cmd.NewParameter("query", "string", "in", -1, "")
		cmd.NewParameter("action", "int", "in", 0, 11)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
		cmd.NewParameter("action_status_id", "int", "inout", 0, 0)
		cmd.NewParameter("action_msg", "string", "inout", 200, "")
	})
}
