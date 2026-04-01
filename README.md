# 🔒 Alcatraz Backend

Backend del gestor de contraseñas **Alcatraz** — una aplicación de bóveda de secretos con arquitectura **Zero Knowledge**. Construido en **Go** con [Echo](https://echo.labstack.com/), [GORM](https://gorm.io/) y [Argon2id](https://pkg.go.dev/golang.org/x/crypto/argon2).

> **Zero Knowledge**: El servidor nunca ve ni tiene acceso a los datos en texto plano del usuario. Actúa únicamente como un almacén ciego de blobs cifrados.

---

## ✨ Características Principales

| Característica | Detalle |
|---|---|
| **Arquitectura Limpia** | Separación clara en capas: Handlers → Services → Repositories |
| **Zero Knowledge** | Los datos sensibles se cifran/descifran exclusivamente en el cliente |
| **Autenticación** | Hashing con Argon2id + JWT en cookies `HttpOnly` |
| **Base de Datos** | PostgreSQL 16 con GORM (auto-migraciones + soft deletes) |
| **Validación** | DTOs con `go-playground/validator` en cada endpoint |
| **Perfiles de Usuario** | Nombre, avatar, idioma configurable |
| **Tipos de Vault** | Passwords, notas, tarjetas, identidades |
| **Carpetas** | Organización de items por carpetas |

---

## 📁 Estructura del Proyecto

```
AlcatrazBack/
├── main.go             # Punto de entrada — inicialización y wiring de dependencias
├── docker-compose.yml  # PostgreSQL 16 Alpine containerizado
├── .env.example        # Variables de entorno requeridas
│
├── db/                 # Conexión a PostgreSQL y auto-migraciones
├── docs/               # Documentación detallada (DATA_FLOW.md)
├── dto/                # Data Transfer Objects — validación de entrada
├── handlers/           # Controladores HTTP (Echo) — capa de transporte
├── middleware/          # Middlewares personalizados (reservado para extensión)
├── models/             # Entidades de dominio y esquema de BD (GORM)
├── repositories/       # Acceso a datos — patrón Repository con interfaces
├── routes/             # Definición de rutas, grupos y JWT middleware
├── security/           # Criptografía: Argon2id hashing + comparación segura
├── services/           # Lógica de negocio — capa de aplicación
└── validator/          # Instancia global de validador + formateo de errores
```

---

## 🚀 Inicio Rápido

### Requisitos Previos

- **Go** 1.25+ (ver `go.mod`)
- **Docker** + **Docker Compose** (para PostgreSQL)

### 1. Clonar el repositorio

```bash
git clone https://github.com/Giankrp/AlcatrazBack.git
cd AlcatrazBack
```

### 2. Configurar variables de entorno

```bash
cp .env.example .env
```

Edita `.env` según tu configuración:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=alcatraz
POSTGRES_PORT=5431
JWT_SECRET=tu_secreto_super_seguro_aqui
DATABASE_URL=postgres://postgres:postgres@localhost:5431/alcatraz?sslmode=disable
ALLOWED_ORIGINS=http://localhost:3000
```

| Variable | Descripción | Default |
|---|---|---|
| `POSTGRES_USER` | Usuario de PostgreSQL | `postgres` |
| `POSTGRES_PASSWORD` | Contraseña de PostgreSQL | `postgres` |
| `POSTGRES_DB` | Nombre de la base de datos | `alcatraz` |
| `POSTGRES_PORT` | Puerto expuesto para PostgreSQL | `5431` |
| `JWT_SECRET` | Clave secreta para firmar tokens JWT | *(requerido)* |
| `DATABASE_URL` | Cadena de conexión completa | Construir con las variables anteriores |
| `ALLOWED_ORIGINS` | Orígenes permitidos CORS (separados por `,`) | `http://localhost:3000` |
| `PORT` | Puerto del servidor HTTP | `8080` |

### 3. Levantar la base de datos

```bash
docker compose up -d
```

El contenedor incluye **healthcheck** automático. Puedes verificar el estado con:

```bash
docker compose ps
```

### 4. Instalar dependencias e iniciar

```bash
go mod download
go run main.go
```

El servidor estará disponible en `http://localhost:8080`.

---

## 🌐 API Endpoints

Base URL: `/api`

### Autenticación (`/api/auth`) — Público

| Método | Ruta | Descripción | Body |
|---|---|---|---|
| `POST` | `/api/auth/register` | Registrar nuevo usuario | `{ email, password }` |
| `POST` | `/api/auth/login` | Iniciar sesión (devuelve cookie JWT) | `{ email, password }` |
| `POST` | `/api/auth/logout` | Cerrar sesión (expira la cookie) | — |

### Perfil de Usuario (`/api/user`) — 🔐 Protegido

| Método | Ruta | Descripción | Body |
|---|---|---|---|
| `GET` | `/api/user/profile` | Obtener perfil del usuario | — |
| `PUT` | `/api/user/profile` | Actualizar perfil | `{ name?, avatar_url?, language? }` |

### Bóveda (`/api/vault`) — 🔐 Protegido

| Método | Ruta | Descripción | Body |
|---|---|---|---|
| `POST` | `/api/vault/items` | Crear item en la bóveda | `{ type, title, icon?, folder_id?, secret }` |
| `GET` | `/api/vault/items` | Listar todos los items | — |
| `GET` | `/api/vault/items/:id` | Obtener item por ID (con secreto) | — |
| `PUT` | `/api/vault/items/:id` | Actualizar item | `{ type?, title?, icon?, folder_id?, trashed?, secret }` |
| `DELETE` | `/api/vault/items/:id` | Eliminar item (soft delete) | — |

> 🔐 Los endpoints protegidos requieren un token JWT válido en la cookie `auth_token` (HttpOnly).

---

## 🏗️ Arquitectura

```mermaid
graph LR
    Client([Cliente/Frontend]) -->|HTTPS| Handlers
    Handlers -->|DTO validado| Services
    Services -->|Modelo| Repositories
    Repositories -->|GORM| DB[(PostgreSQL)]

    subgraph "Alcatraz Backend"
        Handlers
        Services
        Repositories
    end

    style Client fill:#1a1a2e,color:#e94560
    style DB fill:#0f3460,color:#e94560
```

### Flujo de una petición típica

1. **Handler** recibe la petición HTTP, parsea el body a un **DTO**, valida con `validator`
2. **Service** aplica reglas de negocio, transforma DTO → **Model**
3. **Repository** ejecuta la consulta SQL vía GORM
4. La respuesta viaja de vuelta: Repository → Service → Handler → JSON Response

### Inyección de Dependencias

El wiring se realiza en `main.go` siguiendo el patrón de constructor:

```go
// Repositories
userRepo := repositories.NewUserRepository(database)
vaultRepo := repositories.NewVaultRepository(database)

// Services (dependen de repositories)
authService := services.NewAuthService(userRepo)
vaultService := services.NewVaultService(vaultRepo)
userService := services.NewUserService(userRepo)

// Handlers (dependen de services)
authHandler := handlers.NewAuthHandler(authService)
vaultHandler := handlers.NewVaultHandler(vaultService)
userProfileHandler := handlers.NewUserProfileHandler(userService)
```

---

## 🔐 Seguridad

### Modelo Zero Knowledge

El backend **nunca** accede a datos en texto plano. El cliente es responsable de:

1. Derivar una **AuthKey** desde la contraseña maestra (para autenticarse)
2. Derivar una **EncryptionKey** desde la contraseña maestra (para cifrar/descifrar)
3. Cifrar todos los datos sensibles antes de enviarlos al servidor

Para más detalle, ver [docs/DATA_FLOW.md](docs/DATA_FLOW.md).

### Hashing de Contraseñas

Se utiliza **Argon2id** con los siguientes parámetros:

| Parámetro | Valor |
|---|---|
| Memoria | 64 MB |
| Iteraciones | 3 |
| Paralelismo | 2 |
| Salt | 16 bytes aleatorios |
| Key Length | 32 bytes |

### Autenticación JWT

- Los tokens JWT se almacenan en cookies **HttpOnly** (no accesibles desde JavaScript)
- Se configuran con `SameSite=Lax`
- Expiración: **12 horas**
- Claims incluidos: `user_id`, `email`, `exp`

### CORS

Configurado via variable de entorno `ALLOWED_ORIGINS`, soporta `AllowCredentials: true` para envío de cookies.

---

## 🧪 Tests

```bash
# Ejecutar todos los tests
go test ./...

# Tests del paquete security
go test ./security/ -v
```

Actualmente incluye tests para:
- `HashPassword` + `VerifyPassword` (contraseña correcta e incorrecta)
- Validación de formato de hash inválido
- `NeedsRehash` (detección de parámetros cambiados)

---

## 📚 Documentación Detallada

| Documento | Descripción |
|---|---|
| [docs/DATA_FLOW.md](docs/DATA_FLOW.md) | Flujo de datos completo y arquitectura Zero Knowledge |
| [db/README.md](db/README.md) | Conexión a base de datos y migraciones |
| [dto/README.md](dto/README.md) | Data Transfer Objects y validación |
| [handlers/README.md](handlers/README.md) | Controladores HTTP |
| [models/README.md](models/README.md) | Entidades de dominio y esquema |
| [repositories/README.md](repositories/README.md) | Capa de acceso a datos |
| [routes/README.md](routes/README.md) | Definición de rutas y middleware JWT |
| [security/README.md](security/README.md) | Utilidades criptográficas |
| [services/README.md](services/README.md) | Lógica de negocio |
| [validator/README.md](validator/README.md) | Validación de datos |

---

## 🛠️ Stack Tecnológico

| Tecnología | Versión | Uso |
|---|---|---|
| Go | 1.25+ | Lenguaje principal |
| Echo | v4.13 | Framework HTTP |
| GORM | v1.31 | ORM para PostgreSQL |
| PostgreSQL | 16 Alpine | Base de datos relacional |
| Argon2id | — | Hashing de contraseñas |
| JWT (golang-jwt) | v5 | Autenticación por tokens |
| go-playground/validator | v10 | Validación de structs |
| Docker Compose | — | Infraestructura local |
