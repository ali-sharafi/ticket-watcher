package telegram

import (
	"fmt"
	"os"
	"ticket-watcher/domain"
	utils "ticket-watcher/pkg/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func SetupTelegramBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		utils.Logger.Error(err)
		return
	}
	setupCommandListener()
}

func Notify(payload domain.NotifPayload, travelId string) {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("View Ticket", payload.Link),
			tgbotapi.NewInlineKeyboardButtonData("Complete", fmt.Sprintf("complete|%s", travelId)),
		),
	)
	msg := tgbotapi.NewMessageToChannel(os.Getenv("CHANNEL_NAME"), payload.Message)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = inlineKeyboard
	_, err := bot.Send(msg)

	if err != nil {
		utils.Logger.Error("Error sending message to Telegram:", err)
	}
	utils.Logger.Info("notification message:", payload.Message)
	utils.Logger.Info("notification link:", payload.Link)
}

func setupCommandListener() {

}
