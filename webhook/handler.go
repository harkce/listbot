package webhook

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/harkce/listbot/listbot"
	"github.com/julienschmidt/httprouter"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Handler struct{}

func hookResp(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func unsupportedEvent(e *linebot.Event) bool {
	return e.Type != linebot.EventTypeMessage && e.Type != linebot.EventTypeLeave && (e.Source.Type != linebot.EventSourceTypeGroup || e.Source.Type != linebot.EventSourceTypeRoom)
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

		var replyMessage string
		var groupID string
		if event.Source.Type == linebot.EventSourceTypeGroup {
			groupID = event.Source.GroupID
		} else if event.Source.Type == linebot.EventSourceTypeRoom {
			groupID = event.Source.RoomID
		}

		if event.Type == linebot.EventTypeLeave {
			listbot.UnsetEnv(groupID)
			continue
		}

		if event.Source.Type == linebot.EventSourceTypeUser {
			replyMessage = "Halo! Makasih udah chat 😉\n" +
				"Saat ini, List bot ngga bisa chat personal sama kamu, maaf yaa 🙏\n\n" +
				"Kamu harus invite list bot ke grup/multi chat biar bisa nambahin list.\n" +
				"Kalo ada pertanyaan, kesulitan, atau saran, kamu bisa langsung hubungi yang bikin list bot 😃\n\n" +
				"LINE: http://line.me/ti/p/~harkce"
			sendReply(event.ReplyToken, replyMessage)
		}

		message, ok := event.Message.(*linebot.TextMessage)
		if !ok {
			continue
		}

		content := message.Text
		if !strings.HasPrefix(content, "/") {
			continue
		}

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
			replyMessage = "Perintah:\n" +
				"/list - Tampilkan list\n" +
				"/title <title> - Ganti judul list\n" +
				"/add <item> - Tambah item ke list\n" +
				"/edit <nomor> <item> - Edit item di posisi <nomor>\n" +
				"/delete <nomor> - Hapus item dari list\n" +
				"/clear - Hapus semua list\n" +
				"/help - Tampilkan perintah bot"
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
