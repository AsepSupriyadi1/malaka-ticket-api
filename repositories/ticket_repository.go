package repositories

import (
	"case_study_api/entities"

	"gorm.io/gorm"
)

type TicketRepository interface {
	GetByUserID(userID uint) ([]entities.Ticket, error)
	GetByUserIDPaginated(userID uint, offset, limit int) ([]entities.Ticket, int64, error)
	GetByID(id uint) (*entities.Ticket, error)
	Create(ticket *entities.Ticket) error
	Update(ticket *entities.Ticket) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) GetByUserID(userID uint) ([]entities.Ticket, error) {
	var tickets []entities.Ticket
	err := r.db.Preload("Event").Where("user_id = ?", userID).Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) GetByUserIDPaginated(userID uint, offset, limit int) ([]entities.Ticket, int64, error) {
	var tickets []entities.Ticket
	var total int64

	// Get total count
	if err := r.db.Model(&entities.Ticket{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Preload("Event").Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&tickets).Error
	return tickets, total, err
}

func (r *ticketRepository) GetByID(id uint) (*entities.Ticket, error) {
	var ticket entities.Ticket
	err := r.db.Preload("Event").First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) Create(ticket *entities.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *ticketRepository) Update(ticket *entities.Ticket) error {
	return r.db.Save(ticket).Error
}
