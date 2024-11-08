# Ejercicio Técnico - Proceso de Selección Ualá

Este repositorio contiene una implementación simplificada de una plataforma de microblogging similar a Twitter, desarrollada como parte del proceso de selección para Ualá. 
El objetivo principal es permitir a los usuarios publicar mensajes cortos (tweets), seguir a otros usuarios y visualizar un timeline con los tweets de los usuarios que siguen. 
Esta solución ha sido diseñada con un enfoque en la escalabilidad y la optimización para lecturas, con la finalidad de soportar millones de usuarios.

## 1. Objetivo

El objetivo del ejercicio es desarrollar una versión simplificada de una plataforma de microblogging que cumpla con las siguientes funcionalidades:

- **Tweets**: Permitir que los usuarios publiquen mensajes cortos de hasta 280 caracteres.
- **Follow**: Permitir que los usuarios sigan a otros usuarios.
- **Timeline**: Mostrar una línea de tiempo con los tweets de los usuarios que se están siguiendo.

Este proyecto utiliza `Golang` para el backend y está compuesto por tres microservicios principales: 
1. **User Service**: Maneja el registro y las relaciones entre usuarios (seguir/dejar de seguir).
2. **Tweet Service**: Gestiona la creación y eliminación de tweets.
3. **Timeline Service**: Proporciona la línea de tiempo de los tweets publicados por los usuarios.

