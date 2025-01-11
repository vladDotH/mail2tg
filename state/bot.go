package state

import (
	"context"
	"sync"

	"github.com/puzpuzpuz/xsync/v3"
)

type BotMailPack struct {
	Rule *RuleState
	Msg  *ParsedMessage
}

type BotState struct {
	Rules       xsync.MapOf[string, *RuleState]
	DefaultImap ImapParams
	MailsChan   chan *BotMailPack
	Wg          *sync.WaitGroup
	Ctx         context.Context
}
