# Setrya Space Management Service

Microservicio Go del bounded context **Space Management** de Sentrya. Administra espacios de hogares, disponibilidad, referencias a dispositivos vinculados y reseñas. IoT Monitoring conserva la propiedad operativa de los dispositivos, lecturas y alertas.

## Stack y organización

- Go 1.27, `net/http` y MySQL 8 con `database/sql`.
- Capas en `internal/domain`, `internal/application`, `internal/interfaces/rest` e `internal/infrastructure/mysql`.
- OpenAPI 3 en `api/openapi.yaml`, disponible en `/openapi.yaml`.
- Migración MySQL en `migrations/001_init.sql`.

## Ejecutar

Requiere Docker con Compose.

```sh
docker compose up --build
```

- API: `http://localhost:8081`
- Salud: `GET /health`
- Especificación: `GET /openapi.yaml`
- MySQL: `localhost:3307` (`SPACE_MYSQL_PORT` permite cambiarlo).

## API

| Método | Ruta | Función |
|---|---|---|
| `POST` | `/api/v1/spaces` | Crea un espacio y genera su UUID. |
| `GET` | `/api/v1/spaces/{spaceId}` | Obtiene detalle, propietario y estado del espacio. |
| `GET` | `/api/v1/owners/{ownerId}/spaces` | Lista espacios del propietario. |
| `PUT` | `/api/v1/spaces/{spaceId}` | Actualiza nombre, descripción, ubicación y dimensiones. |
| `DELETE` | `/api/v1/spaces/{spaceId}` | Da de baja lógicamente el espacio como no disponible. |
| `PATCH` | `/api/v1/spaces/{spaceId}/pause` | Pausa disponibilidad (ruta del informe). |
| `PATCH` | `/api/v1/spaces/{spaceId}/availability` | Publica/pausa usando `{ "active": true|false }`. |
| `POST` | `/api/v1/spaces/{spaceId}/devices` | Registra el vínculo con un `deviceId` existente. |
| `GET` | `/api/v1/spaces/{spaceId}/devices` | Lista vínculos del espacio. |
| `PATCH` | `/api/v1/spaces/{spaceId}/devices/{deviceId}/deactivate` | Desactiva el vínculo en Space Management. |
| `POST` | `/api/v1/spaces/{spaceId}/reviews` | Registra una reseña con rating de 1 a 5. Solo admite espacios activos. |
| `GET` | `/api/v1/spaces/{spaceId}/reviews` | Lista reseñas del espacio. |
| `GET` | `/api/v1/projects/{projectId}` | Alias compatible con el cliente de Payment Management para validar existencia. |

Ejemplo para crear un espacio:

```json
{
  "ownerId": "UUID-del-usuario-IAM",
  "name": "Sala principal",
  "description": "Espacio monitoreado",
  "location": { "address": "Av. Ejemplo 123", "district": "Miraflores", "city": "Lima" },
  "dimensions": { "width": 4.5, "length": 5.0 }
}
```

Ejemplo para vincular un dispositivo IoT ya registrado:

```json
{ "deviceId": "UUID-del-dispositivo", "deviceType": "temperature" }
```

Ejemplo de reseña:

```json
{ "authorId": "UUID-del-usuario-IAM", "rating": 5, "comment": "Instalación ordenada" }
```

## Límites entre contextos

- Payment Management consulta `/api/v1/projects/{projectId}`, que devuelve el espacio si existe y 404 si no.
- IoT Monitoring consulta `/api/v1/spaces/{spaceId}` y administra el registro/activación del dispositivo y sus datos de monitoreo. Este servicio solo guarda la asociación local al espacio; la integración para validar que `deviceId` exista en IoT requiere acordar el endpoint de servicio a servicio.
- IoT Edge publica lecturas por MQTT a IoT Monitoring; no necesita conectarse a Space Management.
- IAM debe validar JWT y proveer una identidad confiable al propietario/autor. El servicio aún no implementa autenticación directa: `ownerId` y `authorId` del JSON no son prueba de autorización. Mantener detrás del API Gateway hasta cerrar ese contrato.

## Configuración

| Variable | Predeterminado | Uso |
|---|---|---|
| `HTTP_ADDR` | `:8081` | Dirección HTTP del servicio. |
| `MYSQL_DSN` | `sentrya:sentrya@tcp(localhost:3306)/space_management?parseTime=true` | Conexión MySQL. |
| `SPACE_HTTP_PORT` | `8081` | Puerto HTTP publicado por Compose. |
| `SPACE_MYSQL_PORT` | `3307` | Puerto MySQL publicado por Compose. |

La migración se monta para ejecutarse al inicializar un volumen nuevo de MySQL. Si ya existe el volumen, aplica `migrations/001_init.sql` manualmente antes de usar las nuevas tablas.

