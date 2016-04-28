/*
 * ModulesAB server
 * TonyChyi <tonychee1989@gmail.com>
 */
package main

import (
	"fmt"

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
	beego.Info("All is ready, go running...")
	beego.Run()
}

func initDb() {
	o := orm.NewOrm()

	role := []models.Roles{
		models.Roles{
			Id:       uuid.New(),
			Name:     "Administrator",
			RoleFlag: models.RoleFlagAdmin,
		},
		models.Roles{
			Id:       uuid.New(),
			Name:     "Operator",
			RoleFlag: models.RoleFlagOperator,
		},
		models.Roles{
			Id:       uuid.New(),
			Name:     "User",
			RoleFlag: models.RoleFlagUser,
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
