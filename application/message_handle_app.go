package application

import (
	"context"
	"fmt"

	"github.com/KDF5000/nomo/infrastructure/lark_doc"
	"github.com/KDF5000/nomo/infrastructure/notion"
	"github.com/KDF5000/nomo/interfaces/proto"
)

type IMessageHandleApp interface {
	ProcessMessage(ctx context.Context, message *proto.Message) error
}

type messageHandleApp struct {
	notionCli        *notion.NotionClient
	larkDocV2Wrapper *lark_doc.LarkDocV2Wrapper
}

func NewMessageHandleApp() IMessageHandleApp {
	return &messageHandleApp{
		notionCli:        &notion.NotionClient{},
		larkDocV2Wrapper: &lark_doc.LarkDocV2Wrapper{},
	}
}

func (m *messageHandleApp) appendLarkDocx(ctx context.Context, message *proto.Message) error {
	var err error
	info := message.DocxInfo
	switch message.Theme {
	case "flat":
		err = m.larkDocV2Wrapper.InsertBlock(ctx, info.AppID, info.SecretKey, info.DocToken, message.Content)
	default:
		err = fmt.Errorf("theme `%s` not supported for lark docx", message.Theme)
	}

	return err
}

func (m *messageHandleApp) appendNotionPage(ctx context.Context, message *proto.Message) error {
	var err error
	info := message.NotionInfo
	switch message.Theme {
	case "flat":
		err = m.notionCli.AppendBlock(info.SecretKey, info.PageID, message.Content)
	case "gallery":
		err = m.notionCli.AddNewPage2Database(info.SecretKey, info.PageID, message.Content)
	default:
		err = fmt.Errorf("theme `%s` not supported for notion page", message.Theme)
	}

	return err
}

func (m *messageHandleApp) ProcessMessage(ctx context.Context, message *proto.Message) error {
	var err error
	switch message.Type {
	case "docx":
		err = m.appendLarkDocx(ctx, message)
	case "notion":
		err = m.appendNotionPage(ctx, message)
	default:
		err = fmt.Errorf("type `%s` not supported", message.Type)

	}
	return err
}
