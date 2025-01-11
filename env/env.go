package env

import (
	"log"
	"os"
	"strconv"
)

type EnvType struct {
	TgToken       string
	AdminId       int64
	Debug         bool
	LogFile       string
	StoragePrefix string
	HTTPPrefix    string
	HTTPAddr      string
}

var Env = EnvType{}

func LoadEnv() {
	tgToken, exists := os.LookupEnv("TG_TOKEN")
	if !exists {
		log.Fatalf("Invalid EnvType: TG_TOKEN")
	}

	adminId, exists := os.LookupEnv("ADMIN_ID")
	if !exists {
		log.Fatalf("Invalid EnvType: TG_TOKEN")
	}

	adminIdNum, err := strconv.ParseInt(adminId, 10, 64)
	if err != nil {
		log.Fatalf("Invalid AdminId: %v", err)
	}

	var debugVal bool
	if debug, exists := os.LookupEnv("DEBUG"); !exists {
		debugVal = false
	} else {
		debugVal = debug == "true"
	}

	storagePrefix, exists := os.LookupEnv("STORAGE_PREFIX")
	if !exists {
		log.Fatalf("Invalid storage prefix: %v", err)
	}

	httpPrefix, exists := os.LookupEnv("HTTP_PREFIX")
	if !exists {
		log.Fatalf("Invalid http prefix: %v", err)
	}

	httpAddr, exists := os.LookupEnv("HTTP_ADDR")
	if !exists {
		log.Fatalf("Invalid http address: %v", err)
	}

	Env.TgToken = tgToken
	Env.AdminId = adminIdNum
	Env.Debug = debugVal
	Env.LogFile = "./logs"
	Env.StoragePrefix = storagePrefix
	Env.HTTPPrefix = httpPrefix
	Env.HTTPAddr = httpAddr
}
