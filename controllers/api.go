package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"InterfaceAgent/models"
	"strings"
	"encoding/json"
	"time"
	"strconv"
	"regexp"
)

var clientManager *models.ClientManager = models.GetClientManagerInstance()

type HttpApiType int
const (
	ApiReadType HttpApiType = 0
	ApiDealType                = 1
	ApiDepositType            = 2
)

var IdCount = 0

type RequestBody struct {
	method string         `json:"method"`
	params []interface{}  `json:"params"`
	id int                  `json:"id"`
}



type HttpApi struct {
	method string
	apiType HttpApiType
}

 




func init(){
	fmt.Println("= controllers : Api ===-=======-=-=-=-=-=-=-=-=-")

	baseUrl := beego.AppConfig.String("url::baseUrl") + "/"
	if baseUrl == ""{
		fmt.Println("base url is empty string!\n")
	}
	
	ap := models.GetApiPoolInstance()


	httpApi := new(models.HttpApi)

	httpApi.Method = "ticker.do "
	httpApi.ApiType = models.ApiReadType
	httpApi.DoSign = false
	
	beego.Router(baseUrl + httpApi.Method, &ApiTickerController{Api:httpApi})

	ap.AddApi(httpApi)

	httpApi.Method = "depth.do"
	httpApi.ApiType = models.ApiReadType
	httpApi.DoSign = false

	beego.Router(baseUrl + httpApi.Method, &ApiDepthController{Api:httpApi})

	ap.AddApi(httpApi)


	httpApi.Method = "kline.do"
	httpApi.ApiType = models.ApiReadType
	httpApi.DoSign = false

	beego.Router(baseUrl + httpApi.Method, &ApiKlineController{Api:httpApi})

	ap.AddApi(httpApi)

	
	httpApi.Method = "balance.do"
	httpApi.ApiType = models.ApiReadType
	httpApi.DoSign = true

	beego.Router(baseUrl + httpApi.Method, &ApiUserBalanceController{Api:httpApi})

	ap.AddApi(httpApi)


	httpApi.Method = "trade.do"
	httpApi.ApiType = models.ApiDealType
	httpApi.DoSign = true

	beego.Router(baseUrl + httpApi.Method, &ApiTradeController{Api:httpApi})

	ap.AddApi(httpApi)


	
	httpApi.Method = "batch_trade.do"
	httpApi.ApiType = models.ApiDealType
	httpApi.DoSign = true

	beego.Router(baseUrl + httpApi.Method, &ApiBatchTradeController{Api:httpApi})

	ap.AddApi(httpApi)
}


func get_ip(str string) string {
   
	e := strings.Index(str[0:], ":")
	if e < 0 {
		return ""
	}
	return str[0 : e]
}

/////////////////////////////////////////////////    Ticker    /////////////////////////////////////////////////////////////////////


type ApiTickerController struct {
	beego.Controller
	Api *models.HttpApi
}

//=========================== recv matchengine message ===============================
type Result struct{
		Open string                `json:"open"`
		Last string                `json:"last"`
		High string                `json:"high"`
		Low  string                `json:"low"`
		Volume string             `json:"volume"`
		Deal string                `json:"deal"`
}

type MarketStatusToday struct {
	Error models.ErrorDetail    `json:"error"`
	Res Result             `json:"result"`
	Id   int                `json:"id"`
}

//========================== send API response message ===============================


type TickerRes struct {
	Date    int64                `json:"date"`
	Tick    Result                `json:"ticker"`
}


//=============================================================================


