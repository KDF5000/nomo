package utils

import (
	"github.com/KDF5000/pkg/larkbot"
)

type LarkNotify func(msg string)

func ReplyLarkMessage(appid, secretKey, chatID, messageId, msg string) {
	bot := larkbot.NewLarkBot(larkbot.BotOption{
		AppID:     appid,
		AppSecret: secretKey,
	})

	bot.SendTextMessage(larkbot.IDTypeChatID, chatID, messageId, msg)
}
