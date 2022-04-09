package entity

import "gorm.io/gorm"

type Poster struct {
	gorm.Model

	TplID uint   `json:"tpl_id" gorm:"column:tpl_id"`
	Data  string `json:"data" grom:"column:data" comment:"json string"`
}
