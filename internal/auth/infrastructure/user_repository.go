package infrastructure

import (
	"cryptoproject/internal/auth/domain"
	"errors"

	"gorm.io/gorm"
)

// GormUserRepository implementa UserRepository utilizando GORM.
/*
Bueno, esto es el repositorio del usuario, básicamente maneja todas las consultas a la base de datos para usuarios.
Es como el punto de entrada para trabajar con GORM.
*/
type GormUserRepository struct {
	DB *gorm.DB
}

// NewUserRepository crea una nueva instancia de GormUserRepository.
/*
Esto es un constructor sencillo, crea el repositorio de usuario basado en la conexión de la base de datos.
*/
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &GormUserRepository{DB: db}
}

// FindByID busca un usuario por su ID.
/*
Aja, aquí vamos a buscar el usuario con su ID. Si no lo encuentra, devolvemos un error claro.
*/
func (r *GormUserRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No encontramos el usuario, no es grave pero devolvemos este error.
			return nil, errors.New("usuario no encontrado")
		}
		// Aquí algo más grave pasó, devolvemos el error original.
		return nil, err
	}

	// Inicializar CryptoHoldings si está vacío.
	/*
		Por si acaso, si no hay nada en los holdings de criptos, lo inicializamos para evitar problemas después.
	*/
	if user.CryptoHoldings == nil {
		user.CryptoHoldings = make(map[string]float64)
	}

	return &user, nil
}

// FindByUsername busca un usuario por su nombre de usuario.
/*
Igual que la función anterior, pero esta vez buscamos por el username.
*/
func (r *GormUserRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Lo mismo, si no lo encuentra devolvemos un error bien bonito.
			return nil, errors.New("usuario no encontrado")
		}
		// Algo pasó, devolvemos el error para que lo maneje quien llame esta función.
		return nil, err
	}

	// Inicializar CryptoHoldings si está vacío.
	/*
		Otra vez lo mismo, asegurarnos de que no esté vacío. Ojo, esto podría ser optimizado más adelante.
	*/
	if user.CryptoHoldings == nil {
		user.CryptoHoldings = make(map[string]float64)
	}

	return &user, nil
}

// Update actualiza un usuario en la base de datos.
/*
Aquí actualizamos la información del usuario. Si algo falla, devolvemos el error.
*/
func (r *GormUserRepository) Update(user *domain.User) error {
	if err := r.DB.Save(user).Error; err != nil {
		// Ojo, esto puede fallar si hay problemas con la base de datos o el modelo.
		return err
	}
	return nil
}
