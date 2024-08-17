package view

import (
	model "gin_scaffold/models" //导入models包
	"log"

	"github.com/gin-gonic/gin" //导入gin包
)

// 不需要登陆的api，返回items
func Getitem(c *gin.Context) {
	var items []model.Items
	//返回items内容
	model.DB.Find(&items)
	c.JSON(200, gin.H{"items": items})
}

// 定义一个用户信息结构体，因为密码不需要返回，所以不包含密码
type user_info struct {
	ID        uint   `json:"id"` //用户ID
	Email     string `json:"email"`
	Auth      uint8  `json:"auth"`
	CreatedAt string `json:"created_at"`
}

// 需要登陆的api,验证写在路由的中间件,返回用户信息
func GetUser(c *gin.Context) {
	//从中间件取用户并断言为浮点数再转为uint
	useridFloat, ok := c.MustGet("userid").(float64)
	if !ok {
		log.Print("ApiGetUser函数用户ID断言为浮点数错误")
		c.JSON(500, gin.H{"msg": "System encountered an error. Please try again later or contact customer support."})
		return
	}
	userid := uint(useridFloat)
	//查询用户信息
	var user model.Users
	result := model.DB.First(&user, "id = ?", userid)
	if result.RowsAffected == 0 {
		// 没有找到记录
		c.JSON(404, gin.H{"msg": "没有这个用户"})
		return
	}
	var userinfo user_info
	userinfo.Email = user.Email
	userinfo.ID = user.ID
	userinfo.Auth = user.Auth
	userinfo.CreatedAt = user.CreatedAt.String()
	//返回用户信息
	c.JSON(200, gin.H{"data": userinfo})
}
