package easy

import (
	"net/http"
	"ibsi/dbase"
	"ibsi/crud"
)

func init() {

	crud.Handler(crud.CrudHandler {
		Name: "sys-reports",
		Action: "sys-reports",
		Path: "engine",
		KeyName: "id",
		ListDataSource: "DBReporting.GetReportTypes",
		UpdateDataSource: "DBReporting.AddReportType",
		OnInitCrudParams: func(mode string, params dbase.TParameters, w http.ResponseWriter, r *http.Request) {
			// if mode == "edit" || mode == "new" {
				// params["mode"] = 10
			// }
		},
		OnNewRecord: func(mode string, row map[string]interface{}, w http.ResponseWriter, r *http.Request) {
			row["status_code_id"] = 10
		},
	})
	
	dbase.Connections["DBReporting"].NewCommand("GetReportTypes", "SysGetReportTypes", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "in", 0, 0)
		cmd.NewParameter("filter", "string", "in", 100, "")
		cmd.NewParameter("action", "int", "in", 0, 0)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
		cmd.NewParameter("page", "int", "in", 0, 1)
		cmd.NewParameter("pagesize", "int", "in", 0, 50)
		cmd.NewParameter("row_count", "int", "inout", 0, 0)
		cmd.NewParameter("page_count", "int", "inout", 0, 0)
		cmd.NewParameter("sort", "string", "in", 200, "")
		cmd.NewParameter("order", "string", "in", 10, "")
	})

	dbase.Connections["DBReporting"].NewCommand("AddReportType", "SysAddReportType", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "in", 0, 0)
		cmd.NewParameter("report_type", "string", "in", 100, "")
		cmd.NewParameter("status_code_id", "int", "in", 0, 0)
		cmd.NewParameter("action", "int", "in", 0, 10)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
		cmd.NewParameter("action_status_id", "int", "inout", 0, 0)
		cmd.NewParameter("action_msg", "string", "inout", 200, "")
	})
}
