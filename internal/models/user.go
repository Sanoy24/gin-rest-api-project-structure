package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username" validate:"required,min=3,max=20"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"-" bson:"password" validate:"required,min=6"`
	FirstName string             `json:"first_name" bson:"first_name" validate:"required,min=2,max=50"`
	LastName  string             `json:"last_name" bson:"last_name" validate:"required,min=2,max=50"`
	Role      string             `json:"role" bson:"role" validate:"required,oneof=admin user"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type CreateUserRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=20" example:"johndoe"`
	Email     string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password  string `json:"password" validate:"required,min=6" example:"password123"`
	FirstName string `json:"first_name" validate:"required,min=1,max=50" example:"John"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50" example:"Doe"`
	Role      string `json:"role" validate:"required,oneof=admin user" enums:"admin,user" example:"user"`
}

type UpdateUserRequest struct {
	Username  string `json:"username" validate:"omitempty,min=3,max=20" example:"johndoe"`
	Email     string `json:"email" validate:"omitempty,email" example:"johndoe_new@example.com"`
	FirstName string `json:"first_name" validate:"omitempty,min=1,max=50" example:"John"`
	LastName  string `json:"last_name" validate:"omitempty,min=1,max=50" example:"Doe"`
	Role      string `json:"role" validate:"omitempty,oneof=admin user" enums:"admin,user" example:"user"`
	IsActive  *bool  `json:"is_active" example:"true"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id" example:"63a5e3e3e4b0a7e3e3e3e3e3"`
	Username  string             `json:"username" example:"johndoe"`
	Email     string             `json:"email" example:"johndoe@example.com"`
	FirstName string             `json:"first_name" example:"John"`
	LastName  string             `json:"last_name" example:"Doe"`
	Role      string             `json:"role" example:"user"`
	IsActive  bool               `json:"is_active" example:"true"`
	CreatedAt time.Time          `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt time.Time          `json:"updated_at" example:"2023-01-01T12:00:00Z"`
}

// PaginatedUserResponse represents a paginated list of users.
type PaginatedUserResponse struct {
	Data       []UserResponse `json:"data"`
	Total      int64          `json:"total" example:"100"`
	Page       int            `json:"page" example:"1"`
	Limit      int            `json:"limit" example:"10"`
	TotalPages int            `json:"total_pages" example:"10"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
