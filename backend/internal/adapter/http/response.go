package http

// SignUpResponse is the response returned after successful signup
type SignUpResponse struct {
	Message string `json:"message" example:"User registered successfully"`
	Token   string `json:"token" example:"generated token"`
}

// ErrorResponse is returned when an error occurs
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}