# 🔒 Alcatraz Backend

Backend del gestor de contraseñas **Alcatraz** — una aplicación de bóveda de secretos con arquitectura **Zero Knowledge**. Construido en **Go** con [Echo](https://echo.labstack.com/), [GORM](https://gorm.io/) y [Argon2id](https://pkg.go.dev/golang.org/x/crypto/argon2).

> **Zero Knowledge**: El servidor nunca ve ni tiene acceso a los datos en texto plano del usuario. Actúa únicamente como un almacén ciego de blobs cifrados.

---

## ✨ Características Principales

| Característica | Detalle |
|---|---|
| **Arquitectura Limpia** | Separación clara en capas: Handlers → Services → Repositories |
| **Zero Knowledge** | Cifrado/descifrado exclusivo en cliente con Master Key protegida |
| **Autenticación Robusta** | Argon2id + JWT + **2FA (TOTP)** + Backup Codes |
| **Recuperación de Cuenta** | Sistema de Recovery Key para restaurar acceso sin perder datos |
| **Base de Datos** | PostgreSQL 16 con GORM (auto-migraciones + soft deletes) |
| **Seguridad de Datos** | Cálculo de **Security Score** para ítems de la bóveda |
| **Organización** | Gestión de **Carpetas** y Papelera (Trash) |
| **Logs Premium** | Salida estructurada y visual con `charmbracelet/log` |

---

## 📁 Estructura del Proyecto

```
AlcatrazBack/
├── cmd/
│   └── server/
│       └── main.go         # Punto de entrada — inicialización y wiring
├── internal/
│   ├── db/                 # Conexión a PostgreSQL y auto-migraciones
│   ├── dto/                # Data Transfer Objects — validación de entrada
│   ├── handlers/           # Controladores HTTP (Echo) — capa de transporte
│   ├── middleware/          # Middlewares (JWT, Logger, etc.)
│   ├── models/             # Entidades de dominio y esquema de BD (GORM)
│   ├── repositories/       # Acceso a datos — patrón Repository
│   ├── routes/             # Definición de rutas y grupos
│   ├── security/           # Criptografía: Argon2id hashing + comparación
│   ├── services/           # Lógica de negocio — capa de aplicación
│   └── validator/          # Instancia global de validador
├── docs/                   # Documentación detallada
└── docker-compose.yml      # Infraestructura (PostgreSQL 16)
```

---

## 🚀 Inicio Rápido

### Requisitos Previos

- **Go** 1.22+ (ver `go.mod`)
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
PORT=8080
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
| `ENV` | Entorno de ejecución (`production` activa Secure en cookies) | *(vacío = desarrollo)* |

### 3. Levantar la base de datos

```bash
docker compose up -d
```

### 4. Instalar dependencias e iniciar

```bash
go mod download
go run cmd/server/main.go
```

El servidor estará disponible en `http://localhost:8080`.

---

## 🌐 API Endpoints

Base URL: `/api`

### Autenticación (`/api/auth`) — Público

| Método | Ruta | Descripción |
|---|---|---|
| `POST` | `/api/auth/register` | Registrar nuevo usuario |
| `POST` | `/api/auth/login` | Iniciar sesión (devuelve cookie JWT) |
| `POST` | `/api/auth/logout` | Cerrar sesión (expira la cookie) |
| `GET` | `/api/auth/exists` | Verificar si un email ya existe |
| `POST` | `/api/auth/2fa/verify` | Verificar código TOTP durante el login |
| `POST` | `/api/auth/recovery/fetch` | Obtener datos de recuperación (salt/iv) |
| `POST` | `/api/auth/recovery/reset` | Resetear password usando Recovery Key |

### Perfil y Seguridad (`/api/user`) — 🔐 Protegido

| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/api/user/profile` | Obtener perfil del usuario |
| `PUT` | `/api/user/profile` | Actualizar perfil (nombre, avatar, etc.) |
| `POST` | `/api/user/change-password` | Actualizar Master Password (Zero Knowledge) |
| `POST` | `/api/user/2fa/setup` | Generar secreto TOTP para configuración |
| `POST` | `/api/user/2fa/enable` | Activar 2FA tras verificar código |
| `DELETE` | `/api/user/account` | Eliminar cuenta permanentemente |

### Bóveda (`/api/vault`) — 🔐 Protegido

| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/api/vault/items` | Listar todos los items (sin secretos) |
| `POST` | `/api/vault/items` | Crear item en la bóveda |
| `GET` | `/api/vault/items/:id` | Obtener item por ID (con secreto) |
| `PUT` | `/api/vault/items/:id` | Actualizar item |
| `DELETE` | `/api/vault/items/:id` | Mover item a la papelera |
| `GET` | `/api/vault/trash` | Listar items en la papelera |
| `POST` | `/api/vault/items/:id/restore` | Restaurar item de la papelera |
| `DELETE` | `/api/vault/items/:id/permanent` | Eliminar item permanentemente |
| `GET` | `/api/vault/folders` | Listar todas las carpetas |
| `POST` | `/api/vault/folders` | Crear nueva carpeta |
| `PUT` | `/api/vault/folders/:id` | Renombrar/actualizar carpeta |
| `DELETE` | `/api/vault/folders/:id` | Eliminar carpeta |

---

## 🏗️ Arquitectura

```mermaid
graph LR
    Client([Cliente/Frontend]) -->|HTTPS| Handlers
    Handlers -->|DTO validado| Services
    Services -->|Modelo| Repositories
    Repositories -->|GORM| DB[(PostgreSQL)]

    subgraph "Alcatraz Backend (internal/)"
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

---

## 🔐 Seguridad

### Modelo Zero Knowledge

El backend **nunca** accede a datos en texto plano. El cliente es responsable de:

1. Derivar una **AuthKey** desde la contraseña maestra (para autenticarse)
2. Enviar la **Master Key cifrada** al servidor (protegida por el AuthKey o Recovery Key)
3. Cifrar todos los datos sensibles antes de enviarlos al servidor

### Hashing de Contraseñas

Se utiliza **Argon2id** con los siguientes parámetros:

| Parámetro | Valor |
|---|---|
| Memoria | 64 MB |
| Iteraciones | 3 |
| Paralelismo | 2 |

### Autenticación de Dos Factores (2FA)

Soporte nativo para **TOTP** (Time-based One-Time Password). Al activarse, el login requiere un segundo paso de verificación. Se proporcionan backup codes para casos de pérdida del dispositivo.

---

## 🧪 Tests

```bash
# Ejecutar todos los tests
go test ./...

# Tests específicos
go test ./internal/security/ -v
```

---

## 🛠️ Stack Tecnológico

| Tecnología | Versión | Uso |
|---|---|---|
| Go | 1.22+ | Lenguaje principal |
| Echo | v4.13.4 | Framework HTTP |
| GORM | v1.31.1 | ORM para PostgreSQL |
| PostgreSQL | 16 Alpine | Base de datos relacional |
| Argon2id | golang.org/x/crypto | Hashing de contraseñas |
| TOTP (pquerna/otp) | v1.5.0 | Autenticación 2FA |
| JWT (golang-jwt) | v5.3.0 | Autenticación por tokens |
| charmbracelet/log | v1.0.0 | Logging estructurado |
| Docker Compose | v2 | Infraestructura local |
