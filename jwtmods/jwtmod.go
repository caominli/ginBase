// JWT模块，用于生成JWT令牌和验证JWT令牌
package jwtmod

import (
	config "gin_scaffold/config" //导入共享包
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 定义 JWT 签名密钥
var JwtKey = []byte(config.Config.JwtKey)

// JWT 签发函数
func GenerateJWT(userid uint, auth uint8) (string, error) {
	// 创建 JWT Claims 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		//设置超时时间7天
		"exp":    time.Now().Add(time.Hour * 168).Unix(),
		"userid": userid, //这是存入的用户名
		"auth":   auth,   //权限等级
	})
	// 生成 JWT Token
	tokenString, err := token.SignedString(JwtKey)
	// 使用密钥进行签名生成字符串 token
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 必须登陆用户才可以访问的中间件函数，传入参数为权限等级，1+开始为登录用户，如果验证成功返回用户id在上下文中
func AuthUser(auth2 uint8) gin.HandlerFunc {
	return func(c *gin.Context) { //获取cookie中的token字段
		// 获取 Authorization 头
        tokenString := c.GetHeader("Authorization")

		// 如果没有提供信息
		if tokenString == "" {
			//创建消息
			c.JSON(401, gin.H{"msg": "please sign in"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			//这里的参数就是方法，方法直接返回jwtKey令牌
			return JwtKey, nil
		})
		// 检查token是否有效
		if !token.Valid {
			//如果token无效
			c.JSON(401, gin.H{"msg": "The token has expired, please log in"})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userid := claims["userid"] //取得用户userid
			//取得用户权限
			auth3, ok := claims["auth"].(float64)
			if !ok {
				log.Print("jwt模块AuthUser函数用户权限转float64错误")
				c.JSON(500, gin.H{"msg": "A system error occurred. Please try again later or contact customer service."})
				c.Abort()
				return
			}
			//转换为uint8
			auth := uint8(auth3)

			//如果用户权限小于要求的权限
			if auth < auth2 {
				c.JSON(403, gin.H{"msg": "No access permission"})
				c.Abort()
				return
			}
			// 将解析后的claims存储到上下文中，以便后续处理程序可以访问
			c.Set("userid", userid) //返回用户id
		} else {
			c.JSON(500, gin.H{"token_err": err})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 管理员可访问的中间件函数，传入参数为权限等级，3+开始为管理员
func AdminAuth(auth2 uint8) gin.HandlerFunc {
	return func(c *gin.Context) { //获取cookie中的token字段
		// 获取 Authorization 头
        tokenString := c.GetHeader("Authorization")
		// 如果没有提供信息
		if tokenString == "" {
			//直接返回404
			c.Status(404)
			c.Abort()
			log.Print("有用户未登录访问管理员界面")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			//这里的参数就是方法，方法直接返回jwtKey令牌
			return JwtKey, nil
		})
		// 检查token是否有效
		if !token.Valid {
			//如果token无效
			c.Status(404)
			c.Abort()
			log.Print("有用户使用无效token访问管理员界面")
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userid := claims["userid"] //取得用户userid
			//取得用户权限
			auth3, ok := claims["auth"].(float64)
			if !ok {
				c.Status(404)
				c.Abort()
				log.Print("jwt模块Admin_auth函数用户权限转float64错误")
				return
			}
			//转换为uint8
			auth := uint8(auth3)

			//如果用户权限小于要求的权限
			if auth < auth2 {
				c.Status(404)
				log.Print("用户", userid, "权限不足访问管理员界面")
				c.Abort()
				return
			}
			// 将解析后的claims存储到上下文中，以便后续处理程序可以访问
			c.Set("userid", userid) //返回用户id
		} else {
			c.JSON(500, gin.H{"token_err": err})
			c.Abort()
			return
		}
		c.Next()
	}
}
