package main

import (
	"game-server/src/mysql"
	"game-server/src/proceed"
	"game-server/src/user"
	"github.com/gin-gonic/gin"
)

func handleFunc(context *gin.Context) {
	user := user.User{
		Username: context.Query("username"),
		Password: context.Query("password"),
	}

	if mysql.CheckUser(user) {
		proceed.HandleWebSocket(context, user.Username)
	} else {
		context.String(400, "用户名或密码错误")
	}
}

func handleRegister(context *gin.Context) {
	user := user.User{
		Username: context.PostForm("username"),
		Password: context.PostForm("password"),
	}

	if !mysql.CheckUser(user) {
		mysql.AddUser(user)
		context.String(200, "账号注册成功")
	}
}

func handleUpdate(context *gin.Context) {
	user := user.User{
		Username: context.PostForm("username"),
		Password: context.PostForm("password"),
	}

	if user.Password == context.PostForm("newPassword") {
		context.String(400, "新密码不能与原密码相同")
	} else {
		if mysql.UpdateUser(user, context.PostForm("newPassword")) {
			context.String(200, "密码更改成功")
		} else {
			context.String(400, "密码更改失败")
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("/", handleFunc)
	router.POST("/", handleFunc)
	router.POST("/register", handleRegister)
	router.POST("/update", handleUpdate)
	router.Run(":8080")
}
