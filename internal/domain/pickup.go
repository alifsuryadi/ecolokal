package domain

import "time"

type PickupRequest struct {
    ID            int          `json:"id"`
    UserID        int          `json:"user_id"`
    User          *User        `json:"user,omitempty"`
    PetugasID     *int         `json:"petugas_id"`
    Petugas       *User        `json:"petugas,omitempty"`
    ScheduledDate *time.Time   `json:"scheduled_date"`
    ScheduledTime string       `json:"scheduled_time"`
    Status        string       `json:"status"`
    Notes         string       `json:"notes"`
    TotalPoints   int          `json:"total_points"`
    Items         []PickupItem `json:"items,omitempty"`
    CreatedAt     time.Time    `json:"created_at"`
    UpdatedAt     time.Time    `json:"updated_at"`
}

type PickupItem struct {
    ID              int       `json:"id"`
    PickupRequestID int       `json:"pickup_request_id"`
    WasteTypeID     int       `json:"waste_type_id"`
    WasteType       *WasteType `json:"waste_type,omitempty"`
    EstimatedWeight float64   `json:"estimated_weight"`
    ActualWeight    *float64  `json:"actual_weight"`
    PointsEarned    int       `json:"points_earned"`
    CreatedAt       time.Time `json:"created_at"`
}

type CreatePickupRequest struct {
    ScheduledDate string                  `json:"scheduled_date" validate:"required"`
    ScheduledTime string                  `json:"scheduled_time" validate:"required"`
    Notes         string                  `json:"notes"`
    Items         []CreatePickupItem      `json:"items" validate:"required,min=1"`
}

type CreatePickupItem struct {
    WasteTypeID     int     `json:"waste_type_id" validate:"required"`
    EstimatedWeight float64 `json:"estimated_weight" validate:"required,min=0.1"`
}

type UpdatePickupStatus struct {
    Status string `json:"status" validate:"required,oneof=scheduled in_progress completed cancelled"`
}

type UpdatePickupItems struct {
    Items []UpdatePickupItem `json:"items" validate:"required"`
}

type UpdatePickupItem struct {
    ID           int     `json:"id" validate:"required"`
    ActualWeight float64 `json:"actual_weight" validate:"required,min=0.1"`
}