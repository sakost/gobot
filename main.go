package main

import (
	"github.com/nikepan/govkbot"
	"log"
)

const (
	AdminID = 150820042
)

type EnvMap map[string]string

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
		log.Printf("using token: %s", token[:5])
	}
}
