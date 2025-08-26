package repository

import (
	"database/sql"
	"time"

	"github.com/alifsuryadi/ecolokal/internal/domain"
)

type PickupRepository struct {
    db *sql.DB
}

func NewPickupRepository(db *sql.DB) *PickupRepository {
    return &PickupRepository{db: db}
}

func (r *PickupRepository) Create(pickup *domain.PickupRequest) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Insert pickup request
    query := `
        INSERT INTO pickup_requests (user_id, scheduled_date, scheduled_time, notes)
        VALUES ($1, $2, $3, $4)
        RETURNING id, status, created_at, updated_at
    `
    err = tx.QueryRow(
        query,
        pickup.UserID,
        pickup.ScheduledDate,
        pickup.ScheduledTime,
        pickup.Notes,
    ).Scan(&pickup.ID, &pickup.Status, &pickup.CreatedAt, &pickup.UpdatedAt)
    
    if err != nil {
        return err
    }

    // Insert pickup items
    for i := range pickup.Items {
        item := &pickup.Items[i]
        itemQuery := `
            INSERT INTO pickup_items (pickup_request_id, waste_type_id, estimated_weight)
            VALUES ($1, $2, $3)
            RETURNING id, created_at
        `
        err = tx.QueryRow(
            itemQuery,
            pickup.ID,
            item.WasteTypeID,
            item.EstimatedWeight,
        ).Scan(&item.ID, &item.CreatedAt)
        
        if err != nil {
            return err
        }
        item.PickupRequestID = pickup.ID
    }

    return tx.Commit()
}

func (r *PickupRepository) GetByID(id int) (*domain.PickupRequest, error) {
    var pickup domain.PickupRequest
    query := `
        SELECT pr.id, pr.user_id, pr.petugas_id, pr.scheduled_date, pr.scheduled_time,
               pr.status, pr.notes, pr.total_points, pr.created_at, pr.updated_at,
               u.name, u.email, u.phone, u.address
        FROM pickup_requests pr
        JOIN users u ON pr.user_id = u.id
        WHERE pr.id = $1
    `
    
    var petugasID sql.NullInt64
    var scheduledDate sql.NullTime
    user := domain.User{}
    
    err := r.db.QueryRow(query, id).Scan(
        &pickup.ID,
        &pickup.UserID,
        &petugasID,
        &scheduledDate,
        &pickup.ScheduledTime,
        &pickup.Status,
        &pickup.Notes,
        &pickup.TotalPoints,
        &pickup.CreatedAt,
        &pickup.UpdatedAt,
        &user.Name,
        &user.Email,
        &user.Phone,
        &user.Address,
    )
    
    if err != nil {
        return nil, err
    }
    
    if petugasID.Valid {
        pid := int(petugasID.Int64)
        pickup.PetugasID = &pid
    }
    
    if scheduledDate.Valid {
        pickup.ScheduledDate = &scheduledDate.Time
    }
    
    user.ID = pickup.UserID
    pickup.User = &user
    
    // Get items
    items, err := r.getPickupItems(pickup.ID)
    if err != nil {
        return nil, err
    }
    pickup.Items = items
    
    return &pickup, nil
}

