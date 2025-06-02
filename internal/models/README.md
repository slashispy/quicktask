# internal/models

Contiene los **tipos de datos (structs)** usados en la aplicación.

## Responsabilidad

- Definir entidades como `User`, `Task`, etc.
- Incluir tags para serialización JSON o BSON (para MongoDB).
- Representar los datos que viajan entre capas.

## Ejemplo

```go
type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Name     string             `bson:"name"`
    Email    string             `bson:"email"`
    Password string             `bson:"password"`
}
