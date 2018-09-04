// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"InterfaceAgent/controllers"
	"fmt"
	"github.com/astaxie/beego"
	//"InterfaceAgent/models"
	// "time"
	//"context"
)

type sti struct{
	controllers.HttpApi
}

func init() {


	baseUrl := beego.AppConfig.String("url::baseUrl") + "/"

	fmt.Println("+++++++++base url:\n",baseUrl)
	


	//beego.Router("/api/v1/market.list", &controllers.ApiMarketListController{})
	//beego.Router("/api/v1/balance.query",&controllers.ApiBalanceQueryController{})
	//beego.Router("/api/v1/market.list", &controllers.UserController{})
	//beego.Get("/api/v1/hello",func(ctx *context.Context){ctx.Output.Body([]byte("hello world"))})


	fmt.Println("= routers : router ===-===== =-=-=-=-=-=-=-=-=-")
	

	// beego.AddNamespace(ns)
	//beego.AddNamespace(api_ns)
}

