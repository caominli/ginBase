package common

import (
	"fmt"
	config "gin_scaffold/config"
	"log"
	"net/smtp"
	"strings"
)

// 发送验证码，对对方邮箱，验证码
func Sendmail(dfemail string, code string) {
	body := fmt.Sprintf(`<div style="background-color: #f2f2f2; padding: 20px;">
	<h2 style="color: #333; font-size: 24px; margin-bottom: 20px;">这是您的验证码：</h2>
	<div style="background-color: #fff; border: 1px solid green; color:green; border-radius: 5px; padding: 10px; font-size: 24px; font-weight: bold; text-align: center; margin-bottom: 20px;">
	%s
	</div>
	<p style="color: #666; font-size: 16px;">此邮件地址仅供发送邮件，请勿回复此邮件！</p>
	</div>`, code)
	Sendmail_oracle("notify@btcmai.com", dfemail, config.Config.WebName+"的安全验证", body)
}

// 邮箱发送，参数：我方邮箱，对方邮箱，标题，内容
func Sendmail_oracle(from string, to string, title string, leirong string) error {
	//登录账号
	sender_username := config.Config.EmailUser
	//登录密码
	sender_password := config.Config.EmailPassword
	//拼接邮件内容
	body := "To: " + to + "\r\n" +
		"Subject: " + title + "\r\n" +
		"Content-Type: text/html; charset=UTF-8" + "\r\n\r\n" +
		leirong

	smtpServer := config.Config.EmailHost
	auth := smtp.PlainAuth("", sender_username, sender_password, strings.Split(smtpServer, ":")[0])

	err := smtp.SendMail(smtpServer, auth, from, []string{to}, []byte(body))
	if err != nil {
		log.Println("发送邮件失败:", err)
		return err
	}
	//没有错误返回空
	return nil
}

// //使用qq发送邮箱
// // 参数：对方邮箱，标题，内容
// func tuisongemail(dfemail string,title string, body string) {
// 	auth := smtp.PlainAuth("", "13320807", "这里是密码", "smtp.qq.com")
// 	// 邮件内容
// 	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
// 	msg := "From: 13320807@qq.com\n" + "Subject: " + title + "\n" + mime + "\n" + body

// 	// 发送邮件
// 	err := smtp.SendMail("smtp.qq.com:587", auth, "13320807@qq.com", []string{dfemail}, []byte(msg))
// 	if err != nil {
// 		log.Print("QQ邮箱发送失败:", err)
// 		return
// 	}
// }
