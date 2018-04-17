package routers

import (
	"cyftest/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{},"get:Get")
}
