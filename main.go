/*
 * ModulesAB server
 * TonyChyi <tonychee1989@gmail.com>
 */
package main

import (
	"fmt"
	"io/ioutil"
	"runtime"

	"moduleab_server/common"
	_ "moduleab_server/docs"
	"moduleab_server/policies"
	_ "moduleab_server/routers"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const DBS = "%s:%s@tcp(%s)/%s?charset=utf8"

func init() {
	err := orm.RegisterDataBase(
		"default",
		"mysql",
		fmt.Sprintf(DBS,
			beego.AppConfig.String("database::mysqluser"),
			beego.AppConfig.String("database::mysqlpass"),
			beego.AppConfig.String("database::mysqlurl"),
			beego.AppConfig.String("database::mysqldb"),
		),
	)
	if err != nil {
		beego.Alert(err)
	}
}

func main() {
	beego.Info("ModuleAB server", common.Version, "starting...")
	logfile := beego.AppConfig.String("logFile")
	if logfile == "" {
		logfile = "logs/moduleab_server.log"
	}
	err := beego.SetLogger("file", fmt.Sprintf(`{"filename":"%s"}`, logfile))
	if err != nil {
		panic(err)
	}
	beego.SetLevel(beego.LevelInformational)

	beego.Info("Hello!")

	// 别用root运行我！
	if os.Getuid() == 0 {
		beego.Alert("Hey! You're running this server with user root!")
		panic("Don't run me with root!")
	}

	switch beego.BConfig.RunMode {
	case "initdb":
		beego.Info("Got runmode: Initialize database")
		err = orm.RunSyncdb("default", true, true)
		orm.Debug = true
		policies.InitDb()
		beego.Info("Database is ready")
		os.Exit(0)

	case "dev":
		beego.Info("Got runmode: Development")
		beego.SetLevel(beego.LevelDebug)
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

		beego.Debug("Now import database")
		orm.Debug = true
		err = orm.RunSyncdb("default", false, true)

	case "deb":
		beego.Info("Got runmode: Debug")
		beego.SetLevel(beego.LevelDebug)

		beego.Debug("Now import database")
		orm.Debug = true
		err = orm.RunSyncdb("default", false, true)

	default:
		beego.BeeLogger.DelLogger("console")
		err = orm.RunSyncdb("default", false, false)
	}
	if err != nil {
		beego.Alert("Database error:", err, ". go exit.")
		os.Exit(1)
	}
	beego.Debug("Current PID:", os.Getpid())
	ioutil.WriteFile(
		beego.AppConfig.String("pidFile"),
		[]byte(fmt.Sprint(os.Getpid())),
		0600,
	)
	beego.Info("Run check oas job...")
	go policies.CheckOasJob()
	beego.Info("All is ready, go running...")
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "Session_MobuleAB"
	beego.Run()
	defer func() {
		x := recover()
		if x != nil {
			beego.Error("Got fatal error:", x)
			stack := make([]byte, 0)
			runtime.Stack(stack, true)
			beego.Error("Stack trace:\n", string(stack))
			os.Exit(1)
		}
	}()
}
