# 📦 Models Package (`models`)

Define las **Entidades de Dominio** y su representación como tablas en PostgreSQL, gestionadas mediante [GORM](https://gorm.io/).

---

## Responsabilidades

1. **Estructura de Datos**: Define los structs Go que mapean directamente a tablas SQL.
2. **Configuración GORM**: Tipos de columna, índices, foreign keys y constraints.
3. **Relaciones**: Define asociaciones entre tablas (HasOne, BelongsTo, CASCADE).
4. **Seguridad**: Almacena metadatos críticos para la arquitectura **Zero Knowledge**.

---

## Diagrama de Entidades

```mermaid
erDiagram
    User ||--o{ VaultItem : "tiene"
    User ||--|| UserProfile : "tiene"
    User ||--o{ VaultFolder : "tiene"
    VaultItem ||--|| VaultSecret : "tiene"
    VaultFolder |o--o{ VaultItem : "contiene"

    User {
        uuid ID PK
        string Email UK
        string PasswordHash
        string RecoveryKeyHash
        string ProtectedMasterKey
        string MasterKeyIV
        string MasterKeySalt
        string RecoveryProtectedMasterKey
        string RecoveryKeyIV
        string RecoveryKeySalt
        bool TwoFactorEnabled
        string TwoFactorSecret
        json TwoFactorBackupCodes
        timestamp CreatedAt
    }

    UserProfile {
        uuid UserID PK, FK
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
        int SecurityScore
        bool Trashed
        timestamp CreatedAt
        timestamp UpdatedAt
        timestamp DeletedAt
    }

    VaultSecret {
        uuid VaultItemID PK, FK
        string EncryptedData
        string IV
        string Salt
    }

    VaultFolder {
        uuid ID PK
        uuid UserID FK
        string Name
        bool IsDefault
        timestamp CreatedAt
    }
```

---

## Archivos e Identidad

### `user.go` — `User`

El modelo `User` es el pilar de la seguridad Zero Knowledge.

| Campo                        | Propósito                                                                          |
| ---------------------------- | ---------------------------------------------------------------------------------- |
| `PasswordHash`               | Hash de la clave de autenticación (no es la Master Password)                       |
| `RecoveryKeyHash`            | Hash para validar el uso de la clave de recuperación                               |
| `ProtectedMasterKey`         | El "tesoro": la Master Key del usuario cifrada para que el servidor solo la guarde |
| `RecoveryProtectedMasterKey` | Copia de la MK cifrada con la Recovery Key para emergencias                        |
| `TwoFactorEnabled`           | Flag de activación de MFA                                                          |
| `TwoFactorSecret`            | Semilla TOTP cifrada/protegida                                                     |

---

### `vault.go` — `VaultItem` + `VaultSecret`

#### `VaultItem` (Metadatos)

- **Security Score**: Puntuación (0-100) que indica la fortaleza de la contraseña o salud del ítem.
- **Trashed**: Flag para borrado suave (papelera).

#### `VaultSecret` (Carga útil)

Contiene el `EncryptedData`, que es un JSON cifrado en el cliente con la Master Key del usuario. Separado para optimizar listados.

---

### `vault_folder.go` — `VaultFolder`

- **IsDefault**: Indica si es la carpeta raíz ("Personal") que no puede eliminarse.
- Al eliminar una carpeta no-default, sus ítems se reasignan automáticamente a la carpeta default mediante una transacción.

---

### `session.go` — `Session`

Modelo reservado para gestión de sesiones activas. Actualmente se migra a la BD pero no se usa activamente (la autenticación se gestiona via cookies JWT stateless). Está preparado para implementar revocación de sesiones o gestión multi-dispositivo en el futuro.

---

### `user_profile.go` — `UserProfile`

Datos no sensibles del usuario separados del modelo `User` (nombre para mostrar, avatar, idioma). Se crea automáticamente al registrarse y se recupera junto al perfil con `Preload("User")`.

---

## Convenciones de Implementación

- **UUID**: Se utilizan UUIDs (v4) en lugar de IDs incrementales para mayor seguridad y facilidad de sincronización.
- **JSONB**: Se utiliza `datatypes.JSON` de GORM para campos flexibles como los backup codes.
- **Constraints**: Se aplican `OnDelete: CASCADE` para asegurar que al borrar un usuario se eliminen todos sus secretos.