func (a ApiTickerController)Get() {

	
	cli_req := a.Ctx.Request
	addr := cli_req.RemoteAddr // "IP:port" "192.168.1.150:8889"
	ip := get_ip(addr)
	if !clientManager.IsAllowAccess(ip){
		str := models.BuildErrorMsg(models.ReqFrequenTooHigh)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}


	
	var market string
	a.Ctx.Input.Bind(&market,"market")
	if market == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}



	req := httplib.Post(beego.AppConfig.String("MeApiHttpServer"))

	if IdCount < 10000 {
		IdCount ++
	} else {
		IdCount = 0
	}

	bt := fmt.Sprintf("{\"method\": \"market.kline\", \"params\": [\"%s\"], \"id\": %d}",market,IdCount)

	req.Body(bt)

	str, erri := req.String()
	if erri != nil {
	    a.Data["json"] = models.BuildErrorMsg(models.RequestFail)
	} else{
		ch := make(chan string, 1)
		go func(c chan string, s string){
 			c <- s
    		}(ch, string(str))
    		st := <-ch
		var mst MarketStatusToday 
		json.Unmarshal([]byte(st), &mst)

		fmt.Println(mst.Id)
		fmt.Printf("\n%#v\n", mst)

		var ticker TickerRes 
		ticker.Date = time.Now().Unix()
		ticker.Tick = mst.Res


		resp, err := json.Marshal(ticker)
		if err != nil {
			fmt.Println("encoding faild")
		} else {
			fmt.Println("encoded data : ")
			fmt.Println(resp)
			fmt.Println(string(resp))
		}


		a.Data["json"] = string(resp)
	}
	
	a.ServeJSON(false)
}

/////////////////////////////////////////////////// api depth //////////////////////////////////////////////////////////////////
// {
//     "error": null,
//     "result": {
//         "asks": [
//             [
//                 "1000",
//                 "1"
//             ]
//         ],
//         "bids": []
//     },
//     "id": 123
// }



//=========================== recv matchengine message ===============================

//type Order struct {
 //   []string     
//}

