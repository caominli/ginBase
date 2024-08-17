package router

import (
	jwt "gin_scaffold/jwtmods" //导入共享包
	view "gin_scaffold/views"
	"html/template"

	"github.com/gin-contrib/cors" //允许跨域，生产环境需注释
	"github.com/gin-gonic/gin"
)

// 让html不转义
func Html(x string) interface{} {
	return template.HTML(x)
}

func Router() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode) //正式版模式
	r := gin.Default() //初始化gin
	//允许跨域，生产环境需注释
	r.Use(cors.New(cors.Config{ //生产环境需注释
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//自定义过滤器
	r.SetFuncMap(template.FuncMap{
		"html": Html, //让html不转义
	})
	r.LoadHTMLGlob("templates/*") //加载html模板全局，指定templates目录下的所有文件
	//静态文件服务，生产环境需要注释，因为使用反向代理工具来处理静态文件
	r.Static("/static", "./static")

	//下面是曝光路由，由gin判断登录则vue处理，未登录则gin处理
	r.GET("/", view.Index) //主页智能定向语言
	api := r.Group("/api")
	{
		api.POST("/register", view.Register)               //注册
		api.POST("/login", view.Login)                     //登录
		api.POST("/getcode", view.Getcode)                 //获取验证码
		api.POST("/repassword", view.RePassword)           //重置密码
		api.POST("/wxlogin", view.WxLogin)                 //微信登录
		api.GET("/getuser", jwt.AuthUser(1), view.GetUser) //取用户信息，需要登录，核验权限写在中间件
	}
	return r
}
