# Contexto de ingeniería — quicktask

Documento de contexto para entender el porqué de las decisiones del proyecto, no solo
el qué. Complementa al `README.md` (que cubre el uso).

## Propósito

`quicktask` es una API de gestión de usuarios que sirve como **aplicación de demostración
de un circuito GitOps completo**. El foco del portfolio no es la riqueza funcional de la
API, sino cómo se estructura, prueba, empaqueta y despliega siguiendo buenas prácticas.

## Arquitectura y capas

Se sigue una variante pragmática de arquitectura limpia / hexagonal:

- **handlers**: traducen HTTP ↔ dominio. No contienen lógica de negocio.
- **usecases**: la lógica de dominio (`RegisterUser`). Dependen de **interfaces**
  (`UserRepository`, `PasswordEncrypter`), nunca de MongoDB o bcrypt directamente.
- **repositories**: contrato (`UserRepository`) + implementación MongoDB.
- **middleware**: autenticación JWT transversal.
- **models**: entidad de dominio (`User`) y DTO de entrada (`UserDTO`).

Regla de dependencia: el dominio no conoce la infraestructura. Esto hace la lógica
testeable con dobles de prueba y permite sustituir el almacenamiento sin tocar los casos
de uso.

## Decisiones de diseño

- **Inyección por interfaces.** `UserUseCase` recibe un `UserRepository` y un
  `PasswordEncrypter`. En tests se inyectan fakes; en producción, la implementación
  MongoDB y `BcryptEncrypter`.
- **`mongo.ErrNoDocuments` como "no encontrado".** `RegisterUser` trata ese error como
  "el usuario no existe" (caso válido), no como fallo. Hay un test unitario y uno de
  integración que fijan este contrato para que no se rompa en un refactor.
- **Health vs readiness separados.** `/health` (liveness) nunca depende de servicios
  externos: mientras el proceso sirva HTTP, responde 200. `/ready` (readiness) sí hace
  ping a MongoDB y responde 503 si no está disponible. Esta distinción es la correcta
  para Kubernetes: un fallo transitorio de la base no debe reiniciar el pod (liveness),
  pero sí debe sacarlo del balanceo (readiness).
- **Readiness desacoplado del driver.** El handler recibe una `ReadinessFunc`
  (`func(context.Context) error`) en lugar de un `*mongo.Client`, para no acoplar el
  paquete `handlers` a los tipos concretos del driver.
- **Config fail-fast.** Las variables obligatorias se validan al arrancar; si falta
  alguna, el proceso termina de inmediato con un mensaje claro en vez de fallar más tarde
  de forma difusa.
- **Imagen distroless.** La imagen final no tiene shell ni gestor de paquetes: menor
  superficie de ataque y menor tamaño. El binario se compila estático (`CGO_ENABLED=0`).

## Estrategia de testing

Dos niveles, separados por coste y dependencias:

1. **Unitarios** (`go test ./...`): rápidos, sin Docker. Cubren dominio, encriptación,
   config y handlers de salud. Usan dobles de prueba por interfaz.
2. **Integración con Docker** (`go test -tags=integration ...`): levantan MongoDB real
   con testcontainers-go y validan el repositorio contra una base auténtica, incluyendo
   el contrato de `ErrNoDocuments`. Aislados tras un build tag para no exigir Docker en
   el camino rápido ni en CI que no lo tenga.

Esta separación mantiene el feedback rápido por defecto y reserva el coste de levantar
contenedores para cuando se quiere verificar la integración real.

## Migración de Go

El proyecto se migró de **Go 1.23.4 a Go 1.25** (go.mod y Dockerfile) sin cambios de
código: el objetivo fue quedar en una versión soportada y moderna manteniendo la
funcionalidad. Verificado con build, vet y la suite de tests.

## Rol en el circuito GitOps

`quicktask` produce una **imagen de contenedor** que el resto del circuito consume:

- El `helm-chart-generic-app` la despliega como Deployment con probes a `/health` y
  `/ready`, y levanta un MongoDB en el cluster vía values.
- `gitops-app-values` fija los valores concretos (imagen `quicktask:dev`, env, Mongo).
- `gitops-argocd-apps` declara la Application/ApplicationSet que Argo CD reconcilia.
- `kind-gitops-utility-belt` construye la imagen, la carga en KIND y arranca Argo CD.

Por eso la app expone endpoints de salud y se configura 100% por variables de entorno:
son los dos requisitos que la hacen "cloud-native" y operables por el circuito.

## Deuda / próximos pasos

- Validación de `UserDTO` (hay tags `validate` en el DTO pero aún no se ejecuta un
  validador); sería el siguiente incremento natural.
- Endpoint de login que emita JWT (hoy el middleware valida tokens pero no hay emisor).
- Métricas/trazas para observabilidad en el cluster.
