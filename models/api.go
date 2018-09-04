package models

import (
	"fmt"
	"errors"
	"sync"
	//"github.com/astaxie/beego"
	//"InterfaceAgent/controllers"
)
///////////////////////////////////////////////////////////////////////

 var apiPoolOnce sync.Once


 type ApiPool struct {
 	ApiList map[string]*HttpApi
 }

var apiPoolInstance *ApiPool

func GetApiPoolInstance() *ApiPool {
	apiPoolOnce.Do(func() {
		apiPoolInstance = &ApiPool{}
		apiPoolInstance.ApiList = make(map[string]*HttpApi)
	})
	return apiPoolInstance
}



type HttpApiType int
const (
	ApiReadType HttpApiType = 0
	ApiDealType                = 1
	ApiDepositType            = 2
)


type HttpApi struct {
	Method string
	ApiType HttpApiType
	DoSign bool
	//Controller interface{}
}

func init(){
	fmt.Println("= models : Api ===-=======-=-=-=-=-=-=-=-=-")

}


func (this *ApiPool)AddApi(api *HttpApi)error{
	if api==nil{
		return errors.New("api is nil")
	}

	if api.Method == ""{
		return errors.New("empty metohd!")
	}

	this.ApiList[api.Method] = api
	
	return nil
}

func (this *ApiPool)GetApi(method string)*HttpApi{
	if _, ok := this.ApiList[method]; ok {  
		return this.ApiList[method]
	} else {
		return nil
	}

} 
	