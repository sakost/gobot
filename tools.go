package main

import (
	"log"
	"os"
	"strings"
)

func getEnv(prefix string) (env EnvMap) {
	env = EnvMap{}
	_env := os.Environ()
	for _, s := range _env {
		if strings.HasPrefix(s, prefix) {
			s = strings.TrimPrefix(s, prefix)
			separated := strings.Split(s, "=")
			env[separated[0]] = separated[1]
		}
	}
	return
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
