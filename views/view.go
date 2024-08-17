package view

import (
	"log"

	"github.com/gin-gonic/gin" //导入gin包
)

// 服务器渲染
func Index(c *gin.Context) {
	log.Print("运行了index")
	//定义变量
	title := "首页"
	msg := "欢迎来到首页"
	//返回html并传递变量
	c.HTML(200, "index.html", gin.H{title: title, msg: msg})
}
