package main

import (
	"context"
	"io"
	"log"
	"mail2telegram/db"
	"mail2telegram/env"
	"mail2telegram/mailbot"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/peterbourgon/diskv/v3"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	db.Init(diskv.Options{
		BasePath:     "./storage",
		Transform:    func(s string) []string { return []string{} },
		CacheSizeMax: 1024,
	})

	err := godotenv.Load()
	if err != nil {
		log.Printf("Cannot read .env file: %v", err)
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
