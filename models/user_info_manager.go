package models

import (
	//"github.com/astaxie/beego"
	"sync"
	"errors"
	"time"
	"fmt"
)

type UserInfoManager struct{
	info_mutex sync.RWMutex
	Timer        *time.Ticker
	userMap map[string]*UserInfo
}

type UserInfo struct{
	userId int
	appKey string
	secretKey string
	isReadInfo int
	isTrade int
	iswithdraw int
	address UserAddress
	infoTimer uint32
}

func (this *UserInfo)SetUserId(id int){
	this.userId = id
}

func (this *UserInfo)GetUserId()int{
	return this.userId
}

func (this *UserInfo)SetIp(ip string){
	this.address.Ip = ip
}

func (this *UserInfo)SetAppKey(appkey string){
	this.appKey = appkey
}

func (this *UserInfo)SetSecretKey(secretKey string){
	this.secretKey = secretKey
}

func (this *UserInfo)SetIsReadInfo(isReadInfo int){
	this.isReadInfo = isReadInfo
}

func (this *UserInfo)SetIsTrade(isTrade int){
	this.isTrade = isTrade
}

func (this *UserInfo)SetIsWithdraw(iswithdraw int){
	this.iswithdraw = iswithdraw
}

func (this *UserInfo)GetAppKey()string {
	return this.appKey
}

func (this *UserInfo)GetSecretKey()string {
	return this.secretKey
}

func (this *UserInfo)GetIsReadInfo()int {
	return this.isReadInfo
}

func (this *UserInfo)GetIsTrade()int {
	return this.isTrade
}

func (this *UserInfo)GetIswithdraw()int {
	return this.iswithdraw
}

func (this *UserInfo)GetIp()string {
	return this.address.Ip
}


type  UserAddress struct{
	Ip string
	Port string
}

//func init() {
 //     beego.Debug("user info init......")
//}


var userInfoOnce sync.Once
var userInfoManagerInstance *UserInfoManager

func GetUserInfoManager() *UserInfoManager {
	userInfoOnce.Do(func() {
		userInfoManagerInstance = new(UserInfoManager)
		userInfoManagerInstance.userMap = make(map[string]*UserInfo)
		userInfoManagerInstance.Timer = time.NewTicker(10 * time.Second)
		go userInfoManagerInstance.TimerHandle(userInfoManagerInstance.Timer)
	})
	return userInfoManagerInstance
}



func (this *UserInfoManager)PutUserInfo(userInfo *UserInfo)error{

	if userInfo == nil{
		return errors.New("userInfo is nil")
	}
	
	if _, ok := this.userMap[userInfo.appKey]; ok {
		return errors.New("This appKey already exists!")
	}

	userInfo.infoTimer = 0

	this.userMap[userInfo.appKey] = userInfo

	return nil
}

func (this *UserInfoManager)GetUserInfo(appkey string)(*UserInfo,error){
	if _, ok := this.userMap[appkey]; ok {  
		return this.userMap[appkey],nil
	} else {
		return nil,errors.New("appKey is not select!")
	}
}


const MAX_CACHE_TIMER_TICK = 60


func (this *UserInfoManager) TimerHandle(ticker *time.Ticker) {
	for {
		time := <-ticker.C
		/////////// timer handle ////////////
		for appkey, userInfo := range this.userMap {
			userInfo.infoTimer += 1
			//fmt.Printf("(%d).infotimer = %d\n", appkey, userInfo.infoTimer)
			if userInfo.infoTimer > MAX_CACHE_TIMER_TICK {
				fmt.Printf("----- time overappkey=%v, v=%v\n", appkey, userInfo)
				fmt.Printf("----- (%d).infotimer = %d\n", appkey, userInfo.infoTimer)
				this.info_mutex.Lock()
				delete(this.userMap,appkey)
				this.info_mutex.Unlock()
			}

		}
		fmt.Printf("%d  ", time.Second())
	}
}

