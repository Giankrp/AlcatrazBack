# 📦 Routes Package (`routes`)

Centraliza la definición y configuración de todas las rutas HTTP de la API, incluyendo la aplicación del middleware JWT para rutas protegidas.

---

## Responsabilidades

- Agrupar endpoints por funcionalidad (`/api/auth`, `/api/user`, `/api/vault`)
- Asignar cada ruta HTTP a su Handler correspondiente
- Aplicar middleware JWT a las rutas protegidas
- Configurar la extracción del token JWT desde cookies HttpOnly

---

## Archivos

### `routes.go`

Expone una única función `SetupRoutes` que recibe la instancia de Echo y todos los handlers inicializados:

```go
func SetupRoutes(
    e *echo.Echo,
    authHandler *handlers.AuthHandler,
    vaultHandler *handlers.VaultHandler,
    userProfileHandler *handlers.UserProfileHandler,
)
```

---

## Estructura de Rutas

```
/api
├── /auth                    [PÚBLICO]
│   ├── POST /register       → AuthHandler.Register
│   ├── POST /login          → AuthHandler.Login
│   └── POST /logout         → AuthHandler.Logout
│
├── /user                    [🔐 PROTEGIDO]
│   ├── GET  /profile        → UserProfileHandler.GetProfile
│   └── PUT  /profile        → UserProfileHandler.UpdateProfile
│
└── /vault                   [🔐 PROTEGIDO]
    ├── POST   /items        → VaultHandler.CreateItem
    ├── GET    /items        → VaultHandler.GetItems
    ├── GET    /items/:id    → VaultHandler.GetItem
    ├── PUT    /items/:id    → VaultHandler.UpdateItem
    └── DELETE /items/:id    → VaultHandler.DeleteItem
```

---

## Middleware JWT

Las rutas protegidas usan `echo-jwt` configurado para leer el token desde una **cookie HttpOnly**:

```go
protected.Use(echojwt.WithConfig(echojwt.Config{
    SigningKey:  []byte(jwtSecret),
    TokenLookup: "cookie:auth_token",
}))
```

| Configuración | Valor | Descripción |
|---|---|---|
| `SigningKey` | `JWT_SECRET` (env) | Clave para verificar la firma del token |
| `TokenLookup` | `cookie:auth_token` | Lee el JWT desde la cookie `auth_token` |

> 🔒 El token **no** se envía en el header `Authorization`. Se usa una cookie HttpOnly para prevenir acceso desde JavaScript (protección contra XSS).

---

## Grupos de Rutas

| Grupo | Prefijo | Middleware | Descripción |
|---|---|---|---|
| `api` | `/api` | — | Grupo base de la API |
| `auth` | `/api/auth` | — | Autenticación (público) |
| `protected` | `/api` | JWT cookie | Rutas que requieren autenticación |
| `user` | `/api/user` | JWT cookie | Perfil de usuario |
| `vault` | `/api/vault` | JWT cookie | Bóveda de secretos |

---

## Tabla Completa de Endpoints

| Método | Ruta | Auth | Handler | Descripción |
|---|---|---|---|---|
| `POST` | `/api/auth/register` | ❌ | `AuthHandler.Register` | Registro de usuario |
| `POST` | `/api/auth/login` | ❌ | `AuthHandler.Login` | Login (establece cookie) |
| `POST` | `/api/auth/logout` | ❌ | `AuthHandler.Logout` | Logout (expira cookie) |
| `GET` | `/api/user/profile` | ✅ | `UserProfileHandler.GetProfile` | Obtener perfil |
| `PUT` | `/api/user/profile` | ✅ | `UserProfileHandler.UpdateProfile` | Actualizar perfil |
| `POST` | `/api/vault/items` | ✅ | `VaultHandler.CreateItem` | Crear item |
| `GET` | `/api/vault/items` | ✅ | `VaultHandler.GetItems` | Listar items |
| `GET` | `/api/vault/items/:id` | ✅ | `VaultHandler.GetItem` | Obtener item |
| `PUT` | `/api/vault/items/:id` | ✅ | `VaultHandler.UpdateItem` | Actualizar item |
| `DELETE` | `/api/vault/items/:id` | ✅ | `VaultHandler.DeleteItem` | Eliminar item |
