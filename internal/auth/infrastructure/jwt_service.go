package infrastructure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTServiceInterface define las operaciones del servicio JWT.
/*
Aquí definimos lo básico que cualquier servicio JWT debería ofrecer.
Esto nos deja flexibles para cambiar de implementación si hace falta.
*/
type JWTServiceInterface interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	ExtractUserID(token *jwt.Token) (string, error)
}

// JWTService implementa JWTServiceInterface.
/*
Ojo: Esto está bien para ahora, pero si cambiamos la librería de JWT,
tendríamos que ajustar esta implementación.
*/
type JWTService struct {
	secretKey string
	ttl       time.Duration
}

// NewJWTService crea una nueva instancia de JWTService.
/*
Aquí estamos configurando el servicio con la clave secreta y el tiempo de expiración.
Futuro: Tal vez hacer esto más dinámico desde una configuración central.
*/
func NewJWTService(secretKey string, ttl time.Duration) *JWTService {
	return &JWTService{secretKey: secretKey, ttl: ttl}
}

// GenerateToken genera un token JWT para un usuario específico.
/*
Esta función crea un token nuevo. Fácil de entender, pero ojo:
- Si la clave secreta es comprometida, todos los tokens serán vulnerables.
*/
func (s *JWTService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,                       // Aquí guardamos el ID del usuario.
		"exp":     time.Now().Add(s.ttl).Unix(), // Fecha de expiración.
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken valida un token JWT.
/*
Esta función revisa si el token sigue siendo válido.
Cuidado: Si el método de firma no es el esperado, marcamos el token como inválido.
*/
func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validar el método de firma.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido") // Mensaje claro si algo falla.
		}
		return []byte(s.secretKey), nil
	})
}

// ExtractUserID extrae el userID de un token JWT válido.
/*
Aquí sacamos el user_id de los claims. Si algo falla, devolvemos un error.
Esto es útil para usar después de validar el token.
*/
func (s *JWTService) ExtractUserID(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("claims del token inválidos") // Nota: Esto puede ocurrir si el token está corrupto.
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id no encontrado en los claims del token")
	}

	return userID, nil
}
