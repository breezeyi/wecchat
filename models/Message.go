package models

type Message struct {
	ToUserName   string `xml:"ToUserName"`   // 开发者微信号
	FromUserName string `xml:"FromUserName"` // 发送方帐号（一个OpenID）
	CreateTime   int64  `xml:"CreateTime"`   // 消息创建时间 （整型）
	MsgType      string `xml:"MsgType"`      // text
	Content      string `xml:"Content"`      // 文本消息内容
	MsgId        int64  `xml:"MsgId"`        // 消息id，64位整型
}
