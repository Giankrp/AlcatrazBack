# Routes Package (`routes`)

Este paquete centraliza la definición y configuración de las rutas HTTP de la API.

## Responsabilidades
- Agrupar endpoints por funcionalidad (ej. `/api/auth`, `/api/vault`).
- Asignar cada ruta HTTP (GET, POST, etc.) a su Handler correspondiente.
- Aplicar middlewares específicos a grupos de rutas (ej. autenticación JWT).

## Estructura
Generalmente expone una función `SetupRoutes` que recibe la instancia principal de Echo y los Handlers inicializados.

## Ejemplo
```go
func SetupRoutes(e *echo.Echo, authHandler *handlers.AuthHandler) {
    api := e.Group("/api")
    
    // Auth Routes
    auth := api.Group("/auth")
    auth.POST("/register", authHandler.Register)
    auth.POST("/login", authHandler.Login)
}
```
Esto facilita ver de un vistazo toda la superficie de la API.
