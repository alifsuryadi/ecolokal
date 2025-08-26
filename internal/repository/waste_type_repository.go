package repository

import (
	"database/sql"

	"github.com/alifsuryadi/ecolokal/internal/domain"
)

type WasteTypeRepository struct {
    db *sql.DB
}

func NewWasteTypeRepository(db *sql.DB) *WasteTypeRepository {
    return &WasteTypeRepository{db: db}
}

func (r *WasteTypeRepository) Create(wasteType *domain.WasteType) error {
    query := `
        INSERT INTO waste_types (name, point_per_kg, description, is_active)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `
    err := r.db.QueryRow(
        query,
        wasteType.Name,
        wasteType.PointPerKg,
        wasteType.Description,
        wasteType.IsActive,
    ).Scan(&wasteType.ID, &wasteType.CreatedAt, &wasteType.UpdatedAt)
    
    return err
}

func (r *WasteTypeRepository) GetAll() ([]domain.WasteType, error) {
    query := `
        SELECT id, name, point_per_kg, description, is_active, created_at, updated_at
        FROM waste_types
        ORDER BY created_at DESC
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var wasteTypes []domain.WasteType
    for rows.Next() {
        var wt domain.WasteType
        err := rows.Scan(
            &wt.ID,
            &wt.Name,
            &wt.PointPerKg,
            &wt.Description,
            &wt.IsActive,
            &wt.CreatedAt,
            &wt.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        wasteTypes = append(wasteTypes, wt)
    }
    return wasteTypes, nil
}


func (r *WasteTypeRepository) GetByID(id int) (*domain.WasteType, error) {
    var wt domain.WasteType
    query := `
        SELECT id, name, point_per_kg, description, is_active, created_at, updated_at
        FROM waste_types
        WHERE id = $1
    `
    err := r.db.QueryRow(query, id).Scan(
        &wt.ID,
        &wt.Name,
        &wt.PointPerKg,
        &wt.Description,
        &wt.IsActive,
        &wt.CreatedAt,
        &wt.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &wt, nil
}

func (r *WasteTypeRepository) Update(wasteType *domain.WasteType) error {
    query := `
        UPDATE waste_types
        SET name = $1, point_per_kg = $2, description = $3, is_active = $4
        WHERE id = $5
        RETURNING updated_at
    `
    err := r.db.QueryRow(
        query,
        wasteType.Name,
        wasteType.PointPerKg,
        wasteType.Description,
        wasteType.IsActive,
        wasteType.ID,
    ).Scan(&wasteType.UpdatedAt)
    
    return err
}

func (r *WasteTypeRepository) Delete(id int) error {
    query := `DELETE FROM waste_types WHERE id = $1`
    _, err := r.db.Exec(query, id)
    return err
}

func (r *WasteTypeRepository) GetActiveTypes() ([]domain.WasteType, error) {
    query := `
        SELECT id, name, point_per_kg, description, is_active, created_at, updated_at
        FROM waste_types
        WHERE is_active = true
        ORDER BY name
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var wasteTypes []domain.WasteType
    for rows.Next() {
        var wt domain.WasteType
        err := rows.Scan(
            &wt.ID,
            &wt.Name,
            &wt.PointPerKg,
            &wt.Description,
            &wt.IsActive,
            &wt.CreatedAt,
            &wt.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        wasteTypes = append(wasteTypes, wt)
    }
    return wasteTypes, nil
}