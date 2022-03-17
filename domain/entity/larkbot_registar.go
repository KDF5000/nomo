package entity

import "gorm.io/gorm"

type LarkBotRegistar struct {
	gorm.Model

	AppID     string `json:"app_id"  gorm:"column:app_id;uniqIndex;size:255"`
	Token     string `json:"token" gorm:"column:token"`
	TenantKey string `json:"tenant_key" gorm:"column:tenant_key"`
	SecretKey string `json:"secret_key" gorm:"column:secret_key"`
}
