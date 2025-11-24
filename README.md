# SplitFlap

A split-flap display engine. Build configurable, animated split-flap displays and embed them anywhere.

## What is this?

Split-flap displays (also called Solari boards) are those mechanical displays you see in train stations and airports where characters flip to reveal new information. This project is an engine for creating and embedding virtual split-flap displays with real-time data.

**Build** a display in the web app. **Embed** it on any site with a single line of code.

```html
<split-flap display-id="abc123"></split-flap>
<script src="https://splitflap.app/embed.js"></script>
```

## Why?

Employer-mandated GenAI/LLM/Agentic training. Now that global scale corporate plagiarism is possible, it is now a mandatory skill to acquire.Thank YOU humans for all the hard work that was stolen by our corporate overlords to make this silly application possible.

## Project Structure

```

splitflap/
|__ splitflap-api/    # Kotlin + Spring Boot backend
|__ splitflap-web/    # React builder app
|__ splitflap-embed/  # Web component for embedding displays
|__ docs/             # Architecture, API specs, data model
```

## Getting Started

### Prerequisites

- Jave 21+
- Node.js 24+

### Run the API

```sh
cd splitflap-api
./gradlew bootRun
```

### Run the Builder Web App

```sh
cd splitflap-web
npm install
npm run dev
```

### Build the embed component

```sh
cd splitflap-embed
npm install
npm run build
```

## Documentation

- [Architecture](docs/ARCHITECTURE.md)
- [API Specification](docs/API.md)
- [Data Model](docs/DATA_MODEL.md)
- [Roadmap](ROADMAP.md)
- [Current Status](STATUS.md)

## Author

[Clara Brown](https://github.com/clarakko)

## License

[MIT](LICENSE)