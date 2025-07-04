package handler

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/accounts", createAccount)
	r.GET("/accounts/:account_id", getAccount)
	r.POST("/transactions", createTransaction)

	return r
}

func createAccount(c *gin.Context) {
	// Implementation to be added
}

func getAccount(c *gin.Context) {
	// Implementation to be added
}

func createTransaction(c *gin.Context) {
	// Implementation to be added
}
