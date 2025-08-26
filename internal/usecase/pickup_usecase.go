package usecase

import (
	"errors"
	"time"

	"github.com/alifsuryadi/ecolokal/internal/domain"
	"github.com/alifsuryadi/ecolokal/internal/repository"
)

type PickupUsecase struct {
    pickupRepo      *repository.PickupRepository
    userRepo        *repository.UserRepository
    wasteTypeRepo   *repository.WasteTypeRepository
    transactionRepo *repository.TransactionRepository
}

func NewPickupUsecase(
    pickupRepo *repository.PickupRepository,
    userRepo *repository.UserRepository,
    wasteTypeRepo *repository.WasteTypeRepository,
    transactionRepo *repository.TransactionRepository,
) *PickupUsecase {
    return &PickupUsecase{
        pickupRepo:      pickupRepo,
        userRepo:        userRepo,
        wasteTypeRepo:   wasteTypeRepo,
        transactionRepo: transactionRepo,
    }
}

func (u *PickupUsecase) CreatePickupRequest(userID int, req *domain.CreatePickupRequest) (*domain.PickupRequest, error) {
    // Parse scheduled date
    scheduledDate, err := time.Parse("2006-01-02", req.ScheduledDate)
    if err != nil {
        return nil, errors.New("invalid date format")
    }
    
    // Validate scheduled date is not in the past
    if scheduledDate.Before(time.Now().Truncate(24 * time.Hour)) {
        return nil, errors.New("scheduled date cannot be in the past")
    }
    
    // Validate waste types exist
    for _, item := range req.Items {
        _, err := u.wasteTypeRepo.GetByID(item.WasteTypeID)
        if err != nil {
            return nil, errors.New("invalid waste type ID")
        }
    }
    
    pickup := &domain.PickupRequest{
        UserID:        userID,
        ScheduledDate: &scheduledDate,
        ScheduledTime: req.ScheduledTime,
        Notes:         req.Notes,
        Items:         make([]domain.PickupItem, len(req.Items)),
    }
    
    for i, item := range req.Items {
        pickup.Items[i] = domain.PickupItem{
            WasteTypeID:     item.WasteTypeID,
            EstimatedWeight: item.EstimatedWeight,
        }
    }
    
    err = u.pickupRepo.Create(pickup)
    if err != nil {
        return nil, err
    }
    
    return pickup, nil
}

func (u *PickupUsecase) GetPickupByID(id int) (*domain.PickupRequest, error) {
    return u.pickupRepo.GetByID(id)
}

func (u *PickupUsecase) GetUserPickups(userID int) ([]domain.PickupRequest, error) {
    return u.pickupRepo.GetUserPickups(userID)
}

func (u *PickupUsecase) GetPetugasPickups(petugasID int, date string) ([]domain.PickupRequest, error) {
    parsedDate, err := time.Parse("2006-01-02", date)
    if err != nil {
        return nil, errors.New("invalid date format")
    }
    
    return u.pickupRepo.GetPetugasPickups(petugasID, parsedDate)
}

func (u *PickupUsecase) UpdatePickupStatus(id int, req *domain.UpdatePickupStatus, petugasID *int) error {
    // Get pickup to validate
    pickup, err := u.pickupRepo.GetByID(id)
    if err != nil {
        return errors.New("pickup not found")
    }
    
    // Validate status transition
    if err := u.validateStatusTransition(pickup.Status, req.Status); err != nil {
        return err
    }
    
    return u.pickupRepo.UpdateStatus(id, req.Status, petugasID)
}

func (u *PickupUsecase) UpdatePickupItems(pickupID int, userID int, req *domain.UpdatePickupItems) error {
    // Get pickup
    pickup, err := u.pickupRepo.GetByID(pickupID)
    if err != nil {
        return errors.New("pickup not found")
    }
    
    // Only petugas can update items
    user, err := u.userRepo.GetByID(userID)
    if err != nil || user.Role != "petugas" {
        return errors.New("unauthorized")
    }
    
    // Update items
    err = u.pickupRepo.UpdatePickupItems(pickupID, req.Items)
    if err != nil {
        return err
    }
    
    // Get updated pickup to calculate points
    updatedPickup, err := u.pickupRepo.GetByID(pickupID)
    if err != nil {
        return err
    }
    
    // Add points to user
    if updatedPickup.TotalPoints > 0 {
        err = u.userRepo.UpdateUserPoints(pickup.UserID, updatedPickup.TotalPoints)
        if err != nil {
            return err
        }
        
        // Create transaction record
        transaction := &domain.Transaction{
            UserID:      pickup.UserID,
            Type:        "add",
            Points:      updatedPickup.TotalPoints,
            Description: "Points from pickup #" + string(pickupID),
            ReferenceID: &pickupID,
        }
        u.transactionRepo.Create(transaction)
    }
    
    return nil
}

func (u *PickupUsecase) GetPendingPickups() ([]domain.PickupRequest, error) {
    return u.pickupRepo.GetPendingPickups()
}

func (u *PickupUsecase) validateStatusTransition(currentStatus, newStatus string) error {
    validTransitions := map[string][]string{
        "pending":     {"scheduled", "cancelled"},
        "scheduled":   {"in_progress", "cancelled"},
        "in_progress": {"completed", "cancelled"},
        "completed":   {},
        "cancelled":   {},
    }
    
    allowedStatuses, ok := validTransitions[currentStatus]
    if !ok {
        return errors.New("invalid current status")
    }
    
    for _, status := range allowedStatuses {
        if status == newStatus {
            return nil
        }
    }
    
    return errors.New("invalid status transition")
}