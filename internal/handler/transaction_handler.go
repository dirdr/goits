package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dirdr/goits/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	transactionService service.TransactionService
	log                *slog.Logger
	db                 *gorm.DB
}

const (
	maxRetries = 3
	baseDelay  = 10 * time.Millisecond
)

func NewTransactionHandler(transactionService service.TransactionService, log *slog.Logger, db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		log:                log,
		db:                 db,
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

	err := h.processTransactionWithRetry(c, req)
	if err != nil {
		h.log.Error("Failed to process transaction", "source_account_id", req.SourceAccountID, "destination_account_id", req.DestinationAccountID, "amount", req.Amount, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Transaction processed successfully", "source_account_id", req.SourceAccountID, "destination_account_id", req.DestinationAccountID, "amount", req.Amount)
	c.Status(http.StatusCreated)
}

func (h *TransactionHandler) processTransactionWithRetry(c *gin.Context, req CreateTransactionRequest) error {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
			return h.transactionService.ProcessTransfer(c.Request.Context(), tx, req.SourceAccountID, req.DestinationAccountID, req.Amount)
		})

		if err == nil {
			return nil
		}

		lastErr = err

		if !h.isRetryableError(err) {
			return err
		}

		if attempt == maxRetries-1 {
			return fmt.Errorf("transaction failed after %d attempts due to concurrent modifications: %w", maxRetries, lastErr)
		}

		delay := h.calculateBackoffDelay(attempt)
		h.log.Debug("Retrying transaction after optimistic locking conflict",
			"attempt", attempt+1,
			"delay", delay,
			"source_account_id", req.SourceAccountID,
			"destination_account_id", req.DestinationAccountID,
			"error", err)

		select {
		case <-c.Request.Context().Done():
			return c.Request.Context().Err()
		case <-time.After(delay):
		}
	}

	return fmt.Errorf("unexpected error: retry loop exited without return")
}

func (h *TransactionHandler) isRetryableError(err error) bool {
	return strings.Contains(err.Error(), "optimistic locking failed")
}

func (h *TransactionHandler) calculateBackoffDelay(attempt int) time.Duration {
	return time.Duration(1<<attempt) * baseDelay
}
