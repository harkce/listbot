package server

import (
	"net/http"

	"github.com/goware/cors"
	"github.com/harkce/listbot/webhook"
	"github.com/julienschmidt/httprouter"
)

func Router() http.Handler {
	router := httprouter.New()
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		MaxAge:         86400,
	})

	webhookHandler := webhook.Handler{}
	router.POST("/webhook", webhookHandler.WebHook)

	return cors.Handler(router)
}
