package services

import (
	"bytes"
	"case_study_api/dto"
	"case_study_api/repositories"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type ReportService interface {
	GetSummary() (*dto.SummaryReportResponse, error)
	GetEventReport(eventID uint) (*dto.EventReportResponse, error)
	GetSystemReport() (*dto.SystemReportResponse, error)
	GenerateSystemReportPDF() ([]byte, error)
}

type reportService struct {
	reportRepo repositories.ReportRepository
}

func NewReportService(reportRepo repositories.ReportRepository) ReportService {
	return &reportService{
		reportRepo: reportRepo,
	}
}

func (s *reportService) GetSummary() (*dto.SummaryReportResponse, error) {
	report, err := s.reportRepo.GetSummaryReport()
	if err != nil {
		return nil, err
	}

	return &dto.SummaryReportResponse{
		TotalTickets: report.TotalTickets,
		TotalRevenue: report.TotalRevenue,
	}, nil
}

func (s *reportService) GetEventReport(eventID uint) (*dto.EventReportResponse, error) {
	report, err := s.reportRepo.GetEventReport(eventID)
	if err != nil {
		return nil, err
	}

	return &dto.EventReportResponse{
		EventID:     report.EventID,
		Title:       report.Title,
		TicketsSold: report.TicketsSold,
		Revenue:     report.Revenue,
	}, nil
}

func (s *reportService) GetSystemReport() (*dto.SystemReportResponse, error) {
	// Get all the different metrics
	overview, err := s.reportRepo.GetSystemOverview()
	if err != nil {
		return nil, err
	}

	userMetrics, err := s.reportRepo.GetUserMetrics()
	if err != nil {
		return nil, err
	}

	eventMetrics, err := s.reportRepo.GetEventMetrics()
	if err != nil {
		return nil, err
	}

	ticketMetrics, err := s.reportRepo.GetTicketMetrics()
	if err != nil {
		return nil, err
	}

	revenueMetrics, err := s.reportRepo.GetRevenueMetrics()
	if err != nil {
		return nil, err
	}

	topEvents, err := s.reportRepo.GetTopEvents(10)
	if err != nil {
		return nil, err
	}

	categoryBreakdown, err := s.reportRepo.GetCategoryBreakdown()
	if err != nil {
		return nil, err
	}

	monthlyStats, err := s.reportRepo.GetMonthlyStats()
	if err != nil {
		return nil, err
	}

	// Convert repository types to DTO types
	topEventsDTO := make([]dto.TopEventReport, len(topEvents))
	for i, event := range topEvents {
		topEventsDTO[i] = dto.TopEventReport{
			EventID:     event.EventID,
			Title:       event.Title,
			TicketsSold: event.TicketsSold,
			Revenue:     event.Revenue,
			Category:    event.Category,
		}
	}

	categoryBreakdownDTO := make([]dto.CategoryBreakdownReport, len(categoryBreakdown))
	for i, category := range categoryBreakdown {
		categoryBreakdownDTO[i] = dto.CategoryBreakdownReport{
			Category:    category.Category,
			EventCount:  category.EventCount,
			TicketsSold: category.TicketsSold,
			Revenue:     category.Revenue,
		}
	}

	monthlyStatsDTO := make([]dto.MonthlyStatsReport, len(monthlyStats))
	for i, stats := range monthlyStats {
		monthlyStatsDTO[i] = dto.MonthlyStatsReport{
			Month:    stats.Month,
			Events:   stats.Events,
			Tickets:  stats.Tickets,
			Revenue:  stats.Revenue,
			NewUsers: stats.NewUsers,
		}
	}

	return &dto.SystemReportResponse{
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
		SystemName:  "Malaka Ticket",
		Overview: &dto.SystemOverview{
			TotalUsers:       overview.TotalUsers,
			TotalEvents:      overview.TotalEvents,
			TotalTickets:     overview.TotalTickets,
			ActiveEvents:     overview.ActiveEvents,
			CompletedEvents:  overview.CompletedEvents,
			CancelledEvents:  overview.CancelledEvents,
			CancelledTickets: overview.CancelledTickets,
		},
		UserMetrics: &dto.UserMetrics{
			TotalUsers:    userMetrics.TotalUsers,
			AdminUsers:    userMetrics.AdminUsers,
			RegularUsers:  userMetrics.RegularUsers,
			ActiveUsers:   userMetrics.ActiveUsers,
			NewUsersMonth: userMetrics.NewUsersMonth,
		},
		EventMetrics: &dto.EventMetrics{
			TotalEvents:     eventMetrics.TotalEvents,
			ActiveEvents:    eventMetrics.ActiveEvents,
			CompletedEvents: eventMetrics.CompletedEvents,
			CancelledEvents: eventMetrics.CancelledEvents,
			AverageCapacity: eventMetrics.AverageCapacity,
			AveragePrice:    eventMetrics.AveragePrice,
		},
		TicketMetrics: &dto.TicketMetrics{
			TotalTickets:       ticketMetrics.TotalTickets,
			BookedTickets:      ticketMetrics.BookedTickets,
			CancelledTickets:   ticketMetrics.CancelledTickets,
			UsedTickets:        ticketMetrics.UsedTickets,
			AverageTicketPrice: ticketMetrics.AverageTicketPrice,
		},
		RevenueMetrics: &dto.RevenueMetrics{
			TotalRevenue:   revenueMetrics.TotalRevenue,
			MonthlyRevenue: revenueMetrics.MonthlyRevenue,
			AverageRevenue: revenueMetrics.AverageRevenue,
			RefundedAmount: revenueMetrics.RefundedAmount,
		},
		TopEvents:         topEventsDTO,
		CategoryBreakdown: categoryBreakdownDTO,
		MonthlyStats:      monthlyStatsDTO,
	}, nil
}

func (s *reportService) GenerateSystemReportPDF() ([]byte, error) {
	report, err := s.GetSystemReport()
	if err != nil {
		return nil, err
	}

	return s.generatePDF(report)
}

func (s *reportService) generatePDF(report *dto.SystemReportResponse) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, report.SystemName+" - System Report")
	pdf.Ln(15)

	// Generated timestamp
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, "Generated at: "+report.GeneratedAt)
	pdf.Ln(15)

	// System Overview
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "System Overview")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Total Users: %d", report.Overview.TotalUsers))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Total Events: %d", report.Overview.TotalEvents))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Total Tickets: %d", report.Overview.TotalTickets))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Active Events: %d", report.Overview.ActiveEvents))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Completed Events: %d", report.Overview.CompletedEvents))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Cancelled Events: %d", report.Overview.CancelledEvents))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Cancelled Tickets: %d", report.Overview.CancelledTickets))
	pdf.Ln(15)

	// User Metrics
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "User Metrics")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Total Users: %d", report.UserMetrics.TotalUsers))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Admin Users: %d", report.UserMetrics.AdminUsers))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Regular Users: %d", report.UserMetrics.RegularUsers))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Active Users: %d", report.UserMetrics.ActiveUsers))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("New Users This Month: %d", report.UserMetrics.NewUsersMonth))
	pdf.Ln(15)

	// Revenue Metrics
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Revenue Metrics")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Total Revenue: $%.2f", report.RevenueMetrics.TotalRevenue))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Monthly Revenue: $%.2f", report.RevenueMetrics.MonthlyRevenue))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Average Revenue per Event: $%.2f", report.RevenueMetrics.AverageRevenue))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Refunded Amount: $%.2f", report.RevenueMetrics.RefundedAmount))
	pdf.Ln(15)

	// Top Events
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Top Events by Revenue")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	for i, event := range report.TopEvents {
		if i >= 5 { // Show only top 5 in PDF
			break
		}
		pdf.Cell(0, 6, fmt.Sprintf("%d. %s - $%.2f (%d tickets sold)", i+1, event.Title, event.Revenue, event.TicketsSold))
		pdf.Ln(6)
	}

	// Category Breakdown
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Category Breakdown")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	for _, category := range report.CategoryBreakdown {
		pdf.Cell(0, 6, fmt.Sprintf("%s: %d events, %d tickets sold, $%.2f revenue", category.Category, category.EventCount, category.TicketsSold, category.Revenue))
		pdf.Ln(6)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
