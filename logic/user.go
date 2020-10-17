package logic

import (
	"encoding/json"
	"fmt"
	"lin/analysis"
	"lin/framework/session"
	"lin/framework/uuid"
	"lin/msg_dispatcher"
)

type User struct {
	conn      session.IConn
	uid       uuid.UID
	name      string
	loginTIme int64
}

func (u *User) Send(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = u.conn.Write(data)
	return err
}

func (u *User) chat(msg string) {
	if res, ok := dealIfCommand(msg, u.uid); ok {
		u.Send(msg_dispatcher.ChatRsp{
			User: u.name,
			Text: res,
		})
		return
	}
	//replace bad world
	analysis.AddWorldFrequency(msg)
	msg = analysis.ReplaceBadWord(msg, '*')

	WorldRoom.Say(u, msg)
}
