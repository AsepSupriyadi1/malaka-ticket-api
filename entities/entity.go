package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string   `gorm:"type:varchar(255);not null"`
	Email    string   `gorm:"unique;not null;type:varchar(255)"`
	Password string   `gorm:"type:varchar(255);not null"`
	Role     string   `gorm:"type:enum('user','admin');default:'user'"`
	Events   []Event  `gorm:"foreignKey:CreatedBy"`
	Tickets  []Ticket `gorm:"foreignKey:UserID"`
}

type Event struct {
	gorm.Model
	Title       string    `gorm:"unique;not null;type:varchar(255)"`
	Description string    `gorm:"type:text"`
	Location    string    `gorm:"type:varchar(255)"`
	Category    string    `gorm:"type:varchar(100)"`
	Status      string    `gorm:"type:enum('upcoming','ongoing','completed','cancelled');default:'upcoming'"`
	Date        time.Time `gorm:"not null"`
	EndDate     time.Time
	Capacity    int     `gorm:"not null;check:capacity > 0"`
	Price       float64 `gorm:"type:decimal(10,2);not null;check:price >= 0"`
	SoldTickets int     `gorm:"default:0;check:sold_tickets >= 0"`
	CreatedBy   uint    `gorm:"not null"`
	IsActive    bool    `gorm:"default:true"`

	User    User     `gorm:"foreignKey:CreatedBy"`
	Tickets []Ticket `gorm:"foreignKey:EventID"`
}

type Ticket struct {
	gorm.Model
	UserID       uint
	EventID      uint
	Quantity     int       `gorm:"not null;default:1;check:quantity > 0"`
	UnitPrice    float64   `gorm:"type:decimal(10,2);not null"`
	TotalPrice   float64   `gorm:"type:decimal(10,2);not null"`
	Status       string    `gorm:"type:enum('booked','cancelled','used');default:'booked'"`
	BookingCode  string    `gorm:"unique;not null;type:varchar(20)"`
	PurchaseDate time.Time `gorm:"not null"`
	CancelledAt  *time.Time
	CancelReason string `gorm:"type:text"`

	User  User  `gorm:"foreignKey:UserID"`
	Event Event `gorm:"foreignKey:EventID"`
}
