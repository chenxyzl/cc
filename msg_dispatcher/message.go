package msg_dispatcher

type CommonRet struct {
	User string `json:"user"` // 用户名
	Text string `json:"text"` // 事件内容
}

type NameReq struct { //不在user这边处理
	Name string `json:"name"`
}

type ChatReq struct {
	Channel string `json:"channel"`
	Content string `json:"content"`
}

type ChatRsp struct {
	User string `json:"user"` // 用户名
	Text string `json:"text"` // 事件内容
}