type DepthResult struct{
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

type OrderDepth struct {
	Error models.ErrorDetail    `json:"error"`
	Res DepthResult               `json:"result"`
	Id   int                         `json:"id"`
}
//========================== send API response message ===============================
//=============================================================================


type ApiDepthController struct {
	beego.Controller
	Api *models.HttpApi
}

func (a ApiDepthController)Get() {

	cli_req := a.Ctx.Request
	addr := cli_req.RemoteAddr // "IP:port" "192.168.1.150:8889"
	ip := get_ip(addr)
	if !clientManager.IsAllowAccess(ip){
		str := models.BuildErrorMsg(models.ReqFrequenTooHigh)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	
	var market string
	a.Ctx.Input.Bind(&market,"market")
	if market == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}



	req := httplib.Post(beego.AppConfig.String("MeApiHttpServer"))

	if IdCount < 10000 {
		IdCount ++
	} else {
		IdCount = 0
	}

	bt := fmt.Sprintf("{\"method\": \"order.depth\", \"params\": [\"%s\",100,\"0\"], \"id\": %d}",market,IdCount)
	fmt.Println(bt)

	req.Body(bt)

	str, erri := req.String()
	if erri != nil {
	    a.Data["json"] = models.BuildErrorMsg(models.RequestFail)
	} else{
		fmt.Printf("\n===========\n %s\n",str)
		ch := make(chan string, 1)
		go func(c chan string, s string){
 			c <- s
    		}(ch, string(str))
    		st := <-ch

		var odp OrderDepth 
		json.Unmarshal([]byte(st), &odp)

		//fmt.Println(odp.Id)
		//fmt.Printf("\n%#v\n", odp)

		resp, err := json.Marshal(odp.Res)
		if err != nil {
			fmt.Println("encoding faild")
		} 

		a.Data["json"] = string(resp)
	}
	
	a.ServeJSON(false)
}


/////////////////////////////////////////////////  api kline  /////////////////////////////////////////////////////////////////////

//=========================== recv matchengine message ===============================
// {
// 	"method": "market.kline", 
// 	"params": ["BTCBCC",1512449633,1512469633,30], 
// 	"id": 123
// }

//   "result": [
//         [
//             1512449610,
//             "8000",
//             "8000",
//             "8000",
//             "8000",
//             "0",
//             "0",
//             "BTCBCC"
//         ],
//         [
//             1512449640,
//             "8000",
//             "8000",
//             "8000",
//             "8000",
//             "0",
//             "0",
//             "BTCBCC"
//         ]
type Kline struct {
	Error models.ErrorDetail    `json:"error"`
	Res [][]interface{}               `json:"result"`
	Id   int                         `json:"id"`
}

//========================== send API response message ===============================
// ----------- type --------------
// 1min : 1分钟
// 3min : 3分钟
// 5min : 5分钟
// 15min : 15分钟
// 30min : 30分钟
// 1day : 1日
// 3day : 3日
// 1week : 1周
// 1hour : 1小时
// 2hour : 2小时
// 4hour : 4小时
// 6hour : 6小时
// 12hour : 12小时



//=============================================================================
type ApiKlineController struct {
	beego.Controller
	Api *models.HttpApi
}

func (a ApiKlineController)Get() {
	Types := map[string]int64 {
		"1min":60,
		"3min":180,
		"5min":300,
		"15min":900,
		"30min":1800,
		"1day":86400,
		"3day":259200,
		"1week":604800,
		"1hour":3600,
		"2hour":7200,
		"4hour":14400,
		"6hour":21600,
		"12hour":43200,
	} 


	cli_req := a.Ctx.Request
	addr := cli_req.RemoteAddr // "IP:port" "192.168.1.150:8889"
	ip := get_ip(addr)
	if !clientManager.IsAllowAccess(ip){
		str := models.BuildErrorMsg(models.ReqFrequenTooHigh)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	
	var market string
	a.Ctx.Input.Bind(&market,"market")
	if market == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}


	var TimeType string
	a.Ctx.Input.Bind(&TimeType,"type")
	if TimeType == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}



	var since string
	a.Ctx.Input.Bind(&since,"since")
	if since == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}


	since_64, err := strconv.ParseInt(since, 10, 64) 
	if err != nil {
		str := models.BuildErrorMsg(models.ParamError)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	var endtime  int64
	
	if times, ok := Types[TimeType]; ok {  
		endtime = since_64 + times
	} else {
		str := models.BuildErrorMsg(models.ParamError)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}



	req := httplib.Post(beego.AppConfig.String("MeApiHttpServer"))

	if IdCount < 10000 {
		IdCount ++
	} else {
		IdCount = 0
	}




	bt := fmt.Sprintf("{\"method\": \"market.kline\", \"params\": [\"%s\",%d,%d,30], \"id\": %d}",
		market,
		since_64,
		endtime,
		IdCount)
	fmt.Println(bt)

	req.Body(bt)

	str, erri := req.String()
	if erri != nil {
	    a.Data["json"] = models.BuildErrorMsg(models.RequestFail)
	} else{
		fmt.Printf("\n===========\n %s\n",str)
		ch := make(chan string, 1)
		go func(c chan string, s string){
 			c <- s
    		}(ch, string(str))
    		st := <-ch

		var kline Kline 
		json.Unmarshal([]byte(st), &kline)

		//fmt.Println(odp.Id)
		//fmt.Printf("\n%#v\n", odp)

		resp, err := json.Marshal(kline.Res)
		if err != nil {
			fmt.Println("encoding faild")
		} 

		a.Data["json"] = string(resp)
	}
	
	a.ServeJSON(false)
}


/////////////////////////////////////////////////   userinfo  /////////////////////////////////////////////////////////////////////

//=========================== recv matchengine message ===============================
//req:
// {
//         "method": "balance.query",
//         "params": [1],                                    # 获取user_id=1的用户所有资产名称资产。
//         "id": 123
// }
//res:
// Response:
// {
//     "error": null,
//     "result": {
//         "BTC": {
//             "available": "0",                               
//             "freeze": "0"                                           
//         }
//     },
//     "id": 1
// }
type Balance struct {
	Error models.ErrorDetail    `json:"error"`
	Res  map[string]Asset          `json:"result"`
	Id   int                         `json:"id"`
}
//========================== send API response message ===============================
// # Response
// {
//     "assets": [
//           "BTC": {
//                 "available": "0",                               
//                 "freeze": "0"                                           
//            }
//            "EOS":{
//                 "available": "0",                               
//                 "freeze": "0"                                           
//            }
//     ],
//     "error": null
// }
type Asset struct {
	Available   string               `json:"available"`
	Freeze       string               `json:"freeze"`
}

type UserAssets struct {
	Error models.ErrorDetail       `json:"error"`
	Assets   map[string]Asset                   `json:"assets"`
}
//=============================================================================

