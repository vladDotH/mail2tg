package mailbot

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"mail2telegram/db"
	"mail2telegram/state"
)

func (bot *Bot) Set(upd tgbotapi.Update) {
	var data map[string]json.RawMessage
	args := upd.Message.CommandArguments()
	err := json.Unmarshal([]byte(args), &data)
	if err != nil {
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Некорректный json")
		bot.Send(msg)
		return
	}

	if len(data["imapServer"]) > 0 && len(data["imapUser"]) > 0 && len(data["imapToken"]) > 0 {
		var imapData ImapSettingsData
		err := json.Unmarshal([]byte(args), &imapData)
		if err != nil {
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Некорректный json")
			bot.Send(msg)
			return
		}

		bot.setImapData(imapData)

		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Настройки imap обновлены")
		bot.Send(msg)

	} else if len(data["name"]) > 0 && len(data["box"]) > 0 && len(data["chatId"]) > 0 {
		var ruleData = RuleSettingsData{
			Delay: 10,
		}

		err := json.Unmarshal([]byte(args), &ruleData)
		if err != nil {
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Некорректный json")
			bot.Send(msg)
			return
		}

		bot.RunRule(ruleData)

		rules := make([]RuleSettingsData, 0, bot.State.Rules.Size())
		bot.State.Rules.Range(func(_ string, value *state.RuleState) bool {
			rules = append(rules, extractRuleSettings(value.Settings))
			return true
		})

		err = db.Write(RulesDataKey, rules)
		if err != nil {
			log.Printf("Cannot save rules to db: %v", err)
		}

		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Правило применено")
		bot.Send(msg)

	} else {
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Неопознанный json")
		bot.Send(msg)
		return
	}

}

func (bot *Bot) setImapData(imapData ImapSettingsData) {
	bot.State.DefaultImap.Server = imapData.ImapServer
	bot.State.DefaultImap.User = imapData.ImapUser
	bot.State.DefaultImap.Token = imapData.ImapToken
	log.Printf("New imap params: %v", imapData)

	err := db.Write(ImapDataKey, imapData)
	if err != nil {
		log.Printf("Cannot save imapData to db: %v", err)
	}
}
