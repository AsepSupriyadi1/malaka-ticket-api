package services

import (
	"case_study_api/dto"
	"case_study_api/entities"
	"case_study_api/repositories"
	"case_study_api/utils"
	"errors"
	"fmt"
	"time"
)

type TicketService interface {
	GetUserTickets(userID uint) ([]dto.TicketResponse, error)
	GetUserTicketsPaginated(userID uint, pagination utils.PaginationRequest) (*utils.PaginationResponse, error)
	GetByID(ticketID uint) (*dto.TicketResponse, error)
	BookTicket(req dto.CreateTicketRequest, userID uint) (*dto.TicketResponse, error)
	CancelTicket(ticketID uint, userID uint, req dto.CancelTicketRequest) error
}

type ticketService struct {
	ticketRepo repositories.TicketRepository
	eventRepo  repositories.EventRepository
}

func NewTicketService(ticketRepo repositories.TicketRepository, eventRepo repositories.EventRepository) TicketService {
	return &ticketService{
		ticketRepo: ticketRepo,
		eventRepo:  eventRepo,
	}
}

func (s *ticketService) GetUserTickets(userID uint) ([]dto.TicketResponse, error) {
	tickets, err := s.ticketRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var ticketResponses []dto.TicketResponse
	for _, ticket := range tickets {
		ticketResponses = append(ticketResponses, s.entityToResponse(ticket))
	}

	return ticketResponses, nil
}

func (s *ticketService) GetUserTicketsPaginated(userID uint, pagination utils.PaginationRequest) (*utils.PaginationResponse, error) {
	tickets, total, err := s.ticketRepo.GetByUserIDPaginated(userID, pagination.Offset, pagination.PageSize)
	if err != nil {
		return nil, err
	}

	var ticketResponses []dto.TicketResponse
	for _, ticket := range tickets {
		ticketResponses = append(ticketResponses, s.entityToResponse(ticket))
	}

	response := utils.BuildPaginationResponse(ticketResponses, total, pagination)
	return &response, nil
}

func (s *ticketService) GetByID(ticketID uint) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.GetByID(ticketID)
	if err != nil {
		return nil, err
	}

	response := s.entityToResponse(*ticket)
	return &response, nil
}

func (s *ticketService) BookTicket(req dto.CreateTicketRequest, userID uint) (*dto.TicketResponse, error) {
	// Get event details
	event, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	// Check event availability
	if !event.IsActive || time.Now().After(event.EndDate) {
		return nil, errors.New("event is not available")
	}

	// Check ticket availability
	if event.Capacity-event.SoldTickets < req.Quantity {
		return nil, fmt.Errorf("only %d tickets left", event.Capacity-event.SoldTickets)
	}

	// Create ticket
	total := float64(req.Quantity) * event.Price
	ticket := entities.Ticket{
		UserID:       userID,
		EventID:      event.ID,
		Quantity:     req.Quantity,
		UnitPrice:    event.Price,
		TotalPrice:   total,
		BookingCode:  fmt.Sprintf("BK%v%v", userID, time.Now().UnixNano()%100000),
		PurchaseDate: time.Now(),
		Status:       "booked",
	}

	// Save ticket
	if err := s.ticketRepo.Create(&ticket); err != nil {
		return nil, err
	}

	// Update event sold tickets
	event.SoldTickets += req.Quantity
	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}

	response := s.entityToResponse(ticket)
	return &response, nil
}

func (s *ticketService) CancelTicket(ticketID uint, userID uint, req dto.CancelTicketRequest) error {
	// Get ticket
	ticket, err := s.ticketRepo.GetByID(ticketID)
	if err != nil {
		return err
	}

	// Check ownership
	if ticket.UserID != userID {
		return errors.New("unauthorized access")
	}

	// Check status
	if ticket.Status != "booked" {
		return errors.New("only booked tickets can be cancelled")
	}

	// Update ticket status
	now := time.Now()
	ticket.Status = "cancelled"
	ticket.CancelledAt = &now
	ticket.CancelReason = req.Reason

	// Update event sold tickets
	event, err := s.eventRepo.GetByID(ticket.EventID)
	if err == nil && event.SoldTickets >= ticket.Quantity {
		event.SoldTickets -= ticket.Quantity
		s.eventRepo.Update(event)
	}

	return s.ticketRepo.Update(ticket)
}

func (s *ticketService) entityToResponse(ticket entities.Ticket) dto.TicketResponse {
	response := dto.TicketResponse{
		ID:           ticket.ID,
		UserID:       ticket.UserID,
		EventID:      ticket.EventID,
		Quantity:     ticket.Quantity,
		UnitPrice:    ticket.UnitPrice,
		TotalPrice:   ticket.TotalPrice,
		Status:       ticket.Status,
		BookingCode:  ticket.BookingCode,
		PurchaseDate: ticket.PurchaseDate.Format("2006-01-02T15:04:05Z"),
		CancelReason: ticket.CancelReason,
	}

	if ticket.CancelledAt != nil {
		cancelledAt := ticket.CancelledAt.Format("2006-01-02T15:04:05Z")
		response.CancelledAt = &cancelledAt
	}

	// Include event details if available
	if ticket.Event.ID != 0 {
		eventResponse := dto.EventResponse{
			ID:          ticket.Event.ID,
			Title:       ticket.Event.Title,
			Description: ticket.Event.Description,
			Location:    ticket.Event.Location,
			Category:    ticket.Event.Category,
			Status:      ticket.Event.Status,
			Date:        ticket.Event.Date.Format("2006-01-02T15:04:05Z"),
			EndDate:     ticket.Event.EndDate.Format("2006-01-02T15:04:05Z"),
			Capacity:    ticket.Event.Capacity,
			Price:       ticket.Event.Price,
			SoldTickets: ticket.Event.SoldTickets,
			CreatedBy:   ticket.Event.CreatedBy,
			IsActive:    ticket.Event.IsActive,
		}
		response.Event = &eventResponse
	}

	return response
}
