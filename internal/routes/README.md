# 📦 Routes Package (`routes`)

Centraliza la definición y configuración de todas las rutas HTTP de la API, incluyendo la aplicación del middleware JWT para rutas protegidas.

---

## Responsabilidades

- Agrupar endpoints por funcionalidad (`/api/auth`, `/api/user`, `/api/vault`)
- Asignar cada ruta HTTP a su Handler correspondiente
- Aplicar middleware JWT a las rutas protegidas
- Configurar la extracción del token JWT desde cookies HttpOnly

---

## Estructura de Rutas

```
/api
├── /auth                    [PÚBLICO]
│   ├── POST /register
│   ├── POST /login
│   ├── POST /logout
│   ├── GET  /exists
│   ├── POST /2fa/verify
│   ├── POST /recovery/fetch
│   └── POST /recovery/reset
│
├── /user                    [🔐 PROTEGIDO]
│   ├── GET  /profile
│   ├── PUT  /profile
│   ├── POST /change-password
│   ├── POST /2fa/setup
│   ├── POST /2fa/enable
│   └── DELETE /account
│
└── /vault                   [🔐 PROTEGIDO]
    ├── GET    /items
    ├── POST   /items
    ├── GET    /items/:id
    ├── PUT    /items/:id
    ├── DELETE /items/:id
    ├── GET    /trash
    ├── POST   /items/:id/restore
    ├── DELETE /items/:id/permanent
    ├── GET    /folders
    ├── POST   /folders
    ├── PUT    /folders/:id
    └── DELETE /folders/:id
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

> 🔒 **Seguridad**: El token no es accesible desde JavaScript (protección contra XSS). Para el flujo de 2FA, se utiliza un token temporal enviado en el header `X-Temp-Token` durante la verificación final.

---

## Tabla Completa de Endpoints

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `POST` | `/api/auth/register` | ❌ | Registro de usuario |
| `POST` | `/api/auth/login` | ❌ | Login (establece cookie o requiere 2FA) |
| `POST` | `/api/auth/logout` | ❌ | Logout (expira cookie) |
| `GET` | `/api/auth/exists` | ❌ | Comprobar disponibilidad de email |
| `POST` | `/api/auth/2fa/verify` | ❌ | Verificar TOTP (requiere token temporal) |
| `POST` | `/api/auth/recovery/fetch` | ❌ | Obtener metadata de recuperación |
| `POST` | `/api/auth/recovery/reset` | ❌ | Resetear pass con Recovery Key |
| `GET` | `/api/user/profile` | ✅ | Obtener perfil |
| `PUT` | `/api/user/profile` | ✅ | Actualizar perfil |
| `POST` | `/api/user/change-password` | ✅ | Actualizar master password |
| `POST` | `/api/user/2fa/setup` | ✅ | Generar secreto TOTP |
| `POST` | `/api/user/2fa/enable` | ✅ | Activar 2FA |
| `DELETE` | `/api/user/account` | ✅ | Eliminar cuenta |
| `GET` | `/api/vault/items` | ✅ | Listar items activos |
| `POST` | `/api/vault/items` | ✅ | Crear item |
| `GET` | `/api/vault/items/:id` | ✅ | Obtener item (con secreto) |
| `PUT` | `/api/vault/items/:id` | ✅ | Actualizar item |
| `DELETE` | `/api/vault/items/:id` | ✅ | Mover a papelera |
| `GET` | `/api/vault/trash` | ✅ | Listar papelera |
| `POST` | `/api/vault/items/:id/restore`| ✅ | Restaurar item |
| `DELETE` | `/api/vault/items/:id/permanent`| ✅ | Eliminar permanentemente |
| `GET` | `/api/vault/folders` | ✅ | Listar carpetas |
| `POST` | `/api/vault/folders` | ✅ | Crear carpeta |
| `PUT` | `/api/vault/folders/:id` | ✅ | Renombrar carpeta |
| `DELETE` | `/api/vault/folders/:id` | ✅ | Eliminar carpeta |
