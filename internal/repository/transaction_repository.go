package repository

import (
	"database/sql"

	"github.com/alifsuryadi/ecolokal/internal/domain"
)

type TransactionRepository struct {
    db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
    return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *domain.Transaction) error {
    query := `
        INSERT INTO transactions (user_id, type, points, description, reference_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `
    
    var referenceID sql.NullInt64
    if tx.ReferenceID != nil {
        referenceID = sql.NullInt64{Int64: int64(*tx.ReferenceID), Valid: true}
    }
    
    err := r.db.QueryRow(
        query,
        tx.UserID,
        tx.Type,
        tx.Points,
        tx.Description,
        referenceID,
    ).Scan(&tx.ID, &tx.CreatedAt)
    
    return err
}

func (r *TransactionRepository) GetUserTransactions(userID int) ([]domain.Transaction, error) {
    query := `
        SELECT id, user_id, type, points, description, reference_id, created_at
        FROM transactions
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
    
    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var transactions []domain.Transaction
    for rows.Next() {
        var tx domain.Transaction
        var referenceID sql.NullInt64
        
        err := rows.Scan(
            &tx.ID,
            &tx.UserID,
            &tx.Type,
            &tx.Points,
            &tx.Description,
            &referenceID,
            &tx.CreatedAt,
        )
        
        if err != nil {
            return nil, err
        }
        
        if referenceID.Valid {
            rid := int(referenceID.Int64)
            tx.ReferenceID = &rid
        }
        
        transactions = append(transactions, tx)
    }
    
    return transactions, nil
}