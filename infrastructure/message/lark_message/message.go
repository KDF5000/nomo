package lark_message

import (
	"encoding/json"
	"fmt"
	"strings"
)

type EventHeader struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	CreateTime string `json:"create_time"`
	Token      string `json:"token"`
	AppID      string `json:"app_id"`
	TenantKey  string `json:"tenant_key"`
}

type UserID struct {
	UnionID string `json:"union_id"`
	UserID  string `json:"user_id"`
	OpenID  string `json:"open_id"`
}

type EventSender struct {
	SenderID   UserID `json:"sender_id"`
	SenderType string `json:"sender_type"`
	TenantKey  string `json:"tenant_key"`
}

type MentionEvent struct {
	Key       string `json:"key"`
	ID        UserID `json:"id"`
	Name      string `json:"name"`
	TenantKey string `json:"tenant_key"`
}

type Message struct {
	MessageID   string         `json:"message_id"`
	RootID      string         `json:"root_id"`
	ParentID    string         `json:"parent_id"`
	CreatedTime string         `json:"create_time"`
	ChatID      string         `json:"chat_id"`
	ChatType    string         `json:"chat_type"`
	MessageType string         `json:"message_type"`
	Content     string         `json:"content"`
	Mentions    []MentionEvent `json:"mentions"`
}

type Event struct {
	Sender  EventSender `json:"sender"`
	Message Message     `json:"message"`
}

type TextMessage struct {
	Text string `json:"text"`
}

// normal message
type LarkMessageEvent struct {
	Schema string      `json:"schema"`
	Header EventHeader `json:"header"`
	Event  Event       `json:"event"`
}

// robot verification
type UrlVerificationEvent struct {
	Type      string `json:"url_verification"`
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
}

type UrlVerificationResult struct {
	Challenge string `json:"challenge"`
}

func (msg *Message) GetMessageRawContent() (string, error) {
	var text TextMessage
	if err := json.Unmarshal([]byte(msg.Content), &text); err != nil {
		return "", fmt.Errorf("parse content error, content=%s, err=%v", msg.Content, err)
	}

	// trim left @user
	content := text.Text
	for _, mention := range msg.Mentions {
		content = strings.TrimLeft(content, mention.Key)
	}

	return content, nil
}
