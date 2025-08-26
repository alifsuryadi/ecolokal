package handler

import (
	"net/http"

	"github.com/alifsuryadi/ecolokal/internal/usecase"
	"github.com/alifsuryadi/ecolokal/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
    userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
    return &UserHandler{userUsecase: userUsecase}
}

// @Summary Get user profile
// @Description Get current user profile
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=domain.User}
// @Failure 401 {object} utils.Response
// @Router /api/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        utils.ErrorResponse(c, http.StatusUnauthorized, "User ID not found")
        return
    }
    
    user, err := h.userUsecase.GetUserByID(userID.(int))
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "User not found")
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "User profile retrieved", user)
}

// @Summary Get user points
// @Description Get current user points balance
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=domain.UserPoint}
// @Failure 401 {object} utils.Response
// @Router /api/users/points [get]
func (h *UserHandler) GetUserPoints(c *gin.Context) {
    userID, _ := c.Get("userID")
    
    points, err := h.userUsecase.GetUserPoints(userID.(int))
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, "Points not found")
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "User points retrieved", points)
}

// @Summary Get users by role
// @Description Get all users with specific role (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param role path string true "User role (warga/petugas/admin)"
// @Success 200 {object} utils.Response{data=[]domain.User}
// @Failure 403 {object} utils.Response
// @Router /api/users/role/{role} [get]
func (h *UserHandler) GetUsersByRole(c *gin.Context) {
    role := c.Param("role")
    
    users, err := h.userUsecase.GetUsersByRole(role)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get users")
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Users retrieved", users)
}