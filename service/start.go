package service

import (
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
		}
	}
}
