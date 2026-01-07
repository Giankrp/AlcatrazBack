# Database Package (`db`)

Este paquete gestiona la conexión a la base de datos PostgreSQL utilizando **GORM**.

## Responsabilidades
- Establecer la conexión con la base de datos.
- Ejecutar migraciones automáticas de esquemas.
- Proveer la instancia de conexión a otras capas (Repositories).

## Componentes Principales

### `NewConnection() (*gorm.DB, error)`
- Lee la variable de entorno `DATABASE_URL`.
- Inicializa la conexión con PostgreSQL.
- Retorna una instancia de `*gorm.DB` o un error si falla.

### `AutoMigrate(db *gorm.DB) error`
- Recibe una conexión activa.
- Ejecuta `AutoMigrate` para los modelos registrados (User, VaultItem, etc.).
- Mantiene el esquema de base de datos sincronizado con las estructuras de Go.

## Uso
Se inicializa en el `main.go` al arrancar la aplicación:

```go
database, err := db.NewConnection()
if err != nil {
    log.Fatal(err)
}

if err := db.AutoMigrate(database); err != nil {
    log.Fatal(err)
}
```
