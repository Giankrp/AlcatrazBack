# 📦 DTO Package (`dto`)

DTO = **Data Transfer Object**. Este paquete define las estructuras que representan la entrada/salida de datos entre el cliente y la API.

---

## Responsabilidades

- Definir la estructura exacta de los JSON de entrada para cada endpoint
- Desacoplar la estructura interna de la BD (Models) de la interfaz pública de la API
- Incluir reglas de validación (`validate:"..."`) usando [go-playground/validator](https://github.com/go-playground/validator)

---

## Diferencia entre DTOs y Models

| Aspecto | DTO | Model |
|---|---|---|
| **Propósito** | Intención del usuario (lo que envía) | Persistencia en BD (lo que se guarda) |
| **Ejemplo** | `Password` (texto plano para login) | `PasswordHash` (hash Argon2id) |
| **Validación** | Tags `validate:"..."` | Tags `gorm:"..."` |
| **Visibilidad** | Contrato público de la API | Estructura interna |

---

## Archivos y Estructuras

### `auth.go` — Autenticación

```go
type RegisterDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type LoginDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}
```

| Campo | Validación | Usado en |
|---|---|---|
| `Email` | Requerido, formato email válido | Register, Login |
| `Password` | Requerido, mín. 8 caracteres (registro) | Register, Login |

---

### `vault.go` — Bóveda de Secretos

#### Tipos de Item

```go
const (
    ItemTypePassword VaultItemType = "password"
    ItemTypeNote     VaultItemType = "note"
    ItemTypeCard     VaultItemType = "card"
    ItemTypeIdentity VaultItemType = "identity"
)
```

#### `CreateVaultItemDTO`

```go
type CreateVaultItemDTO struct {
    FolderID *string       `json:"folder_id"`
    ItemType VaultItemType `json:"type" validate:"required,oneof=password note card identity"`
    Title    string        `json:"title" validate:"required"`
    Icon     string        `json:"icon"`
    Secret   Secret        `json:"secret" validate:"required"`
}
```

#### `Secret` (Blob cifrado)

```go
type Secret struct {
    Data string `json:"data" validate:"required"`  // Datos cifrados (base64)
    Iv   string `json:"iv" validate:"required"`    // Vector de inicialización
    Salt string `json:"salt" validate:"required"`  // Salt del cifrado
}
```

> ⚠️ Estos campos contienen datos **ya cifrados** por el cliente. El servidor los almacena tal cual (Zero Knowledge).

#### `UpdateVaultItemDTO`

```go
type UpdateVaultItemDTO struct {
    FolderID *string       `json:"folder_id"`
    ItemType VaultItemType `json:"type" validate:"omitempty,oneof=password note card identity"`
    Title    string        `json:"title"`
    Icon     string        `json:"icon"`
    Trashed  *bool         `json:"trashed"`
    Secret   Secret        `json:"secret" validate:"required"`
}
```

#### `CreateVaultFolderDTO`

```go
type CreateVaultFolderDTO struct {
    Name string `json:"name" validate:"required"`
}
```

---

### `user_profileDto.go` — Perfil de Usuario

```go
type UpdateUserProfileDTO struct {
    Name      string `json:"name" validate:"omitempty,min=1,max=50"`
    AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
    Language  string `json:"language" validate:"omitempty,oneof=es en fr de pt"`
}
```

| Campo | Validación | Valores permitidos |
|---|---|---|
| `Name` | Opcional, 1-50 caracteres | Texto libre |
| `AvatarURL` | Opcional, URL válida | URL |
| `Language` | Opcional, valor fijo | `es`, `en`, `fr`, `de`, `pt` |

---

## Ejemplo de JSON de Entrada

### Registro
```json
{
    "email": "user@example.com",
    "password": "mi_password_segura"
}
```

### Crear Item de Bóveda
```json
{
    "type": "password",
    "title": "Mi cuenta de GitHub",
    "icon": "github",
    "folder_id": null,
    "secret": {
        "data": "base64_encrypted_blob...",
        "iv": "base64_iv...",
        "salt": "base64_salt..."
    }
}
```
