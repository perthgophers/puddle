package responses

import (
	"github.com/perthgophers/puddle/messagerouter"
)

var mr *messagerouter.MessageRouter
var commands = make(map[string]messagerouter.MessageHandler)

//Handle Store comamnds before initialisation
func Handle(key string, h messagerouter.MessageHandler) {
	if mr == nil {
		commands[key] = h
		return
	}

	mr.Handle(key, h)
}

//Init initialises responses
func Init(msgr *messagerouter.MessageRouter) {
	mr = msgr
	for k, fn := range commands {
		mr.Handle(k, fn)
	}

	commands = make(map[string]messagerouter.MessageHandler)
}
