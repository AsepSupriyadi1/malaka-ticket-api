package services

import (
	"case_study_api/dto"
	"case_study_api/entities"
	"case_study_api/repositories"
	"case_study_api/utils"
	"errors"
	"time"
)

type EventService interface {
	GetAll() ([]dto.EventResponse, error)
	GetAllPaginated(pagination utils.PaginationRequest) (*utils.PaginationResponse, error)
	GetByID(id uint) (*dto.EventResponse, error)
	Create(req dto.CreateEventRequest, createdBy uint) (*dto.EventResponse, error)
	Update(id uint, req dto.UpdateEventRequest) (*dto.EventResponse, error)
	Delete(id uint) error
}

type eventService struct {
	eventRepo repositories.EventRepository
}

func NewEventService(eventRepo repositories.EventRepository) EventService {
	return &eventService{
		eventRepo: eventRepo,
	}
}

func (s *eventService) GetAll() ([]dto.EventResponse, error) {
	events, err := s.eventRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var eventResponses []dto.EventResponse
	for _, event := range events {
		eventResponses = append(eventResponses, s.entityToResponse(event))
	}

	return eventResponses, nil
}

func (s *eventService) GetAllPaginated(pagination utils.PaginationRequest) (*utils.PaginationResponse, error) {
	events, total, err := s.eventRepo.GetAllPaginated(pagination.Offset, pagination.PageSize)
	if err != nil {
		return nil, err
	}

	var eventResponses []dto.EventResponse
	for _, event := range events {
		eventResponses = append(eventResponses, s.entityToResponse(event))
	}

	response := utils.BuildPaginationResponse(eventResponses, total, pagination)
	return &response, nil
}

func (s *eventService) GetByID(id uint) (*dto.EventResponse, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := s.entityToResponse(*event)
	return &response, nil
}

func (s *eventService) Create(req dto.CreateEventRequest, createdBy uint) (*dto.EventResponse, error) {
	// Parse dates
	date, err := time.Parse("2006-01-02T15:04:05Z", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	var endDate time.Time
	if req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02T15:04:05Z", req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end date format")
		}
		if endDate.Before(date) {
			return nil, errors.New("end date must be after start date")
		}
	}

	event := entities.Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		Category:    req.Category,
		Date:        date,
		EndDate:     endDate,
		Capacity:    req.Capacity,
		Price:       req.Price,
		CreatedBy:   createdBy,
		Status:      "upcoming",
		IsActive:    true,
	}

	if err := s.eventRepo.Create(&event); err != nil {
		return nil, err
	}

	response := s.entityToResponse(event)
	return &response, nil
}

func (s *eventService) Update(id uint, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if time.Now().After(event.Date) {
		return nil, errors.New("cannot update past events")
	}

	// Update fields if provided
	if req.Title != "" {
		event.Title = req.Title
	}
	if req.Description != "" {
		event.Description = req.Description
	}
	if req.Location != "" {
		event.Location = req.Location
	}
	if req.Category != "" {
		event.Category = req.Category
	}
	if req.Date != "" {
		date, err := time.Parse("2006-01-02T15:04:05Z", req.Date)
		if err != nil {
			return nil, errors.New("invalid date format")
		}
		event.Date = date
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02T15:04:05Z", req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end date format")
		}
		event.EndDate = endDate
	}
	if req.Capacity > 0 {
		event.Capacity = req.Capacity
	}
	if req.Price >= 0 {
		event.Price = req.Price
	}
	if req.Status != "" {
		event.Status = req.Status
	}

	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}

	response := s.entityToResponse(*event)
	return &response, nil
}

func (s *eventService) Delete(id uint) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		return err
	}

	if event.SoldTickets > 0 {
		return errors.New("cannot delete event with sold tickets")
	}

	return s.eventRepo.Delete(event)
}

func (s *eventService) entityToResponse(event entities.Event) dto.EventResponse {
	return dto.EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		Category:    event.Category,
		Status:      event.Status,
		Date:        event.Date.Format("2006-01-02T15:04:05Z"),
		EndDate:     event.EndDate.Format("2006-01-02T15:04:05Z"),
		Capacity:    event.Capacity,
		Price:       event.Price,
		SoldTickets: event.SoldTickets,
		CreatedBy:   event.CreatedBy,
		IsActive:    event.IsActive,
	}
}
