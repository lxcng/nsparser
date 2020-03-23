package main

import (
	"log"
	"nsparser/config"
	"nsparser/server"
	"os"
	"os/signal"
)

func main() {
	path, ok := os.LookupEnv("NS_DL_CONF_PATH")
	if !ok {
		path = "config.json"
	}
	log.Println("ver_0.3")
	config.NewConf(path)
	go server.StartServer()
	<-listenForInterrupt()
}

func listenForInterrupt() <-chan bool {
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)
	go shutdown(quit, done)
	return done
}

func shutdown(quit <-chan os.Signal, done chan<- bool) {
	<-quit
	config.Save()
	close(done)
}
