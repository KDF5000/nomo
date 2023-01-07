package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/KDF5000/pkg/log"
	"github.com/patrickmn/go-cache"

	"github.com/KDF5000/nomo/domain/entity"
	"github.com/KDF5000/nomo/domain/repository"
	"github.com/KDF5000/nomo/infrastructure/lark_doc"
	"github.com/KDF5000/nomo/infrastructure/message/lark_message"
	"github.com/KDF5000/nomo/infrastructure/notion"
	. "github.com/KDF5000/nomo/infrastructure/utils"
)

var (
	EnabledThemes = []string{
		"flat", // default
		"gallery",
	}

	DefaultTheme = "flat"

	helpInfo = `
Usage:
  /register app_id secret_key               register a lark bot
  /bind notion secret_key page_id [theme]   bind notion page
  /bind doc page_id [theme]                 bind lark doc page
`
)

type ILarkMessageHandleApp interface {
	ProcessMessage(ctx context.Context, event *lark_message.LarkMessageEvent) error
	VerifyURL(ctx context.Context, event *lark_message.UrlVerificationEvent) (*lark_message.UrlVerificationResult, error)
}

type appendHandler func(ctx context.Context, reg *entity.LarkBotRegistar, pageInfo string, content string) error

type larkMessageHandleApp struct {
	bindRepo        repository.BindInfoRepository
	botRegistarRepo repository.LarkBotRegistarRepository
	larkNotify      LarkNotify
	notionCli       *notion.NotionClient
	larkDocWrapper  *lark_doc.LarkDocWrapper

	// use different handle for diff theme
	handlers map[entity.BindPlatformType]appendHandler
	// case message for deduplication
	eventCache *cache.Cache
}

var _ ILarkMessageHandleApp = &larkMessageHandleApp{}

func NewLarkMessageHandleApp(repo repository.BindInfoRepository,
	registarRepo repository.LarkBotRegistarRepository,
	notifier LarkNotify) *larkMessageHandleApp {
	app := &larkMessageHandleApp{
		bindRepo:        repo,
		botRegistarRepo: registarRepo,
		larkNotify:      notifier,
		notionCli:       &notion.NotionClient{},
		larkDocWrapper:  &lark_doc.LarkDocWrapper{},
		handlers:        make(map[entity.BindPlatformType]appendHandler),
		eventCache:      cache.New(3*time.Minute, 10*time.Minute),
	}

	// register handler for diffrent theme
	app.handlers[entity.BindPlatformTypeLarkDoc] = app.handleLarkAppend
	app.handlers[entity.BindPlatformTypeNotion] = app.handleNotionAppend

	return app
}

func (app *larkMessageHandleApp) isRegisterCommand(content string) bool {
	data := strings.TrimSpace(content)
	return strings.HasPrefix(data, "/register")
}

func (app *larkMessageHandleApp) registerLarkBot(ctx context.Context, header *lark_message.EventHeader, content string) (*entity.LarkBotRegistar, error) {
	data := strings.Fields(strings.TrimSpace(content))
	if len(data) < 3 {
		return nil, fmt.Errorf("command should be like `/register app_id secret_key`")
	}

	if header.AppID != data[1] {
		return nil, fmt.Errorf("mismatch app id with header")
	}

	var reg entity.LarkBotRegistar
	reg.AppID = data[1]
	reg.TenantKey = header.TenantKey
	reg.Token = header.Token
	reg.SecretKey = data[2]

	if err := app.botRegistarRepo.UpdateOrInsert(ctx, &reg); err != nil {
		return nil, err
	}

	return &reg, nil
}

func (app *larkMessageHandleApp) isBindCommand(content string) bool {
	data := strings.TrimSpace(content)
	return strings.HasPrefix(data, "/bind")
}

func (app *larkMessageHandleApp) isValidTheme(theme string) bool {
	for i := range EnabledThemes {
		if EnabledThemes[i] == theme {
			return true
		}
	}

	return false
}

