package handler

import (
	"log/slog"
	"net/http"

	"github.com/dirdr/goits/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IntegrityCheckHandler struct {
	integrityCheckService service.IntegrityCheckService
	log                   *slog.Logger
	db                    *gorm.DB
}

func NewIntegrityCheckHandler(integrityCheckService service.IntegrityCheckService, log *slog.Logger, db *gorm.DB) *IntegrityCheckHandler {
	return &IntegrityCheckHandler{
		integrityCheckService: integrityCheckService,
		log:                   log,
		db:                    db,
	}
}

// CheckIntegrity godoc
// @Summary Check double bookkeeping integrity
// @Description Verifies that the total debits equal total credits in the journal entries to ensure double bookkeeping integrity.
// @Tags integrity
// @Accept json
// @Produce json
// @Success 200 {object} service.IntegrityCheckResult "Integrity check result"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /integrity/check [get]
func (h *IntegrityCheckHandler) CheckIntegrity(c *gin.Context) {
	result, err := h.integrityCheckService.VerifyDoubleBookkeeping(c.Request.Context())
	if err != nil {
		h.log.Error("Failed to verify double bookkeeping integrity", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.IsValid {
		h.log.Info("Double bookkeeping integrity verified successfully",
			"total_debits", result.TotalDebits,
			"total_credits", result.TotalCredits)
	} else {
		h.log.Warn("Double bookkeeping integrity check failed",
			"total_debits", result.TotalDebits,
			"total_credits", result.TotalCredits,
			"difference", result.Difference)
	}

	c.JSON(http.StatusOK, result)
}
