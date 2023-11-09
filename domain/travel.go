package domain

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

type Province struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
