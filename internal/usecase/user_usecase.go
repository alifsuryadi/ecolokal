package usecase

import (
	"github.com/alifsuryadi/ecolokal/internal/domain"
	"github.com/alifsuryadi/ecolokal/internal/repository"
)

type UserUsecase struct {
    userRepo *repository.UserRepository
}

func NewUserUsecase(userRepo *repository.UserRepository) *UserUsecase {
    return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) GetUserByID(id int) (*domain.User, error) {
    return u.userRepo.GetByID(id)
}

func (u *UserUsecase) GetUserPoints(userID int) (*domain.UserPoint, error) {
    return u.userRepo.GetUserPoints(userID)
}

func (u *UserUsecase) GetUsersByRole(role string) ([]domain.User, error) {
    return u.userRepo.GetUsersByRole(role)
}