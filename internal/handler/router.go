package handler

import (
	"log/slog"

	"github.com/dirdr/goits/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRouter(
	accountService service.AccountService,
	transactionService service.TransactionService,
	integrityService service.IntegrityService,
	log *slog.Logger,
	db *gorm.DB,
) *gin.Engine {
	r := gin.Default()

	accountHandler := NewAccountHandler(accountService, log, db)
	transactionHandler := NewTransactionHandler(transactionService, log, db)
	integrityHandler := NewIntegrityHandler(integrityService, log, db)

	r.POST("/accounts", accountHandler.CreateAccount)
	r.GET("/accounts/:account_id", accountHandler.GetAccount)

	r.POST("/transactions", transactionHandler.CreateTransaction)

	r.GET("/integrity/check", integrityHandler.CheckIntegrity)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
