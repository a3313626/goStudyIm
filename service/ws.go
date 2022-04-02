package service

import (
	"chat/cache"
	"chat/conf"
	"chat/pkg/e"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const month = 60 * 60 * 24 * 30 //30天

//发送消息结构体
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

//回复消息结构体
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

//用户结构体
type Client struct {
	ID     string
	SendId string
	Socket *websocket.Conn
	Send   chan []byte
}

//广播类
type BroadCast struct {
	Client  *Client
	Message []byte
	Type    int
}

//用户管理类
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *BroadCast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

//消息转json类
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), //参与链接的用户,性能限制,这里需要设置最大连接数
	Broadcast:  make(chan *BroadCast),
	Reply:      make(chan *Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateID(uid, toUid string) string {
	//这里为了更明显的显示对话关系
	return uid + "->" + toUid
}

func Handler(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")

	conn, err := (&websocket.Upgrader{
		//处理跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil) //升级ws协议

	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	//创建一个用户示例
	client := &Client{
		ID:     CreateID(uid, toUid),
		SendId: CreateID(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}

	//用户注册到用户管理上
	Manager.Register <- client

	go client.Read()
	go client.Write()

}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c //修改状态为关闭状态
		_ = c.Socket.Close()    //关闭socket
	}()

	for {
		c.Socket.PongHandler()
		SendMsg := new(SendMsg)

		err := c.Socket.ReadJSON(&SendMsg)
		if err != nil {
			fmt.Println("数据格式不正确", err)
			Manager.Unregister <- c
			_ = c.Socket.Close() //关闭socket
			break
		}

		if SendMsg.Type == 1 { //发送消息
			r1, _ := cache.RedisClient.Get(c.ID).Result()
			r2, _ := cache.RedisClient.Get(c.SendId).Result()

			//发送消息超过3条,对方没有看到,停止1发送
			if r1 > "3" && r2 == "" {
				ReplyMsg := ReplyMsg{
					Code:    e.WebsocketLimit,
					Content: "发送信息达到限制,请等待对方回复",
				}
				msg, _ := json.Marshal(ReplyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			} else {
				cache.RedisClient.Incr(c.ID)
				_, _ = cache.RedisClient.Expire(c.ID, time.Hour*24*30*3).Result()
			}

		} else if SendMsg.Type == 2 { //获取历史消息
			
			timeT, err := strconv.Atoi(SendMsg.Content) //string To int
			if err != nil {
				timeT = 9999999
			}

			results, _ := FindMany(conf.MongoDBName, c.SendId, c.ID, int64(timeT), 10)

			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				ReplyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "没有更多消息",
				}
				msg, _ := json.Marshal(ReplyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}

			for _, result := range results {
				ReplyMsg := ReplyMsg{
					From:    result.From,
					Content: result.Msg,
				}
				msg, _ := json.Marshal(ReplyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}

		}

		Manager.Broadcast <- &BroadCast{
			Client:  c,
			Message: []byte(SendMsg.Content),
		}

	}

}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			ReplyMsg := ReplyMsg{
				Code:    e.WebsocketSucessMessage,
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(ReplyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)

		}
	}
}
