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

func sendError(update tgbotapi.Update, bot *tgbotapi.BotAPI, err string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, err)
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		log.Warnf("fail to send msg, %v", err)
	}
}

func main() {
	apis := make(map[int]string)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic("Fail to get Telegram updates")
	}
	for {
		select {
		case update := <-updates:
			go func() {
				if update.Message != nil {
					// 命令
					if update.Message.IsCommand() {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
						switch update.Message.Command() {
						case "set":
							apis[update.Message.From.ID] = update.Message.CommandArguments()
							msg.Text = "API token saved."
						case "get":
							text := apis[update.Message.From.ID]
							if text == "" {
								text = "API token empty."
							}
							msg.Text = text
						}
						msg.ReplyToMessageID = update.Message.MessageID
						if _, err := bot.Send(msg); err != nil {
							log.Warnf("fail to send msg, %v", err)
						}
						return
					}
					// 文件形式的图片
					if update.Message.Document != nil {
						if !strings.Contains(update.Message.Document.MimeType, "image/") {
							sendError(update, bot, "File has an invalid extension.")
							return
						}
						fileID := update.Message.Document.FileID
						url, err := bot.GetFileDirectURL(fileID)
						if err != nil {
							sendError(update, bot, "Failed to download the image.")
							return
						}
						msg := uploadHandler(url, apis[update.Message.From.ID], update)
						if _, err := bot.Send(msg); err != nil {
							log.Warnf("fail to send msg, %v", err)
						}
						return
					}
					// 图片
					if update.Message.Photo != nil {
						photo := (*update.Message.Photo)
						fileID := photo[len(photo)-1].FileID
						url, err := bot.GetFileDirectURL(fileID)
						if err != nil {
							sendError(update, bot, "Failed to download the image.")
							return
						}
						msg := uploadHandler(url, apis[update.Message.From.ID], update)
						if _, err := bot.Send(msg); err != nil {
							log.Warnf("fail to send msg, %v", err)
						}
						return
					}
				}
				// Callback 删除图片
				if update.CallbackQuery != nil {
					client := smms.Client{}
					if update.CallbackQuery.Data != "err" {
						if _, err := client.Delete(update.CallbackQuery.Data); err != nil {
							log.Warnf("fail to delete the image, %v", err)
						}
					}
					edit := tgbotapi.EditMessageTextConfig{
						BaseEdit: tgbotapi.BaseEdit{
							ChatID:    int64(update.CallbackQuery.From.ID),
							MessageID: update.CallbackQuery.Message.MessageID,
						},
						Text: "Photo Deleted!",
					}
					if _, err := bot.Send(edit); err != nil {
						log.Warnf("fail to send msg, %v", err)
					}
				}
			}()
		}
	}
}
