package state

import (
	"context"

	"github.com/emersion/go-imap/v2"
)

type RuleSettings struct {
	Imap          *ImapParams
	Name          string
	Delay         int
	Box           string
	ChatId        int64
	OriginalMsgId int
}

type RuleState struct {
	Settings RuleSettings
	Ctx      context.Context
	Cancel   context.CancelFunc
	MailChan chan *ParsedMessage
	UIDNext  imap.UID
}
