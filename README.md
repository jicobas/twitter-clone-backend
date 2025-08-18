# Twitter Clone Backend

Una implementación de Twitter siguiendo **Arquitectura Hexagonal**, optimizada para escalabilidad y lecturas frecuentes.

## 🏗️ Arquitectura y Diseño
- ✅ **Clean Architecture** con separación clara de responsabilidades
- ✅ **Optimizado para lecturas** con estrategia de cache inteligente

### Cualidades del diseño
- ✅ **API REST**
- ✅ **Thread-safe**

### Diseño de packages

```
├── cmd/server/          # Punto de entrada
├── internal/
│   ├── domain/          # Entidades y reglas de negocio
│   ├── ports/           # Interfaces
│   ├── usecases/        # Lógica de aplicación
│   ├── adapters/        # HTTP handlers + repositories
│   └── config/          # Configuración
├── pkg/                 # Logger y utilidades
└── test/               # Tests de integración
```

**Principios aplicados:** Clean Architecture, Repository Pattern, Dependency Injection

## 🎯 Decisiones de Arquitectura

### **Storage Strategy**
- **MVP:** In-memory con thread-safety
- **Escala:** Redis cache + MongoDB sharding para performance y disponibilidad

**Evolución del Storage:**
```
Desarrollo → In-Memory
    ↓
Escala → MongoDB + Redis + Load Balancer
```

### **Cache Strategy: ¿Por qué Redis?**
Para el cache de timelines elegimos **Redis** sobre otras alternativas por razones específicas del dominio de Twitter:

**⚡ Operaciones Atómicas:**
- **MULTI/EXEC**: Para invalidar cache de múltiples seguidores de forma consistente
- **TTL automático**: Expiración de cache sin intervención manual (1 hora timelines)
- **LRU eviction**: Mantiene automáticamente los datos más relevantes

**📈 Escalabilidad:**
- **Redis Cluster**: Sharding automático por user ID
- **Replicación**: Alta disponibilidad con maestro-esclavo

Esta elección es crucial para la **optimización de lecturas** que requiere Twitter, donde cada usuario consulta su timeline frecuentemente. Para entender cuales son las configuraciones adecuadas debemos realizar una prueba de carga para entender efectivamente cuantos recursos se necesitan. Tenemos vencimiento por TTL y por actualizaciones.

### **Database Strategy: ¿Por qué MongoDB?**
Para el almacenamiento principal elegimos **MongoDB** sobre bases de datos relacionales por características específicas del dominio social:

**📊 Ventajas sobre SQL para Social Media:**
- **Modelo de datos natural**: Los tweets y usuarios se mapean directamente a documentos JSON
- **Sin JOINs complejos**: Las relaciones followers/following se almacenan como arrays nativos
- **Esquema evolutivo**: Agregar nuevos campos (hashtags, menciones, media) sin migraciones
- **Sharding automático**: Distribución horizontal por user_id para millones de usuarios
- **Consultas optimizadas**: Índices compuestos para timelines ordenados por fecha

**� Performance Critical para Twitter:**
- **Escrituras masivas**: Millones de tweets simultáneos sin bloqueos de tabla
- **Lecturas frecuentes**: Timeline queries optimizadas con índices específicos
- **Escalabilidad real**: Soporta crecimiento exponencial sin refactoring de base

**⚡ Índices Estratégicos:**
MongoDB permite crear índices compuestos específicos para cada patrón de acceso (timeline personal, global, por usuario) sin las limitaciones de las claves foráneas relacionales.

### **Business Rules**
- Timeline = tweets propios + de usuarios seguidos
- Límite 280 caracteres por tweet
- No auto-seguimiento, no duplicados
- Ordenamiento por fecha descendente

### **Escalabilidad**
- Interfaces preparadas para intercambio fácil (memory → MongoDB)
- Cache layer opcional para timelines
- Repository pattern para diferentes storages

## 📋 Funcionalidades

- ✅ **Crear tweets** (máximo 280 caracteres)
- ✅ **Timeline personalizado** (tweets propios + seguidos)
- ✅ **Seguir/dejar de seguir** usuarios

## 🔧 Stack Tecnológico

