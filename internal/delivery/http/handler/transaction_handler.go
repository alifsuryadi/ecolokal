package handler

import (
	"net/http"

	"github.com/alifsuryadi/ecolokal/internal/domain"
	"github.com/alifsuryadi/ecolokal/internal/usecase"
	"github.com/alifsuryadi/ecolokal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TransactionHandler struct {
    transactionUsecase *usecase.TransactionUsecase
    validator          *validator.Validate
}

func NewTransactionHandler(transactionUsecase *usecase.TransactionUsecase) *TransactionHandler {
    return &TransactionHandler{
        transactionUsecase: transactionUsecase,
        validator:          validator.New(),
    }
}

// @Summary Create transaction
// @Description Create point transaction (admin only)
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateTransaction true "Transaction request"
// @Success 201 {object} utils.Response{data=domain.Transaction}
// @Failure 400 {object} utils.Response
// @Router /api/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
    var req domain.CreateTransaction
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    if err := h.validator.Struct(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
        return
    }
    
    transaction, err := h.transactionUsecase.CreateTransaction(&req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    utils.SuccessResponse(c, http.StatusCreated, "Transaction created", transaction)
}

// @Summary Get user transactions
// @Description Get transaction history for current user
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]domain.Transaction}
// @Router /api/transactions/my [get]
func (h *TransactionHandler) GetMyTransactions(c *gin.Context) {
    userID, _ := c.Get("userID")
    
    transactions, err := h.transactionUsecase.GetUserTransactions(userID.(int))
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get transactions")
        return
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Transactions retrieved", transactions)
}