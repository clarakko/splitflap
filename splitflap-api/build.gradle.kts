plugins {
	kotlin("jvm") version "2.2.21"
	kotlin("plugin.spring") version "2.2.21"
	id("org.springframework.boot") version "4.0.0"
	id("io.spring.dependency-management") version "1.1.7"
	// kotlin("plugin.jpa") version "2.2.21" // Phase 2: Enable when adding database
}

group = "dev.clarakko"
version = "0.0.1-SNAPSHOT"
description = "Split-flap display engine API"

java {
	toolchain {
		languageVersion = JavaLanguageVersion.of(21)
	}
}

repositories {
	mavenCentral()
}

dependencies {
	// implementation("org.springframework.boot:spring-boot-starter-data-jpa") // Phase 2: Enable when adding database
	implementation("org.springframework.boot:spring-boot-starter-webmvc")
	implementation("com.fasterxml.jackson.module:jackson-module-kotlin")
	implementation("org.jetbrains.kotlin:kotlin-reflect")
	developmentOnly("org.springframework.boot:spring-boot-devtools")
	// runtimeOnly("com.h2database:h2") // Phase 2: Enable when adding database
	
	// Test dependencies
	testImplementation("org.springframework.boot:spring-boot-starter-test")
	testImplementation("org.springframework.boot:spring-boot-starter-webmvc-test")
	testImplementation("org.jetbrains.kotlin:kotlin-test-junit5")
	testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

kotlin {
	compilerOptions {
		freeCompilerArgs.addAll("-Xjsr305=strict", "-Xannotation-default-target=param-property")
	}
}

// Phase 2: Enable when adding database
// allOpen {
// 	annotation("jakarta.persistence.Entity")
// 	annotation("jakarta.persistence.MappedSuperclass")
// 	annotation("jakarta.persistence.Embeddable")
// }

tasks.withType<Test> {
	useJUnitPlatform()
}
