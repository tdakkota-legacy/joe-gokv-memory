# Joe Bot - gokv Adapter 
[![Go Report Card](https://goreportcard.com/badge/github.com/tdakkota/joe-gokv-memory)](https://goreportcard.com/report/github.com/tdakkota/joe-gokv-memory)
[![CodeFactor](https://www.codefactor.io/repository/github/tdakkota/joe-gokv-memory/badge)](https://www.codefactor.io/repository/github/tdakkota/joe-gokv-memory)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

Integration [Joe] with [gokv].

## Getting Started

To install:

```
go get github.com/tdakkota/joe-gokv-memory
```

### Note

[gokv] currently does not support GetAll/Keys operations(see this [issue]) 
Package provides `Keys` interface to store list of keys(by default map used).

### Example usage

```go
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

```
## License

[BSD-3-Clause](LICENSE)

[issue]: https://github.com/philippgille/gokv/issues/9
[joe]: https://github.com/go-joe/joe
[gokv]: https://github.com/philippgille/gokv

