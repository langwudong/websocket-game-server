package main

import (
	"game-server/src/mysql"
	"game-server/src/room"
	"game-server/src/user"
	"github.com/gin-gonic/gin"
)

func handleFunc(context *gin.Context) {
	user := user.User{
		Username: context.Query("username"),
		Password: context.Query("password"),
	}

	if mysql.CheckUser(user) {
		room.HandleWebSocket(context, user.Username)
	} else {
		context.String(400, "用户名或密码错误!")
	}
}

//func createRoom(context *gin.Context) {
//	r := &room.Room{
//		Name:       context.PostForm("name"),
//		Players:    make([]*room.Player, 0),
//		Register:   make(chan *room.Player),
//		Unregister: make(chan *room.Player),
//		Gamestate:  make(chan []byte),
//	}
//	room.GetInstance().Rooms = append(room.GetInstance().Rooms, r)
//}
//
//func searchRoom(context *gin.Context) {
//	room.GetInstance().Rooms
//}

func main() {
	router := gin.Default()
	router.GET("/", handleFunc)
	router.POST("/", handleFunc)
	router.Run(":8080")
}
