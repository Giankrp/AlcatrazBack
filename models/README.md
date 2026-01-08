# Models Package

Este paquete define las **Entidades de Dominio** y su representación en la base de datos.

## Responsabilidades

1.  **Estructura de Datos**: Define los structs (`User`, `VaultItem`) que mapean a tablas SQL.
2.  **Etiquetas GORM**: Configura tipos de columnas, índices, claves primarias y relaciones (Foreign Keys) mediante tags (ej: `gorm:"primaryKey"`).
3.  **Hooks**: (Opcional) Define lógica que se ejecuta antes/después de guardar (ej: generar UUIDs si no se hace en DB).

## Modelos Principales

### `User`
Representa al usuario registrado.
*   Guarda `Email` y `PasswordHash` (Hash de Argon2 del hash de autenticación).
*   **No guarda** la contraseña maestra en texto plano.

### `VaultItem`
Representa un elemento guardado en la bóveda (contraseña, tarjeta, nota).
*   **Zero Knowledge**: Los datos sensibles viven en `EncryptedData` (blob cifrado).
*   **Metadatos Visibles**: `Title`, `Type` (password, card, note, identity), `FolderID`.
*   **Metadatos Criptográficos**: `IV`, `Salt` necesarios para el descifrado en el cliente.
