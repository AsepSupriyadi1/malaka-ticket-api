package repositories

import (
	"case_study_api/entities"

	"gorm.io/gorm"
)

type EventRepository interface {
	GetAll() ([]entities.Event, error)
	GetAllPaginated(offset, limit int) ([]entities.Event, int64, error)
	GetByID(id uint) (*entities.Event, error)
	Create(event *entities.Event) error
	Update(event *entities.Event) error
	Delete(event *entities.Event) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) GetAll() ([]entities.Event, error) {
	var events []entities.Event
	err := r.db.Find(&events).Error
	return events, err
}

func (r *eventRepository) GetAllPaginated(offset, limit int) ([]entities.Event, int64, error) {
	var events []entities.Event
	var total int64

	// Get total count
	if err := r.db.Model(&entities.Event{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Offset(offset).Limit(limit).Find(&events).Error
	return events, total, err
}

func (r *eventRepository) GetByID(id uint) (*entities.Event, error) {
	var event entities.Event
	err := r.db.First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) Create(event *entities.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) Update(event *entities.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(event *entities.Event) error {
	return r.db.Delete(event).Error
}
