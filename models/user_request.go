package models

// struct untuk input register user
type RegisterInput struct {
	Username string `json:"username" gorm:"uniqueindex;not null;binding:'required'"`
	Email    string `json:"email" gorm:"uniqueindex;not null;binding:'required,email'"`
	Password string `json:"password" gorm:"not null;binding:'required,min=6'"`
	Role     string `json:"role" gorm:"not null;default:'customer'"`
}

// struct untuk input login user
type LoginInput struct {
	Email    string `json:"email" binding:"required,email" validate:"email"`
	Password string `json:"password" binding:"required"`
}
