package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/postgres"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
//	"github.com/GoAdminGroup/go-admin/examples/datamodel"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"gitlab.atp-fivt.org/fullstack2024a/kondrashovti-project/tables"
	"time"
)

func main() {
	r := gin.Default()

	eng := engine.Default()

	// global config
	cfg := config.Config{
		Databases: config.DatabaseList{
			"default": {
				Host:         "127.0.0.1",
				Port:         "5432",
				User:         "postgres",
				Pwd:          "postgres",
				Name:         "programming_educ",
				MaxIdleConns: 50,
				MaxOpenConns: 150,
				ConnMaxLifetime: time.Hour,
				Driver:       config.DriverPostgresql,
			},
        	},
		UrlPrefix: "admin",
		// STORE is important. And the directory should has permission to write.
		Store: config.Store{
		    Path:   "./uploads", 
		    Prefix: "uploads",
		},
		Language: language.EN,
		// debug mode
		Debug: true,
		// log file absolute path
		InfoLogPath: "./var/logs/info.log",
		AccessLogPath: "./var/logs/access.log",
		ErrorLogPath: "./var/logs/error.log",
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}

	// add component chartjs
	template.AddComp(chartjs.NewChart())

	_ = eng.AddConfig(&cfg).
		AddGenerators(tables.Generators).
	        // add generator, first parameter is the url prefix of table when visit.
    	        // example:
    	        //
    	        // "user" => http://localhost:9033/admin/info/user
    	        //		
	//	AddGenerator("user", datamodel.GetUserTable).
		Use(r)
	
	// customize your pages
	//eng.HTML("GET", "/admin", datamodel.GetContent)

	_ = r.Run(":9033")
}