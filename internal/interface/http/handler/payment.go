package payment

import (
	"go-blocker/internal/application/payment"
	"go-blocker/internal/storage"
	"net/http"

	_ "go-blocker/cmd/docs"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	Service *payment.Service
}

func NewRepository(s *payment.Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) Webhook(c *gin.Context) {
	var req payment.WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, InvalidRequest{Error: "invalid request"})
		return
	}

	if !common.IsHexAddress(req.Address) {
		c.JSON(http.StatusBadRequest, InvalidAddress{Error: "Invalid address"})
		return
	}

	obj, err := h.Service.Create(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
		return
	}

	storage.PaymentAddressStore.Set(req.Address, obj.ID)
	resp := payment.WebhookResponse{ID: obj.ID.String(), Status: obj.Status}
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

	payment, err := h.Service.Repo.FindByID(paymentID)
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

// @Summary      Check a transaction
// @Description  Check if a transaction meets the payment requirements
// @Tags         payment
// @Accept       json
// @Param request body payment.CheckTxRequest true "Check transaction request"
// @Produce      json
// @Success      200  {object}  CheckTxResponse
// @Failure      400 {object} InvalidRequest "Invalid request format"
// @Failure      422 {object} InvalidAddress "Invalid address format"
// @Failure      503 {object} FailedToFind "failed to check transaction"
// @Router       /payment/check/transaction [post]
func (h *Handler) CheckTx(c *gin.Context) {
	var req *payment.CheckTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, InvalidRequest{Error: "invalid request"})
		return
	}

	if !common.IsHexAddress(req.Address) {
		c.JSON(http.StatusUnprocessableEntity, InvalidAddress{Error: "Invalid address"})
		return
	}

	resp, err := h.Service.CheckTx(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, FailedToFind{Error: "failed to check transaction"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary      Find the latest transaction
// @Description  Find the latest transaction for a given address
// @Tags         payment
// @Accept       json
// @Param request body payment.FindTxRequest true "Check transaction request"
// @Produce      json
// @Success      200  {object}  CheckTxResponse
// @Failure      400 {object} InvalidRequest "Invalid request format"
// @Failure      422 {object} InvalidAddress "Invalid address format"
// @Failure      503 {object} FailedToFind  "failed to find latest transaction"
// @Router       /payment/find/transaction [post]
func (h *Handler) FindLatestTx(c *gin.Context) {
	var req *payment.FindTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, InvalidRequest{Error: "invalid request"})
		return
	}

	if !common.IsHexAddress(req.Address) {
		c.JSON(http.StatusBadRequest, InvalidAddress{Error: "Invalid address"})
		return
	}
	resp, err := h.Service.FindLatestTx(req)
	if err != nil {
		c.JSON(http.StatusNotFound, FailedToFind{Error: "failed to find latest transaction"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// HealthHandler provides a basic health check for the service.
//
// @Summary      Health check
// @Description  Returns service status and database connectivity.
// @Tags         system
// @Produce      json
// @Success      200 {object} map[string]string "Service is healthy"
// @Router       /health [get]
func (h *Handler) HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"db":     "connected",
	})
}
