package application

import (
	"cryptoproject/internal/auth/domain"
	"cryptoproject/internal/auth/infrastructure"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthController gestiona las operaciones de autenticación. Aquí están los endpoints como el login.
type AuthController struct {
	jwtService infrastructure.JWTServiceInterface
	userRepo   domain.UserRepository
}

// NewAuthController crea una instancia de AuthController.
// Este es el punto donde inyectamos las dependencias principales.
func NewAuthController(jwtService infrastructure.JWTServiceInterface, userRepo domain.UserRepository) *AuthController {
	return &AuthController{jwtService: jwtService, userRepo: userRepo}
}

// Login permite a los usuarios autenticarse.
/*
Nota: Esta función maneja el proceso de inicio de sesión.
Futuro: Quizás incluir métricas para ver intentos fallidos por usuario o IP.
*/
func (ac *AuthController) Login(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"` // Validar que el username siempre esté presente.
		Password string `json:"password" binding:"required"` // Igual con el password.
	}

	// Aquí validamos los datos enviados en el body.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Intentamos buscar el usuario en la base de datos.
	user, err := ac.userRepo.FindByUsername(request.Username)
	if err != nil || user == nil {
		// Mensaje genérico para evitar revelar si el usuario existe.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña incorrectos"})
		return
	}

	// Validamos la contraseña usando el método en el modelo de usuario.
	if err := user.VerifyPassword(request.Password); err != nil {
		// Este mensaje sigue siendo genérico por seguridad.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña incorrectos"})
		return
	}

	// Generamos el token JWT para el usuario autenticado.
	token, err := ac.jwtService.GenerateToken(user.ID)
	if err != nil {
		/*
			Esto no debería fallar en condiciones normales.
			Futuro: Quizás loggear más detalles del error o agregar una alerta si pasa mucho.
		*/
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}

	// Respondemos con el token al cliente.
	c.JSON(http.StatusOK, gin.H{"token": token})
}
