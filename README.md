# Pathless

A minimal, zero-dependency Go HTTP server that serves a dynamic, keyboard-driven frame viewer as a single pre-compressed HTML page.

## What It Does

Pathless is a lightweight web server that:

1. **Embeds an HTML template** (`pathless.html`) at compile time using Go's `//go:embed` directive
2. **Processes the template** with environment-configurable values (title and API URL)
3. **Minifies the HTML** by removing comments, whitespace, and newlines
4. **Compresses with gzip** for optimal transfer size
5. **Serves from memory** - the processed HTML is stored in RAM, not read from disk

All processing happens **once** during initialization, making subsequent requests extremely fast.

## The Frame Viewer Interface

The embedded HTML provides an interactive frame viewer with:

### Keyboard Controls

- **`1`** - Single panel layout
- **`2`** - Two panel layout (horizontal or vertical)
- **`3`** - Three panel layout (4 variants)
- **`4`** - Rotate frames between panels
- **`Tab`** - Cycle focus between panels
- **`q`** - Previous frame (in focused panel)
- **`e`** - Next frame (in focused panel)

### Layout System

- **Layout 0**: Single panel (1 variant)
- **Layout 1**: Two panels - horizontal or vertical split (2 variants)
- **Layout 2**: Three panels - multiple arrangements (4 variants)

Pressing the same layout key cycles through its variants (e.g., pressing `2` toggles between horizontal and vertical splits).

### Features

- **Frame caching** - Fetched frames are cached in memory
- **Dynamic content loading** - Frames are fetched from the API URL with `X-Frame` header
- **Script execution** - Frames can include `<script>` tags that execute with panel context
- **Responsive panels** - Frames auto-scale with CSS `object-fit: contain`
- **Zero scrollbars** - Clean, distraction-free interface

## Environment Variables

- `TITLE` - The page title (defaults to `"hello_universe"`)
- `API_URL` - The API endpoint URL (defaults to `"http://localhost:1002"`)
  - Frames are fetched via: `GET {API_URL}` with header `X-Frame: {index}`
  - Server should return `X-Index` header with total frame count

## Usage

### Local Development

```bash
# Run with defaults
go run main.go

# Or with custom configuration
TITLE="Frame Viewer" API_URL="https://api.example.com" go run main.go
```

The server listens on port `1001` and serves the page at `http://localhost:1001/`.

### Docker Deployment

#### Build the Image

```bash
# Build for current platform
docker build -t pathless .

# Build for specific platform (e.g., linux/amd64)
docker build --platform linux/amd64 -t pathless .

# Multi-platform build
docker buildx build --platform linux/amd64,linux/arm64 -t pathless .
```

#### Run the Container

```bash
# Run with defaults
docker run -p 1001:1001 pathless

# Run with environment variables
docker run -p 1001:1001 \
  -e TITLE="My Frame Viewer" \
  -e API_URL="https://api.example.com" \
  pathless

# Run with custom port mapping
docker run -p 8080:1001 pathless
```

#### Docker Features

- **Multi-stage build** - Minimal final image size (~5-10MB)
- **Scratch base** - No OS, only the binary
- **Non-root user** - Runs as UID 10001 for security
- **Static binary** - No dependencies, fully portable
- **Build caching** - Go module cache for faster rebuilds
- **Multi-platform** - Supports AMD64, ARM64, etc.

## API Contract

Your frame server (at `API_URL`) should:

1. Accept `X-Frame` header with the requested frame index (integer)
2. Return HTML content for that frame
3. Optionally return `X-Index` header with total number of frames
4. Return `200 OK` for valid frames

Example:
```
GET / HTTP/1.1
Host: api.example.com
X-Frame: 42

HTTP/1.1 200 OK
X-Index: 100
Content-Type: text/html

<img src="frame-42.jpg" />
```

## Routing Behavior

- **Only** serves requests to `/` with no query parameters
- All other paths/queries are redirected to `/` with a `301 Moved Permanently` status
- Responses include `Content-Encoding: gzip` header

## Architecture

```
┌─────────────┐
│  main.go    │ Embeds, minifies, gzips HTML at compile time
└─────────────┘
       │
       ▼
┌─────────────┐
│ pathless.   │ Frame viewer with keyboard controls
│   html      │ Fetches frames from API_URL
└─────────────┘
       │
       ▼
┌─────────────┐
│ Your Frame  │ Serves individual frames
│    API      │ (e.g., images, HTML snippets)
└─────────────┘
```

## Docker Compose Example

```yaml
version: '3.8'
services:
  pathless:
    build: .
    ports:
      - "1001:1001"
    environment:
      TITLE: "Production Viewer"
      API_URL: "http://frame-api:1002"
  
  frame-api:
    image: your-frame-server
    ports:
      - "1002:1002"
```

## Performance

- **Binary size**: ~5-10MB (static, no dependencies)
- **Memory footprint**: Minimal (HTML pre-compressed in RAM)
- **Response time**: Sub-millisecond (serves from memory)
- **Compression**: HTML is gzipped once at startup

---

A self-contained, keyboard-driven frame viewer that works with any frame-serving API backend.