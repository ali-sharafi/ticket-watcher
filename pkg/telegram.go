package telegram

import (
	"fmt"
	"os"
	"strings"
	"ticket-watcher/domain"
	"ticket-watcher/pkg/app"
	utils "ticket-watcher/pkg/utils"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var tempUserInput = make(map[int64]domain.Travel)
var userMessageIDs = make(map[int64]int)

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
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			fmt.Println("callbackData: ", callbackData)
			chatID := update.CallbackQuery.Message.Chat.ID
			userID := update.CallbackQuery.From.ID
			messageID := userMessageIDs[userID]
			input := tempUserInput[userID]
			parts := strings.SplitN(callbackData, "_", 2)
			var msg tgbotapi.EditMessageTextConfig

			if len(parts) == 2 {
				prefix := parts[0]
				data := parts[1]

				switch prefix {
				case "origin":
					input.Origin = data
					keyboard := createProvinceKeyboard("dest")
					msg = tgbotapi.NewEditMessageText(chatID, messageID, "Pls select the destination:")
					msg.ReplyMarkup = &keyboard
				case "dest":
					input.Destination = data
					keyboard := createCalendarInlineKeyboard()
					msg = tgbotapi.NewEditMessageText(chatID, messageID, "Pls select the date:")
					msg.ReplyMarkup = &keyboard
				case "date":
					input.Date = data
					keyboard := createTravelTypeKeyboard()
					msg = tgbotapi.NewEditMessageText(chatID, messageID, "Pls select travel type:")
					msg.ReplyMarkup = &keyboard
				case "type":
					input.Type = data
					keyboard := createConfirmKeyboard()
					msg = tgbotapi.NewEditMessageText(chatID, messageID, getConfirmationText(input))
					msg.ReplyMarkup = &keyboard
				case "confirmation":
					if data == "1" {
						storeUserData(input, userID)
						msg = tgbotapi.NewEditMessageText(chatID, messageID, "Your data has been saved successfully")
					} else {
						msg = tgbotapi.NewEditMessageText(chatID, messageID, "Your data has not been saved!")
					}
				default:
					defaultMsg := tgbotapi.NewMessage(chatID, "Invalid command")
					bot.Send(defaultMsg)
				}
				bot.Send(msg)
				tempUserInput[userID] = input
			}
			return
		}
		if update.Message == nil {
			continue
		}
		userID := update.Message.From.ID
		switch update.Message.Command() {
		case "start":
			initialUserData(userID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Pls select the origin:")
			msg.ReplyMarkup = createProvinceKeyboard("origin")
			sentMessage, err := bot.Send(msg)
			if err == nil {
				userMessageIDs[userID] = sentMessage.MessageID
			}
		case "delete":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Delete command executed")
			bot.Send(msg)
		case "complete":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Complete command executed")
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
			bot.Send(msg)
		}
	}
}

func storeUserData(travel domain.Travel, userID int64) {
	travels := utils.ReadTravelsData()
	travel.ID = utils.GenerateUniqueID()
	travels = append(travels, travel)
	delete(tempUserInput, userID)
	utils.StoreTravelsData(travels)
}

func getConfirmationText(input domain.Travel) string {
	return fmt.Sprintf("You want to save a travel from %s to %s on %s do you confirm?", input.Origin, input.Destination, utils.GetJalaliDate(input.Date))
}

func createTravelTypeKeyboard() tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Train", fmt.Sprintf("type_%s", "train")))
	keyboard = append(keyboard, row)
	row = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Flight", fmt.Sprintf("type_%s", "flight")))
	keyboard = append(keyboard, row)
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

func createConfirmKeyboard() tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Yes", fmt.Sprintf("confirmation_%d", 1)))
	keyboard = append(keyboard, row)
	row = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("No", fmt.Sprintf("confirmation_%d", 0)))
	keyboard = append(keyboard, row)
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

func createCalendarInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	today := time.Now()
	for day := 0; day < 30; day++ {
		currentDate := today.AddDate(0, 0, day)
		dateText := utils.GetJalaliDate(currentDate.Format("2006-01-02"))
		dateCallback := currentDate.Format("2006-01-02")

		row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(dateText, fmt.Sprintf("date_%s", dateCallback)))
		keyboard = append(keyboard, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

func createProvinceKeyboard(prefix string) tgbotapi.InlineKeyboardMarkup {
	columns := 5
	provinces := app.GetProvinces()

	numRows := (len(provinces) + columns - 1) / columns

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < numRows; i++ {
		var row []tgbotapi.InlineKeyboardButton
		for j := 0; j < columns; j++ {
			index := i*columns + j
			if index < len(provinces) {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(provinces[index].Name, fmt.Sprintf("%s_%s", prefix, provinces[index].Code)))
			}
		}
		keyboard = append(keyboard, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

func initialUserData(userID int64) {
	if _, exists := tempUserInput[userID]; exists {
		delete(tempUserInput, userID)
	}
	tempUserInput[userID] = domain.Travel{}
}
