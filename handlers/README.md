# Handlers Package

Este paquete contiene los **Controladores HTTP** (Handlers) de la aplicación. Su única responsabilidad es manejar la capa de transporte HTTP.

## Responsabilidades

1.  **Recibir Peticiones**: Escucha los endpoints definidos en `routes`.
2.  **Parsear Entrada**: Extrae parámetros de URL, query params y el cuerpo (Body) de la petición (JSON).
3.  **Validación**: Usa el paquete `dto` y `validator` para asegurar que los datos de entrada sean sintácticamente correctos.
4.  **Llamar al Servicio**: Invoca a la lógica de negocio (Service) correspondiente.
5.  **Formatear Respuesta**: Convierte los resultados (o errores) del servicio en respuestas HTTP estándar (JSON + Código de Estado).

## Componentes

*   **`AuthHandler`**: Maneja registro (`/register`) y login (`/login`).
*   **`VaultHandler`**: Maneja CRUD de items (`/vault/items`).
*   **`UserHandler`** (Futuro): Gestión de perfil de usuario.

## Ejemplo de Flujo (Create Item)

1.  `VaultHandler.CreateItem` recibe `POST /vault/items`.
2.  Extrae el JWT para obtener el `userID`.
3.  Parsea el JSON body a `dto.CreateVaultItemDTO`.
4.  Valida el DTO.
5.  Llama a `vaultService.CreateItem(userID, dto)`.
6.  Si hay éxito, devuelve `201 Created` con el item.
7.  Si hay error, devuelve `400 Bad Request` o `500 Internal Server Error`.
