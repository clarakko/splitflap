package dev.clarakko.splitflap_api.controller

import dev.clarakko.splitflap_api.dto.Display
import dev.clarakko.splitflap_api.service.DisplayService
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController

@RestController
@RequestMapping("/api/v1/displays")
class DisplayController(
    private val displayService: DisplayService
) {

    @GetMapping("/{id}")
    fun getDisplay(@PathVariable id: String): ResponseEntity<Any> {
        return displayService.getDisplay(id)
            ?.let { ResponseEntity.ok(it) }
            ?: ResponseEntity.notFound().build()
    }
}
