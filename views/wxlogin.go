package view

import (
	// common "gin_scaffold/commons"
	config "gin_scaffold/config"
	// jwt "gin_scaffold/jwtmods"
	// model "gin_scaffold/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin" //导入gin包
)

// 绑定客户端提交的事件code
type wxcode struct {
	Code string `json:"code"`
	// device string `json:"device"` //前端平台
}


// 微信登录
func WxLogin(c *gin.Context) {
	//先检查表单错误
	var form wxcode          //定义一个表单变量
	err := c.BindJSON(&form) //执行绑定
	if err != nil {          //如果验证失败
		//创建消息
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}

	log.Print("客户端提交的code: ", form.Code)
	res, err := http.Get("https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + config.Config.WxAppid + "&secret=" + config.Config.WxSecret + "&code=" + form.Code +  "&grant_type=authorization_code")
	// res, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + config.Config.WxAppid + "&secret=" + config.Config.WxSecret + "&js_code=" + form.Code +  "&grant_type=authorization_code")
	if err != nil {
		log.Printf("微信登录请求错误: %v", err)
	}
	defer res.Body.Close()

	//取返回主主体
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("微信登录取返回主体错误: %v", err)
		c.JSON(500, gin.H{"msg": "遇到错误，请稍后再试或联系客服"})
		return
	}
	log.Printf("微信返回的所有内容: %s,状态码:%d", body, res.StatusCode)

	// 解析响应体为通用的 map
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("微信登录返回json解码错误: %v", err)
		c.JSON(500, gin.H{"msg": "遇到错误，请稍后再试或联系客服"})
		return
	}
	// 检查是否包含错误码
	if errcode, ok := result["errcode"]; ok {
		log.Printf("微信登录返回错误: errcode=%v, errmsg=%v", errcode, result["errmsg"])
		c.JSON(500, gin.H{"msg": fmt.Sprintf("微信登录返回错误: %v", result["errmsg"])})
		return
	}

	// 返回成功响应
	c.JSON(200, gin.H{
		"access_token":  result["access_token"],
		"expires_in":    result["expires_in"],
		"refresh_token": result["refresh_token"],
		"openid":        result["openid"],
		"scope":         result["scope"],
		"unionid":       result["unionid"],
	})
}
