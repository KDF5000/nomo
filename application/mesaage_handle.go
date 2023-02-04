package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/KDF5000/nomo/domain/entity"
	"github.com/KDF5000/nomo/domain/repository"
	"github.com/KDF5000/nomo/infrastructure/lark_doc"
	"github.com/KDF5000/nomo/infrastructure/notion"
	"github.com/KDF5000/pkg/log"
)

var (
	EnabledThemes = []string{
		"flat", // default
		"gallery",
	}

	DefaultTheme = "flat"

	HelpInfo = `
Usage:
  /register app_id secret_key                 register a lark bot
  /bind notion secret_key page_id [theme]     bind notion page
  /bind doc app_id secret_key page_id [theme] bind lark doc page
`
)

// for larkbot
type RegisterCommand struct {
	AppID     string
	SecretKey string
}

type BindCommand struct {
	Platform  entity.BindPlatformType
	SecretKey string
	PageID    string
	Theme     string

	// only used by lark
	AppID string
}

type messageHandler struct {
	bindRepo        repository.BindInfoRepository
	botRegistarRepo repository.LarkBotRegistarRepository
	notionCli       *notion.NotionClient
	larkDocWrapper  *lark_doc.LarkDocWrapper
}

func NewMessageHandler(bind repository.BindInfoRepository, registar repository.LarkBotRegistarRepository) *messageHandler {
	return &messageHandler{
		bindRepo:        bind,
		botRegistarRepo: registar,
		notionCli:       &notion.NotionClient{},
		larkDocWrapper:  &lark_doc.LarkDocWrapper{},
	}
}

func (h *messageHandler) ParseRegisterCommand(content string) (*RegisterCommand, bool, error) {
	data := strings.TrimSpace(content)
	if !strings.HasPrefix(data, "/register") {
		return nil, false, fmt.Errorf("not register command")
	}
	parts := strings.Fields(strings.TrimSpace(content))
	if len(parts) != 3 {
		return nil, true, fmt.Errorf(HelpInfo)
	}

	return &RegisterCommand{
		AppID:     parts[1],
		SecretKey: parts[2],
	}, true, nil
}

func (app *messageHandler) ParseBindCommand(content string) (*BindCommand, bool, error) {
	data := strings.TrimSpace(content)
	if !strings.HasPrefix(data, "/bind") {
		return nil, false, fmt.Errorf("not bind command")
	}

	log.Infof("data: %s", data)
	parts := strings.Fields(strings.TrimSpace(content))
	if len(parts) < 4 {
		return nil, true, fmt.Errorf(HelpInfo)
	}
	log.Infof("parts: %v", parts)

	var cmd BindCommand
	switch parts[1] {
	case "notion":
		cmd.Platform = entity.BindPlatformTypeNotion
		cmd.SecretKey = parts[2]
		cmd.PageID = parts[3]
		if len(parts) > 4 {
			if !app.isValidTheme(parts[4]) {
				return nil, true, fmt.Errorf("invalid theme, must be one of [flat, gallery]")
			}
			cmd.Theme = parts[4]
		}
	case "doc":
		if len(parts) < 5 {
			return nil, true, fmt.Errorf(HelpInfo)
		}
		cmd.Platform = entity.BindPlatformTypeLarkDoc
		cmd.AppID = parts[2]
		cmd.SecretKey = parts[3]
		cmd.PageID = parts[4]
		if len(parts) > 5 {
			if !app.isValidTheme(parts[5]) {
				return nil, true, fmt.Errorf("invalid theme, must be one of [flat, gallery]")
			}
			cmd.Theme = parts[5]
		}
	default:
		return nil, true, fmt.Errorf("invalid platform")
	}

	return &cmd, true, nil
}

// used by lark bot
func (app *messageHandler) RegisterLarkBot(ctx context.Context, tanentKey, tanantToken string, rcmd *RegisterCommand) (*entity.LarkBotRegistar, error) {
	var reg entity.LarkBotRegistar
	reg.AppID = rcmd.AppID
	reg.TenantKey = tanentKey
	reg.Token = tanantToken
	reg.SecretKey = rcmd.SecretKey

	if err := app.botRegistarRepo.UpdateOrInsert(ctx, &reg); err != nil {
		return nil, err
	}

	return &reg, nil
}

func (h *messageHandler) isValidTheme(theme string) bool {
	for i := range EnabledThemes {
		if EnabledThemes[i] == theme {
			return true
		}
	}

	return false
}

// /bind doc page_id [theme]
func (h *messageHandler) BindLarkDocPage(ctx context.Context, platform entity.UserPlatformType, unionId, userInfo string, cmd *BindCommand) (err error) {
	// user info
	var bindInfo entity.BindInfo
	bindInfo.UserPlatform = uint8(platform)
	bindInfo.UnionUserID = unionId
	bindInfo.UserInfo = userInfo

	// bind info
	bindInfo.BindPlatform = uint8(entity.BindPlatformTypeLarkDoc)
	pageInfo := entity.LarkDocPageInfo{
		DocTheme:  cmd.Theme,
		DocToken:  cmd.PageID,
		AppID:     cmd.AppID,
		SecretKey: cmd.SecretKey,
	}

	var info []byte
	if info, err = json.Marshal(&pageInfo); err != nil {
		return err
	}
	bindInfo.PageInfo = string(info)
	return h.bindRepo.UpdateOrInsert(ctx, &bindInfo)
}

// /bind notion secret_key page_id [theme]
func (app *messageHandler) BindNotionPage(ctx context.Context, platform entity.UserPlatformType, unionId, userInfo string, cmd *BindCommand) (err error) {
	// user info
	var bindInfo entity.BindInfo
	bindInfo.UserPlatform = uint8(platform)
	bindInfo.UnionUserID = unionId
	bindInfo.UserInfo = userInfo

	// bind info
	bindInfo.BindPlatform = uint8(entity.BindPlatformTypeNotion)
	pageInfo := entity.NotionPageInfo{
		NotionSecretKey: cmd.SecretKey,
		NotionPageID:    cmd.PageID,
		NotionTheme:     cmd.Theme,
	}

	var info []byte
	if info, err = json.Marshal(&pageInfo); err != nil {
		return err
	}
	bindInfo.PageInfo = string(info)
	return app.bindRepo.UpdateOrInsert(ctx, &bindInfo)
}

func (app *messageHandler) AppendLarkDoc(ctx context.Context, pageInfo *entity.LarkDocPageInfo, content string) error {
	var err error
	switch pageInfo.DocTheme {
	case "flat":
		err = app.larkDocWrapper.InsertBlock(pageInfo.AppID, pageInfo.SecretKey, pageInfo.DocToken, content)
	default:
		err = fmt.Errorf("invalid theme %s", pageInfo.DocTheme)
	}

	return err
}

func (app *messageHandler) AppendNotionPage(ctx context.Context, pageInfo *entity.NotionPageInfo, content string) error {
	var err error
	switch pageInfo.NotionTheme {
	case "flat":
		err = app.notionCli.AppendBlock(pageInfo.NotionSecretKey, pageInfo.NotionPageID, content)
	case "gallery":
		err = app.notionCli.AddNewPage2Database(pageInfo.NotionSecretKey, pageInfo.NotionPageID, content)
	default:
		err = fmt.Errorf("invalid theme %s", pageInfo.NotionTheme)
	}

	return err
}