// /bind doc page_id [theme]
func (app *larkMessageHandleApp) bindLakrDocPage(ctx context.Context, userID *lark_message.UserID, content string) (err error) {
	data := strings.Fields(strings.TrimSpace(content))
	if len(data) < 3 {
		return fmt.Errorf(helpInfo)
	}

	theme := DefaultTheme
	if len(data) > 3 {
		if !app.isValidTheme(data[3]) {
			return fmt.Errorf("invalid theme, must be one of [flat, gallery]")
		}

		theme = data[3]
	}

	// user info
	var bindInfo entity.BindInfo
	bindInfo.UserPlatform = uint8(entity.UserPlatformTypeLark)
	user := entity.LarkUserInfo{
		UserId:  userID.UserID,
		UnionId: userID.UnionID,
		OpenId:  userID.OpenID,
	}
	bindInfo.UnionUserID = user.UnionID()
	var info []byte
	if info, err = json.Marshal(&user); err != nil {
		return
	}
	bindInfo.UserInfo = string(info)

	// bind info
	bindInfo.BindPlatform = uint8(entity.BindPlatformTypeLarkDoc)
	pageInfo := entity.LarkDocPageInfo{
		DocTheme: theme,
		DocToken: data[2],
	}
	if info, err = json.Marshal(&pageInfo); err != nil {
		return err
	}
	bindInfo.PageInfo = string(info)
	return app.bindRepo.UpdateOrInsert(ctx, &bindInfo)
}

// /bind notion secret_key page_id [theme]
func (app *larkMessageHandleApp) bindNotionPage(ctx context.Context, userID *lark_message.UserID, content string) (err error) {
	data := strings.Fields(strings.TrimSpace(content))
	if len(data) < 4 {
		return fmt.Errorf(helpInfo)
	}

	theme := DefaultTheme
	if len(data) > 4 {
		if !app.isValidTheme(data[4]) {
			return fmt.Errorf("invalid theme, must be one of [flat, gallery]")
		}

		theme = data[4]
	}

	// user info
	var bindInfo entity.BindInfo
	bindInfo.UserPlatform = uint8(entity.UserPlatformTypeLark)
	user := entity.LarkUserInfo{
		UserId:  userID.UserID,
		UnionId: userID.UnionID,
		OpenId:  userID.OpenID,
	}
	bindInfo.UnionUserID = user.UnionID()
	var info []byte
	if info, err = json.Marshal(&user); err != nil {
		return
	}
	bindInfo.UserInfo = string(info)

	// bind info
	bindInfo.BindPlatform = uint8(entity.BindPlatformTypeNotion)
	pageInfo := entity.NotionPageInfo{
		NotionSecretKey: data[2],
		NotionPageID:    data[3],
		NotionTheme:     theme,
	}
	if info, err = json.Marshal(&pageInfo); err != nil {
		return err
	}
	bindInfo.PageInfo = string(info)

	return app.bindRepo.UpdateOrInsert(ctx, &bindInfo)
}

func (app *larkMessageHandleApp) handleLarkAppend(ctx context.Context, reg *entity.LarkBotRegistar, pageInfo string, content string) error {
	var docInfo entity.LarkDocPageInfo
	if err := json.Unmarshal([]byte(pageInfo), &docInfo); err != nil {
		return err
	}

	// log.Infof("token: %s, theme: %s, content: %s", docInfo.DocToken, docInfo.DocTheme, content)

	var err error
	switch docInfo.DocTheme {
	case "flat":
		err = app.larkDocWrapper.InsertBlock(reg.AppID, reg.SecretKey, docInfo.DocToken, content)
	default:
		err = fmt.Errorf("invalid theme %s", docInfo.DocTheme)
	}

	return err
}

func (app *larkMessageHandleApp) handleNotionAppend(ctx context.Context, reg *entity.LarkBotRegistar, pageStr string, content string) error {
	var pageInfo entity.NotionPageInfo
	if err := json.Unmarshal([]byte(pageStr), &pageInfo); err != nil {
		return err
	}

	// log.Infof("key: %s, id: %s, theme: %s, content: %s",
	// 	pageInfo.NotionSecretKey, pageInfo.NotionPageID, pageInfo.NotionTheme, content)

	var err error
	switch pageInfo.NotionTheme {
	case "flat":
		err = app.notionCli.AppendBlock(pageInfo.NotionSecretKey, pageInfo.NotionPageID, content)
	case "gallery":
		err = app.notionCli.AddNewPage2Database(pageInfo.NotionSecretKey, pageInfo.NotionPageID, content)
	default:
		err = fmt.Errorf("invalid theme %s", pageInfo.NotionTheme)
	}

	return err
}

func (app *larkMessageHandleApp) appendContent(ctx context.Context, registar *entity.LarkBotRegistar, event *lark_message.Event, content string) error {
	user := entity.LarkUserInfo{
		UserId:  event.Sender.SenderID.UserID,
		UnionId: event.Sender.SenderID.UnionID,
		OpenId:  event.Sender.SenderID.OpenID,
	}

	bindInfo, err := app.bindRepo.GetBindInfoByUnionUserID(ctx, user.UnionID())
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("请先绑定Notion页面! %s", err)
	}

	handler, ok := app.handlers[entity.BindPlatformType(bindInfo.BindPlatform)]
	if !ok {
		return fmt.Errorf("invalid bind platform. platform=%d", bindInfo.BindPlatform)
	}

	return handler(ctx, registar, bindInfo.PageInfo, content)
}

