package user

import "time"

type UserDto struct {
	CustomerName string    `json:"customer_name" validate:"required"`
	PhoneNumber  string    `json:"phone_number" validate:"required,min=11"`
	Email        string    `json:"email" validate:"required,email"`
	DateOfBird   time.Time `json:"dateOfBirth" validate:"required"`
	Sex          bool      `json:"sex" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
}
