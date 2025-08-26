package usecase

import (
	"errors"

	"github.com/alifsuryadi/ecolokal/internal/domain"
	"github.com/alifsuryadi/ecolokal/internal/repository"
)

type WasteTypeUsecase struct {
    wasteRepo *repository.WasteTypeRepository
}

func NewWasteTypeUsecase(wasteRepo *repository.WasteTypeRepository) *WasteTypeUsecase {
    return &WasteTypeUsecase{wasteRepo: wasteRepo}
}

func (u *WasteTypeUsecase) Create(req *domain.WasteTypeRequest) (*domain.WasteType, error) {
    wasteType := &domain.WasteType{
        Name:        req.Name,
        PointPerKg:  req.PointPerKg,
        Description: req.Description,
        IsActive:    req.IsActive,
    }
    
    err := u.wasteRepo.Create(wasteType)
    if err != nil {
        return nil, err
    }
    
    return wasteType, nil
}

func (u *WasteTypeUsecase) GetAll() ([]domain.WasteType, error) {
    return u.wasteRepo.GetAll()
}

func (u *WasteTypeUsecase) GetByID(id int) (*domain.WasteType, error) {
    return u.wasteRepo.GetByID(id)
}

func (u *WasteTypeUsecase) Update(id int, req *domain.WasteTypeRequest) (*domain.WasteType, error) {
    wasteType, err := u.wasteRepo.GetByID(id)
    if err != nil {
        return nil, errors.New("waste type not found")
    }
    
    wasteType.Name = req.Name
    wasteType.PointPerKg = req.PointPerKg
    wasteType.Description = req.Description
    wasteType.IsActive = req.IsActive
    
    err = u.wasteRepo.Update(wasteType)
    if err != nil {
        return nil, err
    }
    
    return wasteType, nil
}

func (u *WasteTypeUsecase) Delete(id int) error {
    _, err := u.wasteRepo.GetByID(id)
    if err != nil {
        return errors.New("waste type not found")
    }
    
    return u.wasteRepo.Delete(id)
}

func (u *WasteTypeUsecase) GetActiveTypes() ([]domain.WasteType, error) {
    return u.wasteRepo.GetActiveTypes()
}