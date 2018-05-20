package webhook

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/harkce/listbot/listbot"
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
	return e.Type != linebot.EventTypeMessage || e.Source.Type != linebot.EventSourceTypeGroup
}

func (h *Handler) WebHook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	events, err := listbot.Client.ParseRequest(r)
	if err != nil {
		log.Println("Error parse:", err)
		hookResp(w)
		return
	}

	for _, event := range events {
		if unsupportedEvent(&event) {
			log.Println("Unsupported event")
			continue
		}

		message, ok := event.Message.(*linebot.TextMessage)
		if !ok {
			log.Println("Not a text message")
			continue
		}

		content := message.Text
		if !strings.HasPrefix(content, "/") {
			log.Println("Not a bot command")
			continue
		}

		var replyMessage string
		groupID := event.Source.GroupID
		replyToken := event.ReplyToken
		args := strings.Split(content, " ")
		if strings.HasPrefix(content, "/list") {
			replyMessage = listbot.LoadList(groupID)
			sendReply(replyToken, replyMessage)
			continue
		}
		if strings.HasPrefix(content, "/title ") {
			if len(args) < 2 {
				replyMessage = ""
			} else {
				replyMessage = listbot.SetTitle(groupID, strings.Join(args[1:], " "))
			}
			sendReply(replyToken, replyMessage)
			continue
		}
		if strings.HasPrefix(content, "/add ") {
			if len(args) < 2 {
				replyMessage = ""
			} else {
				replyMessage = listbot.AddItem(groupID, strings.Join(args[1:], " "))
			}
			sendReply(replyToken, replyMessage)
			continue
		}
		if strings.HasPrefix(content, "/edit ") {
			if len(args) < 3 {
				replyMessage = ""
			} else {
				pos, err := strconv.Atoi(args[1])
				if err != nil {
					replyMessage = ""
				} else {
					replyMessage = listbot.EditItem(groupID, pos, strings.Join(args[2:], " "))
				}
			}
			sendReply(replyToken, replyMessage)
			continue
		}
		if strings.HasPrefix(content, "/delete ") {
			if len(args) < 2 {
				replyMessage = ""
			} else {
				pos, err := strconv.Atoi(args[1])
				if err != nil {
					replyMessage = ""
				} else {
					replyMessage = listbot.DeleteItem(groupID, pos)
				}
			}
			sendReply(replyToken, replyMessage)
			continue
		}
		if strings.HasPrefix(content, "/clear") {
			replyMessage = listbot.ClearItem(groupID)
			sendReply(replyToken, replyMessage)
			continue
		}
		if strings.HasPrefix(content, "/help") {
			replyMessage = "Commands:\n" +
				"/list - Show current list\n" +
				"/title <title> - Set title of the list\n" +
				"/add <item> - Add item to the current list\n" +
				"/edit <item number> <item> - Edit selected item\n" +
				"/delete <item number> - Delete an item\n" +
				"/clear - Reset list\n" +
				"/help - Show bot commands"
			sendReply(replyToken, replyMessage)
			continue
		}
	}

	hookResp(w)
}

func sendReply(replyToken, content string) {
	_, err := listbot.Client.ReplyMessage(replyToken, linebot.NewTextMessage(content)).Do()
	if err != nil {
		log.Println("Error reply:", err)
	}
}
