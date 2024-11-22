## Descripción General:

Crypto Proyecto es una plataforma B2B diseñada para permitir la integración de terceros con el sistema de trading simulado. Permite acceder a datos de mercado en tiempo real, historial de precios de criptomonedas y realizar operaciones simuladas de trading. El proyecto utiliza autenticación segura basada en JWT y está respaldado por PostgreSQL como base de datos principal. Todo el proyecto está dockerizado para facilitar la implementación.

```

.
├── Dockerfile
├── README.md
├── cmd
│   └── server
│       └── main.go
├── configs
│   ├── nginx
│   │   └── nginx.conf
│   └── prometheus
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── account
│   │   ├── application
│   │   │   ├── account_controller.go
│   │   │   └── add_balance_request.go
│   │   ├── domain
│   │   └── infrastructure
│   │       └── db.go
│   ├── auth
│   │   ├── application
│   │   │   ├── auth_controller.go
│   │   │   └── register_controller.go
│   │   ├── domain
│   │   │   └── user.go
│   │   └── infrastructure
│   │       ├── jwt_service.go
│   │       ├── middleware.go
│   │       └── user_repository.go
│   ├── market
│   │   ├── application
│   │   │   └── market_controller.go
│   │   ├── domain
│   │   └── infrastructure
│   │       └── coingecko_service.go
│   ├── server
│   │   └── router.go
│   └── trading
│       ├── application
│       │   └── trading_controller.go
│       ├── domain
│       │   └── transaction.go
│       └── infrastructure
│           └── transaction_repository.go
├── migrations
│   └── init.sql
├── pkg
│   ├── config
│   │   └── config.go
│   ├── logger
│   │   └── logger.go
│   ├── middleware
│   └── utils
└── test
```

## Requerimientos Técnicos:

### Lenguaje y Frameworks:

Lenguaje: Go (versión 1.20+)
Framework Web: Gin Gonic
Base de Datos:

PostgreSQL para almacenamiento de usuarios, transacciones y balances.
APIs Externas:

CoinGecko para obtener precios y datos históricos de criptomonedas.
Autenticación:

JWT (JSON Web Tokens) para proteger los endpoints y validar a los usuarios.
Contenedores:

Docker para contenerización.
Docker Compose para configurar y levantar múltiples servicios.

### **Dependencias Principales**

Estas son las librerías fundamentales que permiten el funcionamiento del proyecto:

1. **`github.com/gin-gonic/gin v1.10.0`**
   * Framework web rápido y minimalista que utilizamos para manejar rutas y middlewares. Es la base para todos nuestros endpoints.
2. **`github.com/golang-jwt/jwt/v5 v5.2.1`**
   * Librería utilizada para la autenticación segura mediante tokens JWT. Nos permite gestionar sesiones y proteger las rutas.
3. **`github.com/google/uuid v1.3.0`**
   * Generación de identificadores únicos universales (UUID) para identificar entidades como usuarios o transacciones.
4. **`github.com/joho/godotenv v1.5.1`**
   * Permite cargar variables de entorno desde archivos `.env`, simplificando la configuración del proyecto.
5. **`github.com/prometheus/client_golang v1.20.5`**
   * Librería para exponer métricas del sistema y facilitar el monitoreo con Prometheus.
6. **`gorm.io/driver/postgres v1.5.9`** y **`gorm.io/gorm v1.25.12`**
   * GORM es nuestro ORM principal para interactuar con la base de datos PostgreSQL. Simplifica las operaciones de persistencia y migraciones.
7. **`github.com/swaggo/gin-swagger v1.6.0`** y **`github.com/swaggo/files v1.0.1`**
   * Herramientas para generar y exponer la documentación de la API basada en Swagger de forma interactiva.

---

### **Dependencias Adicionales**

Estas librerías complementan funcionalidades específicas o indirectamente contribuyen al rendimiento y usabilidad:

