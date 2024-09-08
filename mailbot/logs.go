package mailbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"mail2telegram/env"
	"os/exec"
	"strconv"
)

func (bot *Bot) Logs(upd tgbotapi.Update) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "")

	var lines int64 = -1

	if len(upd.Message.CommandArguments()) > 0 {
		newLines, err := strconv.ParseInt(upd.Message.CommandArguments(), 10, 64)
		if err != nil {
			msg.Text = "Неверное число"
			bot.Send(msg)
		}
		lines = newLines
	}

	var cmd *exec.Cmd
	if lines != -1 {
		cmd = exec.Command("tail", "-n", strconv.FormatInt(lines, 10), env.Env.LogFile)
	} else {
		cmd = exec.Command("tail", env.Env.LogFile)
	}

	output, err := cmd.Output()
	if err != nil {
		log.Printf("Tail error: %v", err)
	}

	msg.Text = string(output)
	bot.Send(msg)
}
