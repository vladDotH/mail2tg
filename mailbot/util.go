package mailbot

import "mail2telegram/state"

func extractRuleSettings(settings state.RuleSettings) RuleSettingsData {
	return RuleSettingsData{
		Name:          settings.Name,
		Box:           settings.Box,
		ChatId:        settings.ChatId,
		OriginalMsgId: settings.OriginalMsgId,
		Delay:         settings.Delay,
	}
}
