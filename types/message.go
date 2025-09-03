package types

type Message struct {
	Type string `json:"type"`           //请求类型,例如测试连接ping,写入数据write... 推荐请求和响应通过code区分
	Code int    `json:"code,omitempty"` //请求结果,推荐 请求:0(或null)  响应: 200成功,500失败... 同http好记一点
	UID  string `json:"uid,omitempty"`  //消息的唯一ID,例如UUID
	Data any    `json:"data,omitempty"` //请求响应的数据
	Msg  string `json:"msg,omitempty"`  //消息
}

// IsRequest 默认code为0是,视作请求
func (this *Message) IsRequest() bool {
	return this.Code == 0
}

// IsResponse 默认code>0时视为响应
func (this *Message) IsResponse() bool {
	return this.Code != 0
}

func (this *Message) Response(code int, data any, msg string) *Message {
	return &Message{
		Type: this.Type,
		Code: code,
		UID:  this.UID,
		Data: data,
		Msg:  msg,
	}
}
