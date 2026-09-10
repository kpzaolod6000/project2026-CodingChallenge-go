# Go API - Matrix Processing Service

Servicio backend desarrollado en **Go** utilizando el framework **Fiber v2**. Este servicio actúa como gateway para recibir matrices de números flotantes, validar su formato rectangular, procesar la descomposición QR y enviar las estadísticas consolidadas al microservicio `node-api` mediante comunicación segura Machine-to-Machine (M2M) respaldada por **JWT**.

---

## 🚀 Características

- **Framework Web**: [Fiber v2](https://gofiber.io/) para alto rendimiento.
- **Autenticación de Clientes**: Middleware de API Key mediante el encabezado `X-API-Key` con comparación en tiempo constante (`crypto/subtle`).
- **Comunicación M2M**: Cliente HTTP interno que genera tokens **JWT** (HS256) dinámicos de corta duración para consumir `node-api`.
- **Seguridad**:
  - Helmet para cabeceras HTTP seguras.
  - CORS configurable.
  - Rate Limiter (limitación de tasa de peticiones).
  - Middleware de recuperación de panics (`recover`).
- **Docker Ready**: `Dockerfile` multi-etapa configurado para desarrollo (`dev` con *live-reloading* vía [Air](https://github.com/air-verse/air)) y producción (`prod` ejecutable minimalista sin privilegios root).

---

## 📁 Estructura del Proyecto

```text
go-api/
├── cmd/
│   └── api/
│       └── main.go           # Punto de entrada de la aplicación
├── internal/
│   ├── client/
│   │   └── node_client.go    # Cliente HTTP M2M con generación de JWT
│   ├── handler/
│   │   └── matrix_handler.go # Handlers y validación de matrices
│   └── middleware/
│       ├── auth.go           # Middleware de autenticación X-API-Key
│       └── security.go       # CORS, Helmet, Rate Limiter y Recover
├── Dockerfile                # Build multi-stage (dev / builder / prod)
├── docker-compose.yml        # Orquestación de servicios
├── go.mod                    # Módulos y dependencias de Go
└── README.md
```

---

## ⚙️ Variables de Entorno

El servicio utiliza las siguientes variables de entorno:

| Variable | Descripción | Valor por Defecto |
| :--- | :--- | :--- |
| `PORT` | Puerto en el que escucha el servidor HTTP | `3000` |
| `API_KEY` | Clave API requerida en la cabecera `X-API-Key` | `coding-challenge` |
| `NODE_API_URL` | URL base del microservicio `node-api` | `http://node-api:4000` |
| `JWT_SECRET` | Clave secreta para firmar tokens JWT M2M | `-` |

---

## 🔌 Endpoints

### 1. Health Check
Comprueba el estado de salud del servicio (no requiere autenticación).

- **Método**: `GET`
- **Ruta**: `/health`
- **Respuesta (`200 OK`)**:
  ```json
  {
    "status": "ok",
    "service": "go-api"
  }
  ```

---

### 2. Procesamiento de Matriz (Factorización QR)
Recibe una matriz bidimensional de números flotantes, valida su estructura y calcula/obtiene las estadísticas consolidadas.

- **Método**: `POST`
- **Ruta**: `/api/factorization`
- **Cabeceras requeridas**:
  - `Content-Type: application/json`
  - `X-API-Key: <TU_API_KEY>` (por defecto `coding-challenge`)

- **Cuerpo de la Petición (Ejemplo)**:
  ```json
  {
    "matrix": [
      [1.0, 2.0, 3.0],
      [4.0, 5.0, 6.0]
    ]
  }
  ```

- **Respuesta (`200 OK`)**:
  ```json
  {
    "qr": {
      "q": [[1.0, 2.0, 3.0], [4.0, 5.0, 6.0]],
      "r": [[1.0, 2.0, 3.0], [4.0, 5.0, 6.0]]
    },
    "stats": {
      "status": "ok",
      "data": {
        "max": 6.0,
        "min": 1.0,
        "sum": 21.0,
        "avg": 3.5
      }
    }
  }
  ```

- **Respuestas de Error**:
  - `401 Unauthorized`: Falta la cabecera `X-API-Key`.
  - `403 Forbidden`: Clave de API inválida.
  - `400 Bad Request`: JSON malformado, matriz vacía o no rectangular.
  - `502 Bad Gateway`: Error de comunicación con `node-api`.

---

## 🛠️ Ejecución y Desarrollo

### Opción 1: Con Docker Compose (Recomendado)

Para levantar el entorno completo con Docker Compose:

```bash
docker compose up --build
```

### Opción 2: Ejecución Local Directa

Asegúrate de contar con **Go 1.22+** instalado.

1. Instalar dependencias:
   ```bash
   go mod tidy
   ```

2. Exportar variables de entorno (opcional):
   ```bash
   export PORT=3000
   export API_KEY=coding-challenge
   export NODE_API_URL=http://localhost:4000
   export JWT_SECRET=super_secret_key
   ```

3. Iniciar la aplicación:
   ```bash
   go run ./cmd/api/main.go
   ```

---

## 🧪 Ejemplo de Petición con `curl`

```bash
curl -X POST http://localhost:3000/api/factorization \
  -H "Content-Type: application/json" \
  -H "X-API-Key: coding-challenge" \
  -d '{
    "matrix": [
      [1.5, 2.5],
      [3.5, 4.5]
    ]
  }'
```

---

## 🧪 Pruebas Unitarias (Tests)

Para ejecutar los tests unitarios del paquete de descomposición matricial (Gram-Schmidt):

### Con Docker (Contenedor en ejecución):

Asegúrate de estar en el directorio raíz del proyecto (`go-api`) y ejecuta:

```bash
docker compose exec go-api go test ./internal/matrix/... -v
```

O para ejecutar todos los tests del proyecto:

```bash
docker compose exec go-api go test ./... -v
```

### Ejecución Local Directa:

```bash
go test ./internal/matrix/... -v
```

