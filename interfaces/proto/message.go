package proto

import "strings"

var validTypeSet = map[string]struct{}{
	"docx":   struct{}{},
	"notion": struct{}{},
}

var validTypeList = []string{"docx", "notion"}

var validThemeSet = map[string]struct{}{
	"flat":    struct{}{},
	"gallery": struct{}{},
}

var validThemeList = []string{"flat, gallery"}

type NotionInfo struct {
	SecretKey string `json:"notion_secret_key"`
	PageID    string `json:"notion_page_id"`
}

func (info *NotionInfo) IsValid() bool {
	return info.SecretKey != "" && info.PageID != ""
}

type LarkDocxInfo struct {
	AppID     string `json:"lark_app_id"`
	SecretKey string `json:"lark_secret_key"`
	DocToken  string `json:"lark_doc_token"`
}

func (info *LarkDocxInfo) IsValid() bool {
	return info.AppID != "" && info.SecretKey != "" || info.DocToken != ""
}

type Message struct {
	Type  string `json:"type" comment:"docx: larkdock, notion: notion page"`
	Theme string `json:"theme" comment:"flat, gallery, default: gallery"`
	// for lark docx
	DocxInfo *LarkDocxInfo `json:"lark_docx_info,omitempty"`
	// for notion page
	NotionInfo *NotionInfo `json:"notion_info,omitempty"`

	// Content to append
	Content string `json:"content"`
}

func IsValidType(t string) (string, bool) {
	_, ok := validTypeSet[t]
	if !ok {
		// return valid type list if type is invalid
		return strings.Join(validTypeList, ","), false
	}

	return "", ok
}

func IsValidTheme(t string) (string, bool) {
	_, ok := validThemeSet[t]
	// return valid theme list if theme is invalid
	if !ok {
		return strings.Join(validThemeList, ","), false
	}
	return "", ok
}
