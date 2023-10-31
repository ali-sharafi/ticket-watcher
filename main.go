package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/erfanmomeniii/jalali"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Travel struct {
	ID          string `json:"id"`
	Origin      string `json:"origin"`
	Destination string `json:"dest"`
	Type        string `json:"type"`
	Date        string `json:"date"`
}

type NotifPayload struct {
	Message string
	Link    string
}

type AlibabaTokenResponse struct {
	Result struct {
		RequestID string `json:"requestId"`
	} `json:"result"`
}

type AlibabaTripsResponse struct {
	Result struct {
		Departing []AlibabaDepartingItem `json:"departing"`
	} `json:"result"`
}

type AlibabaDepartingItem struct {
	Seat          int    `json:"seat"`
	DepartureDate string `json:"departureDateTime"`
}

type AlibabaFlightTripsResponse struct {
	Result struct {
		Departing []AlibabaDepartingFlightItem `json:"departing"`
	} `json:"result"`
}

type AlibabaDepartingFlightItem struct {
	Seat          int    `json:"seat"`
	LeaveDateTime string `json:"leaveDateTime"`
}

var bot *tgbotapi.BotAPI

func main() {
	SetupLogger()
	setupEnv()
	setupTelegramBot()

	run()
}

func run() {
	ticker := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-ticker.C:
			travels := readTravelsData()
			for _, travel := range travels {
				alibaba(travel)
			}
		}
	}
}

func readTravelsData() (travels []Travel) {
	jsonData, err := ioutil.ReadFile("data.json")
	if err != nil {
		logger.Error("Error reading file:", err)
		return
	}

	err = json.Unmarshal(jsonData, &travels)
	if err != nil {
		logger.Error("Error decoding JSON:", err)
	}
	return
}

func setupTelegramBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		logger.Error(err)
	}
}

func setupEnv() {
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file")
	}
}

func SetupLogger() {

	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.ToSlash("./logs/" + time.Now().Format("2006-01-02") + "_app.log"),
		MaxSize:    1, // MB
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	// Fork writing into two outputs
	multiWriter := io.MultiWriter(os.Stderr, lumberjackLogger)

	logFormatter := new(logger.TextFormatter)
	logFormatter.TimestampFormat = time.RFC3339
	logFormatter.FullTimestamp = true

	logger.SetFormatter(logFormatter)
	logger.SetLevel(logger.InfoLevel)
	logger.SetOutput(multiWriter)
}

func notify(payload NotifPayload, travelId string) {
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
		logger.Error("Error sending message to Telegram:", err)
	}
	logger.Info("notification message:", payload.Message)
	logger.Info("notification link:", payload.Link)
}

func alibaba(travel Travel) {
	switch travel.Type {
	case "train":
		checkTrainTickets(travel)
	case "flight":
		checkFlightTickets(travel)
	}
}

func checkTrainTickets(travel Travel) {
	token, _ := getTrainToken(travel)
	tickets, _ := getTrainTrips(token)

	if result, trip := isTrainTicketAvailable(tickets); result != false {
		notifPayload := NotifPayload{
			Message: fmt.Sprintf(`Train Ticket found from: %s To %s on %s`, travel.Origin, travel.Destination, trip.DepartureDate),
			Link:    fmt.Sprintf(`https://www.alibaba.ir/train/%s-%s?adult=1&child=0&ticketType=Family&isExclusive=false&infant=0&departing=%s`, travel.Origin, travel.Destination, getJalaliDate(travel.Date)),
		}
		notify(notifPayload, travel.ID)
	} else {
		logger.Info(fmt.Sprintf("There is not any trips for train from %s to %s on %s in alibaba", travel.Origin, travel.Destination, travel.Date))
	}
}

func checkFlightTickets(travel Travel) {
	token, _ := getFlightToken(travel)
	tickets, _ := getFlightTrips(token)

	if result, trip := isFlightTicketAvailable(tickets); result != false {
		notifPayload := NotifPayload{
			Message: fmt.Sprintf(`Flight Ticket found from %s To %s on %s`, travel.Origin, travel.Destination, trip.LeaveDateTime),
			Link:    fmt.Sprintf(`https://www.alibaba.ir/flights/%s-%s?adult=1&child=0&infant=0&departing=%s`, travel.Origin, travel.Destination, getJalaliDate(travel.Date)),
		}
		notify(notifPayload, travel.ID)
	} else {
		logger.Info(fmt.Sprintf("There is not any trips for flight from %s to %s on %s in alibaba", travel.Origin, travel.Destination, travel.Date))
	}
}

func getFlightTrips(token string) (tripsResponse AlibabaFlightTripsResponse, err error) {
	url := fmt.Sprintf(`https://ws.alibaba.ir/api/v1/flights/domestic/available/%s`, token)
	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		Get(url)

	if err != nil {
		logger.Error("Error in get alibaba trips:", err)
		return
	}
	if err = json.Unmarshal(response.Body(), &tripsResponse); err != nil {
		logger.Error("Error parsing getTrainTrips API response:", err)
		return
	}

	return
}

func getFlightToken(travel Travel) (token string, err error) {
	payload := fmt.Sprintf(`{"departureDate": "%s", "destination": "%s", "origin": "%s", "adult": 1, "child": 0, "infant": 0}`, travel.Date, travel.Destination, travel.Origin)

	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("https://ws.alibaba.ir/api/v1/flights/domestic/available")

	if err != nil {
		logger.Error("Error get alibaba flight token:", err)
		return
	}

	var tokenResponse AlibabaTokenResponse

	if err = json.Unmarshal(response.Body(), &tokenResponse); err != nil {
		logger.Error("Error getFlightToken parsing API response:", err)
		return
	}

	token = tokenResponse.Result.RequestID
	return
}

func isFlightTicketAvailable(trips AlibabaFlightTripsResponse) (bool, AlibabaDepartingFlightItem) {
	for _, trip := range trips.Result.Departing {
		if trip.Seat > 0 {
			return true, trip
		}
	}
	return false, AlibabaDepartingFlightItem{}
}

func isTrainTicketAvailable(trips AlibabaTripsResponse) (bool, AlibabaDepartingItem) {
	for _, trip := range trips.Result.Departing {
		if trip.Seat > 0 {
			return true, trip
		}
	}
	return false, AlibabaDepartingItem{}
}

func getTrainTrips(token string) (tripsResponse AlibabaTripsResponse, err error) {
	url := fmt.Sprintf(`https://ws.alibaba.ir/api/v2/train/available/%s`, token)
	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		Get(url)

	if err != nil {
		logger.Error("Error in get alibaba trips:", err)
		return
	}
	if err = json.Unmarshal(response.Body(), &tripsResponse); err != nil {
		logger.Error("Error parsing getTrainTrips API response:", err)
		return
	}

	return
}

func getTrainToken(travel Travel) (token string, err error) {
	payload := fmt.Sprintf(`{"departureDate": "%s", "destination": "%s", "origin": "%s", "isExclusiveCompartment": false, "passengerCount": 1, "ticketType": "Family"}`, travel.Date, travel.Destination, travel.Origin)

	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("https://ws.alibaba.ir/api/v2/train/available")

	if err != nil {
		logger.Error("Error get alibaba token:", err)
		return
	}

	var tokenResponse AlibabaTokenResponse

	if err = json.Unmarshal(response.Body(), &tokenResponse); err != nil {
		logger.Error("Error getTrainToken parsing API response:", err)
		return
	}

	token = tokenResponse.Result.RequestID
	return
}

func getJalaliDate(gregorianDate string) string {
	t, err := time.Parse("2006-01-02", gregorianDate)
	if err != nil {
		logger.Error("Error:", err)
		return ""
	}

	// Convert to Jalali
	j := jalali.ConvertGregorianToJalali(t)

	jalaliDate := fmt.Sprintf(`%d-%d-%d`, j.Year(), j.Month(), j.Day())

	return jalaliDate
}
