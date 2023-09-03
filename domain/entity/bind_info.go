package entity

import (
	"fmt"

	"gorm.io/gorm"
)

type UserPlatformType uint8
type LarkDocPageVersion uint8

const (
	UserPlatformTypeLark UserPlatformType = iota + 1
	UserPlatformTypeWx
)

const (
	LarkDocPageV1 LarkDocPageVersion = iota + 1
	LarkDocPageV2
)

type BindPlatformType uint8

const (
	BindPlatformTypeNotion BindPlatformType = iota + 1
	BindPlatformTypeLarkDoc
)

type BindInfo struct {
	gorm.Model

	UserPlatform uint8  `json:"user_platform" gorm:"column:user_platform" comment:"0: lark, 1: wechat"`
	UnionUserID  string `json:"union_user_id" gorm:"column:union_user_id; size:255; uniqueIndex;not null"`
	UserInfo     string `json:"user_info" gorm:"column:user_info" comment:"json fromat user info for specified platform"`
	BindPlatform uint8  `json:"bind_platform" gorm:"column:bind_platform" comment:"0: notion, 1: larkdoc"`
	PageInfo     string `json:"page_info" gorm:"column:page_info" comment:"json string for page info"`
}

func (b *BindInfo) BeforeSave(db *gorm.DB) error {
	return nil
}

type LarkUserInfo struct {
	UserId  string `json:"user_id"`
	UnionId string `json:"union_id"`
	OpenId  string `json:"open_id"`
}

func (u *LarkUserInfo) UnionID() string {
	return fmt.Sprintf("lark_%s", u.UnionId)
}

type WXUserInfo struct {
	UserName string `json:"user_name"`
}

func (u *WXUserInfo) UnionID() string {
	return fmt.Sprintf("wx_%s", u.UserName)
}

type NotionPageInfo struct {
	// flat, gallery
	// flat as default
	NotionTheme     string `json:"notion_theme"`
	NotionSecretKey string `json:"notion_secret_key"`
	NotionPageID    string `json:"notion_page_id"`
}

type LarkDocPageInfo struct {
	/* Old struct do not have this member */
	Version   LarkDocPageVersion `json:"version"`
	DocTheme  string             `json:"doc_theme"`
	DocToken  string             `json:"doc_token"`
	AppID     string             `json:"app_id"`
	SecretKey string             `json:"secret_key"`
}
