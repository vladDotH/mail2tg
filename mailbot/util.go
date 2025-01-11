package mailbot

import (
	"mail2telegram/env"
	"mail2telegram/state"
	"os"

	"github.com/google/uuid"
)

func extractRuleSettings(settings state.RuleSettings) RuleSettingsData {
	return RuleSettingsData{
		Name:          settings.Name,
		Box:           settings.Box,
		ChatId:        settings.ChatId,
		OriginalMsgId: settings.OriginalMsgId,
		Delay:         settings.Delay,
	}
}

func UUID2URL(id uuid.UUID) string {
	return env.Env.HTTPPrefix + "/" + env.Env.StoragePrefix + "/" + id.String()
}

func UUID2Path(id uuid.UUID) string {
	return "./" + env.Env.StoragePrefix + "/" + id.String() + ".html"
}

func OpenFile(id uuid.UUID) (*os.File, error) {
	return os.Open(UUID2Path(id))
}
