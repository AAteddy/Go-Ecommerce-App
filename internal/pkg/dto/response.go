package dto

// ApiResponse represents a standardized JSON response for API endpoints.
type ApiResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Error      string      `json:"error"`
	Data       interface{} `json:"data"`
}
