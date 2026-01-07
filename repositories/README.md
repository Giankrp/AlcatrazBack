# Repositories Package (`repositories`)

Este paquete implementa el patrón **Repository** para abstraer el acceso a datos. Es la única capa que interactúa directamente con la base de datos (GORM).

## Responsabilidades
- Realizar operaciones CRUD (Create, Read, Update, Delete) sobre la base de datos.
- Desacoplar la lógica de negocio (Services) de la implementación de base de datos.
- Convertir modelos de base de datos a estructuras utilizables por la aplicación.

## Estructura
Cada repositorio tiene:
1. **Interfaz**: Define los métodos disponibles (contrato).
2. **Implementación**: Estructura concreta que usa `*gorm.DB`.
3. **Constructor**: Función `New...` para inyectar la dependencia de base de datos.

## Ejemplo: `UserRepository`

### Interfaz
```go
type UserRepository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
}
```

### Uso en Servicios
El repositorio se inyecta en los servicios, permitiendo que la lógica de negocio no conozca detalles de SQL o GORM.

```go
func NewAuthService(userRepo repositories.UserRepository) AuthService {
    return &authService{userRepo: userRepo}
}
```