* **`golang.org/x/crypto`**: Proporciona implementaciones de algoritmos criptográficos.
* **`github.com/json-iterator/go`**: Alternativa más rápida para trabajar con JSON.
* **`github.com/jinzhu/inflection`** y **`github.com/jinzhu/now`**: Manejo intuitivo de tiempos y pluralización en nombres de entidades.
* **`github.com/prometheus/procfs`**: Monitoreo de procesos en sistemas operativos.

---

### **Versiones de Go y PostgreSQL**

* **Go**: El proyecto está desarrollado en Go `1.23.3`, aprovechando las últimas optimizaciones y características del lenguaje.
* **PostgreSQL**: Base de datos robusta y escalable para gestionar los datos de usuarios, transacciones y configuraciones.

---

### **Razonamiento Técnico**

1. **Eficiencia y Rendimiento**:
   * Gin y GORM están optimizados para alto rendimiento, esenciales para manejar múltiples solicitudes y operaciones concurrentes.
2. **Facilidad de Extensión**:
   * Las dependencias como `swagger` y `jwt` permiten integrar rápidamente autenticación y documentación sin reinventar la rueda.
3. **Manejo de Errores y Logs**:
   * Uso extensivo de librerías como `logger` para registrar eventos importantes y manejar errores con claridad.
4. **Modularidad**:
   * El uso de GORM y PostgreSQL asegura que podamos extender fácilmente la base de datos y sus interacciones.

Este stack no solo está diseñado para cumplir los requerimientos actuales, sino también para permitir iteraciones futuras de manera ágil y eficiente.

## 

Instalación:
Requisitos Previos:

Docker y Docker Compose instalados en tu máquina.
Configurar variables de entorno en un archivo .env:

```
POSTGRES_USER: Usuario de la base de datos.
POSTGRES_PASSWORD: Contraseña de la base de datos.
POSTGRES_DB: Nombre de la base de datos.
POSTGRES_HOST: Dirección del servidor de la base de datos.
JWT_SECRET: Clave secreta para la generación de tokens JWT.
```

Pasos:

Clona el repositorio:

`git clone https://github.com/tu_usuario/crypto-proyecto.git.`
Entra al directorio: cd crypto-proyecto.
Ejecuta:

`docker-compose up --build.`

## 

Endpoints Disponibles:

### **Registro de Usuario**

**Descripción:**
Permite registrar un nuevo usuario en el sistema proporcionando un nombre de usuario y una contraseña.

**Ruta:**
`POST /register`

**Headers:**

* `Content-Type`: `application/json`

**Parámetros en el cuerpo de la solicitud (JSON):**

* `username`: El nombre de usuario único para el nuevo usuario (obligatorio).
* `password`: La contraseña segura del usuario (obligatorio).

**Respuestas:**

* **201 (Creado):** Usuario registrado correctamente.
* **409 (Conflicto):** El usuario ya existe.
* **400 (Error de validación):** Los datos proporcionados son inválidos.

**Notas:**
No requiere autenticación.

request

```

curl -X POST http://localhost:8080/register \
-H "Content-Type: application/json" \
-d '{
  "username": "usuario123",
  "password": "contraseñaSegura123"
}'

```

response

```
{
  "message": "Usuario registrado exitosamente"
}

```

---

### **Inicio de Sesión**

**Descripción:**
Autentica al usuario y devuelve un token JWT.

**Ruta:**
`POST /auth/login`

**Headers:**

* `Content-Type`: `application/json`

**Parámetros en el cuerpo de la solicitud (JSON):**

* `username`: El nombre de usuario registrado (obligatorio).
* `password`: La contraseña correspondiente (obligatorio).

**Respuestas:**

* **200 (Éxito):** Devuelve un token JWT válido.
* **401 (No autorizado):** Usuario o contraseña incorrectos.
* **400 (Error de validación):** Los datos proporcionados son inválidos.

Request

```

curl -X POST http://localhost:8080/auth/login \
-H "Content-Type: application/json" \
-d '{
  "username": "usuario123",
  "password": "contraseñaSegura123"
}'

```

Response

```
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

```

---

### **Obtener Precio Actual de Criptomonedas**

**Descripción:**
Devuelve el precio actual de una criptomoneda específica en una moneda determinada.

**Ruta:**
`GET /market/:id/price`

**Headers:**

* `Authorization`: `Bearer <token>`

**Parámetros de URL:**

* `id`: Identificador de la criptomoneda (por ejemplo, `bitcoin` o `solana`).

**Parámetros opcionales (Query):**

* `currency`: La moneda en la que se desea el precio (por defecto, `usd`).

**Respuestas:**

* **200 (Éxito):** Devuelve el precio actual.
* **404 (No encontrado):** La criptomoneda no fue encontrada.
* **401 (No autorizado):** El token JWT es inválido o falta.

Request

```
curl -X GET "http://localhost:8080/market/bitcoin/price?currency=usd" \
-H "Authorization: Bearer <token>"

```

Response

```
{
  "crypto": "bitcoin",
  "currency": "usd",
  "price": 95000.12
}

```

---

### **Obtener Precios Históricos de Criptomonedas**

**Descripción:**
Devuelve los precios históricos de una criptomoneda en un rango de fechas.

**Ruta:**
`GET /market/:id/history`

**Headers:**

* `Authorization`: `Bearer <token>`

**Parámetros de URL:**

* `id`: Identificador de la criptomoneda (por ejemplo, `bitcoin` o `solana`).

**Parámetros opcionales (Query):**

* `start`: Fecha de inicio en formato `dd-mm-yyyy`.
* `end`: Fecha de fin en formato `dd-mm-yyyy`.

**Respuestas:**

* **200 (Éxito):** Devuelve los precios históricos.
* **400 (Error de validación):** Las fechas no tienen el formato adecuado.
* **401 (No autorizado):** El token JWT es inválido o falta.

Request

```
curl -X GET "http://localhost:8080/market/bitcoin/history?start=01-11-2024&end=20-11-2024" \
-H "Authorization: Bearer <token>"

```

Response

```
[
  {
    "timestamp": 1698883200000,
    "price": 36521.11
  },
  {
    "timestamp": 1698969600000,
    "price": 36984.12
  }
]

```

---

### **Compra Simulada de Criptomonedas**

**Descripción:**
Permite a un usuario simular la compra de una criptomoneda usando su saldo virtual., saldo inicial es de 1000 usd por temas de velocidad, para desarrollar mas rapido esto esta en la linea

`/Users/wilmer.flores/Documents/wenia/cryptoproject/internal/auth/domain/user.go`

**Ruta:**
`POST /trading/buy`

**Headers:**

* `Authorization`: `Bearer <token>`
* `Content-Type`: `application/x-www-form-urlencoded`

**Parámetros en el cuerpo de la solicitud (Form):**

* `coin`: Identificador de la criptomoneda (por ejemplo, `bitcoin` o `solana`).
* `amount`: Cantidad de criptomoneda a comprar.

**Respuestas:**

* **200 (Éxito):** La compra se realizó con éxito.
* **400 (Error de validación):** La cantidad ingresada no es válida.
* **401 (No autorizado):** El token JWT es inválido o falta.
* **404 (No encontrado):** Usuario no encontrado o criptomoneda no soportada.

Request

```

curl -X POST http://localhost:8080/trading/buy \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/x-www-form-urlencoded" \
-d "coin=bitcoin" \
-d "amount=0.01"

```

Response

