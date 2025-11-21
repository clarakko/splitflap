# Kotlin/Spring Boot API Coding Standards

## General Principles

- Follow Kotlin idioms and conventions
- Prefer immutability (val over var)
- Use data classes for DTOs
- Keep functions small and focused
- Write self-documenting code with clear names

## Code Style

### Naming Conventions

- Classes: PascalCase (`DisplayController`, `DisplayService`)
- Functions: camelCase (`getDisplay`, `findById`)
- Constants: UPPER_SNAKE_CASE (`MAX_ROW_COUNT`)
- Packages: lowercase (`dev.clarakko.splitflap_api.controller`)

### File Organization

```kotlin
// 1. Package declaration
package dev.clarakko.splitflap_api.controller

// 2. Imports (organized by: stdlib, third-party, project)
import org.springframework.web.bind.annotation.*
import dev.clarakko.splitflap_api.service.DisplayService

// 3. Class declaration
@RestController
@RequestMapping("/api/v1/displays")
class DisplayController(
    private val displayService: DisplayService
) {
    // Implementation
}
```

## Spring Boot Patterns

### Controllers

- Use constructor injection (no @Autowired)
- Return `ResponseEntity<T>` for explicit status codes
- Use `@PathVariable`, `@RequestBody` annotations
- Keep controllers thin - delegate to services

```kotlin
@RestController
@RequestMapping("/api/v1/displays")
class DisplayController(private val service: DisplayService) {
    
    @GetMapping("/{id}")
    fun getDisplay(@PathVariable id: String): ResponseEntity<Display> {
        return service.getDisplay(id)
            ?.let { ResponseEntity.ok(it) }
            ?: ResponseEntity.notFound().build()
    }
}
```

### Services

- Annotate with `@Service`
- Contain business logic
- Return domain objects or nulls (not ResponseEntity)

```kotlin
@Service
class DisplayService {
    fun getDisplay(id: String): Display? {
        // Business logic here
    }
}
```

### DTOs

- Use data classes
- Keep them in dedicated `dto` package
- Match JSON structure exactly

```kotlin
data class Display(
    val id: String,
    val content: DisplayContent,
    val config: DisplayConfig
)
```

## Error Handling

### Phase 1

- Return `ResponseEntity.notFound()` for missing resources
- Let Spring handle 500 errors

### Future Phases

- Use `@ControllerAdvice` for global exception handling
- Create custom exception classes
- Return consistent error response format

## Testing

### Unit Tests

- Use JUnit 5
- Test file naming: `{ClassName}Test.kt`
- Use `@WebMvcTest` for controller tests
- Mock dependencies with Mockito

```kotlin
@WebMvcTest(DisplayController::class)
class DisplayControllerTest {
    
    @Autowired
    private lateinit var mockMvc: MockMvc
    
    @MockBean
    private lateinit var service: DisplayService
    
    @Test
    fun `GET existing display returns 200`() {
        // Arrange
        val display = Display(...)
        whenever(service.getDisplay("demo")).thenReturn(display)
        
        // Act & Assert
        mockMvc.get("/api/v1/displays/demo")
            .andExpect { status { isOk() } }
            .andExpect { jsonPath("$.id") { value("demo") } }
    }
}
```

## Configuration

### application.yaml

- Use YAML over properties
- Organize by Spring profile when needed
- Document non-obvious settings with comments

```yaml
spring:
  application:
    name: splitflap-api

server:
  port: 8080
```

### Config Classes

- Use `@Configuration` + `@Bean` for custom beans
- Implement `WebMvcConfigurer` for MVC customization

```kotlin
@Configuration
class CorsConfig : WebMvcConfigurer {
    override fun addCorsMappings(registry: CorsRegistry) {
        registry.addMapping("/api/**")
            .allowedOrigins("http://localhost:5173")
            .allowedMethods("GET")
    }
}
```

## Dependencies

### Phase 1 Only Use

- `spring-boot-starter-web`
- `jackson-module-kotlin`
- `spring-boot-starter-test`

### Don't Add Yet

- ❌ Database dependencies (Phase 2)
- ❌ WebSocket dependencies (Phase 5)
- ❌ Security dependencies (Phase 8)

## Code Comments

- Avoid obvious comments
- Document "why" not "what"
- Use KDoc for public APIs

```kotlin
// ❌ Bad: Obvious
// Get display by ID
fun getDisplay(id: String): Display?

// ✅ Good: Explains decision
// Phase 1: Hardcoded data. Phase 2 will query database.
fun getDisplay(id: String): Display? = demoDisplay.takeIf { it.id == id }
```

## Kotlin-Specific

### Null Safety

- Prefer non-nullable types
- Use `?.let {}` for null-safe operations
- Use `?:` (Elvis) for default values

```kotlin
// ✅ Good
service.getDisplay(id)
    ?.let { ResponseEntity.ok(it) }
    ?: ResponseEntity.notFound().build()

// ❌ Avoid
val display = service.getDisplay(id)
if (display != null) {
    return ResponseEntity.ok(display)
} else {
    return ResponseEntity.notFound().build()
}
```

### Data Classes

- Use for DTOs and value objects
- Automatically get `equals()`, `hashCode()`, `toString()`
- Use `copy()` for immutable updates

### Extension Functions

- Use sparingly
- Only when it truly extends a type's interface
- Don't pollute standard library types

## Git Commits

Use conventional commits:

```
feat(api): add GET /v1/displays/{id} endpoint
fix(api): return 404 for missing displays
test(api): add controller tests for display endpoint
refactor(api): extract display validation logic
docs(api): update API.md with error responses
```

## Phase Discipline

**DO:**
- ✅ Implement only Phase 1 features
- ✅ Add TODO comments referencing future phases
- ✅ Keep code simple and readable

**DON'T:**
- ❌ Add database code (Phase 2)
- ❌ Add POST/PUT/DELETE endpoints (Phase 2)
- ❌ Add authentication (Phase 8)
- ❌ Over-engineer for future needs

```kotlin
// ✅ Good: Simple Phase 1 implementation
private val demoDisplay = Display(
    id = "demo",
    content = DisplayContent(rows = listOf(...))
)

// ❌ Bad: Premature abstraction for Phase 2
interface DisplayRepository {
    fun findById(id: String): Display?
}
```
