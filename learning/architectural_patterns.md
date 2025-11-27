# Architectural Patterns in 2025

This document captures architectural patterns commonly used in backend development as of 2025, with specific context for the SplitFlap project.

---

## Table of Contents

1. [Current Architecture (Layered)](#1-layered-architecture-current)
2. [Hexagonal Architecture (Ports & Adapters)](#2-hexagonal-architecture-ports--adapters)
3. [Vertical Slice Architecture](#3-vertical-slice-architecture)
4. [CQRS (Command Query Responsibility Segregation)](#4-cqrs-command-query-responsibility-segregation)
5. [Functional Core, Imperative Shell](#5-functional-core-imperative-shell)
6. [Event-Driven Architecture](#6-event-driven-architecture)
7. [Recommendations for SplitFlap](#recommendations-for-splitflap)

---

## 1. Layered Architecture (Current)

**What we're using in SplitFlap Phase 1-3**

### Structure

```
Controller → Service → Repository → Database
```

### Code Example (SplitFlap)

```kotlin
// Controller Layer (HTTP concerns)
@RestController
@RequestMapping("/api/v1/displays")
class DisplayController(
    private val displayService: DisplayService
) {
    @GetMapping("/{id}")
    fun getDisplay(@PathVariable id: String): ResponseEntity<Display> {
        return displayService.getDisplay(id)
            ?.let { ResponseEntity.ok(it) }
            ?: ResponseEntity.notFound().build()
    }
}

// Service Layer (Business logic)
@Service
class DisplayService(
    private val displayRepository: DisplayRepository
) {
    fun getDisplay(id: String): Display? {
        return displayRepository.findById(id).orElse(null)
    }
}

// Repository Layer (Data access)
@Repository
interface DisplayRepository : JpaRepository<DisplayEntity, String>

// DTO Layer (API contract)
data class Display(
    val id: String,
    val content: DisplayContent,
    val config: DisplayConfig
)
```

### Pros

- ✅ Simple, well-understood pattern
- ✅ Good for CRUD APIs
- ✅ Easy to test each layer independently
- ✅ Spring Boot's default/recommended pattern
- ✅ Low learning curve for new developers

### Cons

- ❌ Can become "anemic" (controllers just pass-through to services)
- ❌ Business logic can leak across layers
- ❌ Doesn't scale well for complex domains
- ❌ Tight coupling between layers

### Best For

- CRUD APIs (like SplitFlap Phase 1-3)
- Small to medium projects
- Teams familiar with Spring Boot
- Traditional enterprise applications

---

## 2. Hexagonal Architecture (Ports & Adapters)

**Popular in 2025 for domain-driven applications**

### Concept

Isolate core business logic from external concerns (database, HTTP, messaging). The domain defines "ports" (interfaces), and adapters implement them.

### Structure

```
┌─────────────────────────────────────┐
│        Application Core             │
│  ┌──────────────────────────┐      │
│  │   Domain Logic           │      │
│  │   (Business Rules)       │      │
│  └──────────────────────────┘      │
│          ↕️ Ports                    │
└─────────────────────────────────────┘
     ↕️              ↕️              ↕️
┌─────────┐   ┌─────────┐   ┌─────────┐
│ REST    │   │Database │   │ Message │
│ Adapter │   │ Adapter │   │ Adapter │
└─────────┘   └─────────┘   └─────────┘
```

### Code Example

```kotlin
// Domain (core business logic - no framework dependencies)
package dev.clarakko.splitflap_api.domain

class Display(
    val id: String,
    val content: DisplayContent
) {
    fun validate(): Result<Unit> {
        if (content.rows.isEmpty()) {
            return Result.failure(Exception("Display must have content"))
        }
        return Result.success(Unit)
    }
}

// Port (interface - what the domain needs)
package dev.clarakko.splitflap_api.domain.ports

interface DisplayRepository {
    fun findById(id: String): Display?
    fun save(display: Display): Display
}

interface DisplayNotifier {
    fun notifyCreated(displayId: String)
}

// Use Case (application service)
package dev.clarakko.splitflap_api.application

class CreateDisplayUseCase(
    private val repository: DisplayRepository,
    private val notifier: DisplayNotifier
) {
    fun execute(content: DisplayContent): Display {
        val display = Display(id = UUID.randomUUID().toString(), content = content)
        display.validate().getOrThrow()
        
        val saved = repository.save(display)
        notifier.notifyCreated(saved.id)
        
        return saved
    }
}

// Adapter (implementation - database)
package dev.clarakko.splitflap_api.adapters.persistence

@Repository
class PostgresDisplayRepository(
    private val jpaRepository: JpaDisplayRepository
) : DisplayRepository {
    override fun findById(id: String): Display? {
        return jpaRepository.findById(id)
            .map { it.toDomain() }
            .orElse(null)
    }
    
    override fun save(display: Display): Display {
        val entity = DisplayEntity.fromDomain(display)
        return jpaRepository.save(entity).toDomain()
    }
}

// Adapter (implementation - REST)
package dev.clarakko.splitflap_api.adapters.rest

@RestController
@RequestMapping("/api/v1/displays")
class DisplayController(
    private val getDisplayUseCase: GetDisplayUseCase,
    private val createDisplayUseCase: CreateDisplayUseCase
) {
    @GetMapping("/{id}")
    fun getDisplay(@PathVariable id: String): ResponseEntity<DisplayDTO> {
        return getDisplayUseCase.execute(id)
            ?.let { ResponseEntity.ok(it.toDTO()) }
            ?: ResponseEntity.notFound().build()
    }
    
    @PostMapping
    fun createDisplay(@RequestBody request: CreateDisplayRequest): ResponseEntity<DisplayDTO> {
        val display = createDisplayUseCase.execute(request.content)
        return ResponseEntity.ok(display.toDTO())
    }
}
```

### Pros

- ✅ Domain logic isolated from frameworks (Spring, JPA, etc.)
- ✅ Easy to swap implementations (Postgres → MongoDB, REST → GraphQL)
- ✅ Testable without Spring Boot (pure Kotlin tests)
- ✅ Clear separation of concerns
- ✅ Frameworks become plugins, not the foundation

### Cons

- ❌ More boilerplate (interfaces everywhere)
- ❌ Overkill for simple CRUD operations
- ❌ Steeper learning curve
- ❌ More files to navigate

### Best For

- Complex domains with changing requirements
- Long-lived applications
- When database/framework might change
- Domain-driven design (DDD) projects

---

## 3. Vertical Slice Architecture

**Trending in 2025 - organize by feature, not layer**

### Concept

Instead of organizing by technical layer (controllers, services, repositories), organize by feature/use case. Each "slice" contains everything needed for one feature.

### Structure

```
splitflap-api/src/main/kotlin/
├── features/
│   ├── displays/
│   │   ├── GetDisplay.kt          # Everything for GET /displays/{id}
│   │   ├── CreateDisplay.kt       # Everything for POST /displays
│   │   ├── UpdateDisplay.kt       # Everything for PUT /displays/{id}
│   │   └── DeleteDisplay.kt       # Everything for DELETE /displays/{id}
│   └── users/
│       ├── RegisterUser.kt
│       ├── LoginUser.kt
│       └── GetUserProfile.kt
└── shared/
    ├── database/
    ├── security/
    └── validation/
```

### Code Example

```kotlin
// filepath: features/displays/GetDisplay.kt
package dev.clarakko.splitflap_api.features.displays

// Everything for this feature in one file (or folder)

// 1. Request/Response DTOs
data class GetDisplayResponse(
    val id: String,
    val content: DisplayContent,
    val config: DisplayConfig
)

// 2. Handler (business logic)
@Service
class GetDisplayHandler(
    private val repository: DisplayRepository
) {
    fun handle(id: String): GetDisplayResponse? {
        return repository.findById(id)
            ?.let { entity ->
                GetDisplayResponse(
                    id = entity.id,
                    content = parseContent(entity.contentJson),
                    config = DisplayConfig(entity.rowCount, entity.columnCount)
                )
            }
    }
}

// 3. Controller (HTTP endpoint)
@RestController
@RequestMapping("/api/v1/displays")
class GetDisplayController(
    private val handler: GetDisplayHandler
) {
    @GetMapping("/{id}")
    fun getDisplay(@PathVariable id: String): ResponseEntity<GetDisplayResponse> {
        return handler.handle(id)
            ?.let { ResponseEntity.ok(it) }
            ?: ResponseEntity.notFound().build()
    }
}

// 4. Tests (in same package)
@SpringBootTest
class GetDisplayTest {
    @Test
    fun `returns display when it exists`() { /* ... */ }
    
    @Test
    fun `returns 404 when display not found`() { /* ... */ }
}
```

### Alternative: Feature Folders

```kotlin
// filepath: features/displays/get/GetDisplayController.kt
// filepath: features/displays/get/GetDisplayHandler.kt
// filepath: features/displays/get/GetDisplayResponse.kt
// filepath: features/displays/get/GetDisplayTest.kt

// filepath: features/displays/create/CreateDisplayController.kt
// filepath: features/displays/create/CreateDisplayHandler.kt
// filepath: features/displays/create/CreateDisplayRequest.kt
// filepath: features/displays/create/CreateDisplayTest.kt
```

### Pros

- ✅ All related code in one place (easy to find)
- ✅ Low coupling between features
- ✅ Easy to delete entire features
- ✅ Microservices-ready (easy to extract a feature)
- ✅ Clear boundaries (no accidental dependencies between features)
- ✅ New developers can work on features without understanding entire codebase

### Cons

- ❌ Can duplicate code across slices (shared logic needs extraction)
- ❌ Less familiar to traditional Spring developers
- ❌ Harder to enforce cross-cutting concerns
- ❌ Database entities might need to be shared

### Best For

- Feature-rich applications
- Teams working on different features simultaneously
- Preparing for microservices extraction
- When features are truly independent

---

## 4. CQRS (Command Query Responsibility Segregation)

**Common in 2025 for read-heavy or event-sourced systems**

### Concept

Separate read (query) operations from write (command) operations. Can use different data models, databases, or even services for each.

### Structure

```
┌─────────────────────────────────────────────┐
│  Commands (Writes)                          │
│  ┌─────────────────────────────────────┐   │
│  │ CreateDisplayCommand                 │   │
│  │ UpdateDisplayCommand                 │   │
│  │ DeleteDisplayCommand                 │   │
│  └───────────────┬─────────────────────┘   │
│                  ↓                          │
│  ┌─────────────────────────────────────┐   │
│  │ Write Database (Normalized)          │   │
│  │ - Optimized for writes               │   │
│  └───────────────┬─────────────────────┘   │
│                  │ Events                   │
│                  ↓                          │
└─────────────────────────────────────────────┘
                   │
                   ↓
┌─────────────────────────────────────────────┐
│  Queries (Reads)                            │
│  ┌─────────────────────────────────────┐   │
│  │ GetDisplayQuery                      │   │
│  │ ListDisplaysQuery                    │   │
│  │ SearchDisplaysQuery                  │   │
│  └───────────────┬─────────────────────┘   │
│                  ↓                          │
│  ┌─────────────────────────────────────┐   │
│  │ Read Database (Denormalized)         │   │
│  │ - Optimized for queries              │   │
│  │ - Pre-joined data                    │   │
│  └─────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

### Code Example

```kotlin
// Commands (write operations)
sealed interface Command

data class CreateDisplayCommand(
    val content: DisplayContent,
    val config: DisplayConfig
) : Command

data class UpdateDisplayCommand(
    val id: String,
    val content: DisplayContent
) : Command

// Command Handlers
@Service
class CreateDisplayCommandHandler(
    private val writeRepository: DisplayWriteRepository,
    private val eventBus: EventBus
) {
    fun handle(command: CreateDisplayCommand): String {
        val display = Display(
            id = UUID.randomUUID().toString(),
            content = command.content,
            config = command.config
        )
        
        writeRepository.save(display)
        
        // Publish event for read model update
        eventBus.publish(DisplayCreatedEvent(
            displayId = display.id,
            content = display.content,
            timestamp = Instant.now()
        ))
        
        return display.id
    }
}

// Queries (read operations)
sealed interface Query

data class GetDisplayQuery(val id: String) : Query

data class ListDisplaysQuery(
    val page: Int,
    val pageSize: Int,
    val sortBy: String
) : Query

// Query Handlers
@Service
class GetDisplayQueryHandler(
    private val readRepository: DisplayReadRepository  // Different from write repo
) {
    fun handle(query: GetDisplayQuery): DisplayReadModel? {
        return readRepository.findById(query.id)
    }
}

// Read Model (optimized for queries - denormalized)
data class DisplayReadModel(
    val id: String,
    val content: DisplayContent,
    val config: DisplayConfig,
    val createdAt: Instant,
    val updatedAt: Instant,
    val ownerName: String,  // Pre-joined from users table
    val viewCount: Int      // Pre-calculated
)

// Event Handler (updates read model)
@Service
class DisplayReadModelUpdater(
    private val readRepository: DisplayReadRepository
) {
    @EventListener
    fun onDisplayCreated(event: DisplayCreatedEvent) {
        val readModel = DisplayReadModel(
            id = event.displayId,
            content = event.content,
            createdAt = event.timestamp,
            // ... other fields
        )
        readRepository.save(readModel)
    }
}

// Controller (delegates to command/query buses)
@RestController
@RequestMapping("/api/v1/displays")
class DisplayController(
    private val commandBus: CommandBus,
    private val queryBus: QueryBus
) {
    @PostMapping
    fun createDisplay(@RequestBody request: CreateDisplayRequest): ResponseEntity<CreateDisplayResponse> {
        val displayId = commandBus.execute(CreateDisplayCommand(
            content = request.content,
            config = request.config
        ))
        return ResponseEntity.ok(CreateDisplayResponse(displayId))
    }
    
    @GetMapping("/{id}")
    fun getDisplay(@PathVariable id: String): ResponseEntity<DisplayReadModel> {
        return queryBus.execute(GetDisplayQuery(id))
            ?.let { ResponseEntity.ok(it) }
            ?: ResponseEntity.notFound().build()
    }
    
    @GetMapping
    fun listDisplays(
        @RequestParam page: Int,
        @RequestParam pageSize: Int
    ): ResponseEntity<List<DisplayReadModel>> {
        val displays = queryBus.execute(ListDisplaysQuery(page, pageSize, "createdAt"))
        return ResponseEntity.ok(displays)
    }
}
```

### Pros

- ✅ Optimized read/write models separately
- ✅ Scales reads independently from writes (different databases)
- ✅ Clear separation of concerns
- ✅ Great for event sourcing
- ✅ Can handle high read loads (replicated read databases)

### Cons

- ❌ Significant complexity
- ❌ Eventual consistency challenges (read model lags behind writes)
- ❌ More infrastructure (multiple databases, event bus)
- ❌ Overkill for most applications
- ❌ Harder to debug

### Best For

- High-scale systems (Netflix, Amazon scale)
- Event-sourced applications
- Read-heavy workloads (10:1 read/write ratio)
- Systems requiring different query models for different use cases

---

## 5. Functional Core, Imperative Shell

**Growing in Kotlin/Scala communities (2025)**

### Concept

Separate pure functions (no side effects) from impure functions (I/O, database, HTTP). The "core" is pure business logic, the "shell" handles all side effects.

### Structure

```
┌─────────────────────────────────────────┐
│  Functional Core (Pure Functions)       │
│  - No side effects                      │
│  - No dependencies                      │
│  - Easy to test                         │
│  - Business logic only                  │
└─────────────────────────────────────────┘
                  ↑
                  │ Called by
                  │
┌─────────────────────────────────────────┐
│  Imperative Shell (Side Effects)        │
│  - Database calls                       │
│  - HTTP requests                        │
│  - File I/O                             │
│  - External APIs                        │
└─────────────────────────────────────────┘
```

### Code Example

```kotlin
// Pure functions (functional core - no side effects)
object DisplayCore {
    // Validation (pure)
    fun validateContent(content: DisplayContent): Result<DisplayContent> {
        if (content.rows.isEmpty()) {
            return Result.failure(IllegalArgumentException("Content cannot be empty"))
        }
        if (content.rows.size > 100) {
            return Result.failure(IllegalArgumentException("Too many rows (max 100)"))
        }
        return Result.success(content)
    }
    
    // Business logic (pure)
    fun createDisplay(id: String, content: DisplayContent): Display {
        return Display(
            id = id,
            content = content,
            config = DisplayConfig(
                rowCount = content.rows.size,
                columnCount = content.rows.firstOrNull()?.columns?.size ?: 0
            )
        )
    }
    
    // Calculation (pure)
    fun calculateDisplayDimensions(content: DisplayContent): Pair<Int, Int> {
        val rowCount = content.rows.size
        val columnCount = content.rows.maxOfOrNull { it.columns.size } ?: 0
        return Pair(rowCount, columnCount)
    }
    
    // Transformation (pure)
    fun transformToUppercase(content: DisplayContent): DisplayContent {
        return content.copy(
            rows = content.rows.map { row ->
                row.copy(
                    columns = row.columns.map { col ->
                        col.copy(text = col.text.uppercase())
                    }
                )
            }
        )
    }
}

// Impure shell (side effects - I/O, database, HTTP)
@Service
class DisplayService(
    private val repository: DisplayRepository,  // Impure: database
    private val eventPublisher: ApplicationEventPublisher  // Impure: events
) {
    fun getDisplay(id: String): Display? {
        // Side effect: database call
        return repository.findById(id).orElse(null)
    }
    
    fun createDisplay(content: DisplayContent): Display {
        // Pure logic (delegated to core)
        val validated = DisplayCore.validateContent(content).getOrThrow()
        val displayId = UUID.randomUUID().toString()
        val display = DisplayCore.createDisplay(displayId, validated)
        
        // Side effect: save to database
        val entity = DisplayEntity.fromDomain(display)
        repository.save(entity)
        
        // Side effect: publish event
        eventPublisher.publishEvent(DisplayCreatedEvent(display.id))
        
        return display
    }
    
    fun updateDisplayContent(id: String, newContent: DisplayContent): Display {
        // Side effect: fetch from database
        val existing = repository.findById(id).orElseThrow()
        
        // Pure logic
        val validated = DisplayCore.validateContent(newContent).getOrThrow()
        val updated = DisplayCore.createDisplay(existing.id, validated)
        
        // Side effect: save to database
        repository.save(DisplayEntity.fromDomain(updated))
        
        return updated
    }
}

// Tests (pure functions are trivial to test)
class DisplayCoreTest {
    @Test
    fun `validateContent rejects empty content`() {
        val content = DisplayContent(rows = emptyList())
        val result = DisplayCore.validateContent(content)
        
        assertTrue(result.isFailure)
        assertEquals("Content cannot be empty", result.exceptionOrNull()?.message)
    }
    
    @Test
    fun `createDisplay calculates config from content`() {
        val content = DisplayContent(rows = listOf(
            DisplayRow(columns = listOf(Column("A"), Column("B")))
        ))
        
        val display = DisplayCore.createDisplay("test-id", content)
        
        assertEquals(1, display.config.rowCount)
        assertEquals(2, display.config.columnCount)
    }
    
    @Test
    fun `transformToUppercase converts all text`() {
        val content = DisplayContent(rows = listOf(
            DisplayRow(columns = listOf(Column("hello"), Column("world")))
        ))
        
        val transformed = DisplayCore.transformToUppercase(content)
        
        assertEquals("HELLO", transformed.rows[0].columns[0].text)
        assertEquals("WORLD", transformed.rows[0].columns[1].text)
    }
}
```

### Pros

- ✅ Pure functions are extremely easy to test (no mocks, no setup)
- ✅ Business logic has no dependencies (framework-agnostic)
- ✅ Reasoning about code is simpler (no hidden side effects)
- ✅ Functions are reusable across contexts
- ✅ Refactoring is safer (pure functions can't break side effects)

### Cons

- ❌ Requires functional programming mindset
- ❌ Can feel awkward in Spring Boot (framework is imperative)
- ❌ Not all logic can be made pure
- ❌ Team needs to understand the pattern

### Best For

- Complex business rules that need extensive testing
- Fintech applications (calculations, validations)
- Scientific computing
- When business logic changes frequently

---

## 6. Event-Driven Architecture

**Very popular in 2025 for distributed systems and real-time applications**

### Concept

Services communicate through events rather than direct calls. When something happens (event), multiple services can react independently.

### Structure

```
┌──────────────────┐       Event        ┌──────────────────┐
│  Display Service │  ─────────────────> │   Event Bus      │
│  (Publisher)     │   DisplayCreated    │  (Kafka/Redis)   │
└──────────────────┘                     └─────────┬────────┘
                                                   │
                          ┌────────────────────────┼────────────────────┐
                          │                        │                    │
                          ▼                        ▼                    ▼
                 ┌─────────────────┐    ┌─────────────────┐  ┌─────────────────┐
                 │  Analytics      │    │  Notification   │  │  Search Index   │
                 │  Service        │    │  Service        │  │  Service        │
                 │  (Listener)     │    │  (Listener)     │  │  (Listener)     │
                 └─────────────────┘    └─────────────────┘  └─────────────────┘
```

### Code Example

```kotlin
// Domain Event
data class DisplayCreatedEvent(
    val displayId: String,
    val content: DisplayContent,
    val ownerId: String,
    val timestamp: Instant
)

data class DisplayUpdatedEvent(
    val displayId: String,
    val oldContent: DisplayContent,
    val newContent: DisplayContent,
    val timestamp: Instant
)

// Publisher (produces events)
@Service
class DisplayService(
    private val repository: DisplayRepository,
    private val eventPublisher: ApplicationEventPublisher  // Spring's event system
) {
    fun createDisplay(content: DisplayContent, ownerId: String): Display {
        val display = Display(
            id = UUID.randomUUID().toString(),
            content = content,
            config = DisplayConfig(/* ... */)
        )
        
        // Save to database
        repository.save(DisplayEntity.fromDomain(display))
        
        // Publish event (fire-and-forget)
        eventPublisher.publishEvent(
            DisplayCreatedEvent(
                displayId = display.id,
                content = display.content,
                ownerId = ownerId,
                timestamp = Instant.now()
            )
        )
        
        return display
    }
    
    fun updateDisplay(id: String, newContent: DisplayContent): Display {
        val existing = repository.findById(id).orElseThrow()
        val updated = existing.copy(content = newContent)
        
        repository.save(updated)
        
        eventPublisher.publishEvent(
            DisplayUpdatedEvent(
                displayId = id,
                oldContent = existing.content,
                newContent = newContent,
                timestamp = Instant.now()
            )
        )
        
        return updated.toDomain()
    }
}

// Listener 1: Analytics Service
@Service
class AnalyticsService {
    @EventListener
    @Async  // Process asynchronously
    fun onDisplayCreated(event: DisplayCreatedEvent) {
        // Track display creation in analytics dashboard
        analyticsRepository.recordEvent(
            type = "display_created",
            displayId = event.displayId,
            userId = event.ownerId,
            timestamp = event.timestamp
        )
        
        // Update user statistics
        userStatsRepository.incrementDisplayCount(event.ownerId)
    }
    
    @EventListener
    @Async
    fun onDisplayUpdated(event: DisplayUpdatedEvent) {
        // Track display updates
        analyticsRepository.recordEvent(
            type = "display_updated",
            displayId = event.displayId,
            timestamp = event.timestamp
        )
    }
}

// Listener 2: Notification Service
@Service
class NotificationService(
    private val emailService: EmailService
) {
    @EventListener
    @Async
    fun onDisplayCreated(event: DisplayCreatedEvent) {
        // Send confirmation email to display owner
        emailService.send(
            to = getUserEmail(event.ownerId),
            subject = "Display Created Successfully",
            body = "Your display ${event.displayId} has been created!"
        )
    }
}

// Listener 3: Search Index Service
@Service
class SearchIndexService(
    private val elasticsearchClient: ElasticsearchClient
) {
    @EventListener
    @Async
    fun onDisplayCreated(event: DisplayCreatedEvent) {
        // Index display for search
        elasticsearchClient.index(
            index = "displays",
            document = DisplaySearchDocument(
                id = event.displayId,
                content = event.content.toSearchableText(),
                ownerId = event.ownerId,
                createdAt = event.timestamp
            )
        )
    }
    
    @EventListener
    @Async
    fun onDisplayUpdated(event: DisplayUpdatedEvent) {
        // Update search index
        elasticsearchClient.update(
            index = "displays",
            id = event.displayId,
            document = mapOf("content" to event.newContent.toSearchableText())
        )
    }
}

// Listener 4: WebSocket Service (for Phase 5)
@Service
class WebSocketService(
    private val webSocketSessions: WebSocketSessionRegistry
) {
    @EventListener
    fun onDisplayUpdated(event: DisplayUpdatedEvent) {
        // Notify connected clients in real-time
        webSocketSessions.getSessions(event.displayId).forEach { session ->
            session.send(
                TextMessage(objectMapper.writeValueAsString(
                    DisplayUpdateMessage(
                        displayId = event.displayId,
                        content = event.newContent
                    )
                ))
            )
        }
    }
}

// Configuration (enable async processing)
@Configuration
@EnableAsync
class AsyncConfig {
    @Bean
    fun taskExecutor(): Executor {
        val executor = ThreadPoolTaskExecutor()
        executor.corePoolSize = 10
        executor.maxPoolSize = 20
        executor.queueCapacity = 500
        executor.setThreadNamePrefix("event-async-")
        executor.initialize()
        return executor
    }
}
```

### External Event Bus (Kafka Example)

```kotlin
// Using Kafka for distributed events
@Service
class DisplayService(
    private val repository: DisplayRepository,
    private val kafkaTemplate: KafkaTemplate<String, DisplayCreatedEvent>
) {
    fun createDisplay(content: DisplayContent): Display {
        val display = Display(/* ... */)
        repository.save(DisplayEntity.fromDomain(display))
        
        // Publish to Kafka topic
        kafkaTemplate.send(
            "display-events",
            display.id,
            DisplayCreatedEvent(/* ... */)
        )
        
        return display
    }
}

// Listener in different microservice
@Service
class AnalyticsMicroservice {
    @KafkaListener(topics = ["display-events"])
    fun consumeDisplayEvent(event: DisplayCreatedEvent) {
        // This runs in a completely separate service/container
        analyticsDatabase.recordEvent(event)
    }
}
```

### Pros

- ✅ Loose coupling between services (services don't know about each other)
- ✅ Easy to add new features (just add a new listener)
- ✅ Natural fit for microservices
- ✅ Scalable (listeners can run on different machines)
- ✅ Resilient (if one listener fails, others still work)
- ✅ Perfect for real-time updates (WebSocket notifications)

### Cons

- ❌ Debugging is harder (async flows, distributed traces)
- ❌ Eventual consistency (events processed asynchronously)
- ❌ Event versioning challenges (need to support old event formats)
- ❌ Ordering guarantees can be complex
- ❌ Need monitoring/observability tools

### Best For

- Microservices architectures
- Real-time systems (like SplitFlap Phase 5 - WebSockets)
- Systems with many loosely-coupled features
- When you need to react to changes in multiple ways

---

## Recommendations for SplitFlap

### Phase 1-3: Stick with Layered Architecture ✅

**Current architecture is perfect because:**
- Simple CRUD operations
- Small codebase (easy to navigate)
- Fast to implement
- Spring Boot's sweet spot
- Easy for new contributors

```kotlin
// Keep this simple structure
Controller → Service → Repository → Database
     ↓           ↓
    DTO        Domain
```

### Phase 5+: Consider Event-Driven for WebSockets

When implementing real-time updates (Phase 5):

```kotlin
@Service
class DisplayService(
    private val repository: DisplayRepository,
    private val eventPublisher: ApplicationEventPublisher
) {
    fun updateDisplay(id: String, content: DisplayContent): Display {
        val updated = /* ... save to DB ... */
        
        // Publish event for WebSocket notification
        eventPublisher.publishEvent(DisplayUpdatedEvent(id, content))
        
        return updated
    }
}

@Service
class WebSocketNotifier {
    @EventListener
    fun onDisplayUpdated(event: DisplayUpdatedEvent) {
        // Notify all connected WebSocket clients
        webSocketSessions.broadcast(event)
    }
}
```

**Benefits for Phase 5:**
- ✅ Clean separation (business logic doesn't know about WebSockets)
- ✅ Easy to add other listeners (analytics, logging, etc.)
- ✅ Testable (can test display updates without WebSocket code)

### Phase 8+: Consider Vertical Slices if Team Grows

If multiple developers work on different features:

```kotlin
// Each feature is self-contained
features/
  ├── displays/
  │   ├── GetDisplay.kt
  │   ├── CreateDisplay.kt
  │   └── UpdateDisplay.kt
  ├── users/
  │   ├── RegisterUser.kt
  │   └── LoginUser.kt
  └── analytics/
      └── GetDisplayStats.kt
```

**Benefits:**
- ✅ Multiple developers can work without conflicts
- ✅ Easy to understand feature scope
- ✅ Can extract features to microservices later

### Don't Over-Engineer Early

**Avoid these patterns in Phase 1-4:**
- ❌ Hexagonal Architecture (too much boilerplate for CRUD)
- ❌ CQRS (no high-scale read requirements)
- ❌ Functional Core (business logic is simple)

**Start simple, evolve when needed.** The current layered architecture will carry you through Phase 1-4 without issues.

---

## Pattern Selection Matrix

| Pattern | Complexity | Best For | SplitFlap Phase |
|---------|------------|----------|-----------------|
| **Layered** | Low | CRUD APIs | ✅ Phase 1-4 |
| **Hexagonal** | Medium | Domain-driven | ❌ Not needed |
| **Vertical Slice** | Medium | Feature teams | 🤔 Phase 8+ (if team grows) |
| **CQRS** | High | High scale | ❌ Not needed |
| **Functional Core** | Medium | Complex logic | ❌ Not needed |
| **Event-Driven** | Medium | Real-time, microservices | ✅ Phase 5+ (WebSockets) |

---

## Key Takeaways

1. **Layered architecture is not outdated** - Still the right choice for many applications
2. **Choose based on requirements** - Not trends or resume-driven development
3. **Start simple, evolve deliberately** - Don't over-engineer early phases
4. **Event-driven fits real-time well** - Natural for WebSocket updates (Phase 5)
5. **Vertical slices scale teams** - Consider if SplitFlap goes multi-developer

**For SplitFlap:** Stick with layered architecture through Phase 4, add event-driven patterns in Phase 5 for WebSocket support.

---

## Further Reading

- [Spring Boot Best Practices 2025](https://spring.io/guides)
- [Hexagonal Architecture (Alistair Cockburn)](https://alistair.cockburn.us/hexagonal-architecture/)
- [Vertical Slice Architecture (Jimmy Bogard)](https://www.jimmybogard.com/vertical-slice-architecture/)
- [CQRS Pattern (Martin Fowler)](https://martinfowler.com/bliki/CQRS.html)
- [Event-Driven Architecture Patterns](https://www.enterpriseintegrationpatterns.com/)