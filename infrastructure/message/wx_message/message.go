package wx_message

import "encoding/xml"

// <xml>
//
//	<ToUserName><![CDATA[toUser]]></ToUserName>
//	<FromUserName><![CDATA[fromUser]]></FromUserName>
//	<CreateTime>1348831860</CreateTime>
//	<MsgType><![CDATA[text]]></MsgType>
//	<Content><![CDATA[this is a test]]></Content>
//	<MsgId>1234567890123456</MsgId>
//
// </xml>
type WxMessage struct {
	ToUserName   string `json:"ToUserName" xml:"ToUserName"`
	FromUserName string `json:"FromUserName" xml:"FromUserName"`
	CreateTime   uint64 `json:"CreateTime" xml:"CreateTime"`
	MsgType      string `json:"MsgType" xml:"MsgType" comment:"image,text,voice,video,shortvideo,location"`
	Event        string `json:"Event" xml:"Event"`
	PicUrl       string `json:"PicUrl" xml:"PicUrl"`
	MediaId      string `json:"MediaId" xml:"MediaId"`
	ThumbMediaId string `json:"ThumbMediaId" xml:"ThumbMediaId"`
	Format       string `json:"Format" xml:"Format" comment:"amr，speex"`
	Recognition  string `json:"Recognition" xml:"Recognition"`

	// location
	Location  float64 `json:"Location" xml:"Location"`
	LocationX float64 `json:"Location_X" xml:"Location_X"`
	LocationY float64 `json:"Location_Y" xml:"Location_Y"`
	Scale     int     `json:"Scale" xml:"Scale"`
	Label     string  `json:"Lable" xml:"Lable"`

	Content string `json:"Content" xml:"Content"`
	MsgId   string `json:"MsgId" xml:"MsgId"`
}

type WxMessageReply struct {
	XMLName      xml.Name `json:"-" xml:"xml"`
	ToUserName   string   `json:"ToUserName" xml:"ToUserName"`
	FromUserName string   `json:"FromUserName" xml:"FromUserName"`
	CreateTime   uint64   `json:"CreateTime" xml:"CreateTime"`
	MsgType      string   `json:"MsgType" xml:"MsgType" comment:"image,text,voice,video,shortvideo,location"`
	Content      string   `json:"Content" xml:"Content"`
}

type WechatVerifyParam struct {
	Signature string `json:"signature" form:"signature"`
	Timestamp string `json:"timestamp" form:"timestamp"`
	Nonce     string `json:"nonce" form:"nonce"`
	Echostr   string `json:"echostr" form:"echostr"`
}
