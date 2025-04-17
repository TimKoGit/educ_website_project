package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetUsersTable(ctx *context.Context) (userTable table.Table) {

	userTable = table.NewDefaultTable(ctx, table.Config{
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

	info := userTable.GetInfo()

	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("Username", "username", db.Varchar)
	info.AddField("Role", "role", db.Varchar)
	info.AddField("Password", "password", db.Varchar)
	info.AddField("Surname", "surname", db.Varchar)
	info.AddField("Firstname", "firstname", db.Varchar)

	info.SetTable("users").SetTitle("Users").SetDescription("Users")

	formList := userTable.GetForm()

	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit()
	formList.AddField("Username", "username", db.Varchar, form.Text)
	formList.AddField("Role", "role", db.Varchar, form.Text)
	formList.AddField("Password", "password", db.Varchar, form.Text)
	formList.AddField("Surname", "surname", db.Varchar, form.Text)
	formList.AddField("Firstname", "firstname", db.Varchar, form.Text)

	formList.SetTable("users").SetTitle("Users").SetDescription("Users")

	return
}
