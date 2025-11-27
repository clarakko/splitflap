package dev.clarakko.splitflap_api.service

import dev.clarakko.splitflap_api.dto.Display
import dev.clarakko.splitflap_api.dto.DisplayConfig
import dev.clarakko.splitflap_api.dto.DisplayContent
import org.springframework.stereotype.Service

@Service
class DisplayService {

    // Phase 1: Hardcoded demo data. Phase 2 will query database.
    private val demoDisplay = Display(
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

    fun getDisplay(id: String): Display? {
        return demoDisplay.takeIf { it.id == id }
    }
}
