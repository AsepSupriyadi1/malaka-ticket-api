// utils/response.go
package utils

type ErrorResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Errors  []interface{} `json:"errors,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Fungsi-fungsi builder
func BuildErrorResponse(message string, details ...interface{}) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Message: message,
	}
}

// Fungsi-fungsi builder
func BuildErrorResponseWithDetails(message string, details ...interface{}) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Message: message,
		Errors:  details,
	}
}

func BuildSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}
