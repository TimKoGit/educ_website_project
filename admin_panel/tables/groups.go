package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetGroupsTable(ctx *context.Context) (groupTable table.Table) {

	groupTable = table.NewDefaultTable(ctx, table.Config{
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

	info := groupTable.GetInfo()

	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("Name", "name", db.Varchar)
	info.AddField("Code", "code", db.Varchar)

	info.SetTable("groups").SetTitle("Groups").SetDescription("Groups")

	formList := groupTable.GetForm()

	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit()
	formList.AddField("Name", "name", db.Varchar, form.Text)
	formList.AddField("Code", "code", db.Varchar, form.Text)

	formList.SetTable("groups").SetTitle("Groups").SetDescription("Groups")

	return
}
