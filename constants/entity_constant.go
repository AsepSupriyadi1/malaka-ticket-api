package constants

// Event Status Constants
const (
	EventStatusUpcoming  = "upcoming"
	EventStatusOngoing   = "ongoing"
	EventStatusCompleted = "completed"
	EventStatusCancelled = "cancelled"
)

// Ticket Status Constants
const (
	TicketStatusBooked    = "booked"
	TicketStatusCancelled = "cancelled"
	TicketStatusUsed      = "used"
)

// User Role Constants
const (
	UserRoleUser  = "user"
	UserRoleAdmin = "admin"
)

// Sort Order Constants
const (
	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

// Default Pagination Values
const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
	MinLimit     = 1
)

// Event Categories (you can extend this based on your needs)
var ValidEventCategories = []string{
	"conference",
	"workshop",
	"seminar",
	"concert",
	"sports",
	"exhibition",
	"networking",
	"entertainment",
	"education",
	"technology",
	"business",
	"health",
	"food",
	"art",
	"other",
}

// Valid Event Statuses
var ValidEventStatuses = []string{
	EventStatusUpcoming,
	EventStatusOngoing,
	EventStatusCompleted,
	EventStatusCancelled,
}

// Valid Ticket Statuses
var ValidTicketStatuses = []string{
	TicketStatusBooked,
	TicketStatusCancelled,
	TicketStatusUsed,
}

// Valid User Roles
var ValidUserRoles = []string{
	UserRoleUser,
	UserRoleAdmin,
}

// Helper functions to validate enum values
func IsValidEventStatus(status string) bool {
	for _, validStatus := range ValidEventStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func IsValidTicketStatus(status string) bool {
	for _, validStatus := range ValidTicketStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func IsValidUserRole(role string) bool {
	for _, validRole := range ValidUserRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

func IsValidEventCategory(category string) bool {
	for _, validCategory := range ValidEventCategories {
		if category == validCategory {
			return true
		}
	}
	return false
}
