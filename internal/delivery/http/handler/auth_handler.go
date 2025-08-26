package handler

import (
	"net/http"

	"github.com/alifsuryadi/ecolokal/internal/domain"
	"github.com/alifsuryadi/ecolokal/internal/usecase"
	"github.com/alifsuryadi/ecolokal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
    authUsecase *usecase.AuthUsecase
    validator   *validator.Validate
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
    return &AuthHandler{
        authUsecase: authUsecase,
        validator:   validator.New(),
    }
}

// @Summary Register new user
// @Description Register a new user with role warga or petugas
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Register request"
// @Success 201 {object} utils.Response{data=domain.User}
// @Failure 400 {object} utils.Response
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    var req domain.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    if err := h.validator.Struct(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
        return
    }
    
    user, err := h.authUsecase.Register(&req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", user)
}

// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login request"
// @Success 200 {object} utils.Response{data=domain.LoginResponse}
// @Failure 400 {object} utils.Response
// @Router /api/users/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    var req domain.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    if err := h.validator.Struct(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
        return
    }
    
    response, err := h.authUsecase.Login(&req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}