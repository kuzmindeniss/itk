package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kuzmindeniss/itk/internal/handler"
)

func SetupRouter(walletHandler *handler.WalletHandler) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	v1.POST("/wallet", walletHandler.UpdateWalletBalance)
	v1.GET("/wallets/:id", walletHandler.GetWallet)

	return r
}
