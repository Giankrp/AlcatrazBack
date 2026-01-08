# Alcatraz Backend

Backend para el gestor de contraseñas Alcatraz. Construido en Go usando Echo, GORM y Argon2id, diseñado con una arquitectura segura y escalable.

## Características Principales

*   **Arquitectura Limpia (Clean Architecture)**: Separación clara de responsabilidades (Handlers, Services, Repositories).
*   **Seguridad Zero Knowledge**: El backend actúa como un almacén ciego de datos cifrados.
*   **Autenticación Robusta**: Hashing de contraseñas con Argon2id y sesiones vía JWT.
*   **Base de Datos**: PostgreSQL (vía GORM) para persistencia relacional y JSONB.
*   **API RESTful**: Endpoints estandarizados y documentados.

## Estructura del Proyecto

```
AlcatrazBack/
├── db/             # Configuración y conexión a base de datos
├── docs/           # Documentación detallada del proyecto
├── dto/            # Data Transfer Objects (Validación de entrada)
├── handlers/       # Controladores HTTP (Entrada/Salida)
├── models/         # Modelos de base de datos (GORM)
├── repositories/   # Acceso a datos (SQL/ORM)
├── routes/         # Definición de rutas y grupos de API
├── security/       # Utilidades criptográficas (Argon2, etc.)
├── services/       # Lógica de negocio
└── main.go         # Punto de entrada y configuración
```

## Requisitos

*   Go 1.20+
*   PostgreSQL (o MySQL configurado en driver)

## Configuración

Crea un archivo `.env` en la raíz del proyecto:

```env
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=password
DB_NAME=alcatraz_db
DB_PORT=3306
JWT_SECRET=tu_secreto_super_seguro
PORT=8080
```

## Ejecución

1.  **Instalar dependencias**:
    ```bash
    go mod download
    ```

2.  **Ejecutar servidor**:
    ```bash
    go run main.go
    ```
    El servidor iniciará en `http://localhost:8080`.

## Documentación Detallada

*   [Flujo de Datos y Seguridad (Zero Knowledge)](docs/DATA_FLOW.md)
*   Ver `README.md` en cada subdirectorio para detalles de implementación específicos.
