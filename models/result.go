package models

import(
	"encoding/json"
	"fmt"
)

type ErrorCode int

const (
	ReqParamIsEmpty ErrorCode = 10000     //Required parameters can not be empty
	ReqFrequenTooHigh      	= 10001             //用户请求频率过快，超过该接口允许的限额
	SystemError             	= 10002     //
	RequestFail          		= 10004
	ParamError        		= 10008
	SecretKeyNotExit     		= 10005 	//SecretKey不存在
	ApiNotExit				= 10006 	//Api_key不存在
	SignNotMatch			= 10007 	//签名不匹配
	IllegalParameter		= 10008 	//非法参数
	OrderNotExit			= 10009 	//订单不存在
	InsufficienBalance		= 10010 	//余额不足
// 10011 	买卖的数量小于BTC/LTC最小买卖额度
// 10012 	当前网站暂时只支持btc_usd ltc_usd
// 10013 	此接口只支持https请求
// 10014 	下单价格不得≤0或≥1000000
// 10015 	下单价格与最新成交价偏差过大
// 10016 	币数量不足
// 10017 	API鉴权失败
// 10018 	借入不能小于最低限额[usd:100,btc:0.1,ltc:1]
// 10019 	页面没有同意借贷协议
// 10020 	费率不能大于1%
// 10021 	费率不能小于0.01%
// 10023 	获取最新成交价错误
// 10024 	可借金额不足
// 10025 	额度已满，暂时无法借款
// 10026 	借款(含预约借款)及保证金部分不能提出
// 10027 	修改敏感提币验证信息，24小时内不允许提现
// 10028 	提币金额已超过今日提币限额
// 10029 	账户有借款，请撤消借款或者还清借款后提币
// 10031 	存在BTC/LTC充值，该部分等值金额需6个网络确认后方能提出
// 10032 	未绑定手机或谷歌验证
// 10033 	服务费大于最大网络手续费
// 10034 	服务费小于最低网络手续费
// 10035 	可用BTC/LTC不足
// 10036 	提币数量小于最小提币数量
// 10037 	交易密码未设置
// 10040 	取消提币失败
// 10041 	提币地址不存在或未认证
// 10042 	交易密码错误
// 10043 	合约权益错误，提币失败
// 10044 	取消借款失败
// 10047 	当前为子账户，此功能未开放
// 10048 	提币信息不存在
// 10049 	小额委托（<0.15BTC)的未成交委托数量不得大于50个
// 10050 	重复撤单
// 10052 	提币受限
// 10064 	美元充值后的48小时内，该部分资产不能提出
// 10100 	账户被冻结
// 10101 	订单类型错误
// 10102 	不是本用户的订单
// 10103 	私密订单密钥错误
	NoOpenApi				= 10216 	//非开放API
	DataBaseAccessError		= 10217   //
)

type ErrorRes struct {
	ErrorCode int  `json:"error_code"`
	Result bool      `json:"result"`
}


func BuildErrorMsg(errorCode ErrorCode)string{
	er := &ErrorRes{
		int(errorCode),
		false,
	}

	rp,err := json.Marshal(er)

	if err != nil{
		fmt.Println("encoding faild")
	}

	return string(rp)
}


// {
//     "error": {
//         "code": 1,
//         "message": "invalid argument"
//     },
//     "result": null,
//     "id": 123
// }

type ErrorDetail struct{
	ErrorCode int      `json:"code"`
	Msg         string  `json:"message"`
}

