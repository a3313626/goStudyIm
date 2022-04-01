package e

var codeMsg = map[Code]string{
	WebsocketSucessMessage: "解析content内容信息",
	WebsocketSuces:         "发送信息,请求历史记录操作成功",
	WebsocketEnd:           "请求历史记录,但没有更多记录了",
	WebsocketOnlineReply:   "针对回复消息在线应答成功",
	WebsocketOfflineReply:  "针对回复消息离线回答成功",
	WebsocketLimit:         "请求收到限制",
}
