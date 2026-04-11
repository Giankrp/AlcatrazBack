# 📦 Database Package (`db`)

Gestiona la conexión a **PostgreSQL** y las migraciones automáticas de esquema usando [GORM](https://gorm.io/).

---

## Responsabilidades

- Establecer la conexión con PostgreSQL mediante la variable de entorno `DATABASE_URL`
- Ejecutar migraciones automáticas para mantener el esquema sincronizado con los modelos Go
- Proveer la instancia `*gorm.DB` al resto de la aplicación (se inyecta en los Repositories)

---

## Archivos

### `dbConnection.go`

| Función | Firma | Descripción |
|---|---|---|
| `NewConnection` | `() (*gorm.DB, error)` | Lee `DATABASE_URL` del entorno y abre la conexión con PostgreSQL via GORM |
| `AutoMigrate` | `(db *gorm.DB) error` | Crea/actualiza tablas para todos los modelos registrados |

### Modelos Migrados

La función `AutoMigrate` sincroniza las siguientes tablas:

| Modelo | Tabla generada | Descripción |
|---|---|---|
| `User` | `users` | Usuarios registrados (email + hash) |
| `VaultItem` | `vault_items` | Items de la bóveda (metadatos visibles) |
| `VaultSecret` | `vault_secrets` | Datos cifrados de cada item (blob + IV + salt) |
| `VaultFolder` | `vault_folders` | Carpetas para organizar items |
| `Session` | `sessions` | Sesiones activas del usuario |
| `UserProfile` | `user_profiles` | Perfil público del usuario (nombre, avatar, idioma) |

---

## Configuración Requerida

La variable de entorno `DATABASE_URL` debe estar configurada en el archivo `.env`:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5431/alcatraz?sslmode=disable
```

---

## Uso en `main.go`

```go
// 1. Conectar
database, err := db.NewConnection()
if err != nil {
    e.Logger.Fatal("Error connecting to database: ", err)
}

// 2. Migrar esquema
if err := db.AutoMigrate(database); err != nil {
    e.Logger.Fatal("Error migrating database: ", err)
}
```

---

## Docker

La base de datos se levanta con Docker Compose (ver `docker-compose.yml` en la raíz):

```bash
docker compose up -d
```

El contenedor usa **PostgreSQL 16 Alpine** con healthcheck integrado.
