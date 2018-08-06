package server

import (
	"fmt"
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
	router.GET("/me", me)

	return cors.Handler(router)
}

func me(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "ok")
}
