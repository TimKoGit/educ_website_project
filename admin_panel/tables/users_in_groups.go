package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetUsersingroupsTable(ctx *context.Context) (usersInGroupsTable table.Table) {

	usersInGroupsTable = table.NewDefaultTable(ctx, table.Config{
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

	info := usersInGroupsTable.GetInfo()

	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("UserId", "userid", db.Int)
	info.AddField("GroupID", "groupid", db.Int)

	info.SetTable("users_in_groups").SetTitle("Usersingroups").SetDescription("Usersingroups")

	formList := usersInGroupsTable.GetForm()

	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit()
	formList.AddField("UserId", "userid", db.Int, form.Number)
	formList.AddField("GroupID", "groupid", db.Int, form.Number)

	formList.SetTable("users_in_groups").SetTitle("Usersingroups").SetDescription("Usersingroups")

	return
}
