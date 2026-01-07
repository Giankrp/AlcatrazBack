# Handlers Package (`handlers`)

Este paquete actúa como la capa de entrada HTTP (Controladores). Gestiona la interacción con el cliente (request/response).

## Responsabilidades
- Recibir peticiones HTTP (vía Echo framework).
- Validar el formato de entrada (Binding de JSON a DTOs).
- Validar reglas de formato (usando paquete `validator`).
- Llamar a los Servicios correspondientes para ejecutar la lógica.
- Formatear la respuesta HTTP (JSON) y códigos de estado (200, 400, 500, etc.).

## Flujo Típico

1. **Bind**: `c.Bind(&dto)` lee el JSON del body.
2. **Validate**: `validator.Validate.Struct(&dto)` verifica reglas (required, email, min len).
3. **Service Call**: `h.service.Metodo(dto)` ejecuta la acción.
4. **Response**: `c.JSON(status, data)` devuelve el resultado.

## Ejemplo: `AuthHandler`
Maneja rutas como `/api/auth/register` y `/api/auth/login`.

```go
func (h *AuthHandler) Login(c echo.Context) error {
    // 1. Bind
    // 2. Validate
    // 3. Call Service
    token, err := h.authService.Login(loginDTO)
    // 4. Response
    return c.JSON(http.StatusOK, map[string]string{"token": token})
}
```
