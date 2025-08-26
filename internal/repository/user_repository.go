package repository

import (
	"database/sql"

	"github.com/alifsuryadi/ecolokal/internal/domain"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
    query := `
        INSERT INTO users (name, email, password, role, phone, address)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `
    err := r.db.QueryRow(
        query,
        user.Name,
        user.Email,
        user.Password,
        user.Role,
        user.Phone,
        user.Address,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

    if err != nil {
        return err
    }

    // Create user points entry
    _, err = r.db.Exec("INSERT INTO user_points (user_id) VALUES ($1)", user.ID)
    return err
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
    var user domain.User
    query := `
        SELECT id, name, email, password, role, phone, address, created_at, updated_at
        FROM users
        WHERE email = $1
    `
    err := r.db.QueryRow(query, email).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Password,
        &user.Role,
        &user.Phone,
        &user.Address,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {
    var user domain.User
    query := `
        SELECT id, name, email, password, role, phone, address, created_at, updated_at
        FROM users
        WHERE id = $1
    `
    err := r.db.QueryRow(query, id).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Password,
        &user.Role,
        &user.Phone,
        &user.Address,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetUserPoints(userID int) (*domain.UserPoint, error) {
    var points domain.UserPoint
    query := `
        SELECT id, user_id, total_points, last_updated
        FROM user_points
        WHERE user_id = $1
    `
    err := r.db.QueryRow(query, userID).Scan(
        &points.ID,
        &points.UserID,
        &points.TotalPoints,
        &points.LastUpdated,
    )
    if err != nil {
        return nil, err
    }
    return &points, nil
}

func (r *UserRepository) UpdateUserPoints(userID int, points int) error {
    query := `
        UPDATE user_points 
        SET total_points = total_points + $1, last_updated = CURRENT_TIMESTAMP
        WHERE user_id = $2
    `
    _, err := r.db.Exec(query, points, userID)
    return err
}

func (r *UserRepository) GetUsersByRole(role string) ([]domain.User, error) {
    query := `
        SELECT id, name, email, password, role, phone, address, created_at, updated_at
        FROM users
        WHERE role = $1
        ORDER BY created_at DESC
    `
    rows, err := r.db.Query(query, role)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []domain.User
    for rows.Next() {
        var user domain.User
        err := rows.Scan(
            &user.ID,
            &user.Name,
            &user.Email,
            &user.Password,
            &user.Role,
            &user.Phone,
            &user.Address,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}