package main

import (
	"encoding/json"
	"fmt"
	"time"

	"ticket-watcher/domain"
	telegram "ticket-watcher/pkg"

	utils "ticket-watcher/pkg/utils"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

var logger = utils.Logger

func main() {
	setupEnv()
	telegram.SetupTelegramBot()

	run()
}

func run() {
	ticker := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-ticker.C:
			travels := utils.ReadTravelsData()
			for _, travel := range travels {
				alibaba(travel)
				time.Sleep(30 * time.Second)
			}
		}
	}
}

func setupEnv() {
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file")
	}
}

func alibaba(travel domain.Travel) {
	switch travel.Type {
	case "train":
		checkTrainTickets(travel)
	case "flight":
		checkFlightTickets(travel)
	}
}

func checkTrainTickets(travel domain.Travel) {
	token, _ := getTrainToken(travel)
	tickets, _ := getTrainTrips(token)

	if result, trip := isTrainTicketAvailable(tickets); result != false {
		notifPayload := domain.NotifPayload{
			Message: fmt.Sprintf(`Train Ticket found from: %s To %s on %s`, travel.Origin, travel.Destination, trip.DepartureDate),
			Link:    fmt.Sprintf(`https://www.alibaba.ir/train/%s-%s?adult=1&child=0&ticketType=Family&isExclusive=false&infant=0&departing=%s`, travel.Origin, travel.Destination, utils.GetJalaliDate(travel.Date)),
		}
		telegram.Notify(notifPayload, travel.ID)
	} else {
		logger.Info(fmt.Sprintf("There is not any trips for train from %s to %s on %s in alibaba", travel.Origin, travel.Destination, travel.Date))
	}
}

func checkFlightTickets(travel domain.Travel) {
	token, _ := getFlightToken(travel)
	tickets, _ := getFlightTrips(token)

	if result, trip := isFlightTicketAvailable(tickets); result != false {
		notifPayload := domain.NotifPayload{
			Message: fmt.Sprintf(`Flight Ticket found from %s To %s on %s`, travel.Origin, travel.Destination, trip.LeaveDateTime),
			Link:    fmt.Sprintf(`https://www.alibaba.ir/flights/%s-%s?adult=1&child=0&infant=0&departing=%s`, travel.Origin, travel.Destination, utils.GetJalaliDate(travel.Date)),
		}
		telegram.Notify(notifPayload, travel.ID)
	} else {
		logger.Info(fmt.Sprintf("There is not any trips for flight from %s to %s on %s in alibaba", travel.Origin, travel.Destination, travel.Date))
	}
}

func getFlightTrips(token string) (tripsResponse domain.AlibabaFlightTripsResponse, err error) {
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

func getFlightToken(travel domain.Travel) (token string, err error) {
	payload := fmt.Sprintf(`{"departureDate": "%s", "destination": "%s", "origin": "%s", "adult": 1, "child": 0, "infant": 0}`, travel.Date, travel.Destination, travel.Origin)

	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("https://ws.alibaba.ir/api/v1/flights/domestic/available")

	if err != nil {
		logger.Error("Error get alibaba flight token:", err)
		return
	}

	var tokenResponse domain.AlibabaTokenResponse

	if err = json.Unmarshal(response.Body(), &tokenResponse); err != nil {
		logger.Error("Error getFlightToken parsing API response:", err)
		return
	}

	token = tokenResponse.Result.RequestID
	return
}

func isFlightTicketAvailable(trips domain.AlibabaFlightTripsResponse) (bool, domain.AlibabaDepartingFlightItem) {
	for _, trip := range trips.Result.Departing {
		if trip.Seat > 0 {
			return true, trip
		}
	}
	return false, domain.AlibabaDepartingFlightItem{}
}

func isTrainTicketAvailable(trips domain.AlibabaTripsResponse) (bool, domain.AlibabaDepartingItem) {
	for _, trip := range trips.Result.Departing {
		if trip.Seat > 0 {
			return true, trip
		}
	}
	return false, domain.AlibabaDepartingItem{}
}

func getTrainTrips(token string) (tripsResponse domain.AlibabaTripsResponse, err error) {
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

func getTrainToken(travel domain.Travel) (token string, err error) {
	payload := fmt.Sprintf(`{"departureDate": "%s", "destination": "%s", "origin": "%s", "isExclusiveCompartment": false, "passengerCount": 1, "ticketType": "Family"}`, travel.Date, travel.Destination, travel.Origin)

	response, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("https://ws.alibaba.ir/api/v2/train/available")

	if err != nil {
		logger.Error("Error get alibaba token:", err)
		return
	}

	var tokenResponse domain.AlibabaTokenResponse

	if err = json.Unmarshal(response.Body(), &tokenResponse); err != nil {
		logger.Error("Error getTrainToken parsing API response:", err)
		return
	}

	token = tokenResponse.Result.RequestID
	return
}
