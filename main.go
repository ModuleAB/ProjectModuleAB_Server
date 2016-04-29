/*
 * ModulesAB server
 * TonyChyi <tonychee1989@gmail.com>
 */
package main

import (
	"fmt"
	"io/ioutil"
	"os/signal"
	"time"

	_ "moduleab_server/docs"
	"moduleab_server/models"
	_ "moduleab_server/routers"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pborman/uuid"
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
	beego.SetLogger("file", `{"filename":"logs/server.log"}`)
	beego.SetLevel(beego.LevelInformational)

	beego.Info("Hello!")

	// 别用root运行我！
	if os.Getuid() == 0 {
		beego.Alert("Hey! You're running this server with user root!")
		panic("Don't run me with root!")
	}
	err := fmt.Errorf("")
	switch beego.BConfig.RunMode {
	case "initdb":
		beego.Info("Got runmode: Initialize database")
		err = orm.RunSyncdb("default", true, true)
		orm.Debug = true
		initDb()
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
	beego.Info("Run signal notifier...")
	go signalNotifier()
	beego.Info("All is ready, go running...")
	beego.Run()
}

func signalNotifier() {
	beego.Info("Signal notifier started.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	for {
		select {
		case s := <-c:
			beego.Info(
				fmt.Sprintf(
					"Received Signal: %s, stop in 10 seconds...", s,
				),
			)
			time.Sleep(10 * time.Second)
			db, err := orm.GetDB("default")
			if err != nil {
				beego.Warn("Got error:", err)
				os.Exit(1)
			} else {
				db.Close()
				os.Exit(0)
			}
		default:
			continue
		}
	}
}

func initDb() {
	o := orm.NewOrm()

	role := []models.Roles{
		models.Roles{
			Id:        uuid.New(),
			Name:      "Administrator",
			RoleFlag:  models.RoleFlagAdmin,
			Removable: false,
		},
		models.Roles{
			Id:        uuid.New(),
			Name:      "Operator",
			RoleFlag:  models.RoleFlagOperator,
			Removable: false,
		},
		models.Roles{
			Id:        uuid.New(),
			Name:      "User",
			RoleFlag:  models.RoleFlagUser,
			Removable: false,
		},
	}
	o.Begin()
	_, err := o.InsertMulti(1, role)
	if err != nil {
		o.Rollback()
		beego.Alert("Error on inserting roles:", err)
		os.Exit(1)
	}
	o.Commit()

	user := &models.Users{
		Id:       uuid.New(),
		Name:     "admin",
		ShowName: "Administrator",
		Password: "admin",
		Roles: []*models.Roles{
			&role[0],
		},
		Removable: false,
	}
	_, err = models.AddUser(user)
	if err != nil {
		beego.Alert("Error on inserting user:", err)
		os.Exit(1)
	}

	appSet := &models.AppSets{
		Name: "Default",
		Desc: "Default app set",
	}
	_, err = models.AddAppSet(appSet)
	if err != nil {
		beego.Alert("Error on inserting default application set:", err)
		os.Exit(1)
	}

	backupSet := &models.BackupSets{
		Id:   uuid.New(),
		Name: "Default",
		Desc: "Default backup set",
	}
	_, err = models.AddBackupSet(backupSet)
	if err != nil {
		o.Rollback()
		beego.Alert("Error on inserting default backup set:", err)
		os.Exit(1)
	}
}
