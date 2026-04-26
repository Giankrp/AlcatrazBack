# 📦 Security Package (`security`)

Provee utilidades criptográficas críticas para la seguridad del backend. Implementa el hashing de contraseñas con **Argon2id** y generación de entropía segura.

---

## Responsabilidades

- **Hashing**: Protege la AuthKey y la Recovery Key usando Argon2id.
- **Verificación**: Comparación segura en tiempo constante para evitar ataques de temporización.
- **Entropía**: Generación de strings aleatorios seguros para backup codes y sales.

---

## Archivos e Implementación

### `hash.go` — Criptografía

#### Argon2id (Hashing)

Se utiliza el formato PHC string para almacenar los hashes:
`$argon2id$v=19$m=65536,t=3,p=2${salt_base64}${hash_base64}`

| Función | Propósito |
|---|---|
| `HashPassword` | Genera un hash Argon2id con salt aleatorio de 16 bytes. |
| `VerifyPassword` | Compara una contraseña con un hash decodificando sus parámetros originales. |
| `NeedsRehash` | Comprueba si el hash fue generado con parámetros antiguos. |

#### Utilidades Aleatorias

| Función | Propósito |
|---|---|
| `GenerateRandomString(n)` | Genera un string alfanumérico de longitud `n` usando `crypto/rand`. |

> 💡 **Uso**: `GenerateRandomString` se utiliza para generar los 8 códigos de respaldo (Backup Codes) cuando el usuario activa el 2FA.

---

## Comparación de Tiempo Constante

Para evitar **timing attacks**, la comparación de los hashes resultantes se realiza byte a byte sin cortocircuito:

```go
func subtleConstantTimeCompare(a, b []byte) bool {
    // ... XOR bit a bit de todos los bytes
    return diff == 0
}
```

---

## Tests

El paquete incluye tests exhaustivos para asegurar que la verificación es correcta y que el formato del hash es compatible.

```bash
go test ./internal/security/ -v
```
