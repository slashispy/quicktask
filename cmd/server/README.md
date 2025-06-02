# cmd/server

Este directorio contiene el **punto de entrada principal** de la aplicación.

## Contenido

- `main.go`: Archivo que inicializa la configuración, conecta con los servicios externos (MongoDB, etc.), configura el router y arranca el servidor HTTP.

## Convención

En Go es común tener una carpeta `cmd/<nombre>` para cada binario o ejecutable del sistema. En este caso, `cmd/server` representa el ejecutable principal de la API.

## Ejemplo

```bash
go run cmd/server/main.go
