package application

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/c4pt0r/log"
	"github.com/patrickmn/go-cache"

	"github.com/KDF5000/nomo/domain/entity"
	"github.com/KDF5000/nomo/domain/repository"
	"github.com/KDF5000/nomo/infrastructure/message/wx_message"
)

type WXMessageHandleApp struct {
	token string

	bind           repository.BindInfoRepository
	messageHandler *messageHandler
	// case message for deduplication
	eventCache *cache.Cache
}

func NewWXMessageHandleApp(token string, bind repository.BindInfoRepository, registar repository.LarkBotRegistarRepository) *WXMessageHandleApp {
	app := &WXMessageHandleApp{
		token:          token,
		messageHandler: NewMessageHandler(bind, registar),
		eventCache:     cache.New(3*time.Minute, 10*time.Minute),
		bind:           bind,
	}

	return app
}

func (app *WXMessageHandleApp) ProcessMessage(ctx context.Context, message *wx_message.WxMessage) (string, error) {
	if message.MsgType == "event" {
		if message.Event == "subscribe" {
			return MessageWechatWelcome, nil
		}

		return ErrGroupMessageNotSupport, nil
	}

	if message.MsgType != "text" {
		log.Warnf("message type %s not supported", message.MsgType)
		return ErrMessageTypeNotSupport, nil
	}

	content := message.Content
	cmd, isBind, err := app.messageHandler.ParseBindCommand(content)
	if isBind {
		if err != nil {
			return "", err
		}

		userInfo := entity.WXUserInfo{
			UserName: message.FromUserName,
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
			return "", fmt.Errorf("unknown platform %d", cmd.Platform)
		}

		if err != nil {
			return "", err
		}
		return MessageBindSucc, nil
	}

	userInfo := entity.WXUserInfo{
		UserName: message.FromUserName,
	}
	bindInfo, err := app.bind.GetBindInfoByUnionUserID(ctx, userInfo.UnionID())
	if err != nil {
		log.Error(err)
		return MessageWechatWelcome, nil
	}

	switch entity.BindPlatformType(bindInfo.BindPlatform) {
	case entity.BindPlatformTypeNotion:
		var pageInfo entity.NotionPageInfo
		if err := json.Unmarshal([]byte(bindInfo.PageInfo), &pageInfo); err != nil {
			log.Errorf("unmarshal bind page info. info: %s, err: %v", bindInfo.PageInfo, err)
			return ErrInvalidBindPageInfo, nil
		}
		err = app.messageHandler.AppendNotionPage(ctx, &pageInfo, content)
	case entity.BindPlatformTypeLarkDoc:
		var pageInfo entity.LarkDocPageInfo
		if err := json.Unmarshal([]byte(bindInfo.PageInfo), &pageInfo); err != nil {
			log.Errorf("unmarshal bind page info. info: %s, err: %v", bindInfo.PageInfo, err)
			return ErrInvalidBindPageInfo, nil
		}
		err = app.messageHandler.AppendLarkDoc(ctx, &pageInfo, content)
	default:
		return "", fmt.Errorf("unknown bind platform")
	}

	if err != nil {
		return "", fmt.Errorf("append notion error, %v", err)
	}

	return MessageNotionSaveSucc, nil
}

func (app *WXMessageHandleApp) VerifyURL(ctx context.Context, message interface{}) (interface{}, error) {
	msg, ok := message.(*wx_message.WechatVerifyParam)
	if !ok {
		return nil, fmt.Errorf("invalid message type")
	}

	sl := []string{app.token, msg.Timestamp, msg.Echostr}
	sort.Strings(sl)
	sum := sha1.Sum([]byte(sl[0] + sl[1] + sl[2]))
	log.Infof("token: %s, signature: %s, echostr: %s, sum: %s",
		app.token, msg.Signature, msg.Echostr, hex.EncodeToString(sum[:]))
	if msg.Signature != hex.EncodeToString(sum[:]) {
		// return nil, fmt.Errorf("invalid signature")
		log.Info("invalid signature")
	}

	return msg.Echostr, nil
}
