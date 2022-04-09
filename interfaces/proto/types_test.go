package proto

import (
	"encoding/json"
	"testing"
)

func TestTime(t *testing.T) {
	text := `{
		"user_name": "kdf",
		"created_at": "2022-04-08 15:00:01",
		"content": "hello world"
	}`

	var req PosterRequest
	err := json.Unmarshal([]byte(text), &req)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("req: %+v", req)
}
