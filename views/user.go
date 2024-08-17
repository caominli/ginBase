package view

import (
	common "gin_scaffold/commons"
	jwt "gin_scaffold/jwtmods"
	model "gin_scaffold/models"
	"log"
	"time"

	"github.com/gin-gonic/gin" //导入gin包
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 注册
func Register(c *gin.Context) {

	//先检查表单错误
	var form model.Register //定义一个表单变量

	err := c.BindJSON(&form) //执行绑定
	if err != nil {          //如果验证失败
		//创建消息
		c.JSON(400, gin.H{"msg": model.GetValidMsg(err, &form)})
		return
	}

	//执行数据库查询用户是否已存在
	var user model.Users                           //定义数据库查询变量
	model.DB.First(&user, "email = ?", form.Email) //执行查询
	if user.ID != 0 {                              //如果查询到用户
		c.JSON(400, gin.H{"msg": "邮箱已被注册，请直接登录，若忘记密码请修改密码"})
		return
	}

	code := common.DaXie(form.Code) //验证码转大写
	// 核验验证码
	if !model.ValidateCode(form.Email, code) { //核验验证码
		c.JSON(400, gin.H{"msg": "验证码错误"})
		return
	}

	//开始写入数据库逻辑
	err = model.DB.Create(&model.Users{Email: form.Email, Password: setPassword(form.Password)}).Error

	if err != nil {
		c.JSON(500, gin.H{"msg": "Database error, please try again later or report this issue to our customer support team"})
		log.Print("注册用户写入数据库错误：", err)
		return
	}
	c.JSON(200, gin.H{"msg": true})

}

// 登录页页
func Login(c *gin.Context) {
	//先检查表单错误
	var form model.Login           //定义一个表单变量
	err := c.ShouldBindJSON(&form) //执行绑定
	if err != nil {                //如果验证失败
		//返回验证消息
		c.JSON(400, gin.H{"msg": model.GetValidMsg(err, &form)})
		return
	}

	var user model.Users
	//查询数据库
	model.DB.First(&user, "email = ?", form.Email)
	if user.ID == 0 { //如果未查询到用户
		c.JSON(400, gin.H{"msg": "账户或密码错误"})
		return
	}
	// 验证密码如果没通过
	if ok := vPassword(user.Password, form.Password); !ok {
		c.JSON(400, gin.H{"msg": "账户或密码错误"})
		return
	}

	//密码正确后的逻辑
	//调用jwt签发函数并保存到cookie，也可以直接返回json的token让客户端自己处理
	tokenString, err := jwt.GenerateJWT(user.ID, user.Auth)
	if err != nil {
		c.JSON(500, gin.H{"msg": "token生成错误,请稍后再试或联系客服"})
		return
	}
	log.Print("下发的token为:", tokenString)

	c.JSON(200, gin.H{"token": tokenString})

}

// 修改密码
func RePassword(c *gin.Context) {
	//先检查表单错误
	var form model.Register //定义一个表单变量

	err := c.ShouldBindJSON(&form) //执行绑定
	if err != nil {                //如果验证失败
		c.JSON(400, gin.H{"msg": model.GetValidMsg(err, &form)})
		return
	}

	//执行数据库查询用户是否有这个用户
	var user model.Users                           //定义数据库查询变量
	model.DB.First(&user, "email = ?", form.Email) //执行查询
	if user.ID == 0 {                              //如果查询到用户
		c.JSON(400, gin.H{"msg": "未找到用户，请检查邮箱是否正确"})
		return
	}

	code := common.DaXie(form.Code) //验证码转大写
	// 核验验证码
	if !model.ValidateCode(form.Email, code) { //核验验证码
		c.JSON(400, gin.H{"msg": "验证码错误"})
		return
	}

	//更新数据库逻辑
	user.Password = setPassword(form.Password)
	if model.DB.Save(&user).Error != nil {
		c.JSON(500, gin.H{"msg": "Database error, please try again later or report this issue to our customer support team"})
		log.Print("注册用户更新密码数据库错误：", err)
		return
	}

	//更新密码成功
	c.JSON(200, gin.H{"msg": true})

}

// 定义验证码的json结构体
type EmailJson struct {
	Email string `json:"email" binding:"required,email,max=200" msg:"输入正确的邮箱地址,且长度小于200字符"`
}

// 获取验证码API
func Getcode(c *gin.Context) {

	// 解析请求体中的 JSON 数据到结构体
	var emailjson EmailJson
	if err := c.BindJSON(&emailjson); err != nil {
		//返回错误信息
		c.JSON(400, gin.H{"msg": model.GetValidMsg(err, &emailjson)})
		return
	}
	// 查找或创建记录
	var user model.Captchas
	//查询一条数据
	result := model.DB.First(&user, "email = ?", emailjson.Email) //搜索条件为name字段，搜索值为传递的路由参数
	if result.Error != nil {
		// 判断是否查询到数据
		if result.Error != gorm.ErrRecordNotFound {
			//如果不是空值则代表数据库其他错误
			c.JSON(500, gin.H{"msg": "数据库错误"}) //返回查询失败
			log.Print("验证码数据库错误查询错误：", result.Error)
			return
		}
		log.Print("没有这条数据开始创建")
		//如果未查询到数据则创建它
		code := common.Captcha(6) //获得验证码

		//写入数据
		result = model.DB.Create(&model.Captchas{Email: emailjson.Email, Code: code, Date: time.Now()})
		if result.Error != nil {
			c.JSON(500, gin.H{"msg": "The database encountered an error, please try again later or contact customer service"}) //返回查询失败
			log.Print("验证码数据库创建失败：", result.Error)
			return
		}

		//为邮箱发送验证码
		go common.Sendmail(emailjson.Email, code)
		//返回成功
		c.JSON(200, gin.H{"msg": true})
		return
	}

	//计算当前时间和对比时间的时间差
	timeDiff := time.Since(user.Date)
	// 判断时间差是否超过 2 分钟
	if timeDiff < 2*time.Minute {

		// 如果没有超过2分钟,返回成功和通知
		c.JSON(400, gin.H{"msg": "验证码发送需2分钟间隔,请稍后再试"})
		log.Println("时间差没有超过 2 分钟")
		return
	}

	//超过2分钟则更新验证码和时间
	code := common.Captcha(6) //获得验证码

	user.Code = code
	user.Date = time.Now()
	result = model.DB.Save(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{"msg": "验证码重新生成失败"})
		log.Println("验证码重新生成失败：", result.Error)
		return
	}

	//为邮箱发送验证码
	go common.Sendmail(emailjson.Email, code)
	//返回成功
	c.JSON(200, gin.H{"msg": true})
}

// 设置密码，参数：str，返回：[]byte
func setPassword(text string) []byte {
	// 对密码进行哈希加密
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		log.Print("密码加密错误：", err)
		return nil
	}
	return passwordHash
}

// 验证密码，参数：[]byte，str，返回：bool
func vPassword(modelupassword []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(modelupassword, []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}
