package application

import (
	"fmt"
	"net/http"
	"time"

	"cryptoproject/internal/market/infrastructure"
	"cryptoproject/pkg/logger"

	"github.com/gin-gonic/gin"
)

// MarketController maneja todo lo que tiene que ver con datos de mercado.
// Ojo, aquí se usan los métodos del servicio CoinGecko para no inventar la rueda.
type MarketController struct {
	coingeckoService infrastructure.CoingeckoServiceInterface
}

// NewMarketController inicializa el controlador de mercado.
// Pana, este es el constructor, aquí simplemente conectamos con el servicio.
func NewMarketController(service infrastructure.CoingeckoServiceInterface) *MarketController {
	return &MarketController{
		coingeckoService: service,
	}
}

// GetCurrentPriceHandler saca el precio actual de una criptomoneda.
// Aquí usamos el ID de la moneda y la moneda de cambio para obtener el precio.
func (mc *MarketController) GetCurrentPriceHandler(c *gin.Context) {
	crypto := c.Param("id")                       // Esto es el ID de la cripto, tipo "bitcoin" o "ethereum".
	currency := c.DefaultQuery("currency", "usd") // Por defecto trabajamos con USD, pero se puede cambiar.

	price, err := mc.coingeckoService.GetCurrentPrice(crypto, currency)
	if err != nil {
		logger.Error("Error al obtener el precio actual:", err)
		// Pendiente aquí: si falla CoinGecko, devolvemos error, pero quizá podríamos poner un cache para no depender tanto.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el precio actual"})
		return
	}

	// Aquí devolvemos la info al cliente.
	c.JSON(http.StatusOK, gin.H{
		"crypto":   crypto,
		"currency": currency,
		"price":    price,
	})
}

// GetHistoricalPricesHandler obtiene los precios históricos de una cripto.
// Se espera un rango de fechas, pero OJO: el formato tiene que ser específico.
func (mc *MarketController) GetHistoricalPricesHandler(c *gin.Context) {
	cryptoID := c.Param("id") // ID de la cripto, por ejemplo, "bitcoin".
	start := c.Query("start") // Inicio del rango en dd-mm-yyyy.
	end := c.Query("end")     // Fin del rango en el mismo formato.

	// Aquí validamos las fechas. Demasiado importante que el formato sea el correcto.
	startUnix, err := parseDateToUnix(start)
	if err != nil {
		// Pendiente: Deberíamos ser más claros con el formato esperado si esto pasa mucho.
		c.JSON(http.StatusBadRequest, gin.H{"error": "El formato de la fecha de inicio debe ser dd-mm-yyyy"})
		return
	}

	endUnix, err := parseDateToUnix(end)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El formato de la fecha de fin debe ser dd-mm-yyyy"})
		return
	}

	// Pedimos los datos históricos al servicio. Esto podría demorar si son muchos días.
	historicalData, err := mc.coingeckoService.GetHistoricalPrices(cryptoID, fmt.Sprintf("%d", startUnix), fmt.Sprintf("%d", endUnix))
	if err != nil {
		// Ojo: Si hay problemas aquí, seguro es un tema con la API de CoinGecko o con los datos enviados.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mandamos los datos históricos al cliente.
	c.JSON(http.StatusOK, historicalData)
}

// parseDateToUnix convierte una fecha (texto) en un UNIX timestamp.
// ¡Pendiente! Si alguien manda mal el formato, esto devuelve error de una.
func parseDateToUnix(date string) (int64, error) {
	// Intentamos convertir con el formato dd-mm-yyyy.
	parsedDate, err := time.Parse("02-01-2006", date)
	if err == nil {
		return parsedDate.Unix(), nil
	}

	// Si el primer intento falla, probamos con el formato RFC3339 (ISO).
	parsedDate, err = time.Parse(time.RFC3339, date)
	if err == nil {
		return parsedDate.Unix(), nil
	}

	// Aquí no hay más nada que hacer si fallan los dos formatos.
	return 0, fmt.Errorf("la fecha no cumple con los formatos válidos: dd-mm-yyyy o RFC3339")
}
