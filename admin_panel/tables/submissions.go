package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetSubmissionsTable(ctx *context.Context) (submissionsTable table.Table) {

	submissionsTable = table.NewDefaultTable(ctx, table.Config{
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

	info := submissionsTable.GetInfo()

	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("TaskID", "taskid", db.Int)
	info.AddField("UserID", "userid", db.Int)
	info.AddField("Status", "status", db.Varchar)
	info.AddField("URL", "url", db.Varchar)
	info.AddField("CreatedAt", "created_at", db.Timestamp)

	info.SetTable("submissions").SetTitle("Submissions").SetDescription("Submissions")

	formList := submissionsTable.GetForm()

	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit()
	formList.AddField("TaskID", "taskid", db.Int, form.Number)
	formList.AddField("UserID", "userid", db.Int, form.Number)
	formList.AddField("Status", "status", db.Varchar, form.Text)
	formList.AddField("URL", "url", db.Varchar, form.Text)
	formList.AddField("CreatedAt", "created_at", db.Timestamp, form.Datetime)

	formList.SetTable("submissions").SetTitle("Submissions").SetDescription("Submissions")

	return
}
