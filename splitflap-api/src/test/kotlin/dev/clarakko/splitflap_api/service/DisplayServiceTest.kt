package dev.clarakko.splitflap_api.service

import org.junit.jupiter.api.Test
import org.junit.jupiter.api.Assertions.*

/**
 * Unit tests for DisplayService
 * 
 * Tests business logic in isolation without Spring context.
 * Fast execution, no dependencies.
 */
class DisplayServiceTest {

    private val service = DisplayService()

    @Test
    fun `getDisplay returns demo display when id is demo`() {
        // When
        val result = service.getDisplay("demo")

        // Then
        assertNotNull(result)
        assertEquals("demo", result?.id)
        assertEquals(5, result?.config?.rowCount)
        assertEquals(4, result?.config?.columnCount)
        assertEquals(5, result?.content?.rows?.size)
    }

    @Test
    fun `getDisplay returns null when id does not exist`() {
        // When
        val result = service.getDisplay("nonexistent")

        // Then
        assertNull(result)
    }

    @Test
    fun `getDisplay returns null for empty string id`() {
        // When
        val result = service.getDisplay("")

        // Then
        assertNull(result)
    }

    @Test
    fun `demo display has correct content structure`() {
        // When
        val result = service.getDisplay("demo")

        // Then
        assertNotNull(result)
        val rows = result!!.content.rows
        
        // Verify all rows have same column count
        assertTrue(rows.all { it.size == 4 })
        
        // Verify header row
        assertEquals(listOf("TIME", "DESTINATION", "PLATFORM", "STATUS"), rows[0])
        
        // Verify at least one data row exists
        assertTrue(rows.size > 1)
    }

    @Test
    fun `demo display config matches content dimensions`() {
        // When
        val result = service.getDisplay("demo")

        // Then
        assertNotNull(result)
        assertEquals(result!!.content.rows.size, result.config.rowCount)
        assertEquals(result.content.rows[0].size, result.config.columnCount)
    }

    @Test
    fun `demo display contains only valid characters`() {
        // When
        val result = service.getDisplay("demo")

        // Then
        assertNotNull(result)
        val validCharPattern = Regex("^[A-Z0-9 .,:-]*$")
        
        result!!.content.rows.forEach { row ->
            row.forEach { cell ->
                assertTrue(
                    validCharPattern.matches(cell),
                    "Cell '$cell' contains invalid characters"
                )
            }
        }
    }
}
