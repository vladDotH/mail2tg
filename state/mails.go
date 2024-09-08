package state

import (
	"github.com/emersion/go-message/mail"
)

type ImapParams struct {
	Server string
	User   string
	Token  string
}

type MailPart struct {
	Header mail.PartHeader
	Body   []byte
}

type MessageParts struct {
	Header mail.Header
	Parts  []*MailPart
}

type InlinePart struct {
	Header *mail.InlineHeader
	Body   []byte
}

type AttachmentPart struct {
	Header *mail.AttachmentHeader
	Body   []byte
}

type ParsedMessage struct {
	Header      mail.Header
	Inlines     []InlinePart
	Attachments []AttachmentPart
}
