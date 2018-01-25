package system

import (
	"log"
	"ibsi/dbase"
)

func init() {

	connection, ok := dbase.Connections["DBReporting"]
	if !ok {
		log.Println("DBReporting connection was not found")
		return
	}

	connection.NewCommand("AddSavedReportQueryItem", "SysAddSavedReportQueryItem", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "in", 0, 0)
		cmd.NewParameter("name", "string", "in", 100, "")
		cmd.NewParameter("value", "string", "in", 2048, "")
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	})

	connection.NewCommand("GetReport", "SysGetReport", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "in", 0, 0)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
	})

	connection.NewCommand("RunReport", "SysRunReport", "procedure", func(cmd dbase.ICommand) {
		cmd.NewParameter("id", "int", "in", 0, 0)
		cmd.NewParameter("visit_id", "int", "in", 0, 0)
		cmd.NewParameter("page", "int", "in", 0, 1)
		cmd.NewParameter("pagesize", "int", "in", 0, 50)
		cmd.NewParameter("row_count", "int", "inout", 0, 0)
		cmd.NewParameter("page_count", "int", "inout", 0, 0)
		cmd.NewParameter("sort", "string", "in", 200, "")
		cmd.NewParameter("order", "string", "in", 10, "")
	})

	connection.NewQuery("FindSavedReport", func(cmd dbase.ICommand) {
			cmd.NewParameter("id", "int", "inout", 0, 0)
			cmd.NewParameter("report_type_id", "int", "in", 0, 0)
			cmd.NewParameter("user_id", "int", "in", 0, 0)
			cmd.NewParameter("name", "string", "in", 100, "")
		}, 
		func() string {
			return `
				SELECT
					@id = id
				FROM saved_reports
				WHERE report_type_id = @report_type_id AND user_id = @user_id and name = @name
			`
		},
	) 

	// connection.NewQuery("FindSavedReport", "

		// ", "text")
	// })
}
