package application

import (
	"context"
	"os"
	"testing"

	"github.com/KDF5000/nomo/interfaces/proto"
)

func TestNotionPage(t *testing.T) {
	message := &proto.Message{
		Type:  "notion",
		Theme: "flat",
		NotionInfo: &proto.NotionInfo{
			SecretKey: os.Getenv("notion_secret_key"),
			PageID:    os.Getenv("notion_page_id"),
		},
		Content: "#测试 这是一条测试memon",
	}

	app := NewMessageHandleApp()
	err := app.ProcessMessage(context.Background(), message)
	if err != nil {
		t.Fatal(err)
	}

	message.Theme = "gallery"
	message.NotionInfo.PageID = os.Getenv("notion_database_id")
	err = app.ProcessMessage(context.Background(), message)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLarkDocx(t *testing.T) {
	// AppID     = "cli_a10d258394b8500b"
	// AppSecret = "FRKdRyJnOIAHShoImOI1odEcw4m1Di40"
	// TestPage  = "CZJidAkBUonzOTx8cGtcyB6MnOd"
	message := &proto.Message{
		Type:  "docx",
		Theme: "flat",
		DocxInfo: &proto.LarkDocxInfo{
			AppID: os.Getenv("lark_app_id"),
			SecretKey: os.Getenv("lark_app_secret"),
			DocToken:    os.Getenv("lark_page_id"),
		},
		Content: "#测试 这是一条测试memon",
	}
	
	app := NewMessageHandleApp()
	err := app.ProcessMessage(context.Background(), message)
	if err != nil {
		t.Fatal(err)
	}

	message.Theme = "gallery"
	err = app.ProcessMessage(context.Background(), message)
	if err == nil {
		t.Fatal("must be error")
	}
	t.Logf("error: %v", err)
}
