/*ModuleAB server
 * Copyright (C) 2016 TonyChyi <tonychee1989@gmail.com>
 * License: GPL v3 or later.
 */
package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"

	"os"

	_ "github.com/ModuleAB/ModuleAB/server/docs"
	"github.com/ModuleAB/ModuleAB/server/policies"
	_ "github.com/ModuleAB/ModuleAB/server/routers"
	"github.com/ModuleAB/ModuleAB/server/version"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// DBS is template for make database connection string.
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
	defer func() {
		x := recover()
		if x != nil {
			beego.Emergency("Got fatal error:", x)
			var stack = make([]byte, 2<<10)
			runtime.Stack(stack, true)
			beego.Emergency("Stack trace:\n", string(stack))
			os.Exit(1)
		}
	}()

	// Don't run me with root, or will be unexpected security problem.
	if os.Getuid() == 0 {
		beego.Alert("Hey! You're running this server with user root!")
		panic("Don't run me with root!")
	}

	beego.Info("ModuleAB server", version.Version, "starting...")
	logfile := beego.AppConfig.DefaultString(
		"logFile",
		"logs/github.com/ModuleAB/ModuleAB/server.log",
	)
	err := beego.SetLogger("file", fmt.Sprintf(`{"filename":"%s"}`, logfile))
	if err != nil {
		panic(err)
	}
	beego.SetLevel(beego.LevelInformational)

	beego.Info("Hello!")

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
		// Run as Daemon
		if os.Getppid() != 1 {
			exePath, _ := filepath.Abs(os.Args[0])
			cmd := exec.Command(exePath, os.Args[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start()
			beego.Info("ModuleAB will run as daemon.")
			os.Exit(0)
		}

		beego.BeeLogger.DelLogger("console")
		err = orm.RunSyncdb("default", false, false)
	}

	if err != nil {
		beego.Alert("Database error:", err, ". go exit.")
		os.Exit(1)
	}
	beego.Debug("Current PID:", os.Getpid())
	ioutil.WriteFile(
		beego.AppConfig.DefaultString("pidFile", "github.com/ModuleAB/ModuleAB/server.pid"),
		[]byte(fmt.Sprint(os.Getpid())),
		0600,
	)
	beego.Info("Run check oas job...")
	go policies.CheckOasJob()
	beego.Info("All is ready, go running...")
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "Session_MobuleAB"
	beego.BConfig.Listen.ServerTimeOut = beego.AppConfig.DefaultInt64(
		"timeout", 0,
	)
	beego.Run()
}
