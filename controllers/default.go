package controllers

import (
	"fmt"
	"sort"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	c.init()
test:=	make([]int,0)
fmt.Println(test)
var  t= [5]int{}
fmt.Println(t)
	b:=[]int{1,2,5,3,2}
	sort.Ints(b)

	fmt.Print(b,"b")
	c.TplName = "common/mban.html"
}
