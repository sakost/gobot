package main

import (
	"github.com/nikepan/govkbot"
	"log"
	"time"
)

func infOnline() {
	var (
		foo interface{}
		err error
	)
	log.Println("running \"infOnline\" task")
	for running {
		err = govkbot.API.CallMethod("account.setOnline", make(map[string]string), foo)
		if err != nil {
			log.Println("got an error:", err)
			err = nil
		}
		time.Sleep(5 * time.Minute)
	}
	log.Println("\"infOnline\" task exiting..")
}

type ResObj struct {
	Text string `json:"text"`
}

type Response struct {
	Response ResObj `json:"response"`
}

func randomStatus() {
	var (
		status string
		res    interface{}
		getRes Response
		args   = make(map[string]string)
		err    error
	)
	log.Println("running \"randomStatus\" task")
	err = govkbot.API.CallMethod("status.get", make(map[string]string), &getRes)
	status = getRes.Response.Text
	for running {
		args["text"] = status + " | " + time.Now().Format("15:04") + " | автостатус включён"
		err = govkbot.API.CallMethod("status.set", args, res)
		if err != nil {
			log.Println("got an error:", err)
			err = nil
		}
		time.Sleep(time.Minute)
	}
	args["text"] = status
	err = govkbot.API.CallMethod("status.set", args, res)
	if err != nil {
		log.Println("got an error:", err)
	}
	log.Println("\"randomStatus\" task exiting..")
}

func initTasks() {
	go infOnline()
	go randomStatus()
}
