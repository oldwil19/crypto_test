package application

import (
	"cryptoproject/internal/auth/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccountController maneja las operaciones relacionadas con el saldo del usuario.
type AccountController struct {
	userRepo domain.UserRepository
}

// NewAccountController crea una nueva instancia de AccountController.
func NewAccountController(userRepo domain.UserRepository) *AccountController {
	return &AccountController{userRepo: userRepo}
}

/*
HandleAddBalance es el endpoint que se encarga de añadir saldo al usuario.
Aquí validamos que la solicitud sea válida, recuperamos al usuario, ajustamos el saldo y
actualizamos la información en la base de datos. Esto podría mejorar agregando auditorías en el futuro.
*/
func (ac *AccountController) HandleAddBalance(c *gin.Context) {
	// Extraer el userID del JWT. Esto es fundamental para saber a quién estamos modificando.
	userID := c.GetString("user_id")

	/*
	 Ojo con esto:
	 La validación es básica. Si en el futuro se requiere un sistema más robusto,
	 quizás debamos usar una librería especializada o reglas más avanzadas.
	*/
	var request AddBalanceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		// Si el JSON está mal estructurado o falta algo, devolvemos un error claro.
		c.JSON(http.StatusBadRequest, gin.H{"error": "La cantidad ingresada es inválida"})
		return
	}

	// Aquí buscamos al usuario en la base de datos. Si no existe, no podemos continuar.
	user, err := ac.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Añadir el saldo. Pendiente: ¿Y si en el futuro necesitamos límites máximos o mínimos?
	err = user.AddBalance(request.Amount)
	if err != nil {
		/*
		 Puede ser interesante registrar más detalles aquí para auditoría.
		 Por ejemplo, quién intentó añadir saldo, desde qué IP, etc.
		*/
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Guardar los cambios en la base de datos.
	if err := ac.userRepo.Update(user); err != nil {
		// Un error aquí es crítico. Tal vez deberíamos enviar una alerta en un sistema real.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el usuario"})
		return
	}

	/*
	 Respuesta final: le damos al cliente la confirmación y el balance actualizado.
	 Nota: Esto podría incluir más datos, como un historial reciente de operaciones.
	*/
	c.JSON(http.StatusOK, gin.H{
		"message": "Saldo añadido con éxito",
		"user": gin.H{
			"balance": user.Balance,
		},
	})
}
