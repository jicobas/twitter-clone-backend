# Challenge Backend

Este desafío es un punto de partida para crear una plataforma similar a Twitter, pero aún es un proyecto considerable que requerirá tiempo y recursos significativos.

---

## 🎯 Objetivo
Crear una versión simplificada de una plataforma de microblogging similar a Twitter que permita a los usuarios **publicar**, **seguir** y **ver el timeline** de tweets.

---

## 📌 Requerimientos

### Tweets
- Los usuarios deben poder publicar mensajes cortos (tweets) que no excedan un límite de caracteres (por ejemplo, **280 caracteres**).

### Follow
- Los usuarios deben poder seguir a otros usuarios.

### Timeline
- Los usuarios deben poder ver una línea de tiempo que muestre los tweets de los usuarios a los que siguen.

---

## 📄 Assumptions
- Todos los usuarios son válidos, **no es necesario** crear un módulo de *signin* ni manejar sesiones.
- Se puede enviar el identificador de un usuario por **header**, **param**, **body** o cualquier método conveniente.
- Pensar una solución que pueda escalar a **millones de usuarios**.
- La aplicación debe estar **optimizada para lecturas**.

---

## 🧪 Criterios de evaluación
- Documentación **high level** de la arquitectura y componentes usados.
- Elección libre del lenguaje y tecnologías (ejemplos: serverless, Docker, Kubernetes, message brokers, queues, bases de datos, cache, load balancers, gateways, gRPC, websockets, etc.).
- No es necesario desarrollar un **front-end**.
- Se pueden agregar más *assumptions* en un archivo `business.txt`.
- Nos interesa la arquitectura interna y separación de capas:
  - Clean Architecture
  - DDD
  - Arquitectura Hexagonal
  - Ports and Adapters
  - Onion Architecture
  - MVC
- Se puede implementar una **DB in-memory**, pero debe especificarse en la documentación qué motor o tipo de DB se usaría y por qué.
- Valoramos el testing (no es necesario 100% coverage), priorizar casos de uso principales. Tests funcionales, de integración o aceptación son bienvenidos.

---

## 📦 Consideraciones de entrega
- Compartir el repositorio en acceso **público** para corrección o revisión.
- Proporcionar documentación clara para levantar el proyecto en un `README.md`.
- Se puede usar la **wiki** del repositorio para documentación extra.
- Se valora **hostear** o **dockerizar** la(s) aplicación(es).

---
