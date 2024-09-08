package mailbot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/puzpuzpuz/xsync/v3"
	"log"
	"mail2telegram/env"
	"mail2telegram/state"
)

func NewBot(token string) *Bot {
	api, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Fatalf("Unable to start mailbot %v", err)
	}

	return &Bot{BotApi: api, State: state.BotState{
		Rules:       *xsync.NewMapOf[string, *state.RuleState](),
		DefaultImap: state.ImapParams{},
		MailsChan:   make(chan *state.BotMailPack, 16),
	}}
}

const (
	ActHelp = "help"
	ActLogs = "logs"
	ActGet  = "get"
	ActSet  = "set"
	ActDel  = "del"
)

func (bot *Bot) CreateActions() error {
	_, err := bot.BotApi.Send(tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "/" + ActHelp,
			Description: "Список команд",
		},
		tgbotapi.BotCommand{
			Command:     "/" + ActLogs,
			Description: "Получить N последних логов `/logs N`",
		},
		tgbotapi.BotCommand{
			Command:     "/" + ActGet,
			Description: "Получить правила (без аргументов) `/get` или данные о правиле по имени `/get name`",
		},
		tgbotapi.BotCommand{
			Command: "/" + ActSet,
			Description: "Установить правило или параметры в json `/set {\"imapServer\":..., \"imapUser\":..., \"imapToken\":...}`\n" +
				"для правил: `/set {\"name\":...,\"box\":...,\"chatId\":...}` \n" +
				"также опционально originalMessageId (для топиков)",
		},
		tgbotapi.BotCommand{
			Command:     "/" + ActDel,
			Description: "Удалить правило по имени `/del name`",
		},
	))

	return err
}

func (bot *Bot) Run(ctx context.Context) {
	go bot.RunMailsProcessing(ctx)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.BotApi.GetUpdatesChan(u)

mainLoop:
	for {
		select {
		case <-ctx.Done():
			break mainLoop

		case update := <-updates:
			if update.Message == nil ||
				!update.Message.IsCommand() ||
				update.Message.Chat.ID != env.Env.AdminId {
				continue
			}

			switch update.Message.Command() {
			case ActHelp:
				bot.Help(update)
			case ActLogs:
				bot.Logs(update)
			case ActGet:
				bot.Get(update)
			case ActSet:
				bot.Set(update)
			case ActDel:
				bot.Del(update)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неопозннаная команда "+update.Message.Command())
				bot.Send(msg)
			}
		}
	}

	log.Print("Bot stopped")
}

func (bot *Bot) Send(msg tgbotapi.Chattable) {
	if _, err := bot.BotApi.Send(msg); err != nil {
		log.Printf("Error while sengins message: %v", err)
	}
}
