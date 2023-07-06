package main

import (
	tgClient "TelegramBot/clients/telegram"
	"TelegramBot/consumer/eventConsumer"
	"TelegramBot/events/telegram"
	"TelegramBot/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)
	log.Print("service started")

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}

func mustToken() string {
	token := flag.String("token", "", "token for access to tg-bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is empty")
	}
	return *token
}
