package chat

import (
	"container/list"
	"lin/analysis"
	"time"
)

// 保存历史消息的条数
const msgSize = 50
const chanSize = msgSize

const userJoin = "[join room]"
const userLeave = "[leave room]"

type JoinEvent struct {
	uid  UID
	name string
	sub  chan Subscription
}

type User struct {
	uid       UID
	name      string
	event     chan Event
	loginTIme int64
}

// 聊天室
type Room struct {
	users      map[UID]User
	userCount  int
	publishChn chan Event
	msgList    *list.List
	msgChan    chan chan []Event
	joinChan   chan JoinEvent
	leaveChn   chan UID
}

func NewRoom() *Room {
	r := &Room{
		users:      make(map[UID]User),
		publishChn: make(chan Event, chanSize),
		msgChan:    make(chan chan []Event, chanSize),
		msgList:    list.New(),

		joinChan: make(chan JoinEvent, chanSize),
		leaveChn: make(chan UID, chanSize),
	}

	go r.Serve()

	return r
}

// 用来向聊天室发送用户消息
// 这些接口供非websocket连接方式调用
func (r *Room) MsgJoin(user string, uid UID) {
	r.publishChn <- NewEvent(EventTypeJoin, user, uid, userJoin)
}

// 用户订阅聊天室入口函数
// 返回用户订阅的对象，用户根据对象中的属性读取历史消息和即时消息
func (r *Room) Join(username string, uid UID) Subscription {
	resp := make(chan Subscription)
	r.joinChan <- JoinEvent{
		sub:  resp,
		name: username,
		uid:  uid,
	}
	s := <-resp
	s.Username = username
	return s
}

func (r *Room) GetHistoryMsg() []Event {
	ch := make(chan []Event)
	r.msgChan <- ch
	return <-ch
}

// 处理聊天室中的事件
func (r *Room) Serve() {
	for {
		select {
		// 用户加入房间
		case ch := <-r.joinChan:
			chn := make(chan Event, chanSize)
			r.users[ch.uid] = User{
				uid:       ch.uid,
				name:      ch.name,
				event:     chn,
				loginTIme: time.Now().Unix(),
			}
			ch.sub <- Subscription{
				Id:       ch.uid,
				Pipe:     chn,
				EmitCHn:  r.publishChn,
				LeaveChn: r.leaveChn,
			}
			ev := NewEvent(EventTypeSystem, ch.name, ch.uid, userJoin)
			for _, v := range r.users {
				v.event <- ev
			}
		case arch := <-r.msgChan:
			var events []Event
			//历史事件
			for e := r.msgList.Front(); e != nil; e = e.Next() {
				events = append(events, e.Value.(Event))
			}
			arch <- events
		// 有新的消息
		case event := <-r.publishChn:
			if res, ok := dealIfCommand(event.Text, event.Uid); ok {
				r.users[event.Uid].event <- NewEvent(EventTypeSystem, event.User, event.Uid, res)
			}
			//replace bad world
			analysis.AddWorldFrequency(event.Text)
			event.Text = analysis.ReplaceBadWord(event.Text, '*')
			// 推送给所有用户
			for _, v := range r.users {
				v.event <- event
			}
			// 推送消息后，限制本地只保存指定条历史消息
			if r.msgList.Len() >= msgSize {
				r.msgList.Remove(r.msgList.Front())
			}
			r.msgList.PushBack(event)
		// 用户退出房间
		case k := <-r.leaveChn:
			u, ok := r.users[k]
			if !ok {
				return
			}
			delete(r.users, k)
			ev := NewEvent(EventTypeSystem, u.name, u.uid, userLeave)
			for _, v := range r.users {
				v.event <- ev
			}
		}
	}
}
