package domain

import "time"

type WasteType struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    PointPerKg  int       `json:"point_per_kg"`
    Description string    `json:"description"`
    IsActive    bool      `json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type WasteTypeRequest struct {
    Name        string `json:"name" validate:"required"`
    PointPerKg  int    `json:"point_per_kg" validate:"required,min=1"`
    Description string `json:"description"`
    IsActive    bool   `json:"is_active"`
}