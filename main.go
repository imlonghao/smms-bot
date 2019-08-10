package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/imlonghao/smms-bot/smms"
	log "github.com/sirupsen/logrus"
)

func upload(token, url string) (string, string) {
	resp, err := http.Get(url)
	if err != nil {
		return "err", "Fail to download the image"
	}
	defer resp.Body.Close()
	client := smms.Client{Token: token}
	json, err := client.Upload(resp.Body, "1.png")
	if err != nil {
		return "err", "Fail to upload the image"
	}
	if json.Success == false {
		return "err", json.Message
	}
	return json.Data.Hash, json.Data.URL
}

func uploadHandler(url string, token string, update tgbotapi.Update) tgbotapi.MessageConfig {
	hash, url := upload(token, url)
	kb := tgbotapi.InlineKeyboardButton{
		Text:         "Click Here To Delete",
		CallbackData: &hash,
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("`%s`", url))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true
	if hash != "err" {
		msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
			[][]tgbotapi.InlineKeyboardButton{
				[]tgbotapi.InlineKeyboardButton{
					kb,
				}}}
	}
	return msg
}

func main() {
	apis := make(map[int]string)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			// 命令
			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
				case "set":
					apis[update.Message.From.ID] = update.Message.CommandArguments()
					msg.Text = "API token saved."
				case "get":
					msg.Text = apis[update.Message.From.ID]
				}
				msg.ReplyToMessageID = update.Message.MessageID
				if _, err := bot.Send(msg); err != nil {
					log.Warnf("fail to send msg, %v", err)
				}
				continue
			}
			// 文件形式的图片
			if update.Message.Document != nil {
				if !strings.Contains(update.Message.Document.MimeType, "image/") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "File has an invalid extension.")
					msg.ReplyToMessageID = update.Message.MessageID
					if _, err := bot.Send(msg); err != nil {
						log.Warnf("fail to send msg, %v", err)
					}
					continue
				}
				fileID := update.Message.Document.FileID
				url, err := bot.GetFileDirectURL(fileID)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to download the image.")
					msg.ReplyToMessageID = update.Message.MessageID
					if _, err := bot.Send(msg); err != nil {
						log.Warnf("fail to send msg, %v", err)
					}
					continue
				}
				msg := uploadHandler(url, apis[update.Message.From.ID], update)
				if _, err := bot.Send(msg); err != nil {
					log.Warnf("fail to send msg, %v", err)
				}
				continue
			}
			// 图片
			if update.Message.Photo != nil {
				photo := (*update.Message.Photo)
				fileID := photo[len(photo)-1].FileID
				url, err := bot.GetFileDirectURL(fileID)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to download the image.")
					msg.ReplyToMessageID = update.Message.MessageID
					if _, err := bot.Send(msg); err != nil {
						log.Warnf("fail to send msg, %v", err)
					}
					continue
				}
				msg := uploadHandler(url, apis[update.Message.From.ID], update)
				if _, err := bot.Send(msg); err != nil {
					log.Warnf("fail to send msg, %v", err)
				}
				continue
			}
		}
		// Callback 删除图片
		if update.CallbackQuery != nil {
			client := smms.Client{}
			if update.CallbackQuery.Data != "err" {
				client.Delete(update.CallbackQuery.Data)
			}
			edit := tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:    int64(update.CallbackQuery.From.ID),
					MessageID: update.CallbackQuery.Message.MessageID,
				},
				Text: "Photo Deleted!",
			}
			bot.Send(edit)
		}
	}
}
