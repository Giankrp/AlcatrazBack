# 📦 Models Package (`models`)

Define las **Entidades de Dominio** y su representación como tablas en PostgreSQL, gestionadas mediante [GORM](https://gorm.io/).

---

## Responsabilidades

1. **Estructura de Datos**: Define los structs Go que mapean directamente a tablas SQL
2. **Configuración GORM**: Tipos de columna, índices, claves primarias, foreign keys y constraints via tags
3. **Relaciones**: Define asociaciones entre tablas (HasOne, BelongsTo, CASCADE)
4. **Soft Deletes**: Soporte para borrado lógico en `VaultItem`

---

## Diagrama de Entidades

```mermaid
erDiagram
    User ||--o{ VaultItem : "tiene"
    User ||--|| UserProfile : "tiene"
    User ||--o{ Session : "tiene"
    User ||--o{ VaultFolder : "tiene"
    VaultItem ||--|| VaultSecret : "tiene"
    VaultFolder ||--o{ VaultItem : "contiene"

    User {
        uuid ID PK
        string Email UK
        string PasswordHash
        timestamp CreatedAt
    }

    UserProfile {
        uuid UserID PK_FK
        string Name
        string AvatarURL
        string Language
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    VaultItem {
        uuid ID PK
        uuid UserID FK
        uuid FolderID FK
        string ItemType
        string Title
        string Icon
        bool Trashed
        timestamp CreatedAt
        timestamp UpdatedAt
        timestamp DeletedAt
    }

    VaultSecret {
        uuid VaultItemID PK_FK
        string EncryptedData
        string IV
        string Salt
    }

    VaultFolder {
        uuid ID PK
        uuid UserID FK
        string Name
        timestamp CreatedAt
    }

    Session {
        uuid ID PK
        uuid UserID FK
        string DeviceID
        string IP
        string UserAgent
        timestamp ExpiresAt
        timestamp CreatedAt
    }
```

---

## Archivos y Modelos

### `user.go` — `User`

Representa a un usuario registrado en el sistema.

```go
type User struct {
    ID           string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Email        string    `gorm:"unique;not null"`
    PasswordHash string    `gorm:"not null"`
    CreatedAt    time.Time
}
```

| Campo | Tipo BD | Descripción |
|---|---|---|
| `ID` | `UUID` (auto-generado) | Identificador único |
| `Email` | `VARCHAR UNIQUE NOT NULL` | Email del usuario |
| `PasswordHash` | `VARCHAR NOT NULL` | Hash Argon2id del AuthKey |
| `CreatedAt` | `TIMESTAMP` | Fecha de registro |

> 🔒 **Zero Knowledge**: Solo se guarda el hash del AuthKey, nunca la contraseña maestra.

---

### `user_profile.go` — `UserProfile`

Perfil público del usuario. Relación 1:1 con `User`.

```go
type UserProfile struct {
    UserID    string `gorm:"type:uuid;primaryKey;not null"`
    User      User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    Name      string `gorm:"default:''"`
    AvatarURL string `gorm:"default:''"`
    Language  string `gorm:"default:'es'"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

| Campo | Default | Descripción |
|---|---|---|
| `Name` | `""` | Nombre del usuario (extraído del email al registrarse) |
| `AvatarURL` | `""` | URL de la imagen de perfil |
| `Language` | `"es"` | Idioma preferido (`es`, `en`, `fr`, `de`, `pt`) |

---

### `vault.go` — `VaultItem` + `VaultSecret`

Representan un elemento de la bóveda con separación entre metadatos y datos cifrados.

#### `VaultItem` (Metadatos visibles)

```go
type VaultItem struct {
    ID       string        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    UserID   string        `gorm:"index;not null"`
    User     User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
    FolderID *string       `gorm:"index"`
    ItemType VaultItemType `gorm:"column:item_type;not null;index"`
    Title    string        `gorm:"not null"`
    Icon     string        `gorm:"default:'default_icon'"`
    Trashed  bool          `gorm:"default:false;index"`
    Secret   *VaultSecret  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time   `gorm:"index"` // Soft delete
}
```

**Tipos de item soportados:**

| Constante | Valor | Uso |
|---|---|---|
| `ItemTypePassword` | `"password"` | Credenciales de acceso |
| `ItemTypeNote` | `"note"` | Notas seguras |
| `ItemTypeCard` | `"card"` | Tarjetas de crédito/débito |
| `ItemTypeIdentity` | `"identity"` | Datos de identidad personal |

#### `VaultSecret` (Datos cifrados, 1:1)

```go
type VaultSecret struct {
    VaultItemID   string `gorm:"primaryKey"`      // Mismo ID que VaultItem
    EncryptedData string `gorm:"not null"`         // Blob cifrado
    IV            string `gorm:"not null"`         // Vector de Inicialización
    Salt          string `gorm:"not null"`         // Salt del cifrado
}
```

> 💡 **Optimización**: `VaultSecret` está separado de `VaultItem` para que las queries de lista (`GetItems`) no carguen los blobs cifrados. Solo se hace `Preload("Secret")` al obtener un item individual.

#### Structs Auxiliares

- `VaultItemMeta`: Estructura para datos no cifrados futuros (ej. tags)
- `VaultItemPublicData`: Soporte JSONB consultable para datos públicos futuros

---

### `vault_folder.go` — `VaultFolder`

Carpeta para organizar items de la bóveda.

```go
type VaultFolder struct {
    ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    UserID    string    `gorm:"index;not null"`
    Name      string    `gorm:"not null"`
    CreatedAt time.Time
}
```

---

### `session.go` — `Session`

Sesión activa de un usuario (preparado para gestión multi-dispositivo).

```go
type Session struct {
    ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    UserID    string    `gorm:"index;not null"`
    DeviceID  *string   `gorm:"index"`
    IP        string
    UserAgent string
    ExpiresAt time.Time `gorm:"index"`
    CreatedAt time.Time
}
```

| Campo | Descripción |
|---|---|
| `DeviceID` | Identificador del dispositivo (opcional) |
| `IP` | Dirección IP de la sesión |
| `UserAgent` | User-Agent del navegador/cliente |
| `ExpiresAt` | Fecha de expiración de la sesión |

> 📌 **Nota**: El modelo `Session` está definido y migrado, pero la gestión activa de sesiones aún no está implementada en los services. Actualmente la autenticación se gestiona exclusivamente mediante JWT en cookies.

---

## Convenciones GORM

| Tag | Significado |
|---|---|
| `type:uuid` | Columna de tipo UUID de PostgreSQL |
| `default:gen_random_uuid()` | UUID auto-generado por la BD |
| `primaryKey` | Clave primaria |
| `unique` | Constraint UNIQUE |
| `index` | Crea un índice en la columna |
| `not null` | Campo obligatorio |
| `constraint:OnDelete:CASCADE` | Borrado en cascada al eliminar el padre |
| `gorm:"index"` en `DeletedAt` | Habilita soft deletes de GORM |
