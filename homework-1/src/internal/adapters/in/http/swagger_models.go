package handler

type HealthResponse struct {
	Status string `json:"status"`
}

type ValidationDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id"`
	Details   []any  `json:"details"`
}

type ValidationErrorResponse struct {
	Error     string             `json:"error"`
	RequestID string             `json:"request_id"`
	Details   []ValidationDetail `json:"details"`
}
