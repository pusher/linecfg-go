package main

import (
	"log"

	"github.com/pusher/linecfg-go"
)

type Config struct {
	Host     string
	SomePort int `linecfg:"port"`
}

func main() {
	c := new(Config)

	err := linecfg.Getenv("MY_CFG", c)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("host", c.Host)
	log.Println("port", c.SomePort)
}
