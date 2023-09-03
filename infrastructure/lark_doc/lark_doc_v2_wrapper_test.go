package lark_doc

import (
	"context"
	"fmt"
	"testing"
	"time"
)

const (
	AppID     = "cli_a10d258394b8500b"
	AppSecret = "FRKdRyJnOIAHShoImOI1odEcw4m1Di40"
	TestPage  = "CZJidAkBUonzOTx8cGtcyB6MnOd"
)

func TestInsertBlockV2(t *testing.T) {
	wrapper := &LarkDocV2Wrapper{}
	err := wrapper.InsertBlock(context.Background(), AppID, AppSecret, TestPage, fmt.Sprintf("%s, %s", time.Now(), "测试一下"))
	if err != nil {
		t.Fatal(err)
	}
}
