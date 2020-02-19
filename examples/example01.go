package main

import (
	"github.com/go-joe/joe"
	"github.com/philippgille/gokv/redis"
	"github.com/tdakkota/joe-gokv-memory"
	"go.uber.org/zap"
)

type Bot struct {
	*joe.Bot
}

func main() {
	store, err := redis.NewClient(redis.DefaultOptions)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	b := Bot{joe.New(
		"example-bot",
		gokv.Memory(store),
	)}
	b.Respond("remember (.+) is (.+)", b.Remember)
	b.Respond("what is (.+)", b.WhatIs)

	err = b.Run()
	if err != nil {
		zap.L().Fatal(err.Error())
	}
}

func (b *Bot) Remember(msg joe.Message) error {
	key, value := msg.Matches[0], msg.Matches[1]
	msg.Respond("OK, I'll remember %s is %s", key, value)
	return b.Store.Set(key, value)
}

func (b *Bot) WhatIs(msg joe.Message) error {
	key := msg.Matches[0]
	var value string
	ok, err := b.Store.Get(key, &value)
	if err != nil {
		return err
	}
	if ok {
		msg.Respond("%s is %s", key, value)
	} else {
		msg.Respond("I do not remember %q", key)
	}
	return nil
}
