package state

import "github.com/puzpuzpuz/xsync/v3"

type BotMailPack struct {
	Rule *RuleState
	Msg  *ParsedMessage
}

type BotState struct {
	Rules       xsync.MapOf[string, *RuleState]
	DefaultImap ImapParams
	MailsChan   chan *BotMailPack
}
