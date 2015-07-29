package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/tucnak/telebot"
)

type handler struct {
	cmds  []string
	queue chan *runner
	bot *telebot.Bot
}

func Handler(bot *telebot.Bot, cmds []string, size int) *handler {
	ret := &handler{
		cmds,
		make(chan *runner, size),
		bot,
	}

	for i := 0; i < size; i++ {
		ret.queue <- Runner(ret.cmds, bot)
	}

	return ret
}

func (h *handler) Process(bot *telebot.Bot, msg telebot.Message) {
	r := <-h.queue
	r.Run(msg)
	h.queue <- Runner(h.cmds, bot)
}

type reply struct {
	Type string `json:"type"`
	User *telebot.User `json:"user"`
	Content string `json:"content"`
}


type runner struct {
	cmd *exec.Cmd
	bot *telebot.Bot
}

func Runner(cmds []string, bot *telebot.Bot) *runner {
	cmd := cmds[0]
	args := cmds[1:]
	return &runner{exec.Command(cmd, args...), bot}
}

func (r *runner) Run(msg telebot.Message) {
	stdin, err := r.cmd.StdinPipe()
	if err != nil {
		log.Print(err)
	}

	stdout, err := r.cmd.StdoutPipe()
	if err != nil {
		log.Print(err)
	}

	r.cmd.Start()

	if data, err := json.Marshal(msg); err == nil {
		fmt.Fprint(stdin, data)
		stdin.Close()
	} else {
		log.Print(err)
	}

	if data, _ := ioutil.ReadAll(stdout); string(data) != "" {
		var reps []reply
		if err := json.Unmarshal(data, &reps); err == nil {
			for _, rep := range reps {
				r.handleReply(msg, &rep)
			}
		}
	}
}

func (r *runner) handleReply(msg telebot.Message, rep *reply) {
	if rep.User == nil {
		rep.User = &msg.Chat
	}

	switch rep.Type {
	case "doc":
		if file, err := telebot.NewFile(rep.Content); err == nil {
			doc := telebot.Document{File: file}
			r.bot.SendDocument(*rep.User, &doc, nil)
		}
	case "photo":
		if file, err := telebot.NewFile(rep.Content); err == nil {
			photo := telebot.Photo{Thumbnail: telebot.Thumbnail{File: file}}
			r.bot.SendPhoto(*rep.User, &photo, nil)
		}
	case "audio":
		if file, err := telebot.NewFile(rep.Content); err == nil {
			audio := telebot.Audio{File: file}
			r.bot.SendAudio(*rep.User, &audio, nil)
		}
	case "video":
		if file, err := telebot.NewFile(rep.Content); err == nil {
			video := telebot.Video{Audio: telebot.Audio{File: file}}
			r.bot.SendVideo(*rep.User, &video, nil)
		}
	default:
		if rep.Content != "" {
			r.bot.SendMessage(*rep.User, rep.Content, nil)
		}
	}
}
