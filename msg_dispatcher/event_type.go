package msg_dispatcher

//todo 正常这个id是客户端传过来，这里就两个消息，简写未连接曾转过来了
type EventType int32

const (
	KEventTypeAddUser    = 1
	KEventTypeRemoveUser = 2
	KEventTypeChat       = 3
)
