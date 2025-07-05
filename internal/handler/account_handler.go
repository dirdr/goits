package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dirdr/goits/internal/service"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService service.AccountService
	log            *slog.Logger
}

func NewAccountHandler(accountService service.AccountService, log *slog.Logger) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		log:            log,
	}
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Creates a new account with a specified ID and initial balance.
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body CreateAccountRequest true "Account creation request"
// @Success 201 {string} string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Invalid request body for CreateAccount", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountService.CreateAccount(c.Request.Context(), uint(req.AccountID), req.InitialBalance)
	if err != nil {
		h.log.Error("Failed to create account", "account_id", req.AccountID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Account created successfully", "account_id", account.ID)
	c.Status(http.StatusCreated)
}

// GetAccount godoc
// @Summary Get account by ID
// @Description Retrieves an account's details and current balance by its ID.
// @Tags accounts
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID"
// @Success 200 {object} GetAccountResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /accounts/{account_id} [get]
func (h *AccountHandler) GetAccount(c *gin.Context) {
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil || accountID == 0 {
		h.log.Error("Invalid account ID format - must be a positive integer", "account_id", accountIDStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID must be a positive integer"})
		return
	}

	account, err := h.accountService.GetAccountByID(c.Request.Context(), uint(accountID))
	if err != nil {
		h.log.Error("Failed to get account", "account_id", accountIDStr, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account == nil {
		h.log.Info("Account not found", "account_id", accountIDStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	balance, err := h.accountService.GetAccountBalance(c.Request.Context(), uint(accountID))
	if err != nil {
		h.log.Error("Failed to get account balance", "account_id", accountIDStr, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if balance == nil {
		h.log.Info("Account balance not found", "account_id", accountIDStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "Account balance not found"})
		return
	}

	res := GetAccountResponse{
		AccountID: account.ID,
		Balance:   balance.Balance,
		Version:   balance.Version,
		UpdatedAt: balance.UpdatedAt,
	}

	h.log.Info("Account retrieved successfully", "account_id", account.ID)
	c.JSON(http.StatusOK, res)
}
