package handler

import (
	"log/slog"

	"github.com/dirdr/goits/internal/service"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func GetRouter(
	accountService service.AccountService,
	transactionService service.TransactionService,
	log *slog.Logger,
) *gin.Engine {
	r := gin.Default()

	accountHandler := NewAccountHandler(accountService, log)
	transactionHandler := NewTransactionHandler(transactionService, log)

	r.POST("/accounts", accountHandler.CreateAccount)
	r.GET("/accounts/:account_id", accountHandler.GetAccount)

	r.POST("/transactions", transactionHandler.CreateTransaction)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
