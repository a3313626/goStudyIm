package service

import (
	"net/http"

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
	uid := c.Query("id")
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

}

func (c *Client) Write() {

}
