package dev.clarakko.splitflap_api.controller

import dev.clarakko.splitflap_api.dto.Display
import dev.clarakko.splitflap_api.dto.DisplayConfig
import dev.clarakko.splitflap_api.dto.DisplayContent
import dev.clarakko.splitflap_api.service.DisplayService
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.webmvc.test.autoconfigure.AutoConfigureMockMvc
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.context.bean.override.mockito.MockitoBean
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get
import org.springframework.test.web.servlet.result.MockMvcResultMatchers.*
import org.mockito.BDDMockito.given

/**
 * Integration tests for DisplayController
 * 
 * Tests controller with Spring MVC context.
 * Mocks the service layer to isolate controller logic.
 */
@SpringBootTest
@AutoConfigureMockMvc
class DisplayControllerIntegrationTest {

    @Autowired
    private lateinit var mockMvc: MockMvc

    @MockitoBean
    private lateinit var displayService: DisplayService

    @Test
    fun `GET displays by id returns 200 and display when found`() {
        // Given
        val mockDisplay = createMockDisplay()
        given(displayService.getDisplay("demo")).willReturn(mockDisplay)

        // When/Then
        mockMvc.perform(get("/api/v1/displays/demo"))
            .andExpect(status().isOk)
            .andExpect(content().contentType("application/json"))
            .andExpect(jsonPath("$.id").value("demo"))
            .andExpect(jsonPath("$.config.rowCount").value(5))
            .andExpect(jsonPath("$.config.columnCount").value(4))
            .andExpect(jsonPath("$.content.rows").isArray)
            .andExpect(jsonPath("$.content.rows.length()").value(5))
    }

    @Test
    fun `GET displays by id returns 404 when not found`() {
        // Given
        given(displayService.getDisplay("nonexistent")).willReturn(null)

        // When/Then
        mockMvc.perform(get("/api/v1/displays/nonexistent"))
            .andExpect(status().isNotFound)
    }

    @Test
    fun `GET displays endpoint has correct CORS headers`() {
        // Given
        val mockDisplay = createMockDisplay()
        given(displayService.getDisplay("demo")).willReturn(mockDisplay)

        // When/Then
        mockMvc.perform(
            get("/api/v1/displays/demo")
                .header("Origin", "http://localhost:5173")
        )
            .andExpect(status().isOk)
            .andExpect(header().exists("Access-Control-Allow-Origin"))
    }

    @Test
    fun `GET displays with special characters in id returns 404`() {
        // Given
        given(displayService.getDisplay("demo@123")).willReturn(null)

        // When/Then
        mockMvc.perform(get("/api/v1/displays/demo@123"))
            .andExpect(status().isNotFound)
    }

    private fun createMockDisplay() = Display(
        id = "demo",
        content = DisplayContent(
            rows = listOf(
                listOf("TIME", "DESTINATION", "PLATFORM", "STATUS"),
                listOf("10:30", "BOSTON", "3", "ON TIME"),
                listOf("10:45", "NEW YORK", "5", "DELAYED"),
                listOf("11:00", "PHILADELPHIA", "7", "BOARDING"),
                listOf("11:15", "WASHINGTON", "2", "ON TIME")
            )
        ),
        config = DisplayConfig(rowCount = 5, columnCount = 4)
    )
}
