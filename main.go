package main

import (
	_ "cyftest/routers"
	"github.com/astaxie/beego"
)

func main() {
	fmt.print(“aaa”)
	beego.Run()
}

