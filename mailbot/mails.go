package mailbot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

			text := ""
			for _, part := range msg.Msg.Inlines {
				content, _, err := part.Header.ContentType()
				if err == nil && strings.Contains(content, "text/plain") {
					text = string(part.Body)
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

			media := tgbotapi.NewMediaGroup(msg.Rule.Settings.ChatId, files)

			tgMsg := tgbotapi.NewMessage(msg.Rule.Settings.ChatId, text)
			tgMsg.ReplyToMessageID = msg.Rule.Settings.OriginalMsgId

			bot.Send(tgMsg)
			_, err := bot.BotApi.SendMediaGroup(media)
			if err != nil {
				log.Panic(err)
			}
		}
	}

	log.Println("Mails processing stopped")
}
