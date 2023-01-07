package application

import (
	"context"
	"strings"
)

var (
	EnabledThemes = []string{
		"flat", // default
		"gallery",
	}

	DefaultTheme = "flat"

	helpInfo = `
Usage:
  /register app_id secret_key               register a lark bot
  /bind notion secret_key page_id [theme]   bind notion page
  /bind doc page_id [theme]                 bind lark doc page
`
)

type MessageHandleApp interface {
	ProcessMessage(ctx context.Context, message interface{}) error
	VerifyURL(ctx context.Context, message interface{}) (interface{}, error)
}

type MesssageHandle struct {
}

func (app *MesssageHandle) isRegisterCommand(content string) bool {
	data := strings.TrimSpace(content)
	return strings.HasPrefix(data, "/register")
}

func (app *MesssageHandle) isValidTheme(theme string) bool {
	for i := range EnabledThemes {
		if EnabledThemes[i] == theme {
			return true
		}
	}

	return false
}
