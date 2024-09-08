package mailbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *Bot) Help(upd tgbotapi.Update) {
	msg := tgbotapi.NewMessage(
		upd.Message.Chat.ID,
		"Бот для пересылки почты в телеграм.\n"+
			"Пример правила для пересылки из INBOX в чат с id 1234 в ответ на сообщение с id 2 (полезно для топиков): \n"+
			"`/set {\"name\": \"default\", \"box\": \"INBOX\", \"chatId\": 1234, \"originalMsgId\": 2}`\n"+
			"Для установки дефолтного imap клиента: \n"+
			"`/set {\"imapServer\": \"imap.yandex.ru\", \"imapUser\": \"vladDotH\", \"imapToken\": \"123456789\"}`",
	)
	bot.Send(msg)
}
