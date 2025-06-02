# 📁 `internal/handlers/README.md`

```markdown
# internal/handlers

Contiene los **handlers HTTP** de la API.

## Responsabilidad

- Decodificar la request HTTP (body, query params, headers).
- Validar los datos mínimos necesarios.
- Delegar la lógica de negocio a los `usecases`.
- Codificar la respuesta HTTP.

## Importante

Los handlers deben mantenerse lo más delgados posible. Toda la lógica debe estar en los `usecases`.

## Ejemplo

- `POST /users`: crea un nuevo usuario.
- `GET /tasks`: lista las tareas.
