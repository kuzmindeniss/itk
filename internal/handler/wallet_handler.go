package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kuzmindeniss/itk/internal/models"
	"github.com/kuzmindeniss/itk/internal/service"
)

type WalletHandler struct {
	service *service.WalletService
}

func NewWalletHandler(service *service.WalletService) *WalletHandler {
	return &WalletHandler{
		service: service,
	}
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	walletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	wallet, err := h.service.GetWalletByID(c, walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

type UpdateBalanceRequest struct {
	Amount        int32                `json:"amount" binding:"required"`
	WalletID      string               `json:"walletId" binding:"required"`
	OperationType models.OperationType `json:"operationType" binding:"required"`
}

func (h *WalletHandler) UpdateWalletBalance(c *gin.Context) {
	var req UpdateBalanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.OperationType != models.OperationDeposit && req.OperationType != models.OperationWithdraw {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation type"})
		return
	}

	var walletID uuid.UUID

	if req.WalletID != "" {
		var err error
		walletID, err = uuid.Parse(req.WalletID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
			return
		}
	}

	if req.OperationType == models.OperationWithdraw {
		req.Amount = -req.Amount
	}

	wallet, err := h.service.TopUpWalletBalance(c, walletID, int32(req.Amount))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet": gin.H{
			"id":      wallet.ID,
			"balance": wallet.Balance,
		},
	})
}
