package domain

import (
	"cryptoproject/pkg/logger"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User representa un usuario en el sistema. Nota: el campo CryptoHoldings no se persiste en la base de datos.
type User struct {
	ID             string             `gorm:"type:uuid;primaryKey"`
	Username       string             `gorm:"type:varchar(100);unique;not null"`
	PasswordHash   string             `gorm:"type:text;not null"`
	Balance        float64            `gorm:"type:numeric(15,2);default:1000.00"`
	CryptoHoldings map[string]float64 `gorm:"-"` // Aquí usamos `gorm:"-"` para evitar persistir este campo.
	CreatedAt      time.Time          `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt      time.Time          `gorm:"type:timestamp;autoUpdateTime"`
}

// NewUser crea una nueva instancia de usuario con una contraseña encriptada.
// Nota: el saldo inicial y las criptomonedas pueden ajustarse si cambian las reglas de negocio.
func NewUser(username, password string) (*User, error) {
	if username == "" || password == "" {
		return nil, errors.New("nombre de usuario y contraseña son obligatorios")
	}

	id := uuid.New().String()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           id,
		Username:     username,
		PasswordHash: string(passwordHash),
		Balance:      1000.00,
		CryptoHoldings: map[string]float64{
			"btc":  0.0,
			"sol":  0.0,
			"doge": 0.0, // Ejemplo inicial, pero esto puede configurarse de forma dinámica.
		},
	}, nil
}

// VerifyPassword compara una contraseña con el hash almacenado.
// Funciona bien, aunque no se necesita optimización inmediata.
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

// AdjustBalance ajusta el saldo del usuario en USD.
// Ojo: Si el monto es negativo y el balance no alcanza, retorna un error.
func (u *User) AdjustBalance(amount float64) error {
	if u.Balance+amount < 0 {
		return errors.New("saldo insuficiente en USD")
	}
	u.Balance += amount
	logger.Info(fmt.Sprintf("Saldo ajustado. Nuevo balance: %.2f USD", u.Balance))
	return nil
}

// AdjustCryptoBalance ajusta o inicializa el balance de una criptomoneda.
// Pendiente: esto podría registrar más detalles, por ejemplo, la operación que generó el cambio.
func (u *User) AdjustCryptoBalance(crypto string, amount float64) error {
	if currentBalance, exists := u.CryptoHoldings[crypto]; exists {
		// Si la criptomoneda ya existe, ajustamos el balance.
		if currentBalance+amount < 0 {
			return fmt.Errorf("balance insuficiente para %s: actual %.2f, ajuste %.2f", crypto, currentBalance, amount)
		}
		u.CryptoHoldings[crypto] += amount
		logger.Info(fmt.Sprintf("Balance ajustado para %s. Nuevo balance: %.2f", crypto, u.CryptoHoldings[crypto]))
	} else {
		// Si la criptomoneda no existe, la inicializamos.
		if amount < 0 {
			return fmt.Errorf("no se puede inicializar balance negativo para %s", crypto)
		}
		u.CryptoHoldings[crypto] = amount
		logger.Info(fmt.Sprintf("Nueva criptomoneda añadida: %s con balance %.2f", crypto, amount))
	}
	return nil
}

// GetCryptoBalance devuelve el balance de una criptomoneda específica.
// Ojo: Si no existe en el mapa, se retorna un error.
func (u *User) GetCryptoBalance(crypto string) (float64, error) {
	if balance, exists := u.CryptoHoldings[crypto]; exists {
		return balance, nil
	}
	logger.Warn(fmt.Sprintf("Intento de acceso a criptomoneda no existente: %s", crypto))
	return 0, fmt.Errorf("la criptomoneda %s no está en los holdings", crypto)
}

// IsBalanceSufficient verifica si el usuario tiene saldo suficiente en USD.
// Simple, pero cumple con lo necesario por ahora.
func (u *User) IsBalanceSufficient(amount float64) bool {
	return u.Balance >= amount
}

// BeforeCreate es un hook de GORM para generar un UUID antes de insertar un nuevo usuario.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
		logger.Info("Se generó un nuevo UUID para el usuario.")
	}
	return
}

// AddBalance suma saldo al balance en USD del usuario.
// Nota: esto podría incluir un registro más detallado para auditorías.
func (u *User) AddBalance(amount float64) error {
	if amount <= 0 {
		return errors.New("el monto a añadir debe ser positivo")
	}
	u.Balance += amount
	logger.Info(fmt.Sprintf("Se añadió %.2f USD al balance. Nuevo balance: %.2f USD", amount, u.Balance))
	return nil
}

// UserRepository define las operaciones básicas para trabajar con usuarios.
// Esto es bastante estándar, pero se pueden agregar más métodos según sea necesario.
type UserRepository interface {
	FindByID(id string) (*User, error)
	FindByUsername(username string) (*User, error)
	Update(user *User) error
}
