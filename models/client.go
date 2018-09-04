package models

import (
	"time"
	"sync"
	//"errors"
	"container/list"
)

const MaxAccessTimesPerMin int = 3
//const ClearAccessListTimeSec int = 120


 var cliMngOnce sync.Once

type ClientManager struct {
	Cmutex sync.RWMutex
	Timer        *time.Ticker
 	AccessList map[string]*list.List
 }

 func getCurTimeMs()int64{
 	return time.Now().UnixNano()/1e6
 }

var cliMngInstance *ClientManager

func GetClientManagerInstance() *ClientManager {
	cliMngOnce.Do(func() {
		cliMngInstance = &ClientManager{}
		cliMngInstance.AccessList = make(map[string]*list.List)
		cliMngInstance.Timer = time.NewTicker(120 * time.Second)
		go cliMngInstance.TimerHandle(cliMngInstance.Timer)
	})
	return cliMngInstance
}



func (this *ClientManager)probeClientInfo(cliIp string)bool{
	if _, ok := this.AccessList[cliIp]; ok {  
		return true
	} else {
		return false
	}
}

func (this *ClientManager)IsAllowAccess(ip string)bool{
	curTime := getCurTimeMs()
	if this.probeClientInfo(ip){
		for this.AccessList[ip].Len() > 0 {
			oldestAccessTime := this.AccessList[ip].Front()
	
			if curTime - oldestAccessTime.Value.(int64) > 1000{                   //More than one second,delete
				this.Cmutex.Lock()
				this.AccessList[ip].Remove(oldestAccessTime)
				this.Cmutex.Unlock()
			} else {
				this.AccessList[ip].PushBack(curTime)
				if this.AccessList[ip].Len() > MaxAccessTimesPerMin {
					return false
				} else {
					return true
				}
			}
		}
		
		this.AccessList[ip].PushBack(curTime)
		return true
	} else {
		this.AccessList[ip]=list.New()
		this.AccessList[ip].PushBack(curTime)
		return true
	}
}


func (this *ClientManager)TimerHandle(ticker *time.Ticker) {
	curTime := getCurTimeMs()
	for ip,time_list := range this.AccessList {
		newestAccessTime := time_list.Back()
		if curTime - newestAccessTime.Value.(int64) > 1000 {
			delete(this.AccessList,ip)
		}
	}
}