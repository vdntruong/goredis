package ticket

import (
	"time"

	"github.com/shopspring/decimal"
)

type DTORoom struct {
	ExtID  string
	Number string
	Seats  uint
}

type DTORoomCreate struct {
	Number string
	Seats  uint
}

type DTOFilm struct {
	ExtID       string
	Name        string
	ReleaseDate time.Time
}

type DTOFilmCreate struct {
	Name        string
	ReleaseDate time.Time
}

type DTOTicket struct {
	ExtID string

	Date       time.Time
	Price      decimal.Decimal
	SeatNumber uint

	ReceiverEmail *string
}
