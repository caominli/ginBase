package models //定义在models包下

import (
	"reflect" //验证器定义错误需要
	"time"

	"github.com/go-playground/validator/v10" //导入自定义验证器
)

// 定义一个User的数据库模型
type Users struct {
	ID        uint      `gorm:"primaryKey"`                    //id
	Email     string    `gorm:"uniqueIndex;not null;size:200"` //邮箱
	Password  []byte    `gorm:"not null"`                      //密码
	CreatedAt time.Time `gorm:"autoCreateTime"`                //创建时间
	Auth      uint8     `gorm:"default:1"`                     //权限等级0封禁，,1普通用户,2高级用户，3管理员，4超级管理员
	// Carts             []Carts
	// Orders			[]Orders
}

// 验证码数据库模型
type Captchas struct {
	ID    uint      `gorm:"primaryKey"`                    //id
	Email string    `gorm:"uniqueIndex;not null;size:200"` //邮箱
	Code  string    `gorm:"not null;size:6"`               //代码
	Date  time.Time `gorm:"not null"`                      //生成时间
}

// 注册表单验证器
type Register struct {
	Email    string `json:"email" binding:"required,email,max=200" msg:"必须是有效的邮箱,不超过200字符"` //用户名
	Password string `json:"password" binding:"required,min=6,max=50" msg:"密码必须是6-50字符之间"`   //密码,必填，小于5，大于10位数
	Code     string `json:"vcode" binding:"required,alphanum,len=6" msg:"验证码为6个字符"`         //验证码
}

// 登录表单
type Login struct {
	Email    string `json:"email" binding:"required,email,max=200" msg:"不是正确的邮箱地址"`    //用户名
	Password string `json:"password" binding:"required,min=6,max=50" msg:"密码是6-50个字符"` //密码,必填，小于5，大于10位数
}

// 自定义验证器结构体的错误信息，接收err错误对象和结构体的指针，返回str类型
func GetValidMsg(err error, obj interface{}) string {
	//获取结构体的类型信息
	getObj := reflect.TypeOf(obj)
	// 检查错误是否属于validator.ValidationErrors类型
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			// 检查是否可以通过字段名获取结构体字段信息
			if f, exist := getObj.Elem().FieldByName(e.Field()); exist {
				return f.Tag.Get("msg") //获取自定义的结构体msg参数
			}
		}
	}
	// 如果无法通过验证器错误获取自定义消息，返回默认错误消息
	return err.Error()
}

// 验证码核验，传入邮箱，表单验证码，返回bool
func ValidateCode(email string, code string) bool {
	//查询数据
	var user Captchas                   //定义数据库查询变量
	DB.First(&user, "email = ?", email) //执行查询
	if user.ID == 0 {
		//没有查询到用户
		return false
	}
	//判断时间是否在20分钟内
	if time.Since(user.Date) > time.Minute*20 {
		//超过20分钟
		return false
	}
	//核验验证码
	if user.Code == code {
		return true
	}
	return false
}
