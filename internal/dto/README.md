# 📦 DTO Package (`dto`)

DTO = **Data Transfer Object**. Este paquete define las estructuras que representan el contrato de entrada/salida entre el cliente y la API, desacoplando la lógica interna de persistencia de la interfaz pública.

---

## Responsabilidades

- Definir la estructura exacta de los JSON de entrada/salida.
- Incluir reglas de validación (`validate:"..."`) para asegurar la integridad de los datos antes de procesarlos.
- Ocultar campos internos que no deben exponerse directamente (ej. hashes de contraseñas).

---

## Archivos y Estructuras

### `auth.go` — Autenticación y Seguridad

#### `RegisterDTO`
Utilizado para el registro inicial. Incluye todos los metadatos necesarios para el modelo Zero Knowledge.
- **Master Key Metadata**: IV, Salt y la propia MK cifrada con la contraseña.
- **Recovery Metadata**: IV, Salt y MK cifrada con la clave de recuperación.

#### `LoginDTO` & `LoginResponseDTO`
- El login devuelve los metadatos de la Master Key para que el cliente pueda descifrarla localmente.
- Si el 2FA está activo, devuelve `require_2fa: true` y un token temporal.

#### `2FA DTOs`
- `Enable2FADTO`: Para activar el MFA tras verificar el primer código.
- `Verify2FADTO`: Para la verificación durante el proceso de login.

#### `ChangeMasterPasswordDTO`
Permite la rotación de la contraseña maestra actualizando atómicamente el hash de autenticación y los bloques cifrados de la MK.

#### `Recovery DTOs`
- `FetchRecoveryDTO`: Para obtener la metadata de recuperación asociada a un email.
- `ResetPasswordDTO`: Para restaurar el acceso usando la Recovery Key.

---

### `vault.go` — Bóveda y Carpetas

#### `CreateVaultItemDTO` & `UpdateVaultItemDTO`
Representan un ítem de la bóveda.
- **Secret**: Objeto anidado que contiene el blob cifrado (`data`), el `iv` y el `salt`.
- **SecurityScore**: Puntuación de salud del ítem.

#### `Folder DTOs`
- `CreateVaultFolderDTO`: Para nuevas carpetas.
- `UpdateVaultFolderDTO`: Para renombrar carpetas existentes.

---

### `user_profileDto.go` — Perfil

- `UpdateUserProfileDTO`: Para cambios en nombre, avatar e idioma.

---

## Validación

Se utiliza `go-playground/validator`. Algunas reglas comunes aplicadas:
- `required`: El campo no puede estar vacío.
- `email`: Debe ser un formato de correo válido.
- `min=8`: Longitud mínima para contraseñas.
- `oneof=...`: Restringe el valor a una lista permitida (ej. tipos de ítem).
- `url`: Valida que el avatar_url sea una URL válida.

> 🔒 **Importante**: Los DTOs de entrada para datos sensibles (Secret) siempre esperan strings que ya han sido cifrados en el cliente. El servidor nunca recibe el secreto en texto plano.
