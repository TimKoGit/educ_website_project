package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetContestsTable(ctx *context.Context) (contestsTable table.Table) {

	contestsTable = table.NewDefaultTable(ctx, table.Config{
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

	info := contestsTable.GetInfo()

	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("Name", "name", db.Varchar)
	info.AddField("Time", "time", db.Timestamp)
	info.AddField("Duration", "duration", db.Int)
	info.AddField("GroupID", "groupid", db.Int)

	info.SetTable("contests").SetTitle("Contests").SetDescription("Contests")

	formList := contestsTable.GetForm()

	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit()
	formList.AddField("Name", "name", db.Varchar, form.Text)
	formList.AddField("Time", "time", db.Timestamp, form.Datetime)
	formList.AddField("Duration", "duration", db.Int, form.Number)
	formList.AddField("GroupID", "groupid", db.Int, form.Number)

	formList.SetTable("contests").SetTitle("Contests").SetDescription("Contests")

	return
}
