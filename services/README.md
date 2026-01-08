# Services Package

Este paquete encapsula la **Lógica de Negocio** de la aplicación. Es el núcleo funcional que orquesta las operaciones entre los Handlers y los Repositorios.

## Responsabilidades

1.  **Reglas de Negocio**: Aplica validaciones lógicas (ej: "¿El usuario tiene permiso para editar este item?", "El email ya está registrado?").
2.  **Transformación de Datos**: Convierte DTOs (entrada) a Modelos de Dominio (base de datos) y viceversa.
3.  **Coordinación**: Puede llamar a múltiples repositorios u otros servicios para completar una tarea.
4.  **Seguridad de Acceso**: Asegura que las operaciones se realicen sobre los recursos correctos (ej: filtrar siempre por `UserID`).

## Componentes

*   **`AuthService`**: Lógica de autenticación (hashing de contraseñas, generación de tokens JWT).
*   **`VaultService`**: Lógica de gestión de items (creación, actualización segura, borrado lógico).

## Independencia

Los servicios son independientes del transporte (HTTP). Podrían ser invocados por una CLI, gRPC o tareas programadas sin cambios en su código, ya que no dependen de `echo.Context`.
