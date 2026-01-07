# Models Package (`models`)

Este paquete define las **Entidades de Dominio** y su representación en la base de datos.

## Responsabilidades
- Definir la estructura de datos core del negocio (User, VaultItem, etc.).
- Definir etiquetas (tags) de GORM para el mapeo Objeto-Relacional (ORM).
- Definir tipos de datos y restricciones (primary keys, foreign keys, constraints).

## Características
- Las estructuras suelen mapearse 1:1 con tablas en la base de datos.
- Utilizan tags como `gorm:"primaryKey"`, `gorm:"unique"`, `gorm:"index"` para configurar el esquema.

## Ejemplo: `User`
```go
type User struct {
    ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Email        string    `gorm:"unique;not null"`
    PasswordHash string    `gorm:"not null"`
    CreatedAt    time.Time
}
```
Aquí se define que `ID` es un UUID generado automáticamente por la DB y que `Email` debe ser único.
