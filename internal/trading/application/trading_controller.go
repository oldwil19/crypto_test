package application

import (
	authDomain "cryptoproject/internal/auth/domain"
	marketInfra "cryptoproject/internal/market/infrastructure"
	tradingDomain "cryptoproject/internal/trading/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TradingController maneja operaciones simuladas de trading.
type TradingController struct {
	transactionRepo tradingDomain.TransactionRepository
	userRepo        authDomain.UserRepository
	coingecko       marketInfra.CoingeckoServiceInterface
}

// NewTradingController crea una nueva instancia de TradingController.
func NewTradingController(
	transactionRepo tradingDomain.TransactionRepository,
	userRepo authDomain.UserRepository,
	coingecko marketInfra.CoingeckoServiceInterface,
) *TradingController {
	return &TradingController{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		coingecko:       coingecko,
	}
}

// HandleBuy maneja la compra simulada de criptomonedas.
func (tc *TradingController) HandleBuy(c *gin.Context) {
	userID := c.GetString("user_id") // Recuperar ID del usuario desde el contexto JWT
	coin := c.PostForm("coin")
	amountStr := c.PostForm("amount")

	// Validar la cantidad
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La cantidad ingresada es inválida"})
		return
	}

	// Obtener el precio actual de la criptomoneda
	price, err := tc.coingecko.GetCurrentPrice(coin, "usd")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el precio actual"})
		return
	}

	// Calcular el costo total de la transacción
	totalCost := price * amount

	// Recuperar el usuario
	user, err := tc.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Verificar si el usuario tiene saldo suficiente
	if !user.IsBalanceSufficient(totalCost) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Saldo insuficiente"})
		return
	}

	// Actualizar saldo del usuario
	err = user.AdjustBalance(-totalCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error ajustando el saldo"})
		return
	}
	if err := tc.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el usuario"})
		return
	}

	// Convertir el userID a uuid.UUID
	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID de usuario inválido"})
		return
	}

	// Actualizar balance de criptomonedas
	if err := user.AdjustCryptoBalance(coin, amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Registrar la transacción
	transaction := tradingDomain.NewTransaction(userUUID, coin, amount, price)
	if err := tc.transactionRepo.Save(transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar la transacción"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Compra realizada con éxito",
		"user": gin.H{
			"balance":        user.Balance,
			"crypto_balance": user.CryptoHoldings,
		},
		"transaction": transaction,
	})
}

// HandleTransactionHistory devuelve el historial de transacciones de un usuario.
func (tc *TradingController) HandleTransactionHistory(c *gin.Context) {
	userID := c.GetString("user_id") // ID del usuario desde el contexto JWT

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	transactions, err := tc.transactionRepo.FindByUserID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el historial de transacciones"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// HandleBalance devuelve el balance actual del usuario.
func (tc *TradingController) HandleBalance(c *gin.Context) {
	userID := c.GetString("user_id") // Recuperar ID del usuario desde el contexto JWT

	user, err := tc.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Reconstruir balances de criptomonedas a partir de transacciones
	transactions, err := tc.transactionRepo.FindByUserID(uuid.MustParse(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el historial de transacciones"})
		return
	}

	cryptoBalances := make(map[string]float64)
	for _, tx := range transactions {
		cryptoBalances[tx.Coin] += tx.Amount
	}
	user.CryptoHoldings = cryptoBalances

	c.JSON(http.StatusOK, gin.H{
		"usd_balance":     user.Balance,
		"crypto_holdings": user.CryptoHoldings,
	})
}
