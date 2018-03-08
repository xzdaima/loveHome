package routers

import (
	"github.com/astaxie/beego"
	"loveHome/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("api/v1.0/areas", &controllers.AreasController{}, "get:GetAreaInfo")

	beego.Router("api/v1.0/session", &controllers.SessionController{}, "get:GetSessionName;delete:DeleteSession")
	beego.Router("api/v1.0/sessions", &controllers.UserController{}, "post:Login")
	beego.Router("api/v1.0/users", &controllers.UserController{}, "post:Reg;get:GetUserInfo")
	beego.Router("api/v1.0/user", &controllers.UserController{}, "get:GetUserInfo")
	beego.Router("api/v1.0/user/avatar", &controllers.UserController{}, "post:UploadAvatar")

	beego.Router("api/v1.0/houses/index", &controllers.HousesIndexController{}, "get:GetHousesIndex")
}
