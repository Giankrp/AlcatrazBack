# 📦 Middleware Package (`middleware`)

> 📌 Este directorio está reservado para **middlewares personalizados** futuros.

---

## Middlewares Activos (configurados en otras capas)

Aunque este directorio está vacío, la aplicación utiliza los siguientes middlewares configurados en `main.go` y `routes/`:

### Configurados en `main.go`

| Middleware | Paquete | Descripción |
|---|---|---|
| `middleware.Logger()` | Echo built-in | Registra todas las peticiones HTTP en la consola |
| `middleware.Recover()` | Echo built-in | Recupera panics y devuelve 500 en lugar de crashear |
| `middleware.CORSWithConfig(...)` | Echo built-in | Configura orígenes permitidos y credenciales |

#### Configuración CORS

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins:     allowedOrigins,   // Desde ALLOWED_ORIGINS env
    AllowCredentials: true,             // Necesario para cookies HttpOnly
}))
```

### Configurado en `routes/routes.go`

| Middleware | Paquete | Descripción |
|---|---|---|
| `echojwt.WithConfig(...)` | `echo-jwt/v4` | Valida JWT desde cookie `auth_token` |

---

## Custom Error Handler

Además de los middlewares, `main.go` define un **error handler personalizado** que estandariza todas las respuestas de error:

```go
// Respuestas de error siempre en formato JSON:
{
    "error": "mensaje descriptivo del error"
}
```

| Método HTTP | Comportamiento |
|---|---|
| `HEAD` | Devuelve solo el código de estado (sin body) |
| Otros | Devuelve `{ "error": "..." }` con código apropiado |

---

## Extensión Futura

Cuando necesites crear middlewares personalizados (rate limiting, logging avanzado, API keys, etc.), crea archivos `.go` en este directorio siguiendo la convención de Echo:

```go
package middleware

import "github.com/labstack/echo/v4"

func CustomMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Lógica del middleware
            return next(c)
        }
    }
}
```
