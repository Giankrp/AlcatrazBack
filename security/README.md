# Seguridad: hash de contraseñas (Argon2id)

Este proyecto implementa hashing de contraseñas con **Argon2id** en `security/hash.go`.

## Dónde vive y qué se guarda

- La contraseña **nunca** se debe guardar en texto plano.
- El modelo `models.User` guarda el resultado del hash en `PasswordHash` (`models/user.go:5-10`).
- El hash se produce con `security.HashPassword` (`security/hash.go:37-51`) y se valida con `security.VerifyPassword` (`security/hash.go:53-79`).

## Algoritmo usado

Se usa `argon2.IDKey` (Argon2id) del paquete `golang.org/x/crypto/argon2` (`security/hash.go:10-11`, `security/hash.go:44-47`).

Argon2id combina propiedades de Argon2i y Argon2d:

- Resiste ataques por canal lateral (más cercano a Argon2i).
- Resiste ataques con GPU/ASIC y cracking masivo (propio de Argon2 con alta memoria).

## Parámetros (DefaultParams)

Los parámetros por defecto están en `security/hash.go:21-27`:

- `Memory`: `64 * 1024` KiB (64 MiB)
- `Iterations`: `3`
- `Parallelism`: `2`
- `SaltLength`: `16` bytes
- `KeyLength`: `32` bytes

Estos valores controlan el costo del hash:

- **Memoria**: el componente clave para hacer caro el cracking con GPU/ASIC.
- **Iteraciones**: cuántas pasadas hace Argon2 sobre la memoria.
- **Paralelismo**: cuántos “lanes” se computan en paralelo.

## Formato del hash almacenado

`HashPassword` devuelve un string autocontenido que incluye:

1. El algoritmo (`argon2id`)
2. La versión (`v=19`)
3. Los parámetros (`m=...,t=...,p=...`)
4. La sal (salt) en Base64 sin padding
5. El hash derivado en Base64 sin padding

El formato exacto (`security/hash.go:48-50`) es:

```
$argon2id$v=19$m=<memory>,t=<iterations>,p=<parallelism>$<saltB64Raw>$<hashB64Raw>
```

Ejemplo de estructura (valores ilustrativos):

```
$argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
```

Detalles de encoding:

- `salt` y `hash` se codifican con `base64.RawStdEncoding` (`security/hash.go:45-46`).
- “Raw” significa **sin `=` de padding** al final.
- Usa el alfabeto Base64 “standard” (puede contener `+` y `/`), lo cual es perfecto para DB.

## Cómo se genera el hash (paso a paso)

La función `HashPassword(password string)` (`security/hash.go:37-51`) hace:

1. Toma `DefaultParams` (`security/hash.go:38`).
2. Genera una sal aleatoria criptográficamente segura:
   - `generatedRandomBytes(params.SaltLength)` (`security/hash.go:39-42`)
   - Internamente usa `crypto/rand.Read` (`security/hash.go:29-35`).
3. Deriva la clave (hash) con Argon2id:
   - `argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)` (`security/hash.go:44`).
4. Codifica `salt` y `hash` a Base64 “raw” (`security/hash.go:45-46`).
5. Construye el string final autocontenido con `fmt.Sprintf` (`security/hash.go:48-50`).

Propiedad importante: aunque repitas el mismo password, el hash cambia por la sal aleatoria.

## Cómo se verifica un password (paso a paso)

La función `VerifyPassword(password, encoded string)` (`security/hash.go:53-79`) valida así:

1. Separa el string por `$` (`security/hash.go:54`).
2. Valida formato mínimo:
   - Deben ser 6 partes y el algoritmo debe ser `argon2id` (`security/hash.go:55-57`).
3. Extrae:
   - `paramsStr` = `m=..,t=..,p=..` (`security/hash.go:59-61`)
   - `saltStr` y `hashStr` (`security/hash.go:60-61`)
4. Parsea parámetros con `fmt.Sscanf` (`security/hash.go:63-68`).
5. Decodifica `salt` y `expected hash` desde Base64 raw (`security/hash.go:69-76`).
6. Recalcula el hash con los mismos parámetros y la misma sal:
   - La longitud de salida se toma de `len(expected)` para “matchear” exactamente (`security/hash.go:77`).
7. Compara `expected` vs `computed` con comparación en tiempo constante (`security/hash.go:78-79`).

### Por qué comparación en tiempo constante

La comparación `subtleConstantTimeCompare` (`security/hash.go:95-104`) evita filtraciones por timing donde un atacante puede inferir cuántos bytes coincidieron si la comparación “corta” antes.

Nota: la función devuelve `false` si los largos difieren (`security/hash.go:96-98`). Eso es estándar y no es un problema práctico aquí porque el formato decodificado fija el largo del hash esperado.

## Política de rehash (migración de parámetros)

La función `NeedsRehash(encoded string, current ArgonParams)` (`security/hash.go:81-93`) sirve para detectar hashes viejos:

- Parsea el hash almacenado.
- Compara `m`, `t`, `p` del hash con los parámetros “current”.
- Si no coinciden, devuelve `true`.

Uso típico:

1. En login, si `VerifyPassword` da `true`, consultás `NeedsRehash`.
2. Si `NeedsRehash` es `true`, recalculás con `HashPassword` y actualizás `User.PasswordHash`.

Esto permite subir costos con el tiempo sin forzar reset de contraseñas.

## Recomendación de integración (registro y login)

Actualmente `handlers/auth.go:12-23` valida el request pero todavía no persiste usuario ni hashea. El flujo esperado para registrar sería:

1. Validar DTO (`dto/auth.go:3-6`).
2. `encoded, err := security.HashPassword(dto.Password)`.
3. Guardar `encoded` en `models.User.PasswordHash`.

Y para login:

1. Buscar usuario por email.
2. `ok, err := security.VerifyPassword(dto.Password, user.PasswordHash)`.
3. Si `ok`, emitir sesión/token.
4. Opcional: `NeedsRehash` y actualizar hash.

## Pruebas

Hay pruebas unitarias en `security/hash_test.go` que cubren:

- Hash y verificación correcta/incorrecta
- Formato inválido
- Detección de `NeedsRehash`