```
{
  "message": "Compra realizada con éxito",
  "user": {
    "balance": 999.63,
    "crypto_balance": {
      "bitcoin": 0.01
    }
  },
  "transaction": {
    "ID": "12345678-abcd-1234-efgh-567890abcdef",
    "UserID": "02933989-cc89-477b-8711-a0eac1971ecc",
    "Coin": "bitcoin",
    "Amount": 0.01,
    "Price": 36984.12,
    "Timestamp": "2024-11-21T14:00:00Z"
  }
}

```

---

### **Historial de Transacciones**

**Descripción:**
Devuelve todas las transacciones realizadas por un usuario.

**Ruta:**
`GET /trading/history`

**Headers:**

* `Authorization`: `Bearer <token>`

**Respuestas:**

* **200 (Éxito):** Devuelve el historial de transacciones.
* **401 (No autorizado):** El token JWT es inválido o falta.

Request

```
curl -X GET http://localhost:8080/trading/history \
-H "Authorization: Bearer <token>"

```

Response

```
[
  {
    "ID": "12345678-abcd-1234-efgh-567890abcdef",
    "UserID": "02933989-cc89-477b-8711-a0eac1971ecc",
    "Coin": "bitcoin",
    "Amount": 0.01,
    "Price": 36984.12,
    "Timestamp": "2024-11-21T14:00:00Z"
  },
  {
    "ID": "87654321-abcd-4321-efgh-567890abcdef",
    "UserID": "02933989-cc89-477b-8711-a0eac1971ecc",
    "Coin": "doge",
    "Amount": 100,
    "Price": 0.07,
    "Timestamp": "2024-11-20T12:00:00Z"
  }
]

```

---

### **Obtener Balance Actual**

**Descripción:**
Proporciona el saldo actual en USD y las criptomonedas que posee un usuario.

**Ruta:**
`GET /trading/balance`

**Headers:**

* `Authorization`: `Bearer <token>`

**Respuestas:**

* **200 (Éxito):** Devuelve el balance actual.
* **401 (No autorizado):** El token JWT es inválido o falta.

Request

```
curl -X GET http://localhost:8080/trading/balance \
-H "Authorization: Bearer <token>"


```

Response

```
{
  "usd_balance": 999.63,
  "crypto_holdings": {
    "bitcoin": 0.01,
    "doge": 100
  }
}

```

---

### **Añadir Saldo al Usuario**

**Descripción:**
Permite a un usuario añadir saldo virtual a su cuenta.

**Ruta:**
`POST /account/balance/add`

**Headers:**

* `Authorization`: `Bearer <token>`
* `Content-Type`: `application/json`

**Parámetros en el cuerpo de la solicitud (JSON):**

* `amount`: Monto de dinero a añadir al saldo del usuario.

**Respuestas:**

* **200 (Éxito):** El saldo fue añadido exitosamente.
* **400 (Error de validación):** El monto ingresado no es válido.
* **401 (No autorizado):** El token JWT es inválido o falta.

Request

```

curl -X POST http://localhost:8080/account/balance/add \
-H "Authorization: Bearer <token>" \
-H "Content-Type: application/json" \
-d '{
  "amount": 100.00
}'

```

Response

```
{
  "message": "Saldo añadido con éxito",
  "user": {
    "balance": 1099.63
  }
}

```

#### Consideraciones Finales:

Este proyecto fue desarrollado con los principios SOLID, Clean Code y una arquitectura basada en dominios (DDD). Se utilizaron contenedores Docker para simplificar la implementación y CoinGecko para obtener datos de mercado.


## HelthCheck

accede aca, es basico pero extensible.

```
http://localhost/healthz
```

## Documentaciones con swagger

para acceder a la docuementacion
aca

```
http://localhost/swagger/index.html
```

para generar documentacion

```
swag init --generalInfo cmd/server/main.go --output ./docs
```

## **Requerimientos y Cumplimiento**

### **Objetivo del Proyecto**

**Descripción**: Sistema B2B que expone un portal de APIs para acceder a datos de mercado y realizar operaciones simuladas.

