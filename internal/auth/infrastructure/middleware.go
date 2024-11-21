package infrastructure

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware proporciona la funcionalidad del middleware para autenticación JWT.
/*
Aquí controlamos que todas las solicitudes autenticadas pasen primero por este filtro.
Es como un portero, nadie entra sin un token válido.
*/
type JWTMiddleware struct {
	jwtService JWTServiceInterface
}

// NewJWTMiddleware crea una nueva instancia de JWTMiddleware.
/*
Futuro: Si queremos cambiar la implementación de JWT, este constructor lo hace más fácil.
*/
func NewJWTMiddleware(jwtService JWTServiceInterface) *JWTMiddleware {
	return &JWTMiddleware{jwtService: jwtService}
}

// Middleware intercepta las solicitudes para validar el token JWT.
/*
Esto revisa que el token JWT sea válido antes de que la solicitud llegue al controlador.
*/
func (m *JWTMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraer el token del encabezado Authorization.
		/*
			Ojo con esto: Si el cliente no envía el encabezado Authorization, rechazamos de una vez.
		*/
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "El encabezado de autorización es obligatorio"})
			c.Abort()
			return
		}

		// Verificar formato del token.
		/*
			Aquí partimos el encabezado para asegurarnos de que siga el formato: "Bearer {token}".
			Si no tiene ese formato, devolvemos un error claro.
		*/
		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "El formato del encabezado debe ser Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar el token usando el servicio JWT.
		/*
			Esto es como la prueba final. Si el token está vencido o es inválido, no dejamos pasar.
		*/
		token, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Extraer el userID del token.
		/*
			Nota: Aquí es donde sacamos el userID para usarlo después. Si algo sale mal, tampoco seguimos.
		*/
		userID, err := m.jwtService.ExtractUserID(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims del token inválidos"})
			c.Abort()
			return
		}

		// Añadir el userID al contexto.
		/*
			Aquí guardamos el userID en el contexto de Gin. Esto es útil porque lo podemos usar
			en los controladores sin tener que estar  pasando el token otra vez.
		*/
		c.Set("user_id", userID)

		// Continuar con la solicitud.
		/*
			Si todo está bien, dejamos que la solicitud siga su camino.
		*/
		c.Next()
	}
}
