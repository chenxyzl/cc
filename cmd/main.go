package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"lin/chat"
	"lin/framework/uuid"
	"net/http"
	"time"
)

type Message struct {
	Msg string `json:"msg"`
}

func main() {
	fmt.Println("server start ", time.Now())
	http.HandleFunc("/", index)
	http.Handle("/ws", websocket.Handler(webSocket))
	http.ListenAndServe(":9001", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func ReadMessage(ws *websocket.Conn, v *Message) error {
	var data string
	err := websocket.Message.Receive(ws, &data)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), v)
}

func WriteMessage(conn *websocket.Conn, v chat.Event) error {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = conn.Write(data)
	return err
}

func webSocket(ws *websocket.Conn) {
	var uid = chat.UID(uuid.GenUUID())
	WriteMessage(ws, chat.NewEvent("system", "room manager", uid, "please input your name"))
	var firstMessage = Message{}
	err := ReadMessage(ws, &firstMessage)
	if err != nil {
		return
	}

	// get user name
	var name = firstMessage.Msg
	evs := chat.WorldRoom.GetHistoryMsg()
	chat.WorldRoom.MsgJoin(name, uid)
	control := chat.WorldRoom.Join(name, uid)
	defer control.Leave()

	// history msg
	for _, event := range evs {
		if WriteMessage(ws, event) != nil {
			// 用户断开连接
			return
		}
	}

	// read message
	inMsg := make(chan Message)
	go func() {
		for {
			in := Message{}
			err := ReadMessage(ws, &in)
			if err != nil {
				// close chan
				close(inMsg)
				return
			}
			inMsg <- in
		}
	}()

	//send message
	for {
		select {
		case event := <-control.Pipe:
			if WriteMessage(ws, event) != nil {
				// 用户断开连接
				return
			}
		case msg, ok := <-inMsg:
			// close message
			if !ok {
				return
			}
			control.Say(msg.Msg)
		}
	}
}
