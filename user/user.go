package user

import "time"

type User struct {
	CustomerId   string    `json:"id"`
	CustomerName string    `json:"customer_name" validate:"required"`
	PhoneNumber  string    `json:"phone_number" validate:"required,min=11"`
	Email        string    `json:"email" validate:"required,email"`
	DateOfBird   time.Time `json:"dateOfBirth" validate:"required"`
	Sex          bool      `json:"sex" validate:"required"`
	Salt         string    `json:"salt"`
	Password     string    `json:"password" validate:"required,min=6"`
	CreatedAt    time.Time `json:"created_at"`
}
