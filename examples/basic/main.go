package main

import (
	"log"
	"time"

	"github.com/nenavizhuleto/gondler"
)

type MessageType int

const (
	_ MessageType = iota
	INFO
	WARN
	PANIC
	UNKNOWN
)

type Message struct {
	Type MessageType
	Data string
}

func main() {
	messages := []Message{
		{Type: INFO, Data: "hello world"},
		{Type: WARN, Data: "goodbye world"},
		{Type: UNKNOWN, Data: "what am i doing here?"},
		{Type: PANIC, Data: "welcome to hell"},
	}
	source := make(chan Message)
	go func() {
		for _, message := range messages {
			source <- message
			time.Sleep(1 * time.Second)
		}

		close(source)
	}()

	handler := gondler.New(source, func(m Message) MessageType { return m.Type })

	handler.On(INFO, func(m Message) {
		log.Println("just informational message:", m.Data)
	})

	handler.On(WARN, func(m Message) {
		log.Println("warning! take a look at that:", m.Data)
	})

	handler.On(PANIC, func(m Message) {
		log.Println("holy moly!!! we're going to panic right now!", m.Data)
	})

	handler.Default(func(m Message) {
		log.Println("unknown message", m.Data)
	})

	handler.RunSync()
}
