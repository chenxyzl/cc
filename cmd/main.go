package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"lin/framework/uuid"
	"lin/logic"
	"lin/msg_dispatcher"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	Msg string `json:"msg"`
}

func main() {
	fmt.Println("server start ", time.Now())
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.Handle("/ws", websocket.Handler(webSocket))
	http.ListenAndServe(":13001", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func webSocket(ws *websocket.Conn) {
	var uid = uuid.GenUUID()
	WriteMessage(ws, msg_dispatcher.CommonRet{User: "system", Text: "welcome:" + strconv.Itoa(time.Now().Second()),})
	var firstMessage = Message{}
	err := ReadMessage(ws, &firstMessage)
	if err != nil {
		return
	}
	// get user name
	var name = firstMessage.Msg

	logic.UserMgr.ProcessMessage(logic.Event{
		Type: msg_dispatcher.KEventTypeAddUser,
		Uid:  uid,
		Conn: ws,
		Body: name,
	})

	defer logic.UserMgr.ProcessMessage(logic.Event{
		Type: msg_dispatcher.KEventTypeRemoveUser,
		Uid:  uid,
		Conn: ws,
		Body: name,
	})

	// read message
	for {
		in := Message{}
		err := ReadMessage(ws, &in)
		if err != nil {
			return
		}
		logic.UserMgr.ProcessMessage(logic.Event{
			Type: msg_dispatcher.KEventTypeChat,
			Uid:  uid,
			Conn: ws,
			Body: in.Msg,
		})
	}
}
