package application

import (
	"cryptoproject/internal/auth/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterController maneja el registro de usuarios. Aquí metemos toda la lógica para agregar nuevos usuarios.
type RegisterController struct {
	userRepo domain.UserRepository
}

// NewRegisterController crea una instancia de RegisterController. Esto nos permite inyectar dependencias.
func NewRegisterController(userRepo domain.UserRepository) *RegisterController {
	return &RegisterController{userRepo: userRepo}
}

// Register registra un nuevo usuario.
/*
Ojo: Este endpoint debería manejar también validaciones más robustas a futuro,
como verificar fuerza de contraseñas o agregar lógica de verificación por email.
*/
func (c *RegisterController) Register(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"` // Validamos que el nombre de usuario no esté vacío.
		Password string `json:"password" binding:"required"` // Lo mismo con la contraseña.
	}

	// Aquí validamos que la estructura del JSON sea correcta.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Creamos un nuevo usuario usando el modelo del dominio.
	user, err := domain.NewUser(req.Username, req.Password)
	if err != nil {
		// Algo raro pasó al crear el usuario (puede ser un problema de validación interna).
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Revisamos si ya existe un usuario con ese nombre.
	existingUser, err := c.userRepo.FindByUsername(req.Username)
	if err == nil && existingUser != nil {
		/*
			Pendiente: Mejorar el manejo de errores aquí.
			Por ejemplo, enviar un mensaje más descriptivo o registrar por qué falló la creación.
		*/
		ctx.JSON(http.StatusConflict, gin.H{"error": "El usuario ya existe"})
		return
	}

	// Guardamos el nuevo usuario en la base de datos.
	if err := c.userRepo.Update(user); err != nil {
		/*
			Esto podría ser un problema serio si pasa mucho.
			Futuro: Tal vez loggear más información sobre por qué no se pudo crear el usuario.
		*/
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el usuario"})
		return
	}

	// Si todo salió bien, respondemos con un mensaje de éxito.
	ctx.JSON(http.StatusCreated, gin.H{"message": "Usuario registrado exitosamente"})
}
