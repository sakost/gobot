package main

import (
	"github.com/nikepan/govkbot"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	AdminID = 150820042
)

type EnvMap map[string]string

var running bool = true

func initSigHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		s := <-c
		log.Println("got a signal:", s)
		log.Println("gracefully exiting...")
		running = false
	}()
}

func main() {
	var (
		environ = getEnv("SAKOST_BOT_")
		err     error
		token   string
	)

	govkbot.SetDebug(true)

	token, err = GetToken(environ["LOGIN"], environ["PASSWORD"], "", -1)
	checkErr(err)

	if govkbot.API.DEBUG {
		log.Printf("using token: %s...", token[:5])
	}
	govkbot.SetToken(token)
	initSigHandler()
	initTasks()
	for running {
		time.Sleep(time.Minute)
	}
}
