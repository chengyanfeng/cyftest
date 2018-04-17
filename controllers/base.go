package controllers

import (
	"github.com/astaxie/beego"

)

type BaseController struct {
	beego.Controller
}
var userchild=[]map[string]interface{}{
	map[string]interface{}{
		"Path":"/user/getlist",
		"Name":"用户系统管理",
		"On":0,
	},
	map[string]interface{}{
		"Path":"/user/getlist",
		"Name":"团队系统管理",
		"On":0,
	},
}
var Menu=[]map[string]interface{}{
	map[string]interface{}{
		"On":0,
		"Name":"用户管理",
		"Child":userchild,
	},
	map[string]interface{}{
		"On":0,
		"Name":"用户管理",
		"Child":"nil",
	},

}

func (c *BaseController)init(){
		c.Layout="common/layout.html"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["header"]="common/head.html"
		c.LayoutSections["footer"]="common/footer.html"

		c.Data["Menu"]=Menu
}

func (c *BaseController)Page(){

}

