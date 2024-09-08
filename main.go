package main

import (
	"context"
	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"mail2telegram/env"
	"mail2telegram/mailbot"
	"os"
	"os/signal"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Cannot read .env file: %v", err)
	}

	env.LoadEnv()

	logger := lumberjack.Logger{
		Filename:   env.Env.LogFile,
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     31,
	}

	multi := io.MultiWriter(&logger, os.Stdout)
	log.SetOutput(multi)

	bot := mailbot.NewBot(env.Env.TgToken)
	bot.BotApi.Debug = env.Env.Debug

	err = bot.CreateActions()
	if err != nil {
		// Library bug
		log.Printf("Cannot create actions: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go bot.Run(ctx)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)

	for {
		select {
		case <-exit:
			log.Println("Stopping mailbot...")
			cancel()
			time.Sleep(2 * time.Second)
			return
		}
	}
}
