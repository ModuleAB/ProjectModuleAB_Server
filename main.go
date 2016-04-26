/*
 * ModulesAB server
 * TonyChyi <tonychee1989@gmail.com>
 */
package main

import (
	"fmt"
	_ "moduleab_server/docs"
	_ "moduleab_server/routers"
	"os"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const DBS = "%s:%s@tcp(%s)/%s?charset=utf8"

func init() {
	orm.RegisterDataBase(
		"default",
		"mysql",
		fmt.Sprintf(DBS,
			beego.AppConfig.String("database::mysqluser"),
			beego.AppConfig.String("database::mysqlpass"),
			beego.AppConfig.String("database::mysqlurl"),
			beego.AppConfig.String("database::mysqldb"),
		),
	)
}

func main() {
	beego.SetLogger("file", `{"filename":"logs/server.log"}`)
	beego.SetLevel(beego.LevelInformational)

	beego.Info("Hello!")

	err := fmt.Errorf("")
	switch beego.BConfig.RunMode {
	case "newdb":
		err = orm.RunSyncdb("default", true, true)
		beego.Info("Database is clear")
		os.Exit(0)

	case "dev":
		beego.Info("Got runmode: Development")
		beego.SetLevel(beego.LevelDebug)
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

		beego.Debug("Now import database")
		err = orm.RunSyncdb("default", false, true)

	case "deb":
		beego.Info("Got runmode: Debug")
		beego.SetLevel(beego.LevelDebug)

		beego.Debug("Now import database")
		err = orm.RunSyncdb("default", false, true)

	default:
		err = orm.RunSyncdb("default", false, false)
	}
	if err != nil {
		beego.Alert("Database error:", err, ". go exit.")
		os.Exit(1)
	}
	beego.Info("All is ready, go running...")
	beego.Run()
}
