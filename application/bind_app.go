package application

import (
	"github.com/KDF5000/nomo/domain/entity"
	"github.com/KDF5000/nomo/domain/repository"
)

type IBindInfoApp interface {
	BindWXNotion(user *entity.WXUserInfo, page *entity.NotionPageInfo) error
	BindWXLarkDoc(user *entity.WXUserInfo, doc *entity.LarkDocPageInfo) error
	BindLarkNotion(user *entity.LarkUserInfo, page *entity.NotionPageInfo) error
	BindLarkDoc(user *entity.LarkUserInfo, page *entity.NotionPageInfo) error
}

type bindInfoApp struct {
	bindRepo repository.BindInfoRepository
}

func NewBindInfoApp(repo repository.BindInfoRepository) *bindInfoApp {
	return &bindInfoApp{bindRepo: repo}
}

func (app *bindInfoApp) BindWXNotion(user *entity.WXUserInfo, page *entity.NotionPageInfo) error {
	return nil
}

func (app *bindInfoApp) BindWXLarkDoc(user *entity.WXUserInfo, doc *entity.LarkDocPageInfo) error {
	return nil
}

func (app *bindInfoApp) BindLarkNotion(user *entity.LarkUserInfo, page *entity.NotionPageInfo) error {
	return nil
}

func (app *bindInfoApp) BindLarkDoc(user *entity.LarkUserInfo, page *entity.NotionPageInfo) error {
	return nil
}
