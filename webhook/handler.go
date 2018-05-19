package webhook

import (
	"log"
	"net/http"

	"github.com/harkce/listbot"
	"github.com/harkce/listbot/server/response"
	"github.com/julienschmidt/httprouter"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Handler struct{}

func hookResp(w http.ResponseWriter) {
	var resp interface{}
	response.OK(w, resp)
}

func (h *Handler) WebHook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	events, err := listbot.Client.ParseRequest(r)
	if err != nil {
		log.Println(err)
		hookResp(w)
		return
	}

	for _, event := range events {
		if event.Type != linebot.EventTypeMessage {
			hookResp(w)
			return
		}
	}

	hookResp(w)
}
