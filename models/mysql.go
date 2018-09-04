package models

import (
    "github.com/astaxie/beego/orm"
    "time"
    "github.com/astaxie/beego"
     "fmt"
     "strconv"
     "errors"
)


type UserApiInfoTable struct {
	Id         int       `xorm:"fid"`
	Fkey        string    `xorm:"fkey"`
	Fsecret     string    `xorm:"fsecret"`
	Fuser       int       `xorm:"fuser"`
	Label       string    `xorm:"label"`
	Fcreatetime time.Time `xorm:"fcreatetime"`
	Fistrade    int       `xorm:"fistrade"`
	Fiswithdraw int       `xorm:"fiswithdraw"`
	Fisreadinfo int       `xorm:"fisreadinfo"`
	Fip         string    `xorm:"fip"`
	WhiteIps    string    `xorm:"whiteIps"`
}

type UserApiInfo struct {
	AppKey       string   
	SecretKey  string 
	UserId       int
	IsTrade      int
	IsWithdraw  int 
	IsReadInfo  int
}

func (u *UserApiInfoTable) TableName() string {
    // db table name
    return "fapi"
}


func (this *UserApiInfoTable)ToString()string{
	st := fmt.Sprintf("id:%d,Fkey:%s,Fsecret:%d,Fuser:%d,Label:%s,Fcreatetime:%d,Fistrade:%d,Fiswithdraw:%d,Fisreadinfo:%d,Fip:%s,WhiteIps:%s",this.Id,this.Fkey,this.Fsecret,this.Fuser,this.Label,this.Fcreatetime,this.Fistrade,this.Fiswithdraw,this.Fisreadinfo,this.Fip,this.WhiteIps)
	return st
}

func init() {
      orm.RegisterModel(new(UserApiInfoTable))
}



func GetUserApiSecretKey(db string,apikey string) (string,error){
	beego.Debug(db)
	o:=orm.NewOrm()
	o.Using(db)
	var num int64
	var err error
	var maps []orm.Params

	num, err = o.Raw("SELECT fsecret FROM fapi WHERE fkey = ?",apikey).Values(&maps)
	if err == nil && num > 0 {
	    return maps[0]["fsecret"].(string),nil
	} else {
		return "",err
	}
}


func GetUserApiInfo(db string,apikey string) (*UserApiInfo,error){
	beego.Debug(db)
	o:=orm.NewOrm()
	o.Using(db)
	var num int64
	var err error
	var maps []orm.Params

	num, err = o.Raw("SELECT * FROM fapi WHERE fkey = ?",apikey).Values(&maps)
	if err == nil {
		if  num <= 0 {
			return nil,errors.New("this apikey not found!")
		}

		userApiInfo := new(UserApiInfo)
		userApiInfo.AppKey = apikey
		 
		
		UserId,err := strconv.ParseInt(maps[0]["fuser"].(string), 10, 32)
		if err != nil {
			// TODO: log
			
			return nil,err
		}

		userApiInfo.UserId = int(UserId)

		userApiInfo.SecretKey = maps[0]["fsecret"].(string)

		
		IsWithdraw,err := strconv.ParseInt(maps[0]["fiswithdraw"].(string), 10, 32)
		if err != nil {
			// TODO: log
			
			return nil,err
		}

		userApiInfo.IsWithdraw = int(IsWithdraw)
		
		
		IsTrade,err := strconv.ParseInt(maps[0]["fistrade"].(string), 10, 32)
		if err != nil {
			// TODO: log
			
			return nil,err
		}

		userApiInfo.IsTrade = int(IsTrade)
		
		
		IsReadInfo,err := strconv.ParseInt(maps[0]["fisreadinfo"].(string), 10, 32)
		if err != nil {
			// TODO: log
			
			return nil,err
		}

		userApiInfo.IsReadInfo = int(IsReadInfo)
	    //return maps[0]["fsecret"].(string),maps[0]["fisreadinfo"].(int),maps[0]["fistrade"].(int),maps[0]["fiswithdraw"].(int),nil
	   
	     return userApiInfo,nil
	} else { 
		
		return nil,err
	}
}