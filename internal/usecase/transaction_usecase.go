package usecase

import (
	"errors"

	"github.com/alifsuryadi/ecolokal/internal/domain"
	"github.com/alifsuryadi/ecolokal/internal/repository"
)

type TransactionUsecase struct {
    transactionRepo *repository.TransactionRepository
    userRepo        *repository.UserRepository
}

func NewTransactionUsecase(transactionRepo *repository.TransactionRepository, userRepo *repository.UserRepository) *TransactionUsecase {
    return &TransactionUsecase{
        transactionRepo: transactionRepo,
        userRepo:        userRepo,
    }
}

func (u *TransactionUsecase) CreateTransaction(req *domain.CreateTransaction) (*domain.Transaction, error) {
    // Validate user exists
    user, err := u.userRepo.GetByID(req.UserID)
    if err != nil {
        return nil, errors.New("user not found")
    }
    
    // Get current points
    userPoints, err := u.userRepo.GetUserPoints(user.ID)
    if err != nil {
        return nil, err
    }
    
    // If redeem, check if user has enough points
    if req.Type == "redeem" {
        if userPoints.TotalPoints < req.Points {
            return nil, errors.New("insufficient points")
        }
        // Deduct points (negative value for redeem)
        err = u.userRepo.UpdateUserPoints(user.ID, -req.Points)
    } else {
        // Add points
        err = u.userRepo.UpdateUserPoints(user.ID, req.Points)
    }
    
    if err != nil {
        return nil, err
    }
    
    transaction := &domain.Transaction{
        UserID:      req.UserID,
        Type:        req.Type,
        Points:      req.Points,
        Description: req.Description,
        ReferenceID: req.ReferenceID,
    }
    
    err = u.transactionRepo.Create(transaction)
    if err != nil {
        return nil, err
    }
    
    return transaction, nil
}

func (u *TransactionUsecase) GetUserTransactions(userID int) ([]domain.Transaction, error) {
    return u.transactionRepo.GetUserTransactions(userID)
}