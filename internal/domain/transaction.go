package domain

import "time"

type Transaction struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    User        *User     `json:"user,omitempty"`
    Type        string    `json:"type"`
    Points      int       `json:"points"`
    Description string    `json:"description"`
    ReferenceID *int      `json:"reference_id"`
    CreatedAt   time.Time `json:"created_at"`
}

type CreateTransaction struct {
    UserID      int    `json:"user_id" validate:"required"`
    Type        string `json:"type" validate:"required,oneof=add redeem"`
    Points      int    `json:"points" validate:"required,min=1"`
    Description string `json:"description" validate:"required"`
    ReferenceID *int   `json:"reference_id"`
}