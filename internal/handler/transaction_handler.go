package handler

import (
	"log/slog"
	"net/http"

	"github.com/dirdr/goits/internal/service"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService service.TransactionService
	log                *slog.Logger
}

func NewTransactionHandler(transactionService service.TransactionService, log *slog.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		log:                log,
	}
}

// CreateTransaction handles the submission of a new transaction.
// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Processes a transfer of funds between two accounts.
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body CreateTransactionRequest true "Transaction creation request"
// @Success 201 {string} string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Invalid request body for CreateTransaction", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.transactionService.ProcessTransfer(c.Request.Context(), req.SourceAccountID, req.DestinationAccountID, req.Amount)
	if err != nil {
		h.log.Error("Failed to process transaction", "source_account_id", req.SourceAccountID, "destination_account_id", req.DestinationAccountID, "amount", req.Amount, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Transaction processed successfully", "source_account_id", req.SourceAccountID, "destination_account_id", req.DestinationAccountID, "amount", req.Amount)
	c.Status(http.StatusCreated)
}
