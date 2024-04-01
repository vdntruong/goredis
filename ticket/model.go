package ticket

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Film struct {
	gorm.Model

	ExtID string `gorm:"uniqueIndex"`
	Name  string `gorm:"uniqueIndex"`

	ReleaseDate time.Time
}

type Room struct {
	gorm.Model

	ExtID  string `gorm:"uniqueIndex"`
	Number string `gorm:"uniqueIndex"`

	Seats uint
}

type Ticket struct {
	gorm.Model

	ExtID string `gorm:"uniqueIndex"`
	Date  time.Time

	Price      decimal.Decimal
	SeatNumber uint

	ReceiverEmail *string `gorm:"index"`
}
