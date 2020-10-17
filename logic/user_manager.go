package logic

import (
	"lin/framework/session"
	"lin/framework/uuid"
	"lin/msg_dispatcher"
	"time"
)

var UserMgr = NewUserManger()

type UserManager struct {
	users       map[uuid.UID]*User
	nameMapping map[string]uuid.UID
	MsgChan     chan Event
	event       map[msg_dispatcher.EventType]func(*User, string)
}

func NewUserManger() *UserManager {
	r := &UserManager{
		users:       make(map[uuid.UID]*User),
		nameMapping: make(map[string]uuid.UID),
		MsgChan:     make(chan Event, chanSize),
		event:       make(map[msg_dispatcher.EventType]func(*User, string), 0),
	}

	r.event[msg_dispatcher.KEventTypeChat] = r.chat
	r.Start()
	go r.Loop()
	return r
}

//todo 类似这样的方法应该用代码生成器或者反射直接赚到user曾
func (r *UserManager) chat(u *User, msg string) {
	u.chat(msg)
}

func (r *UserManager) Start() {

}

func (r *UserManager) addUser(uid uuid.UID, name string, conn session.IConn) {
	u := &User{conn: conn, uid: uid, name: name, loginTIme: time.Now().Unix()}

	if _, ok := r.users[uid]; ok {
		u.Send(msg_dispatcher.ChatRsp{
			User: u.name,
			Text: "uid 重复",
		})
		u.conn.Close()
		return
	}

	if _, ok := r.nameMapping[name]; ok {
		u.Send(msg_dispatcher.ChatRsp{
			User: u.name,
			Text: "名字重复",
		})
		conn.Close()
		return
	}
	r.nameMapping[name] = uid
	r.users[uid] = u
	WorldRoom.JoinRoom(u)
}

func (r *UserManager) removeUser(uid uuid.UID) {
	u, ok := r.users[uid]
	if !ok {
		return
	}
	delete(r.users, uid)
	delete(r.nameMapping, u.name)
	WorldRoom.LeaveRoom(u)
}

func (r *UserManager) ProcessMessage(event Event) {
	r.MsgChan <- event
}

// 处理聊天室中的事件
func (r *UserManager) Loop() {
	for {
		select {
		case msg := <-r.MsgChan:
			switch msg.Type {
			case msg_dispatcher.KEventTypeAddUser:
				r.addUser(msg.Uid, msg.Body, msg.Conn)
			case msg_dispatcher.KEventTypeRemoveUser:
				r.removeUser(msg.Uid)
			default:
				u, ok := r.users[msg.Uid]
				if !ok {
					return
				}
				if u.conn != msg.Conn {
					u.Send(msg_dispatcher.CommonRet{
						User: u.name,
						Text: "链接被顶了",
					})
					u.conn.Close() //链接重制未新的链接
					u.conn = msg.Conn
				}
				event, ok := r.event[msg.Type]
				if !ok {
					u.Send(msg_dispatcher.CommonRet{
						User: u.name,
						Text: "时间未找到",
					})
					return
				}
				event(u, msg.Body)
			}
		}
	}
}
