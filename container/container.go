package container

import (
	"case_study_api/repositories"
	"case_study_api/services"

	"gorm.io/gorm"
)

type Container struct {
	DB *gorm.DB

	// Repositories
	UserRepo   repositories.UserRepository
	EventRepo  repositories.EventRepository
	TicketRepo repositories.TicketRepository
	ReportRepo repositories.ReportRepository

	// Services
	AuthService   services.AuthService
	EventService  services.EventService
	TicketService services.TicketService
	ReportService services.ReportService
}

func NewContainer(db *gorm.DB) *Container {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	eventRepo := repositories.NewEventRepository(db)
	ticketRepo := repositories.NewTicketRepository(db)
	reportRepo := repositories.NewReportRepository(db)

	// Initialize services with dependency injection
	authService := services.NewAuthService(userRepo)
	eventService := services.NewEventService(eventRepo)
	ticketService := services.NewTicketService(ticketRepo, eventRepo)
	reportService := services.NewReportService(reportRepo)

	return &Container{
		DB:            db,
		UserRepo:      userRepo,
		EventRepo:     eventRepo,
		TicketRepo:    ticketRepo,
		ReportRepo:    reportRepo,
		AuthService:   authService,
		EventService:  eventService,
		TicketService: ticketService,
		ReportService: reportService,
	}
}
