package main

import (
	_ "InterfaceAgent/routers"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	 "InterfaceAgent/models"
	 "fmt"
	 _ "github.com/go-sql-driver/mysql"
)



func main() {

    initDB()

    secret_key,err:=models.GetUserApiSecretKey(beego.AppConfig.String("db::database"),"de2211a189619e1ca43adb700055cb01")
    if err==nil{
        fmt.Println(secret_key)
    }

    secret_key,err=models.GetUserApiSecretKey(beego.AppConfig.String("db::database"),"aac56d9e19644c46d467128b7285960c")
    if err==nil{
        fmt.Println(secret_key)
    }
    
      infoManager := models.GetUserInfoManager()
      ui1 := new(models.UserInfo)
      //ui := new(models.UserInfo)

    ui1.SetAppKey("123")
      //infoManager.PutUserInfo(*ui1)
    ui2,err:=infoManager.GetUserInfo("123")
      if err==nil{
        fmt.Println("=== ui2:",ui2)
      }
        

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}


func initDB() {
    userName := beego.AppConfig.String("db::username")
    password := beego.AppConfig.String("db::password")
    port := beego.AppConfig.String("db::port")
    host := beego.AppConfig.String("db::host")
    database := beego.AppConfig.String("db::database")

    dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", userName, password, host, port, database)
    beego.Debug(dataSource)

    orm.RegisterDriver("mysql", orm.DRMySQL)
    orm.RegisterDataBase("default", "mysql", dataSource)

}