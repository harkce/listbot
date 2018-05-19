package webhook

import (
	"log"
	"net/http"
	"strings"

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

func unsupportedEvent(e *linebot.Event) bool {
	return e.Type != linebot.EventTypeMessage || e.Source != linebot.EventSourceTypeGroup
}

func (h *Handler) WebHook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	events, err := listbot.Client.ParseRequest(r)
	if err != nil {
		log.Println(err)
		hookResp(w)
		return
	}

	for _, event := range events {
		if unsupportedEvent(event) {
			continue
		}

		message, ok := event.Message.(*linebot.TextMessage)
		if !ok {
			continue
		}

		content := message.Text
		if !strings.HasPrefix(content, "/") {
			continue
		}

		if strings.HasPrefix(content, "/list") {
			continue
		}
		if strings.HasPrefix(content, "/title") {
			continue
		}
		if strings.HasPrefix(content, "/add") {
			continue
		}
		if strings.HasPrefix(content, "/delete") {
			continue
		}
		if strings.HasPrefix(content, "/clear") {
			continue
		}
	}

	hookResp(w)
}
