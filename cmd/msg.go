package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"lin/msg_dispatcher"
)

func ReadMessage(ws *websocket.Conn, v *Message) error {
	var data string
	err := websocket.Message.Receive(ws, &data)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), v)
}

func WriteMessage(conn *websocket.Conn, ret msg_dispatcher.CommonRet) error {
	data, err := json.Marshal(ret)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = conn.Write(data)
	return err
}
