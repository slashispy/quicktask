# config

Contiene el código relacionado con la configuración de la aplicación.

## Funcionalidades

- Cargar valores desde `.env`, variables de entorno o archivos de configuración.
- Exponer una estructura `Config` con valores como:
  - Puerto del servidor
  - URI de MongoDB
  - Secretos (como JWT secret)

## Sugerencia

Se puede usar la librería [`spf13/viper`](https://github.com/spf13/viper) para mayor flexibilidad.

## Ejemplo

```go
type Config struct {
    MongoURI string
    Port     string
    JWTSecret string
}
