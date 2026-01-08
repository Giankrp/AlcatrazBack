# Repositories Package

Este paquete maneja el **Acceso a Datos** (DAL - Data Access Layer). Abstrae la base de datos subyacente (PostgreSQL) del resto de la aplicación.

## Responsabilidades

1.  **Consultas SQL**: Ejecuta queries `SELECT`, `INSERT`, `UPDATE`, `DELETE` usando GORM.
2.  **Mapeo Objeto-Relacional**: Convierte registros de base de datos a structs de Go (`models`).
3.  **Transacciones**: (Opcional) Maneja transacciones atómicas si una operación toca múltiples tablas.

## Patrón Repository

Usamos interfaces (`UserRepository`, `VaultRepository`) para definir los contratos. Esto permite:
*   **Testabilidad**: Podemos inyectar repositorios falsos (mocks) en los tests unitarios de los servicios.
*   **Flexibilidad**: Podríamos cambiar GORM por SQLx o MongoDB reimplementando solo esta capa, sin tocar servicios ni handlers.

## Consultas Seguras

Todos los repositorios deben respetar el contexto del usuario (`UserID`).
*   **Mal**: `FindByID(id)` -> Podría devolver un item de otro usuario.
*   **Bien**: `FindByID(id, userID)` -> Garantiza que solo el dueño acceda al dato.
