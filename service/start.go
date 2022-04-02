package service

import (
	"chat/conf"
	"chat/pkg/e"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func (manager *ClientManager) Start() {
	for {
		fmt.Println("---监听管道通信----")
		select {

		case conn := <-manager.Register:
			fmt.Printf("有新链接:%v", conn.ID)
			Manager.Clients[conn.ID] = conn
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuces,
				Content: "登录聊天成功",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)

		case conn := <-Manager.Unregister:
			fmt.Printf("链接失败%s", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "链接终端",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(manager.Clients, conn.ID)
			}

		case broadcast := <-manager.Broadcast:
			message := broadcast.Message
			sendId := broadcast.Client.SendId
			flag := false //默认对方是不在线的

			for id, conn := range manager.Clients {
				if id != sendId {
					continue
				}
				select {

				case conn.Send <- message:
					flag = true

				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}
			id := broadcast.Client.ID

			if flag {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "对方在线,请等待对方回复",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month))
				if err != nil {
					fmt.Println("InsetOne err", err)
				}
			} else {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "对方不在线,可能无法及时回复您的消息",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month))
				if err != nil {
					fmt.Println("InsetOne err", err)
				}
			}
		}
	}
}