* ✅ **Cumplido**: Se desarrolló un sistema con APIs que permite:
  * Acceder a datos en tiempo real de criptomonedas específicas.
  * Realizar operaciones simuladas de compra de criptomonedas.
  * Proteger el acceso a los servicios con autenticación segura.

---

### **1. Autenticación y Autorización**

**Sistema de Autenticación:**

* ✅ **Cumplido**: Implementación de autenticación mediante **JWT**.
  * Middleware para validar tokens.
  * Endpoints protegidos que requieren autenticación.
  * Generación de tokens seguros.

---

### **2. Integración con Coingecko**

**Consumo de APIs Externas:**

* ✅ **Cumplido**:
  * Obtener precios actualizados de **Bitcoin (BTC)** y **Solana (SOL)**.
  * Recuperación de precios históricos para un rango de fechas específico.
  * Manejador robusto con políticas de reintento y control de tasa de llamadas (rate-limiting).

---

### **3. Endpoints Funcionales**

#### **Datos de Mercado**

1. **Obtener Precios en Tiempo Real**:
   * ✅ **Cumplido**: Endpoint `/market/:id/price` devuelve el precio actual en USD de las criptomonedas.
2. **Obtener Datos Históricos**:
   * ✅ **Cumplido**: Endpoint `/market/:id/history` permite obtener precios históricos en un rango de fechas.

#### **Operaciones de Trading Simulado**

1. **Compra Simulada**:
   * ✅ **Cumplido**: Endpoint `/trading/buy` para realizar compras simuladas de criptomonedas.
2. **Validaciones Necesarias**:
   * ✅ **Cumplido**:
     * Verificar saldo suficiente antes de la compra.
     * Registrar cada transacción en el historial.
3. **Registrar Transacciones**:
   * ✅ **Cumplido**: Cada compra se almacena en la base de datos con detalles.

#### **Historial y Estado de Cuenta**

1. **Consultar Historial de Transacciones**:
   * ✅ **Cumplido**: Endpoint `/trading/history` lista las transacciones realizadas por un usuario.
2. **Obtener Balance Actual**:
   * ✅ **Cumplido**: Endpoint `/trading/balance` devuelve el saldo en USD y el balance de criptomonedas.

#### **Saldo del Usuario**

* ✅ **Cumplido**: Endpoint `/account/balance/add` para agregar saldo virtual en USD.

---

### **4. Portal de Documentación del API**

**Documentación Interactiva:**

* ✅ **Cumplido**:
  * Swagger UI disponible para explorar los endpoints.
  * Ejemplos de solicitudes y respuestas incluidos.

---

### **Requerimientos Técnicos**

1. **Lenguaje de Programación**:
   * ✅ **Cumplido**: Proyecto desarrollado en **Go 1.23.3**.
2. **Base de Datos**:
   * ✅ **Cumplido**: Implementación de **PostgreSQL** para almacenar información de usuarios, balances y transacciones.
3. **Control de Versiones**:
   * ✅ **Cumplido**: Proyecto alojado en un repositorio de Git con historial de commits.
4. **Pruebas Unitarias**:
   * ⚠️ **Parcialmente Cumplido**:
     * Pruebas unitarias para validaciones clave implementadas, pero cobertura podría extenderse.
5. **Buenas Prácticas de Desarrollo**:
   * ✅ **Cumplido**:
     * Principios **SOLID** aplicados.
     * Diseño modular y desacoplado basado en **DDD**.
     * Manejo adecuado de errores y excepciones.

---

### **Extras Valorados**

1. **Arquitectura Orientada a Eventos**:
   * ✅ **Cumplido**: Sistema de notificaciones implementado para detectar cambios significativos en precios. Utiliza **canales Go (chan)** con potencial de integración futura con colas como **SQS** o **RabbitMQ**.
2. **Dockerización**:
   * ✅ **Cumplido**: Proyecto dockerizado con **Docker Compose** para facilitar la configuración y ejecución.
