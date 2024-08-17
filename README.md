# gin开发脚手架说明

## 开发时注意
使用时需自行修改mod包名，并替换所有页面中自己实现的包导入路径
修改config.jsonc配置文件,设置为自己需要的

## 编译时注意
编辑router.go，将中间的cors跨域允许，静态文件服务注释掉

## 目录结构

```
📦ginSerBase
 ┣ 📂commons        封装方法文件夹
 ┃ ┣ 📜text.go      文本处理
 ┃ ┗ 📜email.go     邮件发送
 ┣ 📂routers        路由目录
 ┃ ┗ 📜router.go    路由处理
 ┣ 📂views          视图业务目录
 ┃ ┣ 📜view.go      主视图函数
 ┃ ┣ 📜api.go       Api函数
 ┃ ┣ 📜user.go      用户函数
 ┃ ┗ 📜wxlogin.go   微信登录
 ┣ 📂templates      模板目录
 ┃ ┣ 📜index.html   主页
 ┃ ┣ 📜header.html  页头
 ┃ ┗ 📜bottom.html  页尾
 ┣ 📂models         数据模型
 ┃ ┣ 📜init.go      初始化数据库
 ┃ ┣ 📜user.go      用户模型
 ┃ ┗ 📜items.go     项目模型
 ┣ 📂jwtmods        jwt目录
 ┃ ┗ 📜jwtmod.go    jwt处理
 ┣ 📂config         配置文件
 ┃ ┣ 📜getconfig.go 取配置
 ┃ ┗ 📜json.go      自己封装的json读配置，支持json使用注释 
 ┣ 📜go.mod         依赖信息
 ┣ 📜main.go        主文件
 ┣ 📜config.jsonc   配置文件
 ┗ 📜go.sum         依赖文件
```