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

var allowedPrefix = map[string]bool{
	"/list":       true,
	"/title":      true,
	"/add":        true,
	"/edit":       true,
	"/check":      true,
	"/cross":      true,
	"/unmark":     true,
	"/delete":     true,
	"/clear":      true,
	"/help":       true,
	"/multiple":   true,
	"/newlist":    true,
	"/removelist": true,
}

const personalReply = "Halo! Makasih udah chat 😉\n" +
	"Saat ini, list bot ngga bisa chat personal sama kamu, maaf yaa 🙏\n\n" +
	"Kamu harus invite list bot ke grup/multi chat biar bisa nambahin list.\n\n" +
	"Kalo ada pertanyaan, kesulitan, atau saran, kamu bisa langsung hubungi adminnya list bot 😃\n\n" +
	"LINE: http://line.me/ti/p/~harkce"

const singleHelp = "Perintah:\n" +
	"/list - Tampilkan item list\n\n" +
	"/title <judul> - Ganti judul list\n\n" +
	"/add <item> - Tambah item ke list\n\n" +
	"/edit <nomor> <item> - Edit item di posisi <nomor>\n\n" +
	"/check <nomor> - Menandai item dengan ✓\n\n" +
	"/cross <nomor> - Menandai item dengan ✗\n\n" +
	"/unmark <nomor> - Menghilangkan tanda pada item\n\n" +
	"/delete <nomor> - Hapus item dari list\n\n" +
	"/clear - Hapus semua item dari list\n\n" +
	"/multiple on - Mengaktifkan multiple list\n\n" +
	"/help - Tampilkan perintah bot"