func (r *PickupRepository) getPickupItems(pickupID int) ([]domain.PickupItem, error) {
    query := `
        SELECT pi.id, pi.waste_type_id, pi.estimated_weight, pi.actual_weight, 
               pi.points_earned, pi.created_at,
               wt.name, wt.point_per_kg
        FROM pickup_items pi
        JOIN waste_types wt ON pi.waste_type_id = wt.id
        WHERE pi.pickup_request_id = $1
    `
    
    rows, err := r.db.Query(query, pickupID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var items []domain.PickupItem
    for rows.Next() {
        var item domain.PickupItem
        var actualWeight sql.NullFloat64
        wasteType := domain.WasteType{}
        
        err := rows.Scan(
            &item.ID,
            &item.WasteTypeID,
            &item.EstimatedWeight,
            &actualWeight,
            &item.PointsEarned,
            &item.CreatedAt,
            &wasteType.Name,
            &wasteType.PointPerKg,
        )
        
        if err != nil {
            return nil, err
        }
        
        if actualWeight.Valid {
            item.ActualWeight = &actualWeight.Float64
        }
        
        wasteType.ID = item.WasteTypeID
        item.WasteType = &wasteType
        item.PickupRequestID = pickupID
        
        items = append(items, item)
    }
    
    return items, nil
}

func (r *PickupRepository) GetUserPickups(userID int) ([]domain.PickupRequest, error) {
    query := `
        SELECT id, user_id, petugas_id, scheduled_date, scheduled_time,
               status, notes, total_points, created_at, updated_at
        FROM pickup_requests
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var pickups []domain.PickupRequest
    for rows.Next() {
        var pickup domain.PickupRequest
        var petugasID sql.NullInt64
        var scheduledDate sql.NullTime
        
        err := rows.Scan(
            &pickup.ID,
            &pickup.UserID,
            &petugasID,
            &scheduledDate,
            &pickup.ScheduledTime,
            &pickup.Status,
            &pickup.Notes,
            &pickup.TotalPoints,
            &pickup.CreatedAt,
            &pickup.UpdatedAt,
        )
        
        if err != nil {
            return nil, err
        }
        
        if petugasID.Valid {
            pid := int(petugasID.Int64)
            pickup.PetugasID = &pid
        }
        
        if scheduledDate.Valid {
            pickup.ScheduledDate = &scheduledDate.Time
        }
        
        pickups = append(pickups, pickup)
    }
    
    return pickups, nil
}

func (r *PickupRepository) GetPetugasPickups(petugasID int, date time.Time) ([]domain.PickupRequest, error) {
    query := `
        SELECT pr.id, pr.user_id, pr.scheduled_time, pr.status, pr.notes, pr.total_points,
               pr.created_at, pr.updated_at,
               u.name, u.phone, u.address
        FROM pickup_requests pr
        JOIN users u ON pr.user_id = u.id
        WHERE pr.petugas_id = $1 AND pr.scheduled_date = $2
        ORDER BY pr.scheduled_time
    `
    
    rows, err := r.db.Query(query, petugasID, date)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
        var pickups []domain.PickupRequest
    for rows.Next() {
        var pickup domain.PickupRequest
        user := domain.User{}
        
        err := rows.Scan(
            &pickup.ID,
            &pickup.UserID,
            &pickup.ScheduledTime,
            &pickup.Status,
            &pickup.Notes,
            &pickup.TotalPoints,
            &pickup.CreatedAt,
            &pickup.UpdatedAt,
            &user.Name,
            &user.Phone,
            &user.Address,
        )
        
        if err != nil {
            return nil, err
        }
        
        user.ID = pickup.UserID
        pickup.User = &user
        pickup.PetugasID = &petugasID
        pickup.ScheduledDate = &date
        
        pickups = append(pickups, pickup)
    }
    
    return pickups, nil
}

func (r *PickupRepository) UpdateStatus(id int, status string, petugasID *int) error {
    var query string
    var args []interface{}
    
    if petugasID != nil {
        query = `UPDATE pickup_requests SET status = $1, petugas_id = $2 WHERE id = $3`
        args = []interface{}{status, *petugasID, id}
    } else {
        query = `UPDATE pickup_requests SET status = $1 WHERE id = $2`
        args = []interface{}{status, id}
    }
    
    _, err := r.db.Exec(query, args...)
    return err
}

func (r *PickupRepository) UpdatePickupItems(pickupID int, items []domain.UpdatePickupItem) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    totalPoints := 0
    
    for _, item := range items {
        // Get waste type points
        var pointPerKg int
        err = tx.QueryRow("SELECT point_per_kg FROM waste_types WHERE id = (SELECT waste_type_id FROM pickup_items WHERE id = $1)", item.ID).Scan(&pointPerKg)
        if err != nil {
            return err
        }
        
        // Calculate points
        pointsEarned := int(item.ActualWeight * float64(pointPerKg))
        totalPoints += pointsEarned
        
        // Update item
        _, err = tx.Exec(
            "UPDATE pickup_items SET actual_weight = $1, points_earned = $2 WHERE id = $3 AND pickup_request_id = $4",
            item.ActualWeight, pointsEarned, item.ID, pickupID,
        )
        if err != nil {
            return err
        }
    }
    
    // Update total points in pickup request
    _, err = tx.Exec("UPDATE pickup_requests SET total_points = $1 WHERE id = $2", totalPoints, pickupID)
    if err != nil {
        return err
    }
    
    return tx.Commit()
}

func (r *PickupRepository) GetPendingPickups() ([]domain.PickupRequest, error) {
    query := `
        SELECT pr.id, pr.user_id, pr.scheduled_date, pr.scheduled_time, pr.notes,
               pr.created_at, u.name, u.phone, u.address
        FROM pickup_requests pr
        JOIN users u ON pr.user_id = u.id
        WHERE pr.status = 'pending'
        ORDER BY pr.created_at
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var pickups []domain.PickupRequest
    for rows.Next() {
        var pickup domain.PickupRequest
        var scheduledDate sql.NullTime
        user := domain.User{}
        
        err := rows.Scan(
            &pickup.ID,
            &pickup.UserID,
            &scheduledDate,
            &pickup.ScheduledTime,
            &pickup.Notes,
            &pickup.CreatedAt,
            &user.Name,
            &user.Phone,
            &user.Address,
        )
        
        if err != nil {
            return nil, err
        }
        
        if scheduledDate.Valid {
            pickup.ScheduledDate = &scheduledDate.Time
        }
        
        user.ID = pickup.UserID
        pickup.User = &user
        pickup.Status = "pending"
        
        pickups = append(pickups, pickup)
    }
    
    return pickups, nil
}