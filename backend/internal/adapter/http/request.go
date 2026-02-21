package http

type SignupRequest struct {
	Username string `json:"username" binding:"required,min=8,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type BusRequest struct {
	SourceCity string  `json:"source_city" binding:"required"`
	DestCity   string  `json:"dest_city" binding:"required"`
	Date       string  `json:"date" binding:"required"`
	Time       string  `json:"time" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
}
