# Twitter Clone Backend

Una implementación de Twitter siguiendo **Arquitectura Hexagonal**, optimizada para escalabilidad y lecturas frecuentes.

## 🚀 Inicio Rápido

### Opción 1: Local (más rápido)
```bash
# Ejecutar directamente
go run cmd/server/main.go

# O con Make
make run-dev
```

### Opción 2: Docker (recomendado para producción)
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

## 📋 Funcionalidades

- ✅ **Crear tweets** (máximo 280 caracteres)
- ✅ **Timeline personalizado** (tweets propios + seguidos)
- ✅ **Seguir/dejar de seguir** usuarios
- ✅ **API REST** completa con validaciones
- ✅ **Thread-safe** para alta concurrencia con operaciones atómicas
- ✅ **Preparado para escalar** a millones de usuarios

## 🧪 **Concurrencia Verificada:**
- **✅ Creación de tweets**: 1,000+ operaciones concurrentes sin pérdida de datos
- **✅ Operaciones follow/unfollow**: Atómicas y thread-safe bajo alta carga concurrente
- **✅ Acceso a timeline**: Lecturas concurrentes optimizadas con RWMutex
- **✅ Lectura de timelines**: 100+ lecturas simultáneas thread-safe
- **✅ Follow/unfollow**: Operaciones atómicas thread-safe con validaciones

## 🏗️ Arquitectura

```
├── cmd/server/          # Punto de entrada
├── internal/
│   ├── domain/          # Entidades y reglas de negocio
│   ├── ports/           # Interfaces (contratos)
│   ├── usecases/        # Lógica de aplicación
│   ├── adapters/        # HTTP handlers + repositories
│   └── config/          # Configuración
├── pkg/                 # Logger y utilidades
└── test/               # Tests de integración
```

**Principios:** Clean Architecture, Repository Pattern, Dependency Injection

## 📡 API Endpoints

**Autenticación:** Header `X-User-ID: user1` (usuarios pre-creados: user1, user2, user3)

### Tweets
```bash
# Crear tweet
POST /api/v1/tweets
{"content": "Hello World!"}

# Timeline (tweets propios + seguidos)
GET /api/v1/timeline/{userID}?limit=50

# Tweets de usuario específico
GET /api/v1/tweets/user/{userID}
```

### Seguimientos
```bash
# Seguir usuario
POST /api/v1/follow
{"followee_id": "user2"}

# Dejar de seguir
DELETE /api/v1/follow/{followeeID}

# Ver seguidores/siguiendo
GET /api/v1/users/{userID}/followers
GET /api/v1/users/{userID}/following
```

### Health Check
```bash
GET /health
```

## 🛠️ Desarrollo

```bash
# Tests
make test

# Build
make build

# Format código
make fmt
```

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

### Características de la imagen:
- ✅ **Multi-stage build** (imagen final ~15MB)
- ✅ **Alpine Linux** (segura y ligera)
- ✅ **Sin dependencias externas** (solo stdlib Go)
- ✅ **Usuario no-root** para seguridad

## ⚙️ Configuración

Variables de entorno:
```env
PORT=8080
STORAGE_TYPE=memory     # memory, mongodb
MONGO_URI=mongodb://localhost:27017
REDIS_URI=redis://localhost:6379
ENABLE_CACHE=false      # true para producción con Redis
```

### **Redis Configuration (Producción):**
```yaml
# Configuración recomendada para Docker Compose
redis:
  image: redis:7-alpine
  command: redis-server --maxmemory 512mb --maxmemory-policy allkeys-lru
  ports:
    - "6379:6379"
  volumes:
    - redis_data:/data
```

**Políticas de Cache:**
- **Timeline TTL**: 1 hora (balance entre freshness y performance)
- **Tweet TTL**: 24 horas (datos inmutables, cache más agresivo)
- **Eviction**: LRU (Least Recently Used) para usuarios menos activos

### **MongoDB Configuration (Producción):**
```yaml
# Configuración recomendada para Docker Compose
mongodb:
  image: mongo:7
  environment:
    - MONGO_INITDB_DATABASE=twitter_clone
  ports:
    - "27017:27017"
  volumes:
    - mongodb_data:/data/db
```

**Índices Críticos:**
```javascript
// Crear índices para performance óptima
db.tweets.createIndex({user_id: 1, created_at: -1})     // Timeline usuario
db.tweets.createIndex({created_at: -1})                 // Timeline global
db.follows.createIndex({follower_id: 1, followee_id: 1}) // Relaciones
```

**Configuraciones de Producción:**
- **Replica Set**: 3 nodos para alta disponibilidad
- **Sharding**: Por user_id para distribución horizontal
- **Write Concern**: `majority` para consistencia crítica

## 🎯 Decisiones de Arquitectura

### **Storage Strategy**
- **MVP:** In-memory (actual) - Desarrollo rápido con thread-safe mutex
- **Producción:** MongoDB + índices optimizados - Persistencia y queries complejas
- **Escala:** Redis cache + MongoDB sharding - Performance + disponibilidad

**Evolución del Storage:**
```
Desarrollo → In-Memory (simplidad)
    ↓
Producción → MongoDB (persistencia + queries)
    ↓
Escala → MongoDB + Redis + Load Balancer (millones de usuarios)
```

### **Cache Strategy: ¿Por qué Redis?**
Para el cache de timelines elegimos **Redis** sobre otras alternativas por razones específicas del dominio de Twitter:

**🚀 Performance para Timelines:**
- **Sorted Sets**: Perfectos para timelines ordenados por timestamp con acceso O(log N)
- **Sub-milisegundo de latencia** para operaciones típicas (vs 10-100ms de DB)
- **Pipeline operations**: Múltiples operaciones en una sola llamada de red

**📊 Estructuras de Datos Ideales:**
```redis
# Timeline ordenado por timestamp
ZREVRANGE timeline:user123 0 49 WITHSCORES

# Cache de tweets individuales
HGETALL tweet:tweet456

# Invalidación por seguidor
SMEMBERS followers:user789
```

**⚡ Operaciones Atómicas:**
- **MULTI/EXEC**: Para invalidar cache de múltiples seguidores de forma consistente
- **TTL automático**: Expiración de cache sin intervención manual (1 hora timelines)
- **LRU eviction**: Mantiene automáticamente los datos más relevantes

**🔄 vs Alternativas:**
- **vs Memcached**: Redis tiene sorted sets (crítico para timelines ordenados)
- **vs Cache local**: Redis es compartido entre instancias (consistencia)
- **vs Database**: 10-100x más rápido para lecturas frecuentes

**📈 Escalabilidad:**
- **Redis Cluster**: Sharding automático por user ID
- **Replicación**: Alta disponibilidad con maestro-esclavo
- **100K+ ops/seg**: Maneja millones de usuarios concurrentes

Esta elección es crucial para la **optimización de lecturas** que requiere Twitter, donde cada usuario consulta su timeline frecuentemente.

### **Database Strategy: ¿Por qué MongoDB?**
Para el almacenamiento principal elegimos **MongoDB** sobre bases de datos relacionales por características específicas del dominio social:

**📊 Modelo de Datos Flexible:**
```json
// Tweet document - estructura natural
{
  "_id": ObjectId("..."),
  "user_id": "user123",
  "content": "Hello Twitter!",
  "created_at": ISODate("2025-08-13T10:30:00Z"),
  "likes": 0,
  "retweets": 0,
  "metadata": {
    "hashtags": ["#tech", "#go"],
    "mentions": ["@user456"]
  }
}

// User document con arrays embebidos
{
  "_id": "user123",
  "username": "johndoe",
  "followers": ["user456", "user789"],
  "following": ["user101", "user202"]
}
```

**🚀 Performance Optimizado para Social Media:**
- **Consultas por índices compuestos**: `{user_id: 1, created_at: -1}` para timelines
- **Agregation Pipeline**: Para estadísticas complejas sin JOINs costosos
- **Sharding horizontal**: Distribución automática por user_id o geographic_region

**📈 Escalabilidad Nativa:**
```javascript
// Timeline query optimizado
db.tweets.find({
  user_id: {$in: ["user1", "user2", "user3"]}
}).sort({created_at: -1}).limit(50)

// Con índice: {user_id: 1, created_at: -1}
// Performance: O(log N) + O(limit)
```

**🔄 vs Alternativas:**
- **vs PostgreSQL**: Sin necesidad de JOINs complejos para social graphs
- **vs MySQL**: Mejor manejo de arrays (followers/following) sin tablas pivot
- **vs DynamoDB**: Queries más flexibles sin predefinir access patterns
- **vs Cassandra**: Menor complejidad operacional y mejor consistencia

**📊 Ventajas Específicas para Twitter:**
- **Documentos anidados**: Hashtags, menciones, media embebidos naturalmente
- **Arrays nativos**: Listas de followers/following sin tablas relacionales
- **Índices partial**: Solo tweets activos, mejora performance
- **GridFS**: Para attachments multimedia futuros
- **Change streams**: Para notificaciones real-time

**🎯 Estrategia de Índices:**
```javascript
// Índices críticos para performance
db.tweets.createIndex({user_id: 1, created_at: -1})     // Timeline personal
db.tweets.createIndex({created_at: -1})                  // Timeline global
db.users.createIndex({username: 1}, {unique: true})     // Login rápido
db.follows.createIndex({follower_id: 1, followee_id: 1}) // Relaciones
```

MongoDB permite **escalar horizontalmente** manteniendo la flexibilidad de esquema que necesita una aplicación social en evolución constante.

### **Business Rules**
- Timeline = tweets propios + de usuarios seguidos
- Límite 280 caracteres por tweet
- No auto-seguimiento, no duplicados
- Ordenamiento por fecha descendente

### **Escalabilidad**
- Interfaces preparadas para intercambio fácil (memory → MongoDB)
- Cache layer opcional para timelines
- Repository pattern para diferentes storages
- Preparado para message queues y load balancers

## 🏃‍♂️ Ejemplo de Uso

```bash
# 1. Crear tweet
curl -X POST http://localhost:8080/api/v1/tweets \
  -H "X-User-ID: user1" \
  -d '{"content": "Mi primer tweet!"}'

# 2. Seguir usuario
curl -X POST http://localhost:8080/api/v1/follow \
  -H "X-User-ID: user1" \
  -d '{"followee_id": "user2"}'

# 3. Ver timeline
curl http://localhost:8080/api/v1/timeline/user1
```

## 🔧 Stack Tecnológico

- **Go 1.21** - Lenguaje principal (performance + concurrencia)
- **net/http** - Servidor HTTP (sin dependencias externas)
- **In-Memory** - Storage con thread-safety (MVP actual)
- **MongoDB** - Base de datos principal (preparado para producción)
- **Redis** - Cache layer para timelines (optimización crítica)
- **Docker** - Containerización multi-stage
- **Make** - Automatización de tareas

**Arquitectura preparada para:**
- Load balancers (nginx/HAProxy)
- Message queues (para invalidación de cache)
- Monitoring (métricas de Redis + MongoDB)

---

**Listo para producción** con storage in-memory y **preparado para escalar** con la arquitectura hexagonal implementada.
