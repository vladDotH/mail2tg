package mailbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"mail2telegram/state"
	"strings"
)

func (bot *Bot) Get(upd tgbotapi.Update) {
	name := upd.Message.CommandArguments()
	strb := strings.Builder{}

	if len(name) > 0 {
		rule, exists := bot.State.Rules.Load(name)

		if !exists {
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Правило не существует, напишите /get для поиска правил")
			bot.Send(msg)
			return
		}

		_, err := strb.WriteString(
			fmt.Sprintf("%v: %v -> %v delay=%v originalChatId=%v\n",
				rule.Settings.Name, rule.Settings.Box, rule.Settings.ChatId,
				rule.Settings.Delay, rule.Settings.OriginalMsgId,
			),
		)
		if err != nil {
			log.Printf("Cannot write to builder: %v", err)
		}

	} else {
		strb.WriteString(fmt.Sprintf("Настройки imap: %v\n", bot.State.DefaultImap))
		strb.WriteString("Установленные правила:\n")
		bot.State.Rules.Range(func(k string, v *state.RuleState) bool {
			_, err := strb.WriteString(
				fmt.Sprintf("%v: %v -> %v delay=%v originalChatId=%v\n",
					v.Settings.Name, v.Settings.Box, v.Settings.ChatId,
					v.Settings.Delay, v.Settings.OriginalMsgId,
				),
			)
			if err != nil {
				log.Printf("Cannot write to builder: %v", err)
			}
			return true
		})
	}

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, strb.String())
	bot.Send(msg)
}
