package routers

import (
	

	"quickstart/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/user", &controllers.UserController{})
	beego.Router("/register", &controllers.RegisterController{})
	beego.Router("/login", &controllers.LoginController{})
	// beego.Get("/",func(ctx *context.Context){
	// 	ctx.Output.Body([]byte("bob"))
	// })
}
