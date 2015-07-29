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
}

func Handler(cmds []string, size int) *handler {
	ret := &handler{
		cmds,
		make(chan *runner, size),
	}

	for i := 0; i < size; i++ {
		ret.queue <- Runner(ret.cmds)
	}

	return ret
}

func (h *handler) Process(bot *telebot.Bot, msg telebot.Message) {
	r := <-h.queue
	r.Run(bot, msg)
	h.queue <- Runner(h.cmds)
}

type runner struct {
	cmd *exec.Cmd
}

func Runner(cmds []string) *runner {
	cmd := cmds[0]
	args := cmds[1:]
	return &runner{exec.Command(cmd, args...)}
}

func (r *runner) Run(bot *telebot.Bot, msg telebot.Message) {
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
		bot.SendMessage(msg.Chat, string(data), nil)
	}
}