- **Go 1.21** - Lenguaje principal (performance + concurrencia)
- **net/http** - Servidor HTTP (sin dependencias externas)
- **In-Memory** - Storage con thread-safety (MVP)
- **MongoDB** - Base de datos principal
- **Redis** - Cache layer para timelines
- **Docker** - Containerización multi-stage
- **Make** - Automatización de tareas

## 🏃‍♂️ Ejemplo de Uso

```bash
# 1. Crear tweet
curl -X POST http://localhost:8080/tweets \
  -H "X-User-ID: user1" \
  -H "Content-Type: application/json" \
  -d '{"content": "Mi primer tweet!"}'

# 2. Seguir usuario
curl -X POST http://localhost:8080/users/following \
  -H "X-User-ID: user1" \
  -H "Content-Type: application/json" \
  -d '{"followee_id": "user2"}'

# 3. Ver timeline
curl http://localhost:8080/users/user1/timeline

# 4. Ver tweets de un usuario
curl http://localhost:8080/users/user2/tweets

# 5. Ver seguidores
curl http://localhost:8080/users/user2/followers
```

## 🚀 Inicio Rápido

### Opción 1: Local
```bash
# Ejecutar directamente
go run cmd/server/main.go

# O con Make
make run-dev
```

### Opción 2: Docker
```bash
# Construir imagen
docker build -t twitter-clone .

# Ejecutar contenedor
docker run -d -p 8080:8080 --name twitter-api twitter-clone

# Ver logs (opcional)
docker logs twitter-api

# Detener cuando termines
docker stop twitter-api && docker rm twitter-api
```

### Opción 3: Docker Compose (con MongoDB + Redis)
```bash
# Requiere docker-compose instalado
docker-compose up --build
```

El servidor inicia en `http://localhost:8080`

### ✅ Verificar que funciona:
```bash
curl http://localhost:8080/health
# Respuesta: {"status": "ok"}
```

## 📡 API Endpoints

**Autenticación:** Header `X-User-ID: user1` (usuarios pre-creados: user1, user2, user3)

### Tweets
```bash
# Crear tweet
POST /tweets
{"content": "Hello World!"}

# Timeline
GET /users/{userID}/timeline?limit=50

# Tweets de usuario específico
GET /users/{userID}/tweets
```

### Seguimientos
```bash
# Seguir usuario
POST /users/following
{"followee_id": "user2"}

# Dejar de seguir
DELETE /users/following/{followeeID}

# Ver seguidores/siguiendo
GET /users/{userID}/followers
GET /users/{userID}/following
```

### Health Check
```bash
GET /health
```

## ⚙️ Configuración

Variables de entorno:
```env
PORT=8080
STORAGE_TYPE=memory     # memory, mongodb
MONGO_URI=mongodb://localhost:27017
REDIS_URI=redis://localhost:6379
ENABLE_CACHE=false
```

### **Configuración Redis:**

**Políticas de Cache:**
- **Timeline TTL**: 1 hora
- **Tweet TTL**: 24 horas
- **Eviction**: LRU (Least Recently Used) para usuarios menos activos

### **Configuración MongoDB:**

**Índices Críticos:**
```javascript
// Crear índices para performance óptima
db.tweets.createIndex({user_id: 1, created_at: -1})     // Timeline usuario
db.tweets.createIndex({created_at: -1})                 // Timeline global
db.follows.createIndex({follower_id: 1, followee_id: 1}) // Relaciones
```

**Configuraciones:**
- **Replica Set**: Nodos para alta disponibilidad
- **Sharding**: Por user_id para distribución horizontal
- **Write Concern**: `majority` para consistencia crítica

## 🐳 Docker

### Comandos útiles:
```bash
# Construir imagen
docker build -t twitter-clone .

# Ejecutar en background
docker run -d -p 8080:8080 --name twitter-api twitter-clone

# Ver logs en tiempo real
docker logs -f twitter-api

# Entrar al contenedor (debug)
docker exec -it twitter-api sh

# Detener y limpiar
docker stop twitter-api && docker rm twitter-api

# Ver imagen creada
docker images twitter-clone
```
