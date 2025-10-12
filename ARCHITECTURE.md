# Architecture Documentation

## System Overview

JioSaavn API is a REST API wrapper built with Go and Gin framework that provides programmatic access to JioSaavn's music catalog.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         CLIENT                               │
│  (Browser, Mobile App, CLI, Other Services)                  │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/HTTPS
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      MIDDLEWARE LAYER                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  CORS        │  │  Logger      │  │  Recovery    │      │
│  │  Middleware  │  │  Middleware  │  │  Middleware  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                       ROUTES LAYER                           │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Route Registration & HTTP Method Mapping            │   │
│  │  - /song/:id                                         │   │
│  │  - /search?q=query                                   │   │
│  │  - /download/:id                                     │   │
│  │  - /lyrics/:id                                       │   │
│  │  - /artist/:id                                       │   │
│  │  - /album/:id                                        │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      SERVICES LAYER                          │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Business Logic & API Handlers                      │    │
│  │  - GetSongHandler()                                 │    │
│  │  - SearchSongsHandler()                             │    │
│  │  - DownloadSongHandler()                            │    │
│  │  - GetLyricsHandler()                               │    │
│  │  - GetArtistHandler()                               │    │
│  │  - GetAlbumHandler()                                │    │
│  └─────────────────────────────────────────────────────┘    │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                       UTILS LAYER                            │
│  ┌────────────────┐  ┌────────────────┐  ┌──────────────┐  │
│  │  URL Decrypt   │  │  Format Song   │  │  Escape      │  │
│  │  (DES)         │  │  Data          │  │  String      │  │
│  └────────────────┘  └────────────────┘  └──────────────┘  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    EXTERNAL API                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  JioSaavn API (https://www.jiosaavn.com/api.php)    │   │
│  │  - song.getDetails                                   │   │
│  │  - search.getResults                                 │   │
│  │  - lyrics.getLyrics                                  │   │
│  │  - artist.getArtistPageDetails                       │   │
│  │  - content.getAlbumDetails                           │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Component Breakdown

### 1. Entry Point (main.go)

**Responsibilities:**
- Application initialization
- Configuration loading
- Middleware registration
- Route registration
- Server startup

**Key Functions:**
```go
func main()
```

### 2. Configuration Layer (config/)

**Responsibilities:**
- Environment variable management
- Default value provision
- Configuration validation

**Components:**
```go
type Config struct {
    ServerPort      string
    JioSaavnBaseURL string
    DecryptionKey   string
}
```

### 3. Middleware Layer (middleware/)

#### CORS Middleware
**Purpose:** Handle Cross-Origin Resource Sharing
**Headers Set:**
- Access-Control-Allow-Origin
- Access-Control-Allow-Methods
- Access-Control-Allow-Headers

#### Logger Middleware
**Purpose:** Request/Response logging
**Logs:**
- HTTP method
- Request path
- Client IP
- Status code
- Request duration

### 4. Routes Layer (routes/)

**Responsibilities:**
- Endpoint definition
- Method mapping
- Handler assignment

**Structure:**
```go
func RegisterRoutes(r *gin.Engine) {
    r.GET("/song/:id", services.GetSongHandler)
    r.GET("/search", services.SearchAllHandler)
    // ... more routes
}
```

### 5. Services Layer (services/)

**Responsibilities:**
- Business logic implementation
- API request handling
- Response formatting
- Error handling

**Key Handlers:**
- `GetSongHandler` - Fetch song details
- `SearchAllHandler` - Search across categories
- `DownloadSongHandler` - Generate download URLs
- `GetLyricsHandler` - Fetch lyrics
- `GetArtistHandler` - Artist information
- `GetAlbumHandler` - Album details

### 6. Utils Layer (utils/)

**Responsibilities:**
- Helper functions
- Data transformation
- Encryption/Decryption

**Key Functions:**
- `DecryptURL()` - Decrypt media URLs using DES
- `FormatSong()` - Format song data
- `EscapeString()` - URL encode strings

## Data Flow

### Request Flow Example: Get Song Details

```
1. Client Request
   └─> GET /song/abc123

2. Middleware Chain
   ├─> CORS Middleware (headers)
   └─> Logger Middleware (logging)

3. Route Matching
   └─> Match: /song/:id → GetSongHandler

4. Service Layer
   ├─> Extract ID from URL
   ├─> Build JioSaavn API URL
   ├─> Make HTTP GET request
   └─> Parse JSON response

5. Utils Layer
   ├─> DecryptURL (encrypted_media_url)
   ├─> FormatSong (enhance data)
   └─> Return formatted data

6. Response
   └─> JSON {success: true, data: [...]}

7. Logger Middleware
   └─> Log request details
```

## Error Handling Flow

```
┌─────────────┐
│   Request   │
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│  Validation     │───────┐
│  - Check params │       │ Error
│  - Verify input │       ▼
└──────┬──────────┘  ┌─────────────┐
       │             │ Return 400  │
       │ Valid       │ Bad Request │
       ▼             └─────────────┘
┌─────────────────┐
│  API Call       │───────┐
│  - HTTP Request │       │ Error
│  - Parse JSON   │       ▼
└──────┬──────────┘  ┌─────────────┐
       │             │ Return 500  │
       │ Success     │ Server Error│
       ▼             └─────────────┘
┌─────────────────┐
│  Data Process   │───────┐
│  - Decrypt      │       │ Error
│  - Format       │       ▼
└──────┬──────────┘  ┌─────────────┐
       │             │ Return 500  │
       │ Success     │ Server Error│
       ▼             └─────────────┘
┌─────────────────┐
│  Return 200 OK  │
│  with data      │
└─────────────────┘
```

## Security Considerations

### 1. Input Validation
- All user inputs are validated before processing
- Parameter existence checks
- Type validation

### 2. Error Messages
- Generic error messages to prevent information leakage
- Detailed errors logged server-side only

### 3. CORS Configuration
- Currently allows all origins (configurable)
- Can be restricted for production use

### 4. Rate Limiting (Future)
- Not implemented yet
- Should be added for production

### 5. Authentication (Future)
- No authentication currently
- API keys can be added if needed

## Scalability Considerations

### Current Architecture
- Stateless design (scales horizontally)
- No session management
- Pure request/response model

### Potential Improvements

1. **Caching Layer**
   ```
   Client → Load Balancer → API Server → Redis Cache → JioSaavn API
   ```

2. **Database Layer**
   ```
   Client → API Server → PostgreSQL → JioSaavn API
   ```

3. **Message Queue**
   ```
   Client → API Server → RabbitMQ → Workers → JioSaavn API
   ```

## Performance Considerations

### Current Implementation
- Synchronous API calls
- No caching
- Direct pass-through to JioSaavn

### Optimization Opportunities
1. **Response Caching**: Cache frequent queries
2. **Connection Pooling**: Reuse HTTP connections
3. **Batch Requests**: Combine multiple requests
4. **CDN Integration**: Cache static resources

## Monitoring & Observability

### Current Logging
- Request/response logging via middleware
- HTTP method, path, IP, status, duration

### Recommended Additions
1. **Metrics**: Prometheus metrics
2. **Tracing**: Distributed tracing (Jaeger/Zipkin)
3. **Alerting**: Error rate monitoring
4. **Health Checks**: Detailed health endpoints

## Deployment Architecture

### Development
```
Local Machine → Go Runtime → Port 8080
```

### Production (Recommended)
```
Internet → Load Balancer → Docker Container (API) → JioSaavn
                         ↓
                      Redis Cache
                         ↓
                    PostgreSQL (Logs)
```

## Technology Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.21+ |
| Framework | Gin Web Framework |
| HTTP Client | net/http (stdlib) |
| Encryption | crypto/des (stdlib) |
| Containerization | Docker |
| Orchestration | Docker Compose |
| Testing | Go testing (stdlib) |

## API Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response
```json
{
  "error": "Error message"
}
```

## Configuration Management

```
Environment Variables → Config Loader → Application Config
         ↓
    Default Values
```

**Priority:**
1. Environment variables (highest)
2. Default values (fallback)

## Testing Strategy

```
┌─────────────────┐
│  Unit Tests     │ ← Test individual functions
├─────────────────┤
│  Integration    │ ← Test component interaction
├─────────────────┤
│  E2E Tests      │ ← Test complete flows (Future)
└─────────────────┘
```

**Current Coverage:**
- Utils package: ✅
- Config package: ✅
- Services: ⏳ (To be added)
- Routes: ⏳ (To be added)

## Future Architecture Enhancements

1. **Microservices**: Split into smaller services
2. **GraphQL**: Add GraphQL endpoint
3. **WebSocket**: Real-time features
4. **Message Queue**: Async processing
5. **API Gateway**: Centralized entry point
6. **Service Mesh**: Inter-service communication

## Conclusion

The current architecture follows clean architecture principles with clear separation of concerns, making it maintainable, testable, and scalable. The modular design allows for easy extension and modification as requirements evolve.
