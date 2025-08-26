package usecase

import (
	"errors"

	"github.com/alifsuryadi/ecolokal/config"
	"github.com/alifsuryadi/ecolokal/internal/domain"
	"github.com/alifsuryadi/ecolokal/internal/repository"
	"github.com/alifsuryadi/ecolokal/pkg/utils"
)

type AuthUsecase struct {
    userRepo *repository.UserRepository
    cfg      *config.Config
}

func NewAuthUsecase(userRepo *repository.UserRepository, cfg *config.Config) *AuthUsecase {
    return &AuthUsecase{
        userRepo: userRepo,
        cfg:      cfg,
    }
}

func (u *AuthUsecase) Register(req *domain.RegisterRequest) (*domain.User, error) {
    // Check if email already exists
    existingUser, _ := u.userRepo.GetByEmail(req.Email)
    if existingUser != nil {
        return nil, errors.New("email already registered")
    }
    
    // Hash password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }
    
    // Set default role if not provided
    if req.Role == "" {
        req.Role = "warga"
    }
    
    // Validate role
    if req.Role != "warga" && req.Role != "petugas" {
        return nil, errors.New("invalid role")
    }
    
    user := &domain.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
        Role:     req.Role,
        Phone:    req.Phone,
        Address:  req.Address,
    }
    
    err = u.userRepo.Create(user)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

func (u *AuthUsecase) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
    // Get user by email
    user, err := u.userRepo.GetByEmail(req.Email)
    if err != nil {
        return nil, errors.New("invalid email or password")
    }
    
    // Check password
    if !utils.CheckPassword(req.Password, user.Password) {
        return nil, errors.New("invalid email or password")
    }
    
    // Generate JWT token
    token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, u.cfg)
    if err != nil {
        return nil, err
    }
    
    return &domain.LoginResponse{
        Token: token,
        User:  *user,
    }, nil
}