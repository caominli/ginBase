package main //主程序

import (
	route "gin_scaffold/routers" //导入gin运行路由
	config "gin_scaffold/config" //导入配置包
)



func main() {
	r := route.Router()
	//开始启动
	r.Run(":"+config.Config.Port)
}
