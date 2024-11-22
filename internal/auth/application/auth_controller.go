package application

import (
	"cryptoproject/internal/auth/domain"
	"cryptoproject/internal/auth/infrastructure"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginRequest representa las credenciales requeridas para iniciar sesión.
// @Description Estructura del cuerpo de la solicitud para el endpoint de login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // El nombre de usuario debe ser obligatorio.
	Password string `json:"password" binding:"required"` // La contraseña debe ser obligatoria.
}

// AuthController gestiona las operaciones de autenticación.
type AuthController struct {
	jwtService infrastructure.JWTServiceInterface
	userRepo   domain.UserRepository
}

// NewAuthController crea una instancia de AuthController.
func NewAuthController(jwtService infrastructure.JWTServiceInterface, userRepo domain.UserRepository) *AuthController {
	return &AuthController{jwtService: jwtService, userRepo: userRepo}
}

// Login godoc
// @Summary Autenticación de usuarios
// @Description Permite a los usuarios autenticarse y obtener un token JWT para acceder a los endpoints protegidos.
// @Tags Auth
// @Accept json
// @Produce json
// @Param LoginRequest body LoginRequest true "Credenciales del usuario para autenticación"
// @Success 200 {object} map[string]string "Token JWT generado para el usuario autenticado"
// @Failure 400 {object} map[string]string "Error en la validación de datos enviados"
// @Failure 401 {object} map[string]string "Credenciales inválidas o usuario no encontrado"
// @Failure 500 {object} map[string]string "Error interno al generar el token"
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var request LoginRequest

	// Validar los datos enviados en el cuerpo de la solicitud.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Intentar buscar el usuario en la base de datos.
	user, err := ac.userRepo.FindByUsername(request.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña incorrectos"})
		return
	}

	// Verificar la contraseña.
	if err := user.VerifyPassword(request.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña incorrectos"})
		return
	}

	// Generar el token JWT.
	token, err := ac.jwtService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