func (app *larkMessageHandleApp) getBotRegistar(ctx context.Context, appId string) (*entity.LarkBotRegistar, error) {
	return app.botRegistarRepo.GetLarkBotRegistarByUnionUserID(ctx, appId)
}

func (app *larkMessageHandleApp) ProcessMessage(ctx context.Context, event *lark_message.LarkMessageEvent) error {
	if _, ok := app.eventCache.Get(event.Header.EventID); ok {
		return fmt.Errorf("repeated lark message +%v", *event)
	}
	app.eventCache.Set(event.Header.EventID, true, cache.DefaultExpiration)

	message := &event.Event.Message
	if message.MessageType != "text" {
		msg := fmt.Sprintf("unsupported message type: %s, app_id: %s  chat_id: %s, messageid: %s",
			event.Event.Message.MessageType, event.Header.AppID, message.ChatID, message.MessageID)
		app.larkNotify(msg)
		if reg, err := app.getBotRegistar(ctx, event.Header.AppID); err == nil {
			ReplyLarkMessage(reg.AppID, reg.SecretKey, message.ChatID, message.MessageID,
				fmt.Sprintf("目前只支持文本消息，当前类型为 %s", event.Event.Message.MessageType))
		}
		return fmt.Errorf("%s", msg)
	}

	content, err := message.GetMessageRawContent()
	if err != nil {
		log.Errorf("failed to get message content. err=%v", err)
		app.larkNotify(fmt.Sprintf("event: %+v, err: %v", *event, err.Error()))
		return err
	}

	// /register app_secret_key
	if app.isRegisterCommand(content) {
		reg, err := app.registerLarkBot(ctx, &event.Header, content)
		if err != nil {
			msg := fmt.Sprintf("register lark bot error. event: %+v, err: %v", event, err)
			app.larkNotify(msg)
			return err
		}

		ReplyLarkMessage(reg.AppID, reg.SecretKey, message.ChatID, message.MessageID, "注册成功!")
		return nil
	}

	reg, err := app.getBotRegistar(ctx, event.Header.AppID)
	if err != nil {
		app.larkNotify(fmt.Sprintf("Failed to get bot registar. event: %+v, err: %v", *event, err))
		return err
	}

	// /bind notion secret_key page_id [theme] or /bind doc page_id [theme]
	if app.isBindCommand(content) {
		parts := strings.Fields(strings.TrimSpace(content))
		if len(parts) < 2 {
			log.Errorf("invalid bind command. %s", content)
			ReplyLarkMessage(reg.AppID, reg.SecretKey, message.ChatID, message.MessageID, helpInfo)
			return fmt.Errorf("invalid bind command, %s", content)
		}

		switch parts[1] {
		case "notion":
			err = app.bindNotionPage(ctx, &event.Event.Sender.SenderID, content)
		case "doc":
			err = app.bindLakrDocPage(ctx, &event.Event.Sender.SenderID, content)
		default:
			log.Errorf("invalid bind command. %s", content)
			ReplyLarkMessage(reg.AppID, reg.SecretKey, message.ChatID, message.MessageID, helpInfo)
			return fmt.Errorf("invalid bind command, %s", content)
		}

		if err != nil {
			log.Errorf("failed to bind page. err=%v", err)
			ReplyLarkMessage(reg.AppID, reg.SecretKey, message.ChatID, message.MessageID, err.Error())
			return err
		}

		ReplyLarkMessage(reg.AppID, reg.SecretKey, message.ChatID, message.MessageID, "绑定成功~")
		return nil
	}

	// log.Infof("content==> %s", content)
	if err := app.appendContent(ctx, reg, &event.Event, content); err != nil {
		msg := fmt.Sprintf("向Notion页面写入失败, %v", err)
		log.Errorf(msg)
		ReplyLarkMessage(reg.AppID, reg.SecretKey, message.ChatID, message.MessageID, err.Error())
		return err
	}

	ReplyLarkMessage(reg.AppID, reg.SecretKey, message.ChatID, message.MessageID, "已保存，可以前往Notion页面查看~")
	return nil
}

func (app *larkMessageHandleApp) VerifyURL(ctx context.Context, event *lark_message.UrlVerificationEvent) (*lark_message.UrlVerificationResult, error) {
	return &lark_message.UrlVerificationResult{Challenge: event.Challenge}, nil
}
