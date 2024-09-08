package mailbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mail2telegram/state"
)

type Bot struct {
	State  state.BotState
	BotApi *tgbotapi.BotAPI
}

type ImapSettingsData struct {
	ImapServer string `json:"imapServer"`
	ImapUser   string `json:"imapUser"`
	ImapToken  string `json:"imapToken"`
}

type RuleSettingsData struct {
	Name          string `json:"name"`
	Box           string `json:"box"`
	ChatId        int64  `json:"chatId"`
	OriginalMsgId int    `json:"originalMsgId"`
	Delay         int    `json:"delay"`
}
