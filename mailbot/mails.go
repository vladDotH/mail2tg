package mailbot

import (
	"bytes"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"jaytaylor.com/html2text"
	"log"
	"strings"
)

func (bot *Bot) RunMailsProcessing(ctx context.Context) {
mainLoop:
	for {
		select {
		case <-ctx.Done():
			break mainLoop

		case msg := <-bot.State.MailsChan:
			log.Printf("Recieved: %v -> %v", msg.Rule.Settings.Box, msg.Rule.Settings.ChatId)

			msgText := strings.Builder{}

			subject, err := msg.Msg.Header.Subject()
			if err != nil {
				fmt.Printf("Cannot read subject: %v", err)
			} else {
				msgText.WriteString("Subject: ")
				msgText.WriteString(subject)
				msgText.WriteString("\n")
			}

			from, err := msg.Msg.Header.AddressList("From")
			if err != nil {
				fmt.Printf("Cannot read from: %v", err)
			} else {
				msgText.WriteString("From: ")
				for _, address := range from {
					msgText.WriteString(address.Name)
					msgText.WriteString(" <" + address.Address + "> ")
				}
				msgText.WriteString("\n")
			}

			for _, part := range msg.Msg.Inlines {
				content, _, err := part.Header.ContentType()
				if err == nil && strings.Contains(content, "text/html") {
					text, err := html2text.FromReader(bytes.NewReader(part.Body), html2text.Options{PrettyTables: false})
					if err != nil {
						log.Printf("Cannot parse html: %v", err)
					}
					msgText.WriteString(text)
				}
			}

			files := make([]interface{}, len(msg.Msg.Attachments))
			i := 0
			for _, part := range msg.Msg.Attachments {
				fileName := fmt.Sprintf("File_%d", i)
				_fileName, err := part.Header.Filename()
				if err == nil {
					fileName = _fileName
				}

				fileReader := tgbotapi.FileBytes{
					Name:  fileName,
					Bytes: part.Body,
				}

				files[i] = tgbotapi.NewInputMediaDocument(fileReader)
				i++
			}

			tgMsg := tgbotapi.NewMessage(msg.Rule.Settings.ChatId, msgText.String())
			tgMsg.ReplyToMessageID = msg.Rule.Settings.OriginalMsgId
			if len(tgMsg.Text) > 0 {
				bot.Send(tgMsg)
			}

			if len(files) > 0 {
				media := tgbotapi.NewMediaGroup(msg.Rule.Settings.ChatId, files)
				media.ReplyToMessageID = msg.Rule.Settings.OriginalMsgId
				_, err = bot.BotApi.SendMediaGroup(media)
				if err != nil {
					log.Printf("Cannot send files: %v", err)
					errMsg := tgbotapi.NewMessage(msg.Rule.Settings.ChatId, "Не удалось отправить приложения к письму :(")
					errMsg.ReplyToMessageID = msg.Rule.Settings.OriginalMsgId
					bot.Send(errMsg)
				}
			}
		}
	}

	log.Println("Mails processing stopped")
}
