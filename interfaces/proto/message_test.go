package proto

import (
	"encoding/json"
	"testing"
)

func TestMessageBasic(t *testing.T) {
	msg := Message {
		Type :"docx",
		Theme: "flat",
		DocxInfo: &LarkDocxInfo{
			AppID: "appidxxxx",
			SecretKey: "app_secret_key",
			DocToken: "doc_token",
		},
	}

	data, _ := json.Marshal(msg)
	t.Logf("%v", string(data))

	str := `{
	    "type":"docx",
	    "theme":"flat",
		"lark_docx_info": {
			"lark_app_id":"appidxxxx",
			"lark_secret_key":"app_secret_key",
			"lark_doc_token":"doc_token"
		}
	}`

	var msg2 Message
	json.Unmarshal([]byte(str), &msg2)
	t.Logf("%v", msg2)
}