package application

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/KDF5000/nomo/domain/entity"
	"github.com/KDF5000/nomo/infrastructure/lark_doc"
	"github.com/KDF5000/nomo/infrastructure/message/wx_message"
	"github.com/KDF5000/nomo/infrastructure/notion"
)

type wxMessageHandleApp struct {
	MesssageHandle

	token string

	notionCli      *notion.NotionClient
	larkDocWrapper *lark_doc.LarkDocWrapper

	// use different handle for diff theme
	handlers map[entity.BindPlatformType]appendHandler
	// case message for deduplication
	eventCache *cache.Cache
}

func NewwxMessageHandleApp(token string) *wxMessageHandleApp {
	app := &wxMessageHandleApp{
		notionCli:      &notion.NotionClient{},
		larkDocWrapper: &lark_doc.LarkDocWrapper{},
		handlers:       make(map[entity.BindPlatformType]appendHandler),
		eventCache:     cache.New(3*time.Minute, 10*time.Minute),
	}

	// register handler for diffrent theme
	// app.handlers[entity.BindPlatformTypeLarkDoc] = app.handleLarkAppend
	// app.handlers[entity.BindPlatformTypeNotion] = app.handleNotionAppend

	return app
}

func (app *wxMessageHandleApp) ProcessMessage(ctx context.Context, message interface{}) error {
	return nil
}

func (app *wxMessageHandleApp) VerifyURL(ctx context.Context, message interface{}) (interface{}, error) {
	msg, ok := message.(wx_message.WechatVerifyParam)
	if !ok {
		return nil, fmt.Errorf("invalid message type")
	}

	sl := []string{app.token, msg.Timestamp, msg.Echostr}
	sort.Strings(sl)
	sum := sha1.Sum([]byte(sl[0] + sl[1] + sl[2]))
	if msg.Signature != hex.EncodeToString(sum[:]) {
		return nil, fmt.Errorf("invalid signature")
	}

	return msg.Echostr, nil
}
