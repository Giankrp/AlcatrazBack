# 📦 Repositories Package (`repositories`)

Implementa la **Capa de Acceso a Datos** (Data Access Layer). Abstrae toda interacción con PostgreSQL del resto de la aplicación usando el **Patrón Repository** con interfaces.

---

## Responsabilidades

1. **Ejecutar consultas** SQL via GORM (`SELECT`, `INSERT`, `UPDATE`, `DELETE`)
2. **Mapear** registros de BD a structs Go (`models`)
3. **Aislar** la tecnología de persistencia — se podría reemplazar GORM por SQLx o MongoDB reimplementando solo esta capa
4. **Garantizar aislamiento de datos** — siempre filtrar por `UserID`

---

## Patrón Repository

Cada repository define una **interfaz** (contrato) y una **implementación** privada:

```go
// Contrato público
type VaultRepository interface {
    Create(item *models.VaultItem) error
    FindByID(id string, userID string) (*models.VaultItem, error)
    // ...
}

// Implementación privada
type vaultRepository struct {
    db *gorm.DB
}

// Constructor
func NewVaultRepository(db *gorm.DB) VaultRepository {
    return &vaultRepository{db: db}
}
```

### Ventajas

| Beneficio | Descripción |
|---|---|
| **Testabilidad** | Se pueden inyectar mocks en tests de services |
| **Flexibilidad** | Cambiar de ORM sin tocar services ni handlers |
| **Desacoplamiento** | Los services no conocen GORM ni SQL |

---

## Archivos e Interfaces

### `user_repository.go` — `UserRepository`

Gestiona usuarios y perfiles.

```go
type UserRepository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
    CreateProfile(profile *models.UserProfile) error
    FindProfileByUserID(userID string) (*models.UserProfile, error)
    UpdateProfile(profile *models.UserProfile) error
}
```

| Método | Query | Descripción |
|---|---|---|
| `Create` | `INSERT INTO users` | Crea un nuevo usuario |
| `FindByEmail` | `WHERE email = ?` | Busca usuario por email (para login/registro) |
| `CreateProfile` | `INSERT INTO user_profiles` | Crea perfil inicial al registrarse |
| `FindProfileByUserID` | `WHERE user_id = ?` | Obtiene el perfil de un usuario |
| `UpdateProfile` | `SAVE (upsert)` | Actualiza el perfil del usuario |

---

### `vault_repository.go` — `VaultRepository`

Gestiona los items de la bóveda con consideraciones de rendimiento y seguridad.

```go
type VaultRepository interface {
    Create(item *models.VaultItem) error
    FindByID(id string, userID string) (*models.VaultItem, error)
    FindAllByUserID(userID string) ([]models.VaultItem, error)
    Update(item *models.VaultItem) error
    Delete(id string, userID string) error
}
```

| Método | Query | Notas |
|---|---|---|
| `Create` | `INSERT INTO vault_items + vault_secrets` | Crea item y secreto en cascada |
| `FindByID` | `WHERE id = ? AND user_id = ?` + `Preload("Secret")` | **Incluye** datos cifrados |
| `FindAllByUserID` | `WHERE user_id = ? AND deleted_at IS NULL` | **No incluye** secretos (optimización) |
| `Update` | `Session(FullSaveAssociations: true).Save()` | Actualiza item + secreto en transacción |
| `Delete` | `WHERE id = ? AND user_id = ?` + soft delete | Borrado lógico (GORM `DeletedAt`) |

#### Optimización de Rendimiento

```go
// ❌ FindAllByUserID NO hace Preload("Secret")
// → Los blobs cifrados no se cargan en la lista (ahorro de ancho de banda)

// ✅ FindByID SÍ hace Preload("Secret")
// → Los datos cifrados solo se cargan al ver el detalle de un item
```

---

## Seguridad: Aislamiento por Usuario

Todas las queries que acceden a datos del vault incluyen un filtro por `userID`:

```go
// ✅ Correcto — Solo el dueño puede acceder
db.Where("id = ? AND user_id = ?", id, userID).First(&item)

// ❌ Incorrecto — Podría devolver datos de otro usuario
db.Where("id = ?", id).First(&item)
```

Este patrón previene vulnerabilidades de tipo **Insecure Direct Object Reference (IDOR)**.

---

## Inyección en `main.go`

```go
userRepo := repositories.NewUserRepository(database)
vaultRepo := repositories.NewVaultRepository(database)
```

Los repositories se inyectan en los Services como dependencias.
