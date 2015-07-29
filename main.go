package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/tucnak/telebot"
)

var (
	token_file = flag.String("token", "", "File contains telegram bot token")
	max_worker = flag.Int("worker", 5, "Max worker threads")
)

func main() {
	flag.Parse()
	if *token_file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	args := flag.Args()
	log.Print(args)
	if len(args) < 1 {
		log.Fatal(`Usage: ` + os.Args[0] + ` -worker 5 -token=token_file cmd_to_run arg1 arg2 ...`)
	}

	token, err := ioutil.ReadFile(*token_file)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telebot.NewBot(string(token))
	if err != nil {
		log.Fatal(err)
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	h := CreateHandler(bot, args, *max_worker)
	log.Printf("Workers: %d\n", *max_worker)

	for msg := range messages {
		go h.Process(bot, msg)
	}
}