La documentación detallada de la arquitectura y los componentes utilizados está disponible en la [wiki](https://github.com/DevOpsLP/microblogging-platform/wiki/Overview)  del repositorio.

## 2. Requerimientos del Proyecto

Para poder levantar el proyecto, se necesitan los siguientes requisitos:

- **Golang**: Es necesario tener `Golang` instalado para poder compilar y ejecutar el código del backend.
- **Docker**: Se utilizará `Docker` para contenerizar los servicios y simplificar su despliegue.
- **Docker Compose**: `Docker Compose` facilitará la gestión y la orquestación de los contenedores.

## 3. Guía para Levantar el Proyecto

### 3.1 Clonar el Repositorio
Primero, clone el repositorio desde GitHub utilizando el siguiente comando:

```sh
git clone https://github.com/DevOpsLP/microblogging-platform.git
cd microblogging-platform
```

### 3.2 Configurar Variables de Entorno
El proyecto requiere algunas variables de entorno para configurarse adecuadamente. Asegúrese de crear un archivo `.env` en la raíz del proyecto con las siguientes variables:
Por ejemplo:
```env
DB_HOST=localhost
DB_USER=devuser
DB_PASSWORD=devpassword
DB_NAME=userdb
DB_PORT=5432
USER_SERVICE_URL=http://user-service:8080
TWEET_SERVICE_URL=http://tweet-service:8081
```

> Para efectos de este proyecto, todas estan colocadas y expuestas en el archivo `docker-compose.yml` para un levantamiento mas sencillo de este proyecto

### 3.3 Levantar los Servicios con Docker Compose
El proyecto incluye un archivo `docker-compose.yml` que contiene la configuración para todos los microservicios necesarios (user-service, tweet-service y timeline-service), así como la base de datos.
EL proyecto esta pensado para ser compilado usando `FROM ubuntu:latest`, el cual facilita el despliegue de ser necesario

Para levantar todos los servicios, ejecute el siguiente comando:

```sh
docker-compose up --build
```
Este comando descargará las imágenes necesarias, construirá los contenedores y levantará los servicios especificados.


### 3.4 Acceso a los Endpoints
Una vez que los servicios estén ejecutándose, podrá acceder a los endpoints disponibles desde el navegador o utilizando una herramienta como `cURL` o `Postman`. A continuación se presenta un resumen de los endpoints principales:

- **User Service**: 
  - Registrar usuarios, seguir y dejar de seguir usuarios.
  - Ejemplo: `http://localhost:8080/register`

- **Tweet Service**:
  - Crear y eliminar tweets.
  - Ejemplo: `http://localhost:8081/tweets`

- **Timeline Service**:
  - Obtener el timeline de tweets.
  - Ejemplo: `http://localhost:8082/timeline`

## 4. Consideraciones de Arquitectura

La arquitectura de la plataforma está orientada a la escalabilidad y está dividida en múltiples microservicios para garantizar una buena separación de responsabilidades. Cada microservicio tiene su propia responsabilidad y comunica con los demás a través de peticiones HTTP.

- **User Service**: Maneja la gestión de usuarios, incluyendo la funcionalidad de seguir y dejar de seguir.
- **Tweet Service**: Se encarga de la creación, eliminación y almacenamiento de tweets.
- **Timeline Service**: Permite obtener un resumen consolidado de los tweets de los usuarios seguidos.

La arquitectura utilizada sigue el enfoque de `MICROSERVICIOS` para garantizar la modularidad y la facilidad de mantenimiento. La documentación detallada sobre cómo se dividen los servicios y los componentes está disponible en la [wiki del repositorio](https://github.com/DevOpsLP/microblogging-platform/wiki/Overview).

# 4.1 Consideraciones de base de datos

Se utilizó una configuración de base de datos para PostgreSQL, la cual es compatible con AWS RDS y facilita la migración de desarrollo a producción. Esto permite utilizar una base de datos que está enfocada en cost-on-demand, considerando el factor costo-beneficio y la simplicidad de poder mejorar aún más el servicio si se usara AWS AURORA en caso de necesitar aún más velocidad.
Hay que destacar que DynamoDB no es una base de datos recomendada para esta tarea por la complejidad que llega a tener las relaciones many-to-many que pueden tener los followers/following. Sin embargo, se puede utilizar por separado para llevar una relación one-to-many, pero, para simplicidad, todo está dentro de PostgreSQL.

## 5. Testing

El proyecto incluye pruebas unitarias y de integración para asegurar la calidad del código y la correcta implementación de los casos de uso principales. Para ejecutar las pruebas, puede utilizar los siguientes comandos, según los archivos específicos de prueba:

```sh
go test ./tweet-service/internal/infrastructure/api/tweet_handler_test.go
```
- **tweet_handler_test.go**: Verifica el flujo completo de creación, obtención y eliminación de tweets, incluyendo la integración con el `user-service` para validar a los usuarios. Incluye pruebas como crear un tweet, obtener tweets de un usuario específico, y eliminar un tweet.

```sh
go test ./tweet-service/internal/infrastructure/persistence/db_test.go
```
- **db_test.go (tweet-service)**: Prueba la conexión a la base de datos para el `tweet-service`, asegurando que se pueda conectar y realizar operaciones básicas como `ping` para verificar la disponibilidad.

```sh
go test ./user-service/internal/infrastructure/persistence/db_test.go
```
- **db_test.go (user-service)**: Similar a la prueba de `tweet-service`, verifica que el `user-service` puede conectarse correctamente a la base de datos y realizar operaciones básicas de verificación.

```sh
go test ./user-service/internal/infrastructure/api/user_handler_test.go
```
- **user_handler_test.go**: Prueba el flujo de seguir y dejar de seguir usuarios. Incluye verificaciones como que un usuario pueda seguir a otro, validar que la relación se haya establecido correctamente, y luego dejar de seguir para confirmar que la relación se elimina.

Ejemplo del flujo de prueba **TestUserFollowUnfollowFlow**:
- **Paso 1**: Usuario 1 sigue a Usuario 2 y se verifica la respuesta exitosa.
- **Paso 2**: Se verifica que Usuario 1 tiene a Usuario 2 en su lista de seguidos.
- **Paso 3**: Se verifica que Usuario 2 tiene a Usuario 1 como seguidor.
- **Paso 4**: Usuario 1 deja de seguir a Usuario 2 y se verifica la respuesta.
- **Paso 5**: Confirmar que la lista de seguidos de Usuario 1 está vacía.
- **Paso 6**: Confirmar que la lista de seguidores de Usuario 2 está vacía.

Ejemplo del flujo de prueba **TestTweetFlowWithHTTPUserRepo**:
- **Paso 1**: Crear un tweet para User1 y verificar que se crea correctamente.
- **Paso 2**: Crear un tweet para User2 y verificar que se crea correctamente.
- **Paso 3**: Obtener los tweets de User1 y validar que el contenido es correcto.
- **Paso 4**: Obtener todos los tweets y verificar que ambos tweets están presentes.
- **Paso 5**: Eliminar el tweet de User1 y confirmar que se elimina exitosamente.

Estas pruebas cubren los componentes más importantes del `tweet-service` y `user-service` para garantizar el correcto funcionamiento de las funcionalidades principales.

## 6. Conclusión

El `microblogging-platform` es una solución funcional que cumple con los requisitos definidos en el ejercicio técnico para Ualá, proporcionando una plataforma robusta para la publicación, seguimiento y visualización de tweets. La documentación más detallada y los diagramas de arquitectura están disponibles en la wiki del repositorio, donde se explica la estructura interna y los componentes utilizados.

El resto de la informacion esta condensada en la wiki de este repositorio, tal como se indicó en la tarea [wiki del repositorio](https://github.com/DevOpsLP/microblogging-platform/wiki/Overview) 

