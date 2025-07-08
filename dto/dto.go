package dto

// Auth DTOs
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// Event DTOs
type CreateEventRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Location    string  `json:"location" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Date        string  `json:"date" binding:"required"`
	EndDate     string  `json:"end_date"`
	Capacity    int     `json:"capacity" binding:"required,min=1"`
	Price       float64 `json:"price" binding:"required,min=0"`
}

type UpdateEventRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
	EndDate     string  `json:"end_date"`
	Capacity    int     `json:"capacity" binding:"min=1"`
	Price       float64 `json:"price" binding:"min=0"`
	Status      string  `json:"status"`
}

type EventResponse struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	Category    string  `json:"category"`
	Status      string  `json:"status"`
	Date        string  `json:"date"`
	EndDate     string  `json:"end_date"`
	Capacity    int     `json:"capacity"`
	Price       float64 `json:"price"`
	SoldTickets int     `json:"sold_tickets"`
	CreatedBy   uint    `json:"created_by"`
	IsActive    bool    `json:"is_active"`
}

// Ticket DTOs
type CreateTicketRequest struct {
	EventID  uint `json:"event_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,min=1"`
}

type CancelTicketRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type TicketResponse struct {
	ID           uint           `json:"id"`
	UserID       uint           `json:"user_id"`
	EventID      uint           `json:"event_id"`
	Quantity     int            `json:"quantity"`
	UnitPrice    float64        `json:"unit_price"`
	TotalPrice   float64        `json:"total_price"`
	Status       string         `json:"status"`
	BookingCode  string         `json:"booking_code"`
	PurchaseDate string         `json:"purchase_date"`
	CancelledAt  *string        `json:"cancelled_at,omitempty"`
	CancelReason string         `json:"cancel_reason,omitempty"`
	Event        *EventResponse `json:"event,omitempty"`
}

// Report DTOs
type SummaryReportResponse struct {
	TotalTickets int     `json:"total_tickets"`
	TotalRevenue float64 `json:"total_revenue"`
}

type EventReportResponse struct {
	EventID     uint    `json:"event_id"`
	Title       string  `json:"title"`
	TicketsSold int     `json:"tickets_sold"`
	Revenue     float64 `json:"revenue"`
}

// Comprehensive System Report DTOs
type SystemReportResponse struct {
	GeneratedAt       string                    `json:"generated_at"`
	SystemName        string                    `json:"system_name"`
	Overview          *SystemOverview           `json:"overview"`
	UserMetrics       *UserMetrics              `json:"user_metrics"`
	EventMetrics      *EventMetrics             `json:"event_metrics"`
	TicketMetrics     *TicketMetrics            `json:"ticket_metrics"`
	RevenueMetrics    *RevenueMetrics           `json:"revenue_metrics"`
	TopEvents         []TopEventReport          `json:"top_events"`
	CategoryBreakdown []CategoryBreakdownReport `json:"category_breakdown"`
	MonthlyStats      []MonthlyStatsReport      `json:"monthly_stats"`
}

type SystemOverview struct {
	TotalUsers       int64 `json:"total_users"`
	TotalEvents      int64 `json:"total_events"`
	TotalTickets     int64 `json:"total_tickets"`
	ActiveEvents     int64 `json:"active_events"`
	CompletedEvents  int64 `json:"completed_events"`
	CancelledEvents  int64 `json:"cancelled_events"`
	CancelledTickets int64 `json:"cancelled_tickets"`
}

type UserMetrics struct {
	TotalUsers    int64 `json:"total_users"`
	AdminUsers    int64 `json:"admin_users"`
	RegularUsers  int64 `json:"regular_users"`
	ActiveUsers   int64 `json:"active_users"`
	NewUsersMonth int64 `json:"new_users_this_month"`
}

type EventMetrics struct {
	TotalEvents     int64   `json:"total_events"`
	ActiveEvents    int64   `json:"active_events"`
	CompletedEvents int64   `json:"completed_events"`
	CancelledEvents int64   `json:"cancelled_events"`
	AverageCapacity float64 `json:"average_capacity"`
	AveragePrice    float64 `json:"average_price"`
}

type TicketMetrics struct {
	TotalTickets       int64   `json:"total_tickets"`
	BookedTickets      int64   `json:"booked_tickets"`
	CancelledTickets   int64   `json:"cancelled_tickets"`
	UsedTickets        int64   `json:"used_tickets"`
	AverageTicketPrice float64 `json:"average_ticket_price"`
}

type RevenueMetrics struct {
	TotalRevenue   float64 `json:"total_revenue"`
	MonthlyRevenue float64 `json:"monthly_revenue"`
	AverageRevenue float64 `json:"average_revenue_per_event"`
	RefundedAmount float64 `json:"refunded_amount"`
}

type TopEventReport struct {
	EventID     uint    `json:"event_id"`
	Title       string  `json:"title"`
	TicketsSold int     `json:"tickets_sold"`
	Revenue     float64 `json:"revenue"`
	Category    string  `json:"category"`
}

type CategoryBreakdownReport struct {
	Category    string  `json:"category"`
	EventCount  int     `json:"event_count"`
	TicketsSold int     `json:"tickets_sold"`
	Revenue     float64 `json:"revenue"`
}

type MonthlyStatsReport struct {
	Month    string  `json:"month"`
	Events   int     `json:"events"`
	Tickets  int     `json:"tickets"`
	Revenue  float64 `json:"revenue"`
	NewUsers int     `json:"new_users"`
}
