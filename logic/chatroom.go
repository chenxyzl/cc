package logic

import (
	"container/list"
	"lin/framework/uuid"
	"lin/msg_dispatcher"
)

// 保存历史消息的条数
const msgSize = 50
const chanSize = msgSize

const userJoin = "[join room]"
const userLeave = "[leave room]"

// 聊天室
type Room struct {
	users   map[uuid.UID]*User //
	msgList *list.List
}

func NewRoom() *Room {
	r := &Room{
		users:   make(map[uuid.UID]*User),
		msgList: list.New(),
	}
	return r
}

func (r *Room) JoinRoom(u *User) {
	//发送历史消息
	data := r.msgList.Front()
	for {
		if data != nil {
			u.Send(data.Value)
			data = data.Next()
		} else {
			break
		}

	}
	for _, o := range r.users {
		o.Send(msg_dispatcher.ChatRsp{u.name, userJoin})
	}
	r.users[u.uid] = u
}

func (r *Room) LeaveRoom(u *User) {
	delete(r.users, u.uid)
	for _, u := range r.users {
		u.Send(msg_dispatcher.ChatRsp{u.name, userLeave})
	}
}

func (r *Room) Say(u *User, msg string) {
	//过滤不在本房间的消息。 //todo 正常应该走监听自己感兴趣的消息类型，频道，这里时间来不及了
	if _, ok := r.users[u.uid]; !ok {
		return
	}

	// 推送消息后，限制本地只保存指定条历史消息
	if r.msgList.Len() >= msgSize {
		r.msgList.Remove(r.msgList.Front())
	}
	message := msg_dispatcher.ChatRsp{u.name, msg}
	r.msgList.PushBack(message)

	//广播
	for _, u := range r.users {
		u.Send(message)
	}
}
