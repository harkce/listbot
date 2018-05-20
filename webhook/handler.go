package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	events, err := ParseRequest(os.Getenv("CHANNEL_SECRET"), r)
	if err != nil {
		log.Println("Error parse:", err)
		hookResp(w)
		return
	}

	for _, event := range events {
		if unsupportedEvent(event) {
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
	}

	hookResp(w)
}

func sendReply(replyToken, content string) {
	_, err := listbot.Client.ReplyMessage(replyToken, linebot.NewTextMessage(content)).Do()
	if err != nil {
		log.Println("Error reply:", err)
	}
}

func ParseRequest(channelSecret string, r *http.Request) ([]*linebot.Event, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if !validateSignature(channelSecret, r.Header.Get("X-Line-Signature"), body) {
		return nil, linebot.ErrInvalidSignature
	}

	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err = json.Unmarshal(body, request); err != nil {
		return nil, err
	}
	return request.Events, nil
}

func validateSignature(channelSecret, signature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))
	hash.Write(body)
	return hmac.Equal(decoded, hash.Sum(nil))
}
