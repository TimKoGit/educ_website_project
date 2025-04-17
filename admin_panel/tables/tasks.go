package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetTasksTable(ctx *context.Context) (tasksTable table.Table) {

	tasksTable = table.NewDefaultTable(ctx, table.Config{
		Driver:     db.DriverPostgresql,
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		Connection: table.DefaultConnectionName,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: table.DefaultPrimaryKeyName,
		},
	})

	info := tasksTable.GetInfo()

	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("Name", "name", db.Varchar)
	info.AddField("URL", "url", db.Varchar)
	info.AddField("ContestID", "contestid", db.Int)

	info.SetTable("tasks").SetTitle("Tasks").SetDescription("Tasks")

	formList := tasksTable.GetForm()

	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit()
	formList.AddField("Name", "name", db.Varchar, form.Text)
	formList.AddField("URL", "url", db.Varchar, form.Text)
	formList.AddField("ContestID", "contestid", db.Int, form.Number)

	formList.SetTable("tasks").SetTitle("Tasks").SetDescription("Tasks")

	return
}
