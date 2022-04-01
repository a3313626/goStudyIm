package e

type Code int

const (
	WebsocketSucessMessage = 50001
	WebsocketSuces         = 50002
	WebsocketEnd           = 50003
	WebsocketOnlineReply   = 50004
	WebsocketOfflineReply  = 50005
	WebsocketLimit         = 50006
)

func (c Code) Msg() string {
	return codeMsg[c]
}
