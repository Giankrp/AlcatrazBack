# Validator Package (`validator`)

Este paquete configura y expone la instancia global del validador de datos.

## Responsabilidades
- Proveer una instancia única (Singleton) de `go-playground/validator`.
- Centralizar la lógica de validación de estructuras.
- (Opcional) Proveer utilidades para formatear errores de validación en mensajes legibles.

## Uso
Se utiliza principalmente en la capa de **Handlers** para validar DTOs antes de procesarlos.

```go
if err := validator.Validate.Struct(&dto); err != nil {
    return c.JSON(http.StatusBadRequest, ...)
}
```

Esto asegura que solo datos válidos lleguen a la capa de Servicio.
