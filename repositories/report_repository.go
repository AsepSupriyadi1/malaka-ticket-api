package repositories

import (
	"case_study_api/entities"
	"time"

	"gorm.io/gorm"
)

type SummaryReport struct {
	TotalTickets int     `json:"total_tickets"`
	TotalRevenue float64 `json:"total_revenue"`
}

type EventReport struct {
	EventID     uint    `json:"event_id"`
	Title       string  `json:"title"`
	TicketsSold int     `json:"tickets_sold"`
	Revenue     float64 `json:"revenue"`
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

type ReportRepository interface {
	GetSummaryReport() (*SummaryReport, error)
	GetEventReport(eventID uint) (*EventReport, error)
	GetSystemOverview() (*SystemOverview, error)
	GetUserMetrics() (*UserMetrics, error)
	GetEventMetrics() (*EventMetrics, error)
	GetTicketMetrics() (*TicketMetrics, error)
	GetRevenueMetrics() (*RevenueMetrics, error)
	GetTopEvents(limit int) ([]TopEventReport, error)
	GetCategoryBreakdown() ([]CategoryBreakdownReport, error)
	GetMonthlyStats() ([]MonthlyStatsReport, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetSummaryReport() (*SummaryReport, error) {
	var result SummaryReport
	err := r.db.Model(&entities.Ticket{}).
		Where("status = ?", "booked").
		Select("SUM(quantity) as total_tickets, SUM(total_price) as total_revenue").
		Scan(&result).Error

	return &result, err
}

func (r *reportRepository) GetEventReport(eventID uint) (*EventReport, error) {
	var result EventReport

	err := r.db.Table("tickets").
		Select("tickets.event_id, events.title, SUM(tickets.quantity) as tickets_sold, SUM(tickets.total_price) as revenue").
		Joins("JOIN events ON tickets.event_id = events.id").
		Where("tickets.event_id = ? AND tickets.status = ?", eventID, "booked").
		Group("tickets.event_id, events.title").
		Scan(&result).Error

	return &result, err
}

func (r *reportRepository) GetSystemOverview() (*SystemOverview, error) {
	var result SystemOverview

	// Get total users
	r.db.Model(&entities.User{}).Count(&result.TotalUsers)

	// Get total events
	r.db.Model(&entities.Event{}).Count(&result.TotalEvents)

	// Get total tickets
	r.db.Model(&entities.Ticket{}).Count(&result.TotalTickets)

	// Get active events
	r.db.Model(&entities.Event{}).Where("status = ?", "upcoming").Count(&result.ActiveEvents)

	// Get completed events
	r.db.Model(&entities.Event{}).Where("status = ?", "completed").Count(&result.CompletedEvents)

	// Get cancelled events
	r.db.Model(&entities.Event{}).Where("status = ?", "cancelled").Count(&result.CancelledEvents)

	// Get cancelled tickets
	r.db.Model(&entities.Ticket{}).Where("status = ?", "cancelled").Count(&result.CancelledTickets)

	return &result, nil
}

func (r *reportRepository) GetUserMetrics() (*UserMetrics, error) {
	var result UserMetrics

	// Get total users
	r.db.Model(&entities.User{}).Count(&result.TotalUsers)

	// Get admin users
	r.db.Model(&entities.User{}).Where("role = ?", "admin").Count(&result.AdminUsers)

	// Get regular users
	r.db.Model(&entities.User{}).Where("role = ?", "user").Count(&result.RegularUsers)

	// Get active users (users who have booked tickets)
	r.db.Model(&entities.User{}).
		Joins("JOIN tickets ON users.id = tickets.user_id").
		Where("tickets.status = ?", "booked").
		Group("users.id").
		Count(&result.ActiveUsers)

	// Get new users this month
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	r.db.Model(&entities.User{}).
		Where("created_at >= ?", startOfMonth).
		Count(&result.NewUsersMonth)

	return &result, nil
}

func (r *reportRepository) GetEventMetrics() (*EventMetrics, error) {
	var result EventMetrics

	// Get total events
	r.db.Model(&entities.Event{}).Count(&result.TotalEvents)

	// Get active events
	r.db.Model(&entities.Event{}).Where("status = ?", "upcoming").Count(&result.ActiveEvents)

	// Get completed events
	r.db.Model(&entities.Event{}).Where("status = ?", "completed").Count(&result.CompletedEvents)

	// Get cancelled events
	r.db.Model(&entities.Event{}).Where("status = ?", "cancelled").Count(&result.CancelledEvents)

	// Get average capacity
	r.db.Model(&entities.Event{}).Select("AVG(capacity)").Scan(&result.AverageCapacity)

	// Get average price
	r.db.Model(&entities.Event{}).Select("AVG(price)").Scan(&result.AveragePrice)

	return &result, nil
}

func (r *reportRepository) GetTicketMetrics() (*TicketMetrics, error) {
	var result TicketMetrics

	// Get total tickets
	r.db.Model(&entities.Ticket{}).Count(&result.TotalTickets)

	// Get booked tickets
	r.db.Model(&entities.Ticket{}).Where("status = ?", "booked").Count(&result.BookedTickets)

	// Get cancelled tickets
	r.db.Model(&entities.Ticket{}).Where("status = ?", "cancelled").Count(&result.CancelledTickets)

	// Get used tickets
	r.db.Model(&entities.Ticket{}).Where("status = ?", "used").Count(&result.UsedTickets)

	// Get average ticket price
	r.db.Model(&entities.Ticket{}).Select("AVG(unit_price)").Scan(&result.AverageTicketPrice)

	return &result, nil
}

func (r *reportRepository) GetRevenueMetrics() (*RevenueMetrics, error) {
	var result RevenueMetrics

	// Get total revenue from booked tickets
	r.db.Model(&entities.Ticket{}).
		Where("status = ?", "booked").
		Select("SUM(total_price)").
		Scan(&result.TotalRevenue)

	// Get monthly revenue
	startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	r.db.Model(&entities.Ticket{}).
		Where("status = ? AND created_at >= ?", "booked", startOfMonth).
		Select("SUM(total_price)").
		Scan(&result.MonthlyRevenue)

	// Get average revenue per event
	r.db.Table("tickets").
		Select("AVG(event_revenue)").
		Joins("JOIN (SELECT event_id, SUM(total_price) as event_revenue FROM tickets WHERE status = 'booked' GROUP BY event_id) as event_totals ON tickets.event_id = event_totals.event_id").
		Scan(&result.AverageRevenue)

	// Get refunded amount (cancelled tickets)
	r.db.Model(&entities.Ticket{}).
		Where("status = ?", "cancelled").
		Select("SUM(total_price)").
		Scan(&result.RefundedAmount)

	return &result, nil
}

func (r *reportRepository) GetTopEvents(limit int) ([]TopEventReport, error) {
	var results []TopEventReport

	err := r.db.Table("tickets").
		Select("tickets.event_id, events.title, events.category, SUM(tickets.quantity) as tickets_sold, SUM(tickets.total_price) as revenue").
		Joins("JOIN events ON tickets.event_id = events.id").
		Where("tickets.status = ?", "booked").
		Group("tickets.event_id, events.title, events.category").
		Order("revenue DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

func (r *reportRepository) GetCategoryBreakdown() ([]CategoryBreakdownReport, error) {
	var results []CategoryBreakdownReport

	err := r.db.Table("events").
		Select("events.category, COUNT(events.id) as event_count, COALESCE(SUM(tickets.quantity), 0) as tickets_sold, COALESCE(SUM(tickets.total_price), 0) as revenue").
		Joins("LEFT JOIN tickets ON events.id = tickets.event_id AND tickets.status = 'booked'").
		Group("events.category").
		Order("revenue DESC").
		Scan(&results).Error

	return results, err
}

func (r *reportRepository) GetMonthlyStats() ([]MonthlyStatsReport, error) {
	var results []MonthlyStatsReport

	// Get stats for the last 6 months
	err := r.db.Raw(`
		SELECT 
			DATE_FORMAT(months.month, '%Y-%m') as month,
			COALESCE(event_stats.events, 0) as events,
			COALESCE(ticket_stats.tickets, 0) as tickets,
			COALESCE(ticket_stats.revenue, 0) as revenue,
			COALESCE(user_stats.new_users, 0) as new_users
		FROM (
			SELECT DATE_FORMAT(DATE_SUB(CURDATE(), INTERVAL n MONTH), '%Y-%m-01') as month
			FROM (SELECT 0 as n UNION SELECT 1 UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5) numbers
		) months
		LEFT JOIN (
			SELECT DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as events
			FROM events
			WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 MONTH)
			GROUP BY DATE_FORMAT(created_at, '%Y-%m')
		) event_stats ON months.month = event_stats.month
		LEFT JOIN (
			SELECT DATE_FORMAT(created_at, '%Y-%m') as month, SUM(quantity) as tickets, SUM(total_price) as revenue
			FROM tickets
			WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 MONTH) AND status = 'booked'
			GROUP BY DATE_FORMAT(created_at, '%Y-%m')
		) ticket_stats ON months.month = ticket_stats.month
		LEFT JOIN (
			SELECT DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as new_users
			FROM users
			WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 MONTH)
			GROUP BY DATE_FORMAT(created_at, '%Y-%m')
		) user_stats ON months.month = user_stats.month
		ORDER BY months.month DESC
	`).Scan(&results).Error

	return results, err
}
