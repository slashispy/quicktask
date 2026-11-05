# quicktask

[![CI](https://github.com/slashispy/quicktask/actions/workflows/ci.yml/badge.svg)](https://github.com/slashispy/quicktask/actions/workflows/ci.yml)

API REST en Go para gestión de usuarios, construida con una arquitectura limpia
(handlers → casos de uso → repositorios) y pensada para desplegarse en Kubernetes
mediante un circuito GitOps con Argo CD.

Este repositorio es la **aplicación**. El despliegue (Helm chart, values de GitOps,
Argo CD y el tooling de KIND) vive en repos separados —ver
[El circuito GitOps](#el-circuito-gitops).

## Stack

- **Go 1.25**
- [chi](https://github.com/go-chi/chi) como router HTTP
- **MongoDB** como almacenamiento (driver oficial)
- **JWT** (`golang-jwt/v5`) para autenticación
- **bcrypt** para el hashing de contraseñas
- Imagen final [distroless](https://github.com/GoogleContainerTools/distroless)

## Arquitectura

```
cmd/server            punto de entrada: carga config, conecta Mongo, monta el router
config                carga de configuración desde variables de entorno
internal/
  handlers            adaptadores HTTP (register user, health, readiness)
  middleware          autenticación JWT
  models              entidades (User) y DTOs (UserDTO)
  repositories        interfaz UserRepository + implementación MongoDB
  usecases            lógica de dominio (RegisterUser) y encriptación (bcrypt)
  database            conexión a MongoDB
```

La lógica de dominio (`usecases`) depende de **interfaces** (`UserRepository`,
`PasswordEncrypter`), no de implementaciones concretas, lo que permite testearla con
dobles de prueba y cambiar el almacenamiento sin tocar el dominio.

## Endpoints

| Método | Ruta         | Auth | Descripción                                  |
|--------|--------------|------|----------------------------------------------|
| GET    | `/health`    | No   | Liveness. Responde `200 {"status":"ok"}`.    |
| GET    | `/ready`     | No   | Readiness. Hace ping a MongoDB; `503` si falla. |
| POST   | `/v1/users`  | JWT  | Registra un usuario.                         |

`/health` y `/ready` existen para los probes de Kubernetes: liveness nunca depende de
dependencias externas; readiness sí verifica MongoDB.

## Configuración (variables de entorno)

| Variable        | Requerida | Default | Descripción                      |
|-----------------|-----------|---------|----------------------------------|
| `MONGO_URI`     | Sí        | —       | URI de conexión a MongoDB.       |
| `MONGO_DB_NAME` | Sí        | —       | Nombre de la base de datos.      |
| `JWT_SECRET`    | Sí        | —       | Secreto para validar los JWT.    |
| `PORT`          | No        | `8080`  | Puerto HTTP del servidor.        |

Las variables requeridas se validan al arrancar: si falta alguna, el proceso termina.

## Correr en local

### Con Go y un MongoDB local

```bash
export MONGO_URI="mongodb://localhost:27017"
export MONGO_DB_NAME="quicktask"
export JWT_SECRET="dev-secret"
go run ./cmd/server
```

### Con Docker

```bash
docker build -t quicktask:dev .
docker run --rm -p 8080:8080 \
  -e MONGO_URI="mongodb://host.docker.internal:27017" \
  -e MONGO_DB_NAME="quicktask" \
  -e JWT_SECRET="dev-secret" \
  quicktask:dev
```

## Tests

### Unitarios (sin dependencias externas)

```bash
go test ./...
```

Cubren la lógica de dominio (`RegisterUser` con dobles de prueba), el encriptador
bcrypt, la carga de configuración y los handlers de health/readiness.

### Integración con Docker

Un test de integración levanta un **MongoDB real en un contenedor** con
[testcontainers-go](https://golang.testcontainers.org/) y ejercita el repositorio
contra esa instancia. Requiere un demonio Docker en ejecución y está aislado tras el
build tag `integration`:

```bash
go test -tags=integration ./internal/repositories/...
```

## Build y verificación

```bash
go build ./...
go vet ./...
go test ./...
```

## Integración continua

El workflow de GitHub Actions ([`.github/workflows/ci.yml`](.github/workflows/ci.yml))
se ejecuta de forma manual (`workflow_dispatch`, desde la pestaña *Actions*):

- **Build & Test**: verifica formato (`gofmt`), `go vet`, compila, corre los tests
  unitarios con `-race` y coverage, y los tests de integración con Docker
  (`-tags=integration`, MongoDB vía testcontainers-go).
- **Docker image build**: construye la imagen del `Dockerfile` (sin publicarla) con
  cache de capas, validando que el empaquetado funciona.

El pipeline representa el primer tramo del circuito; la publicación de la imagen a un
registry y el despliegue vía Argo CD se documentan en `kind-gitops-utility-belt`.

## El circuito GitOps

`quicktask` es el primer componente de un circuito GitOps reproducible en local:

| Repo | Rol |
|------|-----|
| **quicktask** (este) | La aplicación. |
| `helm-chart-generic-app` | Chart Helm genérico reutilizable para desplegar apps. |
| `gitops-app-values` | Values por app que parametrizan el chart genérico. |
| `gitops-argocd-apps` | App of Apps + ApplicationSet de Argo CD. |
| `kind-gitops-utility-belt` | Scripts y tooling para levantar todo en KIND. |

Ver [`context.md`](context.md) para el detalle de arquitectura y decisiones de diseño.
