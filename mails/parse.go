package mails

import (
	"fmt"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
	"io"
	"log"
	"mail2telegram/state"
)

func ParseMessageToParts(msg *imapclient.FetchMessageData) (state.MessageParts, error) {
	var bodySection imapclient.FetchItemDataBodySection
	ok := false
	for {
		item := msg.Next()
		if item == nil {
			break
		}
		bodySection, ok = item.(imapclient.FetchItemDataBodySection)
		if ok {
			break
		}
	}

	if !ok {
		return state.MessageParts{}, fmt.Errorf("failed to parse message body")
	}

	mr, err := mail.CreateReader(bodySection.Literal)
	if err != nil {
		return state.MessageParts{}, err
	}

	parsed := state.MessageParts{}
	parsed.Header = mr.Header
	parsed.Parts = make([]*state.MailPart, 0, 16)

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return state.MessageParts{}, err
		}
		body, err := io.ReadAll(part.Body)

		if err != nil {
			log.Printf("failed to read body: %v", err)
			continue
		}

		parsed.Parts = append(parsed.Parts, &state.MailPart{
			Header: part.Header,
			Body:   body,
		})

	}
	return parsed, nil
}

func ParseMessageParts(msg state.MessageParts) (state.ParsedMessage, error) {
	parsed := state.ParsedMessage{
		Header:      msg.Header,
		Inlines:     make([]state.InlinePart, 0, 2),
		Attachments: make([]state.AttachmentPart, 0, 8),
	}

	for _, part := range msg.Parts {
		switch h := part.Header.(type) {
		case *mail.InlineHeader:
			parsed.Inlines = append(parsed.Inlines, state.InlinePart{Header: h, Body: part.Body})
		case *mail.AttachmentHeader:
			parsed.Attachments = append(parsed.Attachments, state.AttachmentPart{Header: h, Body: part.Body})
		default:
			return state.ParsedMessage{}, fmt.Errorf("failed to parse part headers")
		}
	}

	return parsed, nil
}
