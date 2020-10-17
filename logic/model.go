package logic

import (
	"lin/framework/session"
	"lin/framework/uuid"
	"lin/msg_dispatcher"
)


// 聊天室事件定义
type Event struct {
	Type msg_dispatcher.EventType
	Uid  uuid.UID
	Conn session.IConn
	Body string
}

