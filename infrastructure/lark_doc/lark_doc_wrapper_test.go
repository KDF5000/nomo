package lark_doc

import "testing"

func TestInsertBlock(t *testing.T) {
	wrapper := &LarkDocWrapper{}
	err := wrapper.InsertBlock("xxxx", "xxxx",
		"xxx", "test")
	if err != nil {
		t.Fatal(err)
	}
}
