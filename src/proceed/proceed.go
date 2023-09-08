package proceed

import (
	"encoding/json"
	"game-server/src/message"
	"game-server/src/upgrader"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"sync"
)

type rooms struct {
	Rooms []*Room `json:"rooms"`
}

type Room struct {
	ID string `json:"room_id,omitempty"`
	//玩家列表
	Players []*Player `json:"players"`
}

type Player struct {
	Username  string          `json:"username,omitempty"`
	RoomID    string          `json:"roomID,omitempty"`
	Conn      *websocket.Conn `json:"conn,omitempty"`
	GameState string          `json:"gameState,omitempty"`
	Index     int             `json:"index,omitempty"`
}

var (
	instance *rooms
	once     sync.Once
	players  = make(map[*websocket.Conn]*Player)
)

func GetInstance() *rooms {
	//使用单例模式创建rooms,注意需要初始化rooms.Rooms
	once.Do(func() {
		instance = &rooms{}
		instance.Rooms = make([]*Room, 0)
	})
	return instance
}

func HandleWebSocket(context *gin.Context, username string) {
	//将http连接升级为WebSocket连接
	conn, err := upgrader.Upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		logrus.Error("Failed to send message:", err)
		return
	}

	//初始化玩家相关数据
	player := &Player{
		Username: username,
		Conn:     conn,
	}

	//将玩家添加到玩家列表中
	players[conn] = player

	//持续接收玩家的消息
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.Error("Failed to send message:", err)
			break
		}

		t := message.Message{}
		json.Unmarshal(msg, &t)
		//处理收到的消息
		processMessage(player, t.Type, msg)
	}
	//当玩家连接断开时,执行相关清理操作
	disconnectPlayer(player)
}

func processMessage(player *Player, t string, msg []byte) {
	switch t {
	case "display_rooms":
		sendMessage(player, GetInstance())
	case "create_room":
		quitRoom(player)
		//初始化房间相关数据
		room := &Room{
			ID:      generateUniqueID(player),
			Players: []*Player{player},
		}
		player.RoomID = room.ID
		//GetInstance().Rooms[proceed.ID] = proceed
		GetInstance().Rooms = append(GetInstance().Rooms, room)

		//向玩家发送房间号等信息
		response := map[string]interface{}{
			"room_id": room.ID,
		}
		sendMessage(player, response)
	case "join_room":
		//解析消息中的房间号信息
		roomID := parseRoomIDFromMessage(msg)

		//解析到Room
		var r *Room
		err := json.Unmarshal(msg, &r)

		if err != nil {
			logrus.Error(err)
		}

		if r.Players != nil && len(r.Players) > 0 {
			//获取Room中保存的玩家信息并解析到player
			b, _ := json.Marshal(r.Players[0])
			if b != nil {
				json.Unmarshal(b, &player)
			}
		}

		//如果要加入的房间号,不等于玩家已加入的房间号,则先退出已加入的房间
		if roomID != player.RoomID {
			quitRoom(player)
			if ok, room := Contain(GetInstance().Rooms, roomID); ok {
				player.RoomID = roomID
				room.Players = append(room.Players, player)

				//向玩家发送房间信息、其他玩家列表
				response := map[string]interface{}{
					"room_id": room.ID,
					"players": room.Players,
				}
				sendMessage(player, response)
			} else {
				//房间不存在或已满,发送错误消息给玩家
				response := map[string]interface{}{
					"error": "Room not found or full.",
				}
				sendMessage(player, response)
			}
		} else {
			if ok, room := Contain(GetInstance().Rooms, roomID); ok {
				//向玩家发送房间信息、其他玩家列表
				response := map[string]interface{}{
					"room_id": room.ID,
					"players": room.Players,
				}
				sendMessage(player, response)
			}
		}
	case "normal":
		var data map[string]interface{}
		err := json.Unmarshal(msg, &data)
		if err != nil {
			logrus.Error("Failed to send message:", err)
		}

		//err := json.Unmarshal(message, &other)
		//if err != nil {
		//	logrus.Error("Failed to convert json:", err)
		//}

		broadcastMessage(player, QueryRoom(GetInstance().Rooms, player), data)
	default:
		logrus.Error("The type of message not found.")
	}
}

func disconnectPlayer(player *Player) {
	//玩家断开连接后,从玩家列表中移除玩家,并关闭玩家的连接
	quitRoom(player)
	delete(players, player.Conn)
	player.Conn.Close()
}

func generateUniqueID(player *Player) string {
	return player.Username
}

func parseRoomIDFromMessage(msg []byte) string {
	var data map[string]interface{}
	err := json.Unmarshal(msg, &data)
	if err != nil {
		//返回空房间号,表示解析失败
		logrus.Error("Failed to send message:", err)
		return ""
	}

	roomID, ok := data["room_id"].(string)
	if !ok {
		//返回空房间号,表示解析失败
		logrus.Error("Failed to send message: RoomID not found.")
		return ""
	}
	return roomID
}

func quitRoom(player *Player) {
	if player.RoomID != "" {
		room := QueryRoom(GetInstance().Rooms, player)

		if room != nil {
			//删除房间中玩家的信息
			if room.Players != nil || len(room.Players) > 0 {
				for index, value := range room.Players {
					if value.Conn == player.Conn {
						//从房间的玩家列表里移除玩家
						room.Players = append(room.Players[:index], room.Players[index+1:]...)
					}
				}
			}
		}

		if room != nil {
			//如果房间中没有玩家,则删除该房间信息
			if room.Players == nil || len(room.Players) == 0 {
				//delete(GetInstance().Rooms, player.RoomID)
				RemoveRoom(&GetInstance().Rooms, room)
			}
		}

		player.RoomID = ""
	}
}

func sendMessage(player *Player, msg interface{}) {
	err := player.Conn.WriteJSON(msg)
	if err != nil {
		logrus.Error("Failed to send message:", err)
	}
}

func broadcastMessage(player *Player, room *Room, msg interface{}) {
	for _, p := range room.Players {
		if player.Conn != p.Conn {
			err := p.Conn.WriteJSON(msg)
			if err != nil {
				logrus.Error("Failed to send message:", err)
			}
		}
	}
}

func broadcastAllMessage(room *Room, msg interface{}) {
	for _, player := range room.Players {
		err := player.Conn.WriteJSON(msg)
		if err != nil {
			logrus.Error("Failed to send message:", err)
		}
	}
}
