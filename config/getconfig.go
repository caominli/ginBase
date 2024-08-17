package config

//这个包用来读取jsonc配置，并存储在全局变量中
import (
	"log"
	"os"
)

// 定义读取配置变量的结构体
// 取配置变量
type config struct {
	WebName       string `json:"webname"`       //网站名称
	Port          string `json:"port"`          //gin运行端口
	Host          string `json:"host"`          //数据库地址
	JwtKey        string `json:"jwtkey"`        //32位安全密钥,jwt需要
	DBhost        string `json:"dbhost"`        //数据库地址
	DBport        string `json:"dbport"`        //数据库端口，字符串格式
	DBname        string `json:"dbname"`        //数据库名称
	DBuser        string `json:"dbuser"`        //数据库用户名
	DBpassword    string `json:"dbpassword"`    //数据库密码
	EmailUser     string `json:"emailuser"`     //邮箱用户名
	EmailPassword string `json:"emailpassword"` //邮箱密码
	EmailHost     string `json:"emailhost"`     //邮箱服务器地址
	WxAppid       string `json:"wxappid"`       //微信小程序appid
	WxSecret      string `json:"wxsecret"`      //微信小程序secret
}

// 定义一个全局变量，用来存储配置信息
var Config config

func init() {
	// 读取JSON设置
	data, err := os.ReadFile("config.jsonc")
	if err != nil {
		log.Print("json配置文件无法读取", err)
		panic("json配置文件无法读取，关闭程序")
	}
	//解析带注释的JSON文件
	err = JsonGet(data, &Config)
	if err != nil {
		panic("json配置文件无法解析，关闭程序:" + err.Error())
	}
}
