package handler

import (
	"github.com/alifsuryadi/ecolokal/internal/usecase"
	"github.com/go-playground/validator/v10"
)

type PickupHandler struct {
    pickupUsecase *usecase.PickupUsecase
    validator     *validator.Validate
}

func NewPickupHandler(pickupUsecase *usecase.PickupUsecase) *PickupHandler {
    return &PickupHandler{
        pickupUsecase: pickupUsecase,
        validator:     validator.New(),
    }
}

// @Summary Create pickup request
// @Description Create new pickup request (warga only)
// @Tags pickups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreatePickupRequest true "Pickup request"
// @Success 201 {object} utils.Response{data=domain.PickupRequest}
// @Failure 400 {object} utils.Response
// @Router /api/pickups [post]
func (h *PickupHandler) CreatePickupRequest(c *gin.Context) {
    userID, _ := c.Get("userID")
    role, _ := c.Get("userRole")
    
    if role != "warga" {
        utils.ErrorResponse(c, http.StatusForbidden, "Only warga can create pickup requests")
        return
    }
    
    var req domain.CreatePickupRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    if err := h.validator.Struct(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
        return
    }
    
    pickup, err := h.pickupUsecase.CreatePickupRequest(userID.(int), &req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusCreated, "Pickup request created", pickup)
}

// @Summary Get pickup by ID
// @Description Get pickup request details
// @Tags pickups
// @Produce json
// @Security BearerAuth
// @Param id path int true "Pickup ID"
// @Success 200 {object} utils.Response{data=domain.PickupRequest}
// @Failure 404 {object} utils.Response
// @Router /api/pickups/{id} [get]
func (h *PickupHandler) GetPickupByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
        return
    }
    
    pickup, err := h.pickupUsecase.GetPickupByID(id)
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "Pickup not found")
        return
    }
    
    // Check access permission
    userID, _ := c.Get("userID")
    role, _ := c.Get("userRole")
    
    if role != "admin" && role != "petugas" && pickup.UserID != userID.(int) {
        utils.ErrorResponse(c, http.StatusForbidden, "Access denied")
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Pickup retrieved", pickup)
}

// @Summary Get user pickups
// @Description Get all pickup requests for current user
// @Tags pickups
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]domain.PickupRequest}
// @Router /api/pickups/my [get]
func (h *PickupHandler) GetMyPickups(c *gin.Context) {
    userID, _ := c.Get("userID")
    
    pickups, err := h.pickupUsecase.GetUserPickups(userID.(int))
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get pickups")
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Pickups retrieved", pickups)
}

// @Summary Get petugas pickups
// @Description Get pickups assigned to petugas for specific date
// @Tags pickups
// @Produce json
// @Security BearerAuth
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} utils.Response{data=[]domain.PickupRequest}
// @Router /api/pickups/petugas [get]
func (h *PickupHandler) GetPetugasPickups(c *gin.Context) {
    userID, _ := c.Get("userID")
    role, _ := c.Get("userRole")
    
    if role != "petugas" && role != "admin" {
        utils.ErrorResponse(c, http.StatusForbidden, "Access denied")
        return
    }
    
    date := c.Query("date")
    if date == "" {
        utils.ErrorResponse(c, http.StatusBadRequest, "Date parameter required")
        return
    }
    
    pickups, err := h.pickupUsecase.GetPetugasPickups(userID.(int), date)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Pickups retrieved", pickups)
}

// @Summary Update pickup status
// @Description Update pickup request status (admin/petugas)
// @Tags pickups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Pickup ID"
// @Param request body domain.UpdatePickupStatus true "Status update"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/pickups/{id}/status [put]
func (h *PickupHandler) UpdatePickupStatus(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
        return
    }
    
    var req domain.UpdatePickupStatus
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    if err := h.validator.Struct(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
        return
    }
    
    userID, _ := c.Get("userID")
    role, _ := c.Get("userRole")
    
    var petugasID *int
    if req.Status == "scheduled" && role == "admin" {
        // Admin assigning to petugas
        petugasID = nil
    } else if role == "petugas" {
        // Petugas updating their own pickup
        pid := userID.(int)
        petugasID = &pid
    }
    
    err = h.pickupUsecase.UpdatePickupStatus(id, &req, petugasID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Status updated", nil)
}

// @Summary Update pickup items
// @Description Update actual weight of pickup items (petugas only)
// @Tags pickups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Pickup ID"
// @Param request body domain.UpdatePickupItems true "Items update"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/pickups/{id}/items [put]
func (h *PickupHandler) UpdatePickupItems(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
        return
    }
    
    var req domain.UpdatePickupItems
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    if err := h.validator.Struct(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
        return
    }
    
    userID, _ := c.Get("userID")
    
    err = h.pickupUsecase.UpdatePickupItems(id, userID.(int), &req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Items updated", nil)
}

// @Summary Get pending pickups
// @Description Get all pending pickups (admin only)
// @Tags pickups
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]domain.PickupRequest}
// @Router /api/pickups/pending [get]
func (h *PickupHandler) GetPendingPickups(c *gin.Context) {
    pickups, err := h.pickupUsecase.GetPendingPickups()
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get pickups")
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Pending pickups retrieved", pickups)
}

// @Summary Assign pickup to petugas
// @Description Assign pickup request to petugas (admin only)
// @Tags pickups
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Pickup ID"
// @Param request body map[string]int true "Petugas assignment"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/pickups/{id}/assign [put]
func (h *PickupHandler) AssignPickupToPetugas(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
        return
    }
    
    var req map[string]int
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    petugasID, ok := req["petugas_id"]
    if !ok {
        utils.ErrorResponse(c, http.StatusBadRequest, "petugas_id required")
        return
    }
    
    status := domain.UpdatePickupStatus{Status: "scheduled"}
    err = h.pickupUsecase.UpdatePickupStatus(id, &status, &petugasID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Pickup assigned to petugas", nil)
}