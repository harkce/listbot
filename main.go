package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/harkce/listbot/listbot"
	"github.com/harkce/listbot/server"
	"github.com/subosito/gotenv"
)

func main() {
	log.Println("Starting listbot...")
	gotenv.Load(os.Getenv("GOPATH") + "/src/github.com/harkce/listbot/.env")

	err := listbot.InitBot()
	if err != nil {
		log.Fatalln(err)
		return
	}

	router := server.Router()

	log.Println("Listbot started @:8063")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		http.ListenAndServe(":8063", router)
	}()

	<-sigChan
	log.Println("Shutting down listbot...")
	log.Println("listbot stopped")
}
