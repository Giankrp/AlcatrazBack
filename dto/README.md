# DTO Package (`dto`)

DTO significa **Data Transfer Object**. Este paquete define las estructuras utilizadas para transportar datos entre el cliente (frontend) y la aplicación, o entre capas de la aplicación.

## Responsabilidades
- Definir la estructura exacta de los JSON de entrada y salida.
- Desacoplar la estructura interna de la base de datos (Modelos) de la estructura pública de la API.
- Incluir etiquetas de validación (`validate:"..."`) para asegurar la integridad de los datos de entrada.

## Diferencia con Models
- **Models**: Representan tablas de base de datos. Tienen campos internos como `PasswordHash` o `ID`.
- **DTOs**: Representan intenciones de usuario. Tienen campos como `Password` (texto plano para login) que no existen en el modelo persistido.

## Ejemplo: `RegisterDTO`
```go
type RegisterDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```
Define que para registrarse se requiere un email válido y una contraseña de mínimo 8 caracteres.
