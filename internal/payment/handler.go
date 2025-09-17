package payment

import (
	constants "go-blocker/internal/const"
	"go-blocker/internal/storage"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service *PaymentService
}

func (h *Handler) Webhook(c *gin.Context) {
	var req constants.WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if !common.IsHexAddress(req.Address) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address"})
		return
	}

	obj, err := h.Service.Create(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
		return
	}

	storage.PaymentAddressStore.Set(req.Address, obj.ID)
	resp := constants.WebhookResponse{ID: obj.ID.String(), Status: obj.Status}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Status(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing payment ID"})
		return
	}

	paymentID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		return
	}

	payment, err := h.Service.repo.FindByID(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":            payment.ID.String(),
		"status":        payment.Status,
		"amount":        payment.Amount,
		"actual_amount": payment.ReceivedAmount,
		"currency":      payment.Currency,
		"address":       payment.Address,
		"created_at":    payment.CreatedAt,
		"expires_at":    payment.ExpiresAt,
		"txid":          payment.TxID,
		"stuck":         payment.IsStuck,
		"callback_url":  payment.CallbackURL,
	})
}

func (h *Handler) CheckTx(c *gin.Context) {
	var req constants.CheckTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if !common.IsHexAddress(req.Address) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address"})
		return
	}

	// service here

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
	// 	return
	// }
}
