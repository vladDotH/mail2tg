package mailbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) Del(upd tgbotapi.Update) {
	name := upd.Message.CommandArguments()
	rule, exists := bot.State.Rules.Load(name)

	if !exists {
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Правило не существует, напишите /get для поиска правил")
		bot.Send(msg)
		return
	}

	rule.Cancel()
	bot.State.Rules.Delete(name)

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Правило удалено")
	bot.Send(msg)
}
