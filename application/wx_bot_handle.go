package application

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/KDF5000/nomo/domain/entity"
	"github.com/KDF5000/nomo/domain/repository"
	"github.com/KDF5000/pkg/log"
	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

type HandlerType string

type Notifer func(msg string) error

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

const (
	FriendAddReplyMessage = `Hi，欢迎使用Nomo小助手~
--
Hi, welcome to use Nomo assistant~
`
	ErrMessageTypeNotSupport  = "当前只支持文本消息，更多消息类型稍后会进行支持~"
	ErrGroupMessageNotSupport = "当前还不支持群聊，请私信~"
	ErrInvalidBindPageInfo    = "绑定信息有异常，请重新绑定Notion或者飞书Doc页面"
	ErrAppendFailed           = "保存失败, 请稍后重试"
	MessageNotionSaveSucc     = "已保存，可以前往Notion页面查看~"
	MessageBindSucc           = "绑定成功~"
	MessageNotBind            = "请先绑定Notion页面!"

	MessageWechatWelcome = `
谢谢关注43号广场~
---
这里会定期分享分布式、存储、数据库相关的技术以及论文解读。同时因为本人也是一个效率工具爱好者，所以也会根据个人需求开发一些效率工具在此分享，比如：本公众号也实现了一个类Flomo的功能，可以在聊天框以对话形式保存个人日常记录到Notion，详细使用方式见https://blog.openhex.cn/posts/35d22c04-5518-4871-9812-832af9e8d5fa
`
)

type wxBotHandleApp struct {
	messageHandler *messageHandler

	bind     repository.BindInfoRepository
	registar repository.LarkBotRegistarRepository
}

func NewWXBotHandleApp(bind repository.BindInfoRepository, registar repository.LarkBotRegistarRepository) *wxBotHandleApp {
	return &wxBotHandleApp{
		messageHandler: NewMessageHandler(bind, registar),
		bind:           bind,
		registar:       registar,
	}
}

func (app *wxBotHandleApp) QrCodeCallBack(uuid string) {
	if runtime.GOOS == "windows" {
		openwechat.PrintlnQrcodeUrl(uuid)
		return
	}

	// macos or linux
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Printf("Please scan the QR code to login your wechat account: %s", q.ToString(true))
}

func (app *wxBotHandleApp) processMessage(ctx context.Context, notify Notifer, message *openwechat.Message) error {
	sender, err := message.Sender()
	if err != nil {
		notify("微信内部信息格式错误, 请联系管理员~")
		return err
	}

	content := message.Content
	cmd, isBind, err := app.messageHandler.ParseBindCommand(content)
	if isBind {
		if err != nil {
			notify(fmt.Sprintf("错误的bind格式,  %v\n------------\n%s", err, HelpInfo))
			return err
		}

		userInfo := entity.WXUserInfo{
			UserName: sender.PYInitial,
		}
		data, _ := json.Marshal(&userInfo)
		switch cmd.Platform {
		case entity.BindPlatformTypeLarkDoc:
			err = app.messageHandler.BindLarkDocPage(ctx, entity.UserPlatformTypeWx,
				userInfo.UnionID(), string(data), cmd)
		case entity.BindPlatformTypeNotion:
			err = app.messageHandler.BindNotionPage(ctx, entity.UserPlatformTypeWx,
				userInfo.UnionID(), string(data), cmd)
		default:
			return fmt.Errorf("unknown platform %d", cmd.Platform)
		}

		if err != nil {
			notify(fmt.Sprintf("绑定账号失败，%s", err))
			return err
		}
		notify(MessageBindSucc)
		return nil
	}

	userInfo := entity.WXUserInfo{
		UserName: sender.PYInitial,
	}
	bindInfo, err := app.bind.GetBindInfoByUnionUserID(ctx, userInfo.UnionID())
	if err != nil {
		notify(MessageNotBind)
		return fmt.Errorf("%s, %s", MessageNotBind, err)
	}

	switch entity.BindPlatformType(bindInfo.BindPlatform) {
	case entity.BindPlatformTypeNotion:
		var pageInfo entity.NotionPageInfo
		if err := json.Unmarshal([]byte(bindInfo.PageInfo), &pageInfo); err != nil {
			notify(ErrInvalidBindPageInfo)
			return fmt.Errorf("unmarshal bind page info. info: %s, err: %v", bindInfo.PageInfo, err)
		}
		err = app.messageHandler.AppendNotionPage(ctx, &pageInfo, content)
	case entity.BindPlatformTypeLarkDoc:
		var pageInfo entity.LarkDocPageInfo
		if err := json.Unmarshal([]byte(bindInfo.PageInfo), &pageInfo); err != nil {
			notify(ErrInvalidBindPageInfo)
			return fmt.Errorf("unmarshal bind page info. info: %s, err: %v", bindInfo.PageInfo, err)
		}
		err = app.messageHandler.AppendLarkDoc(ctx, &pageInfo, content)
	default:
		return fmt.Errorf("unknown bind platform %d", bindInfo.BindPlatform)
	}

	if err != nil {
		notify(ErrAppendFailed)
		return err
	}

	notify(MessageNotionSaveSucc)
	return nil
}

func (app *wxBotHandleApp) handlePrivateUserMessage(message *openwechat.Message) error {
	notify := func(msg string) error {
		_, err := message.ReplyText(msg)
		if err != nil {
			log.Errorf("failed to reply wechat message, %v", err)
			return err
		}
		return nil
	}

	if !message.IsText() {
		return notify(ErrMessageTypeNotSupport)
	}

	sender, _ := message.Sender()
	log.Infof("receive message: %+v, sender: %+v", *message, *sender)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return app.processMessage(ctx, notify, message)
}

func (app *wxBotHandleApp) handleGroupMessage(msg *openwechat.Message) error {
	_, err := msg.Agree(FriendAddReplyMessage)
	if err != nil {
		log.Errorf("failed to agree friend add request, %v", err)
	}

	return err
}

// Handler 全局处理入口
func (app *wxBotHandleApp) Handler(msg *openwechat.Message) {
	// group message
	if msg.IsSendByGroup() {
		app.handleGroupMessage(msg)
		return
	}

	// 好友申请
	if msg.IsFriendAdd() {
		if _, err := msg.Agree(FriendAddReplyMessage); err != nil {
			log.Errorf("failed to agree friend add request, %v", err)
			return
		}
		return
	}

	// p2p message
	app.handlePrivateUserMessage(msg)
}
