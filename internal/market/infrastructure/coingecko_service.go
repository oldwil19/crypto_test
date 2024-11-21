package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"cryptoproject/pkg/logger"
)

// CoingeckoServiceInterface define los métodos que usamos para interactuar con CoinGecko.
// Honestamente, esto asegura flexibilidad, pero no significa que sea 100% perfecto.
type CoingeckoServiceInterface interface {
	GetCurrentPrice(crypto string, currency string) (float64, error)
	GetHistoricalPrices(crypto string, start string, end string) ([]map[string]interface{}, error)
	CheckAPIStatus() bool
}

// CoingeckoService estructura el servicio de integración con CoinGecko.
type CoingeckoService struct {
	baseURL      string
	client       *http.Client
	rateLimitMax int
	rateLimit    time.Duration
	lastRequest  time.Time
	mu           sync.Mutex
	cachedCoins  []string
	cacheTime    time.Time
	cacheTTL     time.Duration
}

// Estas variables globales son funcionales, pero no me encantan.
// Quizás podamos mover esto a algo más limpio en el futuro.
var (
	once                 sync.Once
	coingeckoServiceInst *CoingeckoService
)

// NewCoingeckoService inicializa el servicio de CoinGecko usando el patrón Singleton.
// Ojo: Este patrón funciona aquí, pero no lo abuses en otras partes, se pone feo si no es necesario.
func NewCoingeckoService() CoingeckoServiceInterface {
	once.Do(func() {
		timeout, err := time.ParseDuration(getEnv("COINGECKO_TIMEOUT", "5s"))
		if err != nil {
			logger.Warn("Error al parsear COINGECKO_TIMEOUT, usando valor por defecto: 5s")
			timeout = 5 * time.Second
		}

		rateLimitMax, err := strconv.Atoi(getEnv("COINGECKO_RATE_LIMIT", "50"))
		if err != nil {
			logger.Warn("Error al parsear COINGECKO_RATE_LIMIT, usando valor por defecto: 50")
			rateLimitMax = 50
		}

		coingeckoServiceInst = &CoingeckoService{
			baseURL:      getEnv("COINGECKO_BASE_URL", "https://api.coingecko.com/api/v3"),
			client:       &http.Client{Timeout: timeout},
			rateLimitMax: rateLimitMax,
			rateLimit:    time.Minute / time.Duration(rateLimitMax),
			lastRequest:  time.Time{},
			cachedCoins:  nil,
			cacheTime:    time.Time{},
			cacheTTL:     24 * time.Hour, // Esto es mucho tiempo, pero para el caso está bien.
		}
	})
	return coingeckoServiceInst
}

// enforceRateLimit asegura que respetemos los límites de la API.
// Ojo aquí: Si empiezan a aparecer retrasos por el sleep, quizás tengamos que revisar.
func (s *CoingeckoService) enforceRateLimit() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if s.lastRequest.IsZero() || now.Sub(s.lastRequest) >= s.rateLimit {
		s.lastRequest = now
		return
	}

	timeToWait := s.rateLimit - now.Sub(s.lastRequest)
	time.Sleep(timeToWait)
	s.lastRequest = time.Now()
}

// retryPolicy maneja reintentos para las solicitudes a la API.
// Honestamente, este número de reintentos es arbitrario. Quizás sea mejor configurable.
func (s *CoingeckoService) retryPolicy(requestFunc func() (*http.Response, error)) (*http.Response, error) {
	const maxRetries = 3
	var err error
	var resp *http.Response

	for attempts := 0; attempts < maxRetries; attempts++ {
		resp, err = requestFunc()
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		// Ojo: El sleep aquí hace que los reintentos sean lentos si hay muchas solicitudes fallidas.
		time.Sleep(time.Second * time.Duration(attempts+1))
	}
	return nil, errors.New("error tras múltiples intentos: " + err.Error())
}

// CheckAPIStatus revisa si la API de CoinGecko está operativa.
// Este método funciona, pero no es el más eficiente.
// Quizás en el futuro podamos implementar algo más ligero.
func (s *CoingeckoService) CheckAPIStatus() bool {
	s.enforceRateLimit()
	url := fmt.Sprintf("%s/ping", s.baseURL)

	resp, err := s.retryPolicy(func() (*http.Response, error) {
		return s.client.Get(url)
	})
	if err != nil {
		logger.Error("Error al verificar el estado de CoinGecko:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// GetCurrentPrice obtiene el precio actual de una criptomoneda.
// Aquí no hay magia: si la API falla, no hay mucho que podamos hacer.
func (s *CoingeckoService) GetCurrentPrice(crypto, currency string) (float64, error) {
	if !s.CheckAPIStatus() {
		return 0, fmt.Errorf("API de CoinGecko no está disponible")
	}

	s.enforceRateLimit()
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=%s", s.baseURL, crypto, currency)

	resp, err := s.retryPolicy(func() (*http.Response, error) {
		return s.client.Get(url)
	})
	if err != nil {
		logger.Error("Error al realizar solicitud a CoinGecko:", err)
		return 0, fmt.Errorf("fallo en la solicitud a CoinGecko: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API de CoinGecko devolvió estado: %d", resp.StatusCode)
	}

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Error("Error al decodificar respuesta de CoinGecko:", err)
		return 0, fmt.Errorf("error al decodificar JSON: %w", err)
	}

	if _, exists := data[crypto]; !exists {
		return 0, fmt.Errorf("moneda '%s' no soportada o no disponible en CoinGecko", crypto)
	}
	if price, ok := data[crypto][currency]; ok {
		return price, nil
	}

	return 0, fmt.Errorf("precio para '%s' en '%s' no encontrado", crypto, currency)
}

// GetHistoricalPrices obtiene precios históricos de una criptomoneda.
// Esto está bien para ahora, pero si las fechas son largas, los datos se vuelven enormes.
func (s *CoingeckoService) GetHistoricalPrices(crypto, start, end string) ([]map[string]interface{}, error) {
	if !s.CheckAPIStatus() {
		return nil, fmt.Errorf("API de CoinGecko no está disponible")
	}

	s.enforceRateLimit()
	url := fmt.Sprintf("%s/coins/%s/market_chart/range?vs_currency=usd&from=%s&to=%s", s.baseURL, crypto, start, end)

	resp, err := s.retryPolicy(func() (*http.Response, error) {
		return s.client.Get(url)
	})
	if err != nil {
		logger.Error("Error al realizar solicitud a CoinGecko:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Prices [][]float64 `json:"prices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logger.Error("Error al decodificar respuesta de CoinGecko:", err)
		return nil, err
	}

	historicalPrices := []map[string]interface{}{}
	for _, price := range response.Prices {
		if len(price) == 2 {
			historicalPrices = append(historicalPrices, map[string]interface{}{
				"timestamp": price[0],
				"price":     price[1],
			})
		}
	}
	return historicalPrices, nil
}

// getEnv devuelve un valor de variable de entorno o un valor por defecto.
// Lo más simple, pero evita muchas sorpresas.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
