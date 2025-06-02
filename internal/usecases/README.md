# internal/usecases

Contiene la **lógica de negocio** del sistema.

## Responsabilidad

- Orquestar las reglas de negocio.
- Validar entidades.
- Llamar a los repositorios (pero no conocer detalles de implementación).
- Mantenerse aislado del mundo exterior (no conoce HTTP ni Mongo).

## Ventaja

Esta capa puede ser fácilmente testeada con mocks y reutilizada en otros contextos (REST, gRPC, CLI, etc).

## Ejemplo

- `RegisterUser()`: valida datos, encripta password, guarda usuario.
- `CompleteTask()`: marca una tarea como completada.
