package mailbot

import (
	"context"
	"mail2telegram/mails"
	"mail2telegram/state"
)

func (bot *Bot) RunRule(ruleData RuleSettingsData) {
	rule, exists := bot.State.Rules.Load(ruleData.Name)
	if exists {
		rule.Cancel()
		bot.State.Rules.Delete(ruleData.Name)
	}

	ctx, cancel := context.WithCancel(context.Background())
	newRule := state.RuleState{
		Settings: state.RuleSettings{
			Imap:          &bot.State.DefaultImap,
			Name:          ruleData.Name,
			Delay:         ruleData.Delay,
			Box:           ruleData.Box,
			ChatId:        ruleData.ChatId,
			OriginalMsgId: ruleData.OriginalMsgId,
		},
		Ctx:      ctx,
		Cancel:   cancel,
		MailChan: make(chan *state.ParsedMessage, 16),
		UIDNext:  0,
	}

	bot.State.Rules.Store(ruleData.Name, &newRule)

	go mails.RunMailerRule(&newRule)

	go func() {
		for msg := range newRule.MailChan {
			bot.State.MailsChan <- &state.BotMailPack{
				Rule: &newRule,
				Msg:  msg,
			}
		}
	}()
}
