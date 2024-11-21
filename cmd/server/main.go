package main

import (
	accountApp "cryptoproject/internal/account/application"
	accountInfra "cryptoproject/internal/account/infrastructure"
	"cryptoproject/internal/auth/application"
	"cryptoproject/internal/auth/domain"
	"cryptoproject/internal/auth/infrastructure"
	marketApp "cryptoproject/internal/market/application"
	marketInfra "cryptoproject/internal/market/infrastructure"
	"cryptoproject/internal/server"
	tradingApp "cryptoproject/internal/trading/application"
	tradingDomain "cryptoproject/internal/trading/domain"
	tradingInfra "cryptoproject/internal/trading/infrastructure"
	"cryptoproject/pkg/config"
	"cryptoproject/pkg/logger"
	"os"
	"time"

	"gorm.io/gorm"
)

func main() {
	config.LoadConfig()
	logger.InitLogger()

	db := initializeDatabase()
	defer closeDatabase(db)

	if err := runMigrations(db); err != nil {
		logger.Error("Error ejecutando migraciones:", err)
		return
	}

	jwtService := infrastructure.NewJWTService(os.Getenv("JWT_SECRET"), 24*time.Hour)
	jwtMiddleware := infrastructure.NewJWTMiddleware(jwtService)

	authController := initializeAuthController(db, jwtService)
	registerController := initializeRegisterController(db)
	marketController := initializeMarketController()
	tradingController := initializeTradingController(db)
	accountController := initializeAccountController(db)

	router := server.SetupRouter(authController, marketController, registerController, tradingController, accountController, jwtMiddleware)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Servidor iniciado en el puerto:", port)
	if err := router.Run(":" + port); err != nil {
		logger.Error("Error al iniciar el servidor:", err)
	}
}

// Inicializa la conexión a la base de datos. Asegúrate de configurar bien tu conexión.
func initializeDatabase() *gorm.DB {
	return accountInfra.ConnectDatabase()
}

// Cierra la conexión con la base de datos de manera segura.
func closeDatabase(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Error al obtener la conexión de la base de datos:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		logger.Error("Error al cerrar la conexión de la base de datos:", err)
	}
}

// Ejecuta migraciones. Por ahora, agregamos las necesarias para las entidades clave.
func runMigrations(db *gorm.DB) error {
	logger.Info("Ejecutando migraciones...")
	// Esta lógica depende de la base de datos que estés usando. Asegúrate de que esté configurada correctamente.
	return db.AutoMigrate(&domain.User{}, &tradingDomain.Transaction{})
}

// Configura el controlador de autenticación.
func initializeAuthController(db *gorm.DB, jwtService infrastructure.JWTServiceInterface) *application.AuthController {
	userRepo := infrastructure.NewUserRepository(db)
	return application.NewAuthController(jwtService, userRepo)
}

// Configura el controlador de registro.
func initializeRegisterController(db *gorm.DB) *application.RegisterController {
	userRepo := infrastructure.NewUserRepository(db)
	return application.NewRegisterController(userRepo)
}

// Configura el controlador de mercado.
func initializeMarketController() *marketApp.MarketController {
	coingeckoService := marketInfra.NewCoingeckoService()
	return marketApp.NewMarketController(coingeckoService)
}

// Configura el controlador de trading.
func initializeTradingController(db *gorm.DB) *tradingApp.TradingController {
	transactionRepo := tradingInfra.NewTransactionRepository(db)
	userRepo := infrastructure.NewUserRepository(db)
	coingeckoService := marketInfra.NewCoingeckoService()
	return tradingApp.NewTradingController(transactionRepo, userRepo, coingeckoService)
}

// Configura el controlador de cuentas.
func initializeAccountController(db *gorm.DB) *accountApp.AccountController {
	userRepo := infrastructure.NewUserRepository(db)
	return accountApp.NewAccountController(userRepo)
}
