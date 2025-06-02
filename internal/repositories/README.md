# internal/repositories

Contiene las interfaces de persistencia (repositorios) y sus implementaciones concretas (MongoDB en este caso).

## Responsabilidad

- Abstraer la lógica de acceso a datos.
- Proveer implementaciones que se adapten a MongoDB u otras bases.

## Ventaja

Gracias a las interfaces, la capa de negocio no depende de Mongo directamente, y se facilita el testeo.

## Ejemplo

- `UserRepository`: interfaz con métodos como `Save()`, `FindByEmail()`.
- `MongoUserRepository`: implementación concreta usando `mongo-go-driver`.
