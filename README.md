# Setrya Space Management Service

Microservicio Go del bounded context **Space Management** de Sentrya. Este repositorio contiene únicamente este servicio; se creó desde el repositorio del microservicio, por lo que el código vive en la raíz y no dentro de una carpeta anidada `space-management/`.

Administra espacios de hogares: propietario, nombre, descripción, ubicación, dimensiones y disponibilidad. IoT Monitoring conserva la propiedad de dispositivos, lecturas y alertas.

## Stack y organización

- Go 1.27 y `net/http`, siguiendo el stack del servicio IoT Monitoring de referencia.
- MySQL 8 con `database/sql` y el driver MySQL del mismo servicio de referencia.
- Separación por capas: `internal/domain`, `internal/application`, `internal/interfaces/rest` e `internal/infrastructure/mysql`.
- OpenAPI 3 en `api/openapi.yaml`, servido en `/openapi.yaml`.
- Docker y Compose para ejecución local.

## Ejecutar

Requiere Docker con Compose.

```sh
docker compose up --build
```

- API: `http://localhost:8081`
- Salud: `GET http://localhost:8081/health`
- OpenAPI: `http://localhost:8081/openapi.yaml`
- MySQL expuesto en `localhost:3307` (ajustable con `SPACE_MYSQL_PORT`).

Los datos iniciales usan el usuario/clave `sentrya` para el entorno local. No reutilizar esas credenciales fuera del desarrollo.

## API

| Método | Ruta | Función |
|---|---|---|
| `POST` | `/api/v1/spaces` | Crea un espacio y genera su UUID. |
| `GET` | `/api/v1/spaces/{spaceId}` | Obtiene el espacio. IoT puede usar esta respuesta para confirmar existencia y estado. |
| `GET` | `/api/v1/owners/{ownerId}/spaces` | Lista espacios del propietario. |
| `PUT` | `/api/v1/spaces/{spaceId}` | Actualiza nombre, descripción, ubicación y dimensiones. |
| `PATCH` | `/api/v1/spaces/{spaceId}/availability` | Publica o pausa disponibilidad con `{ "active": true|false }`. |
| `GET` | `/api/v1/projects/{projectId}` | Alias de compatibilidad para Payment Management; consulta por UUID y responde 404 si no existe. |

Ejemplo para crear:

```json
{
  "ownerId": "UUID-del-usuario-IAM",
  "name": "Sala principal",
  "description": "Espacio monitoreado",
  "location": {
    "address": "Av. Ejemplo 123",
    "district": "Miraflores",
    "city": "Lima"
  },
  "dimensions": { "width": 4.5, "length": 5.0 }
}
```

El campo `ownerId` debe ser el identificador acordado con IAM. El ID del espacio lo genera el servicio.

## Integración entre contextos

- **Payment Management:** el cliente HTTP de referencia consulta `GET /api/v1/projects/{projectId}`. Se expone ese alias y usa UUID, como espera el contrato de pagos.
- **IoT Monitoring:** consulta `GET /api/v1/spaces/{spaceId}`. Un espacio inexistente responde 404; el JSON incluye propietario y disponibilidad.
- **IAM:** los ZIP proporcionados muestran JWT en los servicios consumidores, pero no definen un contrato de middleware compartido para este servicio. La validación del token y la derivación del propietario deben hacerse en el gateway o acordarse antes de exponer endpoints a clientes. El servicio no acepta un `ownerId` como prueba de autorización.
- **IoT Edge:** publica lecturas por MQTT hacia IoT Monitoring y no llama directamente a Space Management; no se agrega una dependencia innecesaria.

## Configuración

| Variable | Predeterminado | Uso |
|---|---|---|
| `HTTP_ADDR` | `:8081` | Dirección HTTP del servicio. |
| `MYSQL_DSN` | `sentrya:sentrya@tcp(localhost:3306)/space_management?parseTime=true` | Conexión MySQL. |
| `SPACE_HTTP_PORT` | `8081` | Puerto HTTP publicado por Compose. |
| `SPACE_MYSQL_PORT` | `3307` | Puerto MySQL publicado por Compose. |

`migrations/001_init.sql` crea la tabla al iniciar un volumen MySQL nuevo. Para una base existente, aplícala una sola vez antes de iniciar el servicio.
