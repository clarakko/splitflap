package dev.clarakko.splitflap_api.contract

import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.webmvc.test.autoconfigure.AutoConfigureMockMvc
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get
import org.springframework.test.web.servlet.result.MockMvcResultMatchers.*
import org.hamcrest.Matchers

/**
 * Contract tests for Display API
 * 
 * Validates that the API implementation matches the specification in docs/API.md.
 * Tests the full stack with real service implementation.
 * 
 * Contract verification:
 * - Response structure matches JSON schema
 * - HTTP status codes match specification
 * - Content-Type headers are correct
 * - Field types and names are exact
 */
@SpringBootTest
@AutoConfigureMockMvc
class DisplayApiContractTest {

    @Autowired
    private lateinit var mockMvc: MockMvc

    @Test
    fun `GET displays demo returns exact schema from API spec`() {
        mockMvc.perform(get("/api/v1/displays/demo"))
            .andExpect(status().isOk)
            .andExpect(content().contentType("application/json"))
            
            // Verify top-level structure
            .andExpect(jsonPath("$.id").exists())
            .andExpect(jsonPath("$.content").exists())
            .andExpect(jsonPath("$.config").exists())
            
            // Verify id field
            .andExpect(jsonPath("$.id").isString)
            .andExpect(jsonPath("$.id").value("demo"))
            
            // Verify content structure
            .andExpect(jsonPath("$.content.rows").isArray)
            .andExpect(jsonPath("$.content.rows[0]").isArray)
            .andExpect(jsonPath("$.content.rows[0][0]").isString)
            
            // Verify config structure
            .andExpect(jsonPath("$.config.rowCount").isNumber)
            .andExpect(jsonPath("$.config.columnCount").isNumber)
            .andExpect(jsonPath("$.config.rowCount").value(5))
            .andExpect(jsonPath("$.config.columnCount").value(4))
            
            // Verify no extra fields at root level
            .andExpect(jsonPath("$.*").value(Matchers.hasSize<Any>(3)))
    }

    @Test
    fun `GET displays with nonexistent id returns 404 as per spec`() {
        mockMvc.perform(get("/api/v1/displays/nonexistent"))
            .andExpect(status().isNotFound)
            .andExpect(content().string("")) // Spring default for 404 with no body
    }

    @Test
    fun `GET displays response has correct Content-Type header`() {
        mockMvc.perform(get("/api/v1/displays/demo"))
            .andExpect(status().isOk)
            .andExpect(header().string("Content-Type", "application/json"))
    }

    @Test
    fun `demo display has exactly 5 rows as per API spec`() {
        mockMvc.perform(get("/api/v1/displays/demo"))
            .andExpect(status().isOk)
            .andExpect(jsonPath("$.content.rows.length()").value(5))
    }

    @Test
    fun `demo display has exactly 4 columns per row as per API spec`() {
        mockMvc.perform(get("/api/v1/displays/demo"))
            .andExpect(status().isOk)
            .andExpect(jsonPath("$.content.rows[0].length()").value(4))
            .andExpect(jsonPath("$.content.rows[1].length()").value(4))
            .andExpect(jsonPath("$.content.rows[2].length()").value(4))
            .andExpect(jsonPath("$.content.rows[3].length()").value(4))
            .andExpect(jsonPath("$.content.rows[4].length()").value(4))
    }

    @Test
    fun `demo display header row matches API spec example`() {
        mockMvc.perform(get("/api/v1/displays/demo"))
            .andExpect(status().isOk)
            .andExpect(jsonPath("$.content.rows[0][0]").value("TIME"))
            .andExpect(jsonPath("$.content.rows[0][1]").value("DESTINATION"))
            .andExpect(jsonPath("$.content.rows[0][2]").value("PLATFORM"))
            .andExpect(jsonPath("$.content.rows[0][3]").value("STATUS"))
    }

    @Test
    fun `config dimensions match actual content dimensions`() {
        mockMvc.perform(get("/api/v1/displays/demo"))
            .andExpect(status().isOk)
            // This validates the contract that config.rowCount matches actual row count
            .andExpect(jsonPath("$.config.rowCount").value(5))
            .andExpect(jsonPath("$.content.rows.length()").value(5))
            // And columnCount matches actual column count
            .andExpect(jsonPath("$.config.columnCount").value(4))
            .andExpect(jsonPath("$.content.rows[0].length()").value(4))
    }

    @Test
    fun `GET displays endpoint follows versioned API pattern`() {
        // Verify the endpoint follows /api/v1/... pattern as specified
        mockMvc.perform(get("/api/v1/displays/demo"))
            .andExpect(status().isOk)
        
        // Verify non-versioned path doesn't work
        mockMvc.perform(get("/displays/demo"))
            .andExpect(status().isNotFound)
    }

    @Test
    fun `CORS headers allow specified origins for GET requests`() {
        // Test Vite default port (5173)
        mockMvc.perform(
            get("/api/v1/displays/demo")
                .header("Origin", "http://localhost:5173")
        )
            .andExpect(status().isOk)
            .andExpect(header().string("Access-Control-Allow-Origin", "http://localhost:5173"))

        // Test alternative dev port (3000)
        mockMvc.perform(
            get("/api/v1/displays/demo")
                .header("Origin", "http://localhost:3000")
        )
            .andExpect(status().isOk)
            .andExpect(header().string("Access-Control-Allow-Origin", "http://localhost:3000"))
    }
}
