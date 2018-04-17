package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
	"net/url"
	"fmt"

)
//初始化数据库
func init() {
	//读取配置文件信息
	host := beego.AppConfig.String("host")
	port := beego.AppConfig.String("port")
	name := beego.AppConfig.String("name")
	username := beego.AppConfig.String("username")
	password := beego.AppConfig.String("password")
	timezone := beego.AppConfig.DefaultString("mysql_timezone", "Asia/Shanghai")
	//数据库连接信息
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&loc=%s", username, password, host, port, name, url.QueryEscape(timezone))
	if host == "" {
		return
	}
	//数据库驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//数据库连接
	orm.RegisterDataBase("default", "mysql", connection)
	//最大连接数
	orm.SetMaxIdleConns("default", 30)

	orm.SetMaxOpenConns("default", 30)
	runmode := beego.AppConfig.DefaultString("runmode", "dev")
	if runmode == "dev" {
		orm.Debug = true
	}
}
//创建一个基类和基类方法
type Base struct{
}
//创建公共的Orm的查询句柄接口
func (c *Base) Orm ()(orm.Ormer){
	o := orm.NewOrm()
	//与刚刚上面的数据库连接名和句柄接口
	err:=o.Using("default")
	if err!=nil{
		panic(err)
	}
	return o

}

//通过Orm的接口，创建QueryTable接口，Querytables 接口
func (c *Base) Query(object interface{}) (orm.QuerySeter) {
	//传进来一个model，但是记住一定要注册过的！！！不然会报错！！！可以通过初始化注册
	qs:=c.Orm().QueryTable(object)
//返回一个查询QuerySeter
	return qs
	}

