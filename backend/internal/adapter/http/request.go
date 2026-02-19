package http

type SignupRequest struct {
    Username string `json:"username" binding:"required,min=8,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}



