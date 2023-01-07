package application

import (
	"fmt"
	"runtime"

	"github.com/KDF5000/pkg/log"
	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

type messageHandler func(*openwechat.Message) error

type wXBotHandleApp struct {
	msgHandlers map[HandlerType]messageHandler
}

func NewWXBotHandleApp() *wXBotHandleApp {
	app := &wXBotHandleApp{
		msgHandlers: make(map[HandlerType]messageHandler),
	}

	app.registMessageHandler()
	return app
}

func (app *wXBotHandleApp) registMessageHandler() {
	app.msgHandlers[UserHandler] = app.handlePrivateUserMessage
	app.msgHandlers[GroupHandler] = app.handleGroupMessage
}

// QrCodeCallBack 登录扫码回调，
func (app *wXBotHandleApp) QrCodeCallBack(uuid string) {
	if runtime.GOOS == "windows" {
		// 运行在Windows系统上
		openwechat.PrintlnQrcodeUrl(uuid)
	} else {
		log.Info("login in linux")
		q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
		fmt.Println(q.ToString(true))
	}
}

func (app *wXBotHandleApp) handlePrivateUserMessage(msg *openwechat.Message) error {
	if !msg.IsText() {
		return nil
	}

	log.Infof("msg: %+v", msg)
	_, err := msg.ReplyText(msg.Content)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return err
}

func (app *wXBotHandleApp) handleGroupMessage(msg *openwechat.Message) error {
	return nil
}

// Handler 全局处理入口
func (app *wXBotHandleApp) Handler(msg *openwechat.Message) {
	// 处理群消息
	if msg.IsSendByGroup() {
		app.msgHandlers[GroupHandler](msg)
		return
	}

	// 好友申请
	if msg.IsFriendAdd() {
		_, err := msg.Agree("hi, I am your personal memo recorder")
		if err != nil {
			log.Errorf("failed to agree friend add request, %v", err)
			return
		}
	}

	app.msgHandlers[UserHandler](msg)
}
