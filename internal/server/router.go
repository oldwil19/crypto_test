package server

import (
	"cryptoproject/docs"
	accountApp "cryptoproject/internal/account/application" // Añadimos esta línea
	"cryptoproject/internal/auth/application"
	"cryptoproject/internal/auth/infrastructure"
	marketApp "cryptoproject/internal/market/application"
	tradingApp "cryptoproject/internal/trading/application"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	authController *application.AuthController,
	marketController *marketApp.MarketController,
	registerController *application.RegisterController,
	tradingController *tradingApp.TradingController,
	accountController *accountApp.AccountController, // Añadimos AccountController aquí
	jwtMiddleware *infrastructure.JWTMiddleware,
) *gin.Engine {
	docs.SwaggerInfo.Title = "Crypto API"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	r := gin.New()
	r.Use(gin.Recovery())

	// Endpoints públicos
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/register", registerController.Register)
	r.POST("/auth/login", authController.Login)

	// Endpoints protegidos
	protected := r.Group("/")
	protected.Use(jwtMiddleware.Middleware())

	// Mercado
	protected.GET("/market/:id/price", marketController.GetCurrentPriceHandler)
	protected.GET("/market/:id/history", marketController.GetHistoricalPricesHandler)

	// Trading
	protected.POST("/trading/buy", tradingController.HandleBuy)
	protected.GET("/trading/history", tradingController.HandleTransactionHistory)
	protected.GET("/trading/balance", tradingController.HandleBalance)

	// Account
	protected.POST("/account/balance/add", accountController.HandleAddBalance) // Añadimos este endpoint

	return r
}
