package handler

import (
	"net/http"
	"strconv"

	"github.com/alifsuryadi/ecolokal/internal/domain"
	"github.com/alifsuryadi/ecolokal/internal/usecase"
	"github.com/alifsuryadi/ecolokal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type WasteTypeHandler struct {
    wasteTypeUsecase *usecase.WasteTypeUsecase
    validator        *validator.Validate
}

func NewWasteTypeHandler(wasteTypeUsecase *usecase.WasteTypeUsecase) *WasteTypeHandler {
    return &WasteTypeHandler{
        wasteTypeUsecase: wasteTypeUsecase,
        validator:        validator.New(),
    }
}

// @Summary Create waste type
// @Description Create new waste type (admin only)
// @Tags waste-types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.WasteTypeRequest true "Waste type request"
// @Success 201 {object} utils.Response{data=domain.WasteType}
// @Failure 400 {object} utils.Response
// @Router /api/waste-types [post]
func (h *WasteTypeHandler) Create(c *gin.Context) {
    var req domain.WasteTypeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    if err := h.validator.Struct(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
        return
    }
    
    wasteType, err := h.wasteTypeUsecase.Create(&req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusCreated, "Waste type created", wasteType)
}

// @Summary Get all waste types
// @Description Get all waste types
// @Tags waste-types
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]domain.WasteType}
// @Router /api/waste-types [get]
func (h *WasteTypeHandler) GetAll(c *gin.Context) {
    wasteTypes, err := h.wasteTypeUsecase.GetAll()
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get waste types")
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Waste types retrieved", wasteTypes)
}

// @Summary Get waste type by ID
// @Description Get waste type by ID
// @Tags waste-types
// @Produce json
// @Security BearerAuth
// @Param id path int true "Waste type ID"
// @Success 200 {object} utils.Response{data=domain.WasteType}
// @Failure 404 {object} utils.Response
// @Router /api/waste-types/{id} [get]
func (h *WasteTypeHandler) GetByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
        return
    }
    
    wasteType, err := h.wasteTypeUsecase.GetByID(id)
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "Waste type not found")
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Waste type retrieved", wasteType)
}

// @Summary Update waste type
// @Description Update waste type (admin only)
// @Tags waste-types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Waste type ID"
// @Param request body domain.WasteTypeRequest true "Waste type request"
// @Success 200 {object} utils.Response{data=domain.WasteType}
// @Failure 400 {object} utils.Response
// @Router /api/waste-types/{id} [put]
func (h *WasteTypeHandler) Update(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
        return
    }
    
    var req domain.WasteTypeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    if err := h.validator.Struct(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
        return
    }
    
    wasteType, err := h.wasteTypeUsecase.Update(id, &req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Waste type updated", wasteType)
}

// @Summary Delete waste type
// @Description Delete waste type (admin only)
// @Tags waste-types
// @Produce json
// @Security BearerAuth
// @Param id path int true "Waste type ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/waste-types/{id} [delete]
func (h *WasteTypeHandler) Delete(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
        return
    }
    
    err = h.wasteTypeUsecase.Delete(id)
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Waste type deleted", nil)
}

// @Summary Get active waste types
// @Description Get only active waste types
// @Tags waste-types
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]domain.WasteType}
// @Router /api/waste-types/active [get]
func (h *WasteTypeHandler) GetActiveTypes(c *gin.Context) {
    wasteTypes, err := h.wasteTypeUsecase.GetActiveTypes()
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get waste types")
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Active waste types retrieved", wasteTypes)
}