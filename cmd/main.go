package main

import (
	"TelegramBot/clients/telegram"
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	// get token
	token := mustToken()

	// create client
	tgClient = telegram.New(tgBotHost, token)

	// fetcher = fetcher.New(tgClient)

	// processor = processor.New(tgClient)

	// consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String("token", "", "token for access to tg-bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is empty")
	}
	return *token
}
