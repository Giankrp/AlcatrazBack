# Services Package (`services`)

Este paquete contiene la **Lógica de Negocio** de la aplicación. Es el núcleo funcional del sistema.

## Responsabilidades
- Implementar las reglas de negocio (ej. validaciones complejas, flujos de registro).
- Orquestar llamadas a repositorios y otros componentes (como seguridad).
- Transformar DTOs (Data Transfer Objects) en Modelos de dominio y viceversa.
- Generar tokens JWT y manejar autenticación.

## Estructura
Al igual que los repositorios, los servicios se definen por interfaces para facilitar el testing y desacoplamiento.

## Ejemplo: `AuthService`

### Funcionalidades
- **Register**:
  - Verifica si el email ya existe usando `UserRepository`.
  - Hashea la contraseña usando el paquete `security`.
  - Crea el usuario en la base de datos.
- **Login**:
  - Busca el usuario por email.
  - Verifica la contraseña hasheada.
  - Genera y firma un token JWT.

### Inyección de Dependencias
Los servicios reciben los repositorios necesarios en su constructor:

```go
type authService struct {
    userRepo repositories.UserRepository
}
```
