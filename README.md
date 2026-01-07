# AlcatrazBack API

Backend para el gestor de contraseñas Alcatraz. Construido con **Go (Golang)**, **Echo**, **GORM** y **PostgreSQL**, siguiendo una arquitectura limpia y modular.

## Arquitectura

El proyecto sigue una estructura de capas (Clean Architecture simplificada) para separar responsabilidades, facilitar el testing y el mantenimiento.

### Capas Principales
1.  **Routes (`routes/`)**: Define los endpoints y conecta URLs con Handlers.
2.  **Handlers (`handlers/`)**: Controladores HTTP. Validan entrada (DTOs) y llaman a Servicios.
3.  **Services (`services/`)**: Lógica de negocio pura. Orquestan operaciones y validan reglas de dominio.
4.  **Repositories (`repositories/`)**: Acceso a datos. Interactúan con la base de datos vía GORM.
5.  **Models (`models/`)**: Entidades de dominio mapeadas a la base de datos.

### Componentes Transversales
-   **DB (`db/`)**: Configuración y conexión a PostgreSQL.
-   **DTO (`dto/`)**: Objetos de transferencia de datos (JSON shapes).
-   **Security (`security/`)**: Utilidades criptográficas (Argon2id para contraseñas).
-   **Validator (`validator/`)**: Reglas de validación de entrada.

## Requisitos

-   Go 1.22+
-   PostgreSQL

## Configuración

Crear un archivo `.env` en la raíz del proyecto con las siguientes variables:

```env
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
JWT_SECRET=tu_secreto_seguro_para_jwt
PORT=8080
```

## Ejecución

1.  **Instalar dependencias**:
    ```bash
    go mod download
    ```

2.  **Ejecutar en desarrollo**:
    ```bash
    go run .
    ```
    El servidor iniciará en `http://localhost:8080` (o el puerto definido).

3.  **Compilar**:
    ```bash
    go build -o alcatraz-api .
    ./alcatraz-api
    ```

## Flujo de Desarrollo

1.  Definir modelo en `models/` (si es nuevo).
2.  Definir operaciones de DB en `repositories/`.
3.  Definir DTOs de entrada/salida en `dto/`.
4.  Implementar lógica en `services/`.
5.  Crear handler en `handlers/`.
6.  Registrar ruta en `routes/`.
7.  Inyectar dependencias en `main.go`.

## Seguridad

-   **Contraseñas**: Hashing robusto con **Argon2id** + Salt aleatoria única por usuario.
-   **Autenticación**: JWT (JSON Web Tokens) firmados con HS256.
