package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/harkce/listbot/listbot"
	"github.com/harkce/listbot/server"
	"github.com/subosito/gotenv"
)

const (
	CHANNEL_SECRET = "1953ed2b01fbb896a7cd067804628184"
	CHANNEL_TOKEN  = "7I0e9hSYEmdWRmXr/bK9x12WtLXF2vlWArGCe0yRVuuyF022ZJAxabuGWJ9u1SjvWnFspmpass12IvKWhj+oNxse1nXKQ+4l460RqveXMFa3pYAGJIYmGgp9+nkp3r8jxyXz1UIzU6Gkg54IrHrgKAdB04t89/1O/w1cDnyilFU="
)

func main() {
	log.Println("Starting listbot...")
	os.Setenv("CHANNEL_SECRET", CHANNEL_SECRET)
	os.Setenv("CHANNEL_TOKEN", CHANNEL_TOKEN)
	gotenv.Load(os.Getenv("GOPATH") + "/src/github.com/harkce/listbot/.env")

	//test
	dbURL := os.Getenv("DATABASE_URL")
	fmt.Println(dbURL)

	err := listbot.InitBot()
	if err != nil {
		log.Fatalln(err)
		return
	}

	router := server.Router()
	port := os.Getenv("PORT")

	log.Println(fmt.Sprintf("Listbot started @:%s", port))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	}()

	<-sigChan
	log.Println("Shutting down listbot...")
	log.Println("listbot stopped")
}
