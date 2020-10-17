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
	uid  UUID
	name string
	sub  chan Subscription
}

type User struct {
	uid       UUID
	name      string
	event     chan Event
	loginTIme int64
}

// 聊天室
type Room struct {
	users      map[UUID]User     // 当前房间订阅者
	userCount  int               // 当前房间总人数
	publishChn chan Event        // 聊天室的消息推送入口
	msgList    *list.List        // 历史记录 todo 未持久化 重启失效
	msgChan    chan chan []Event // 通过接受chan来同步聊天内容
	joinChan   chan JoinEvent    // 接收订阅事件的通道 用户加入聊天室后要把历史事件推送给用户
	leaveChn   chan UUID         // 用户取消订阅通道 把通道中的历史事件释放并把用户从聊天室用户列表中删除
}

func NewRoom() *Room {
	r := &Room{
		users:      make(map[UUID]User),
		publishChn: make(chan Event, chanSize),
		msgChan:    make(chan chan []Event, chanSize),
		msgList:    list.New(),

		joinChan: make(chan JoinEvent, chanSize),
		leaveChn: make(chan UUID, chanSize),
	}

	go r.Serve()

	return r
}

// 用来向聊天室发送用户消息
// 这些接口供非websocket连接方式调用
func (r *Room) MsgJoin(user string) {
	r.publishChn <- NewEvent(EventTypeJoin, user, userJoin)
}

func (r *Room) MsgSay(user, message string) {
	r.publishChn <- NewEvent(EventTypeMsg, user, message)
}

func (r *Room) MsgLeave(user string) {
	r.publishChn <- NewEvent(EventTypeLeave, user, userLeave)
}

func (r *Room) Remove(id UUID) {
	r.leaveChn <- id // 将用户从聊天室列表中移除
}

// 用户订阅聊天室入口函数
// 返回用户订阅的对象，用户根据对象中的属性读取历史消息和即时消息
func (r *Room) Join(username string, uid UUID) Subscription {
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
			ev := NewEvent(EventTypeSystem, "", "")
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
			ev := NewEvent(EventTypeSystem, u.name, userLeave)
			for _, v := range r.users {
				v.event <- ev
			}
		}
	}
}