const multipleHelp = "Perintah:\n" +
	"/newlist <judul> - Buat list baru\n\n" +
	"/list - Tampilkan semua list\n\n" +
	"/list <nomorlist> - Tampilkan item di list <nomorlist>\n\n" +
	"/title <nomorlist> <judul> - Ganti judul list\n\n" +
	"/add <nomorlist> <item> - Tambah item ke list <nomorlist>\n\n" +
	"/edit <nomorlist> <nomoritem> <item> - Edit item list\n\n" +
	"/check <nomorlist> <nomoritem> - Menandai item dengan ✓\n\n" +
	"/cross <nomorlist> <nomoritem> - Menandai item dengan ✗\n\n" +
	"/unmark <nomorlist> <nomoritem> - Menghilangkan tanda pada item\n\n" +
	"/delete <nomorlist> <nomoritem> - Hapus item dari list <nomorlist>\n\n" +
	"/removelist <nomorlist> - Hapus list nomor <nomorlist>\n\n" +
	"/multiple off - Menonaktifkan multiple list\n\n" +
	"/help - Tampilkan perintah bot"

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

		if event.Type == linebot.EventTypeLeave && event.Source.Type == linebot.EventSourceTypeRoom {
			listbot.UnsetEnv(groupID)
			continue
		}

		if event.Source.Type == linebot.EventSourceTypeUser {
			replyMessage = personalReply
			sendReply(event.ReplyToken, replyMessage)
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

		replyToken := event.ReplyToken
		args := strings.Split(content, " ")

		var l *listbot.List
		if allowedPrefix[args[0]] {
			l, _ = listbot.Retrieve(groupID)
		}

		if strings.HasPrefix(content, "/multiple") {
			replyMessage = l.SetMultiple(args[1])
			if l.Multiple {
				replyMessage += "\n" + "Gunakan '/multiple off' untuk menonaktifkan multiple list"
			} else {
				replyMessage += "\n" + "Gunakan '/multiple on' untuk mengaktifkan multiple list"
			}
			sendReply(replyToken, replyMessage)
			continue
		}

		if l.Multiple {
			if strings.HasPrefix(content, "/newlist") {
				if len(args) < 2 {
					replyMessage = ""
				} else {
					replyMessage = l.CreateList(strings.Join(args[1:], " "))
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/list") {
				var pos int
				var err error
				if len(args) > 1 {
					pos, err = strconv.Atoi(args[1])
					if err != nil {
						replyMessage = l.LoadMultiple()
					} else {
						replyMessage = l.LoadElement(pos)
					}
				} else {
					replyMessage = l.LoadMultiple()
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/title") {
				if len(args) < 3 {
					replyMessage = ""
				} else {
					pos, err := strconv.Atoi(args[1])
					if err != nil {
						replyMessage = ""
					} else {
						replyMessage = l.SetElementTitle(pos, strings.Join(args[2:], " "))
					}
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/add") {
				if len(args) < 3 {
					replyMessage = ""
				} else {
					pos, err := strconv.Atoi(args[1])
					if err != nil {
						replyMessage = ""
					} else {
						replyMessage = l.AddElementItem(pos, strings.Join(args[2:], " "))
					}
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/edit") {
				if len(args) < 4 {
					replyMessage = ""
				} else {
					listpos, err1 := strconv.Atoi(args[1])
					pos, err2 := strconv.Atoi(args[2])
					if err1 != nil || err2 != nil {
						replyMessage = ""
					} else {
						replyMessage = l.EditElementItem(listpos, pos, strings.Join(args[3:], " "))
					}
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/check") {
				if len(args) < 3 {
					replyMessage = ""
				} else {
					listpos, err1 := strconv.Atoi(args[1])
					pos, err2 := strconv.Atoi(args[2])
					if err1 != nil || err2 != nil {
						replyMessage = ""
					} else {
						replyMessage = l.CheckElementItem(listpos, pos)
					}
				}
				sendReply(replyToken, replyMessage)
			}
			if strings.HasPrefix(content, "/cross") {
				if len(args) < 3 {
					replyMessage = ""
				} else {
					listpos, err1 := strconv.Atoi(args[1])
					pos, err2 := strconv.Atoi(args[2])
					if err1 != nil || err2 != nil {
						replyMessage = ""
					} else {
						replyMessage = l.CrossElementItem(listpos, pos)
					}
				}
				sendReply(replyToken, replyMessage)
			}
			if strings.HasPrefix(content, "/unmark") {
				if len(args) < 3 {
					replyMessage = ""
				} else {
					listpos, err1 := strconv.Atoi(args[1])
					pos, err2 := strconv.Atoi(args[2])
					if err1 != nil || err2 != nil {
						replyMessage = ""
					} else {
						replyMessage = l.UncheckElementItem(listpos, pos)
					}
				}
				sendReply(replyToken, replyMessage)
			}
			if strings.HasPrefix(content, "/delete") {
				if len(args) < 3 {
					replyMessage = ""
				} else {
					listpos, err1 := strconv.Atoi(args[1])
					pos, err2 := strconv.Atoi(args[2])
					if err1 != nil || err2 != nil {
						replyMessage = ""
					} else {
						replyMessage = l.DeleteElementItem(listpos, pos)
					}
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/removelist") {
				if len(args) < 2 {
					replyMessage = ""
				} else {
					pos, err := strconv.Atoi(args[1])
					if err != nil {
						replyMessage = ""
					} else {
						replyMessage = l.RemoveList(pos)
					}
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/help") {
				replyMessage = multipleHelp
				sendReply(replyToken, replyMessage)
				continue
			}
		} else {
			if strings.HasPrefix(content, "/list") {
				replyMessage = l.LoadList(groupID)
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/title ") {
				if len(args) < 2 {
					replyMessage = ""
				} else {
					replyMessage = l.SetTitle(strings.Join(args[1:], " "))
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/add ") {
				if len(args) < 2 {
					replyMessage = ""
				} else {
					replyMessage = l.AddItem(strings.Join(args[1:], " "))
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
						replyMessage = l.EditItem(pos, strings.Join(args[2:], " "))
					}
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/check ") {
				if len(args) < 2 {
					replyMessage = ""
				} else {
					pos, err := strconv.Atoi(args[1])
					if err != nil {
						replyMessage = ""
					} else {
						replyMessage = l.CheckItem(pos)
					}
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/cross ") {
				if len(args) < 2 {
					replyMessage = ""
				} else {
					pos, err := strconv.Atoi(args[1])
					if err != nil {
						replyMessage = ""
					} else {
						replyMessage = l.CrossItem(pos)
					}
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/unmark ") {
				if len(args) < 2 {
					replyMessage = ""
				} else {
					pos, err := strconv.Atoi(args[1])
					if err != nil {
						replyMessage = ""
					} else {
						replyMessage = l.UncheckItem(pos)
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
						replyMessage = l.DeleteItem(pos)
					}
				}
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/clear") {
				replyMessage = l.ClearItem()
				sendReply(replyToken, replyMessage)
				continue
			}
			if strings.HasPrefix(content, "/help") {
				replyMessage = singleHelp
				sendReply(replyToken, replyMessage)
				continue
			}
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
