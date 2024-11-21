package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Transaction representa una operación de compra o venta.
type Transaction struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Coin      string    `gorm:"type:text;not null"` // Ejemplo: BTC, SOL
	Amount    float64   `gorm:"type:numeric;not null"`
	Price     float64   `gorm:"type:numeric;not null"`
	Timestamp time.Time `gorm:"autoCreateTime"`
}

// NewTransaction crea una nueva transacción.
func NewTransaction(userID uuid.UUID, coin string, amount, price float64) *Transaction {
	return &Transaction{
		ID:        uuid.New(),
		UserID:    userID,
		Coin:      coin,
		Amount:    amount,
		Price:     price,
		Timestamp: time.Now(),
	}
}

// BeforeCreate es un hook de GORM que se ejecuta antes de insertar una nueva transacción.
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}

// TransactionRepository define las operaciones para trabajar con transacciones.
type TransactionRepository interface {
	Save(transaction *Transaction) error
	FindByUserID(userID uuid.UUID) ([]Transaction, error)
}

// GormTransactionRepository implementa TransactionRepository usando GORM.
type GormTransactionRepository struct {
	DB *gorm.DB
}

// NewTransactionRepository crea una nueva instancia de GormTransactionRepository.
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &GormTransactionRepository{DB: db}
}

// Save guarda una nueva transacción en la base de datos.
func (r *GormTransactionRepository) Save(transaction *Transaction) error {
	if err := r.DB.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

// FindByUserID encuentra todas las transacciones asociadas a un usuario dado.
func (r *GormTransactionRepository) FindByUserID(userID uuid.UUID) ([]Transaction, error) {
	var transactions []Transaction
	if err := r.DB.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
