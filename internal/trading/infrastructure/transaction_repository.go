package infrastructure

import (
	"cryptoproject/internal/trading/domain"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormTransactionRepository implementa la interfaz TransactionRepository usando GORM.
type GormTransactionRepository struct {
	DB *gorm.DB
}

// NewTransactionRepository crea una nueva instancia de GormTransactionRepository.
func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &GormTransactionRepository{DB: db}
}

// Save guarda una nueva transacción en la base de datos.
func (r *GormTransactionRepository) Save(transaction *domain.Transaction) error {
	if err := r.DB.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

// FindByUserID recupera las transacciones realizadas por un usuario específico.
func (r *GormTransactionRepository) FindByUserID(userID uuid.UUID) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.DB.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return transactions, nil
}

// FindAll devuelve todas las transacciones (opcional para extensibilidad futura).
func (r *GormTransactionRepository) FindAll() ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	if err := r.DB.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// CalculateCryptoHoldings calcula el balance total por criptomoneda para un usuario.
func (r *GormTransactionRepository) CalculateCryptoHoldings(userID uuid.UUID) (map[string]float64, error) {
	var transactions []domain.Transaction
	if err := r.DB.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}

	// Mapa para acumular las cantidades de cada criptomoneda
	cryptoHoldings := make(map[string]float64)
	for _, transaction := range transactions {
		cryptoHoldings[transaction.Coin] += transaction.Amount
	}

	return cryptoHoldings, nil
}