type ApiUserBalanceController struct {
	beego.Controller
	Api *models.HttpApi
}

func (a ApiUserBalanceController)Post() {

	params := make(map[string]string)
	cli_req := a.Ctx.Request
	addr := cli_req.RemoteAddr // "IP:port" "192.168.1.150:8889"
	ip := get_ip(addr)
	if !clientManager.IsAllowAccess(ip){
		str := models.BuildErrorMsg(models.ReqFrequenTooHigh)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	
	var apikey string
	a.Ctx.Input.Bind(&apikey,"apikey")
	if apikey == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
	params["apikey"] = apikey

	var sign string
	a.Ctx.Input.Bind(&sign,"sign")
	if sign == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
	params["sign"] = sign

	userInfoManager := models.GetUserInfoManager()

	var userId int


	userInfo,_ := userInfoManager.GetUserInfo(apikey)
	if userInfo == nil {
		userApiInfo,Err := models.GetUserApiInfo(beego.AppConfig.String("db::database"),apikey)
		if Err != nil {
			str := models.BuildErrorMsg(models.SystemError)
			a.Data["json"] = str
			a.ServeJSON()
			return
		}
		
		userId = 7//userApiInfo.UserId

		userInfo = new(models.UserInfo)
		userInfo.SetAppKey(userApiInfo.AppKey)
		userInfo.SetSecretKey(userApiInfo.SecretKey)
		userInfo.SetIsReadInfo(userApiInfo.IsReadInfo)
		userInfo.SetIsTrade(userApiInfo.IsTrade)
		userInfo.SetIsWithdraw(userApiInfo.IsWithdraw)
	}


	if !models.IsCorrectSign(userInfo,params) {
		str := models.BuildErrorMsg(models.SignNotMatch)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	req := httplib.Post(beego.AppConfig.String("MeApiHttpServer"))

	if IdCount < 10000 {
		IdCount ++
	} else {
		IdCount = 0
	}

	bt := fmt.Sprintf("{\"method\": \"balance.query\", \"params\": [%d], \"id\": %d}",userId,IdCount)
	fmt.Println(bt)

	req.Body(bt)

	str, erri := req.String()
	if erri != nil {
	    a.Data["json"] = models.BuildErrorMsg(models.RequestFail)
	} else{
		fmt.Printf("\n===========\n %s\n",str)
		ch := make(chan string, 1)
		go func(c chan string, s string){
 			c <- s
    		}(ch, string(str))
    		st := <-ch

		var bl Balance 
		json.Unmarshal([]byte(st), &bl)

		//fmt.Println(odp.Id)
		//fmt.Printf("\n%#v\n", odp)
		var ua UserAssets 
		ua.Error = bl.Error
		ua.Assets = bl.Res

		resp, err := json.Marshal(ua)
		if err != nil {
			fmt.Println("encoding faild")
		} 

		a.Data["json"] = string(resp)
	}
	
	a.ServeJSON(false)

}


/////////////////////////////////////////////////  api trade  /////////////////////////////////////////////////////////////////////

//=========================== recv matchengine message ===============================
// {
//         "error": null,
//         "result": {
//         "id": 1,
//         "market": "BTCUSD",
//         "source": "",
//         "type": 1,
//         "side": 1,
//         "user": 1,
//         "ctime": 1512375992.515996,
//         "mtime": 1512375992.515996,
//         "price": "8000",
//         "amount": "10",
//         "taker_fee": "0.002",
//         "maker_fee": "0.001",
//         "left": "10",
//         "deal_stock": "0",
//         "deal_money": "0",
//         "deal_fee": "0"
//     },
//     "id": 123
// }
type PutOrderRes struct {
	Error *models.ErrorDetail    `json:"error"`
	Res  *Order                      `json:"result"`
	Id   int                         `json:"id"`
}

type Order struct {
	Order_id      int               `json:"id"`
	Market         string
	Source         string
	Type_order   int                `json:"type"`
	Side            int
	User_id        int               `json:"user"`
	Ctime           float32
	Mtime           float32
	Price           string
	Amount          string
	Taker_fee      string
	Maker_fee      string
	Left             string
	Deal_stock     string
	Deal_money     string
	Deal_fee        string
}
//========================== send API response message ===============================
type TradeRes struct {
	Result  bool			`json:"result"`
	Order_id int 		`json:"order_id"`
}
//=============================================================================
type ApiTradeController struct {
	beego.Controller
	Api *models.HttpApi
}

//


func (a ApiTradeController)Post() {

	params := make(map[string]string)

	//=========== access filter =================
	cli_req := a.Ctx.Request
	addr := cli_req.RemoteAddr // "IP:port" "192.168.1.150:8889"
	ip := get_ip(addr)
	if !clientManager.IsAllowAccess(ip){
		str := models.BuildErrorMsg(models.ReqFrequenTooHigh)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	
	//============ param check ==================
	var market string
	a.Ctx.Input.Bind(&market,"market")
	if market == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
	

	params["market"] = market

	var type_para string
	a.Ctx.Input.Bind(&type_para,"type")
	if type_para == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
	

	params["type"] = type_para

	var price string
	a.Ctx.Input.Bind(&price,"price")
	if price == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
	
	params["price"] = price




	var amount string
	a.Ctx.Input.Bind(&amount,"amount")
	if amount == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
	params["amount"] = amount
	


	var apikey string




	a.Ctx.Input.Bind(&apikey,"apikey")
	if apikey == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
	params["apikey"] = apikey
	

	var sign string
	a.Ctx.Input.Bind(&sign,"sign")
	if sign == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	params["sign"] = sign
	//beego.Debug(" *** ")



//=============== sign auth ===================  

	userInfoManager := models.GetUserInfoManager()

	var userId int


	userInfo,_ := userInfoManager.GetUserInfo(apikey)
	if userInfo == nil {
		userApiInfo,Err := models.GetUserApiInfo(beego.AppConfig.String("db::database"),apikey)
		if Err != nil {
			str := models.BuildErrorMsg(models.SystemError)
			a.Data["json"] = str
			a.ServeJSON()
			return
		}
		//beego.Debug(" *** ")
		userId = userApiInfo.UserId

		userInfo = new(models.UserInfo)
		userInfo.SetAppKey(userApiInfo.AppKey)
		userInfo.SetSecretKey(userApiInfo.SecretKey)
		userInfo.SetIsReadInfo(userApiInfo.IsReadInfo)
		userInfo.SetIsTrade(userApiInfo.IsTrade)
		userInfo.SetIsWithdraw(userApiInfo.IsWithdraw)
	}


	if !models.IsCorrectSign(userInfo,params) {
		str := models.BuildErrorMsg(models.SignNotMatch)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
//beego.Debug(" *** ")

	const taker_fee string = "0.0002"
	const maker_fee string = "0.0001"
	//============== internal api request ==================
	var method string
	var side int
	if (params["type"] == "buy") || (params["type"] == "sell"){
		method = "order.put_limit"
	} else if (params["type"] == "market_buy") || (params["type"] == "market_sell"){
		method = "order.put_market"
	} else {
		str := models.BuildErrorMsg(models.IllegalParameter)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	if (params["type"] == "buy") || (params["type"] == "market_buy"){
		side = 2
	} else if (params["type"] == "sell") || (params["type"] == "market_sell"){
		side = 1
	}



	req := httplib.Post(beego.AppConfig.String("MeApiHttpServer"))

	if IdCount < 10000 {
		IdCount ++
	} else {
		IdCount = 0
	}


	bt := fmt.Sprintf("{\"method\": \"%s\", \"params\": [%d,\"%s\",%d,\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"], \"id\": %d}",method,userId,market,side,amount,price,taker_fee,maker_fee,"empty source",IdCount)
	fmt.Println(bt)

	req.Body(bt)

	//===================== app response ===================
	str, erri := req.String()
	if erri != nil {
	    a.Data["json"] = models.BuildErrorMsg(models.RequestFail)
	} else{
		ch := make(chan string, 1)
		go func(c chan string, s string){
 			c <- s
    		}(ch, string(str))
    		st := <-ch

		var por PutOrderRes 
		json.Unmarshal([]byte(st), &por)

		//fmt.Println(odp.Id)
		//fmt.Printf("\n%#v\n", por)

		var tr TradeRes

		if por.Error != nil{
			tr.Result = false
		} else if por.Res != nil{
			tr.Order_id = (*por.Res).Order_id
			tr.Result = true
		}
	

		resp, err := json.Marshal(tr)
		if err != nil {
			fmt.Println("encoding faild")
		} 

		a.Data["json"] = string(resp)
	}
	
	a.ServeJSON(false)

}



//////////////////////////////////////////////////  batch_trade  /////////////////////////////////////////////////////////////////////

//=========================== recv matchengine message ===============================

//========================== send API response message ===============================
// {
// 	"order_info":[
// 		{"order_id":41724206},
// 		{"error_code":10011,"order_id":-1},
// 		{"error_code":10014,"order_id":-1}
// 	],
// 	"result":true
// }

type OrderResult struct {
	Order_id    int        `json:"order_id"`
	Error_code int         `json:"error_code"`
}

type OrdersInfo struct {
	Order_info     []*OrderResult    `json:"order_info"`
	Result           bool             `json:"result"`
}
//=============================================================================






func GetOrderRequires(orders string)[]string{
	reg := regexp.MustCompile(`\[.+\]`)
	strmatch := reg.FindAllString(orders, -1)

	if strmatch == nil{
		fmt.Println("strslince nil")
		return nil
	}
	
	pt_order := "{(.*?)}"
	
	if ok, _ := regexp.Match(pt_order, []byte(orders)); ok {  
	
		fmt.Println("match found")  
		return nil
	}

	re, err := regexp.Compile(pt_order)
	if err != nil || re == nil{
		return nil
	}

	strs := re.FindAllString(orders,-1)
	if strs == nil {
		return nil
	}
	
	return strs
}

func GetOneOrderRequire(req string)[]string{            // price    amount   type
	r := regexp.MustCompile("price\\:(\\w*.\\w*)\\,amount\\:(\\w*.\\w*)\\,type\\:\\'(\\w*)\\'")
	if r == nil {
		return nil
	}
	strings := r.FindAllStringSubmatch(req, -1)
	if strings == nil {
		return nil
	}

	return strings[1]
}


type ApiBatchTradeController struct {
	beego.Controller
	Api *models.HttpApi
}

//


func (a ApiBatchTradeController)Post() {

	params := make(map[string]string)

	//=========== access filter =================
	cli_req := a.Ctx.Request
	addr := cli_req.RemoteAddr // "IP:port" "192.168.1.150:8889"
	ip := get_ip(addr)
	if !clientManager.IsAllowAccess(ip){
		str := models.BuildErrorMsg(models.ReqFrequenTooHigh)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	
	//============ param check ==================
	var market string
	a.Ctx.Input.Bind(&market,"market")
	if market == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
	

	params["market"] = market

	var type_para string
	a.Ctx.Input.Bind(&type_para,"type")
	if type_para == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
	
	params["type"] = type_para






	var orders_data  string
	a.Ctx.Input.Bind(&orders_data,"orders_data ")
	if orders_data  == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
	params["orders_data "] = orders_data 
	


	var apikey string


	a.Ctx.Input.Bind(&apikey,"apikey")
	if apikey == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
	params["apikey"] = apikey
	

	var sign string
	a.Ctx.Input.Bind(&sign,"sign")
	if sign == ""{
		str := models.BuildErrorMsg(models.ReqParamIsEmpty)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	params["sign"] = sign
	//beego.Debug(" *** ")



//=============== sign auth ===================  

	userInfoManager := models.GetUserInfoManager()

	var userId int


	userInfo,_ := userInfoManager.GetUserInfo(apikey)
	if userInfo == nil {
		userApiInfo,Err := models.GetUserApiInfo(beego.AppConfig.String("db::database"),apikey)
		if Err != nil {
			str := models.BuildErrorMsg(models.SystemError)
			a.Data["json"] = str
			a.ServeJSON()
			return
		}
		//beego.Debug(" *** ")
		userId = userApiInfo.UserId

		userInfo = new(models.UserInfo)
		userInfo.SetAppKey(userApiInfo.AppKey)
		userInfo.SetSecretKey(userApiInfo.SecretKey)
		userInfo.SetIsReadInfo(userApiInfo.IsReadInfo)
		userInfo.SetIsTrade(userApiInfo.IsTrade)
		userInfo.SetIsWithdraw(userApiInfo.IsWithdraw)
	}


	if !models.IsCorrectSign(userInfo,params) {
		str := models.BuildErrorMsg(models.SignNotMatch)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}
//beego.Debug(" *** ")

	const taker_fee string = "0.0002"
	const maker_fee string = "0.0001"
	//============== internal api request ==================
	
	var side int

	req := httplib.Post(beego.AppConfig.String("MeApiHttpServer"))

	if IdCount < 10000 {
		IdCount ++
	} else {
		IdCount = 0
	}


	orders_req := GetOrderRequires(orders_data)
	if orders_req == nil {
		str := models.BuildErrorMsg(models.IllegalParameter)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	if  len(orders_req) > 5 {
		str := models.BuildErrorMsg(models.IllegalParameter)
		a.Data["json"] = str
		a.ServeJSON()
		return
	}

	oi := new(OrdersInfo)
	oi.Result = false
	oi.Order_info = make([]*OrderResult,len(orders_req))

	for i := 0; i < len(orders_req); i ++ {
		order_req_one := GetOneOrderRequire(orders_req[i])
		if order_req_one == nil {
			str := models.BuildErrorMsg(models.IllegalParameter)
			a.Data["json"] = str
			a.ServeJSON()
			return
		}

		price := order_req_one[0]
		fmt.Printf("************   price: %s\n",price)
		amount := order_req_one[1]
		fmt.Printf("************   amount: %s\n",amount)

		this_type := order_req_one[3]

		if this_type == "buy" {
			side = 2
		} else if this_type == "sell" {
			side = 1
		} else {
			if params["type"] == "buy"{
				side = 2
			} else if params["type"] == "sell" {
				side = 1
			} else {
				str := models.BuildErrorMsg(models.IllegalParameter)
				a.Data["json"] = str
				a.ServeJSON()
				return
			}
		}


		bt := fmt.Sprintf("{\"method\": \"%s\", \"params\": [%d,\"%s\",%d,\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"], \"id\": %d}","put_limit",userId,market,side,amount,price,taker_fee,maker_fee,"empty source",IdCount)
		fmt.Println(bt)

		req.Body(bt)

		str, erri := req.String()
		if erri != nil {
			a.Data["json"] = models.BuildErrorMsg(models.RequestFail)
			a.ServeJSON(false)
			return

		} else{
			ch := make(chan string, 1)
			go func(c chan string, s string){
 				c <- s
    			}(ch, string(str))
    			st := <-ch

			var por PutOrderRes 
			json.Unmarshal([]byte(st), &por)

		//fmt.Println(odp.Id)
		//fmt.Printf("\n%#v\n", por)

			or := new(OrderResult)


			if por.Error != nil{
				or.Error_code = (*por.Error).ErrorCode
				or.Order_id = -1
			} else if por.Res != nil{
				or.Order_id = (*por.Res).Order_id
				or.Error_code = 0
				oi.Result = true
			}


			oi.Order_info[i] = or

		}
	
	}

	resp, err := json.Marshal(oi)
	if err != nil {
		fmt.Println("encoding faild")
	} 

	a.Data["json"] = string(resp)

	a.ServeJSON(false)
	//===================== app response ===================
}



/////////////////////////////////////////////////        /////////////////////////////////////////////////////////////////////

//=========================== recv matchengine message ===============================

//========================== send API response message ===============================
//=============================================================================


/////////////////////////////////////////////////        /////////////////////////////////////////////////////////////////////

//=========================== recv matchengine message ===============================

//========================== send API response message ===============================
//=============================================================================


/////////////////////////////////////////////////        /////////////////////////////////////////////////////////////////////

//=========================== recv matchengine message ===============================

//========================== send API response message ===============================
//=============================================================================