# 📦 Repositories Package (`repositories`)

Implementa la **Capa de Acceso a Datos** (Data Access Layer). Abstrae toda interacción con PostgreSQL del resto de la aplicación usando el **Patrón Repository** con interfaces.

---

## Responsabilidades

1. **Ejecutar consultas** SQL via GORM (`SELECT`, `INSERT`, `UPDATE`, `DELETE`)
2. **Mapear** registros de BD a structs Go (`models`)
3. **Aislar** la tecnología de persistencia — se podría reemplazar GORM por SQLx reimplementando solo esta capa
4. **Garantizar aislamiento de datos** — siempre filtrar por `UserID`

---

## Patrón Repository

Cada repository define una **interfaz** (contrato) y una **implementación** privada:

```go
// Contrato público
type VaultRepository interface {
    Create(item *models.VaultItem) error
    FindByID(id uuid.UUID, userID uuid.UUID) (*models.VaultItem, error)
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
    FindProfileByUserID(userID uuid.UUID) (*models.UserProfile, error)
    UpdateProfile(profile *models.UserProfile) error
    FindByID(id uuid.UUID) (*models.User, error)
    Update(user *models.User) error
    Delete(id uuid.UUID) error
}
```

| Método | Query | Descripción |
|---|---|---|
| `Create` | `INSERT INTO users` | Crea un nuevo usuario |
| `FindByEmail` | `WHERE email = ?` | Busca usuario por email (para login/registro) |
| `CreateProfile` | `INSERT INTO user_profiles` | Crea perfil inicial al registrarse |
| `FindProfileByUserID` | `WHERE user_id = ?` + `Preload("User")` | Obtiene el perfil con datos del usuario |
| `UpdateProfile` | `SAVE (upsert)` | Actualiza el perfil del usuario |
| `FindByID` | `WHERE id = ?` | Busca usuario por UUID (para 2FA y cambio de contraseña) |
| `Update` | `SAVE (upsert)` | Actualiza datos del usuario (contraseña, 2FA, etc.) |
| `Delete` | Transacción | Elimina en cascada: items → carpetas → perfil → usuario |

---

### `vault_repository.go` — `VaultRepository`

Gestiona los items y carpetas de la bóveda con consideraciones de rendimiento y seguridad.

```go
type VaultRepository interface {
    // Item methods
    Create(item *models.VaultItem) error
    FindByID(id uuid.UUID, userID uuid.UUID) (*models.VaultItem, error)
    FindAllByUserID(userID uuid.UUID) ([]models.VaultItem, error)
    FindTrashedByUserID(userID uuid.UUID) ([]models.VaultItem, error)
    Update(item *models.VaultItem) error
    MoveToTrash(id uuid.UUID, userID uuid.UUID) error
    PermanentlyDelete(id uuid.UUID, userID uuid.UUID) error

    // Folder methods
    CreateFolder(folder *models.VaultFolder) error
    FindFoldersByUserID(userID uuid.UUID) ([]models.VaultFolder, error)
    FindFolderByID(id uuid.UUID, userID uuid.UUID) (*models.VaultFolder, error)
    FindDefaultFolder(userID uuid.UUID) (*models.VaultFolder, error)
    UpdateFolder(folder *models.VaultFolder) error
    DeleteFolder(id uuid.UUID, userID uuid.UUID, defaultFolderID uuid.UUID) error
}
```

#### Métodos de Items

| Método | Query | Notas |
|---|---|---|
| `Create` | `INSERT INTO vault_items + vault_secrets` | Crea item y secreto en cascada |
| `FindByID` | `WHERE id = ? AND user_id = ?` + `Preload("Secret")` | **Incluye** datos cifrados |
| `FindAllByUserID` | `WHERE user_id = ? AND trashed = false` | **No incluye** secretos (optimización) |
| `FindTrashedByUserID` | `WHERE user_id = ? AND trashed = true` | Items en la papelera |
| `Update` | `Session(FullSaveAssociations: true).Save()` | Actualiza item + secreto |
| `MoveToTrash` | `UPDATE trashed = true` | Soft-delete (papelera) |
| `PermanentlyDelete` | `Unscoped().Delete()` | Eliminación física definitiva |

#### Métodos de Carpetas

| Método | Descripción |
|---|---|
| `CreateFolder` | Crea una nueva carpeta para el usuario |
| `FindFoldersByUserID` | Lista todas las carpetas del usuario |
| `FindFolderByID` | Busca una carpeta verificando que sea del usuario (anti-IDOR) |
| `FindDefaultFolder` | Obtiene la carpeta "Personal" (IsDefault=true) para reasignación |
| `UpdateFolder` | Actualiza nombre u otros campos de la carpeta |
| `DeleteFolder` | Transacción: reasigna items → elimina carpeta |

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
