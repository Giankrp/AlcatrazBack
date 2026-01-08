# Security Package

Este paquete provee utilidades criptográficas críticas para la seguridad del backend.

## Funcionalidades

### Argon2id Hashing (`hash.go`)
Implementa el algoritmo de hashing de contraseñas **Argon2id**, recomendado por OWASP.

*   **Uso**: Se usa para hashear el "Hash de Autenticación" que envía el cliente durante el registro y login.
*   **Parámetros**:
    *   Time: 1
    *   Memory: 64MB
    *   Threads: 4
    *   KeyLen: 32 bytes
*   **Sal (Salt)**: Genera sales aleatorias criptográficamente seguras de 16 bytes para cada hash.

### Comparación de Tiempo Constante
Para prevenir ataques de tiempo (timing attacks), la comparación de hashes y contraseñas nunca debe usar operadores estándar (`==`). Este paquete asegura comparaciones seguras.
