package mails

import (
	"fmt"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"log"
	"mail2telegram/state"
	"time"
)

func RunMailerRule(rule *state.RuleState) {
	settings := &rule.Settings

	logger := log.New(
		log.Writer(),
		fmt.Sprintf("[Rule %v: %v -> %v] ", rule.Settings.Name, rule.Settings.Box, rule.Settings.ChatId),
		log.LstdFlags,
	)

	logger.Printf("Starting new rule... Delay = %v, OriginalMsgId = %v", rule.Settings.Delay, rule.Settings.OriginalMsgId)

ruleLoop:
	for {
		select {
		case <-rule.Ctx.Done():
			logger.Printf("Stopping rule...")
			break ruleLoop
		case <-time.After(time.Duration(settings.Delay) * time.Second):
		}

		// II RuleFunction
		func() {
			mailClient, err := imapclient.DialTLS(settings.Imap.Server, nil)

			defer func(mailClient *imapclient.Client) {
				err := mailClient.Close()
				if err != nil {
					logger.Printf("Failed to close IMAP connection: %v", err)
				}
			}(mailClient)

			if err != nil {
				logger.Printf("Failed to dial %v IMAP server: %v", settings.Imap.Server, err)
				return
			}

			if err := mailClient.Login(settings.Imap.User, settings.Imap.Token).Wait(); err != nil {
				logger.Printf("Failed to login in %v: %v", settings.Imap.Server, err)
				return
			}

			selectedMbox, err := mailClient.Select(settings.Box, nil).Wait()
			if err != nil {
				logger.Printf("Failed to select %v: %v", settings.Box, err)
				return
			}

			if selectedMbox.UIDNext <= rule.UIDNext {
				return
			}

			logger.Printf("New messages. New UIDNext=%v", selectedMbox.UIDNext)

			oldUid := rule.UIDNext
			rule.UIDNext = selectedMbox.UIDNext

			seqSet := imap.UIDSet{imap.UIDRange{Start: oldUid, Stop: rule.UIDNext}}

			fetchOptions := &imap.FetchOptions{
				Envelope: true,
				Flags:    true,
				UID:      true,
				BodySection: []*imap.FetchItemBodySection{
					{},
				},
			}

			messages := mailClient.Fetch(seqSet, fetchOptions)
			if err != nil {
				logger.Printf("FETCH command failed: %v", err)
				return
			}

			defer func(messages *imapclient.FetchCommand) {
				err := messages.Close()
				if err != nil {
					logger.Printf("Failed to close messages: %v", err)
				}
			}(messages)

			count := 0
			for msg := messages.Next(); msg != nil; msg = messages.Next() {
				parts, err := ParseMessageToParts(msg)
				if err != nil {
					logger.Printf("Cannot parse message body: %v", err)
				}

				parsed, err := ParseMessageParts(parts)
				if err != nil {
					logger.Printf("Cannot parse message parts: %v", err)
				}

				rule.MailChan <- &parsed

				count++
			}

			logger.Printf("Fetched %v new messages", count)
		}()
	}

	logger.Printf("Rule stopped")
	close(rule.MailChan)
}
