# pathless

Zero-dependency viewport allocator. 


## Installation


## Overview

**pathless** delivers rich, interactive experiences from a single domain (e.g., `timefactory.io`), without exposing internal content paths. Requests can only be served from `/`. All requests to paths/queries are automatically redirected back to `/`, ensuring the peer facing URL remains clean. Content is dynamically allocated to responsive panels, providing a distraction-free interface with no scrollbars. sit behind a reverse proxy, delivers a full 

- **Responsive panels** - Frames auto-scale with CSS `object-fit: contain`
- **Zero scrollbars** - Clean, distraction-free interface

frames are a finite pool of simulataneously observable html content, cached after first fetch. state is managed panel -> frame -> state. 

## Architecture

Pathless is a lightweight Go HTTP server that embeds a sophisticated client-side frame viewer. The server processes, minifies, and compresses an HTML template once at startup, then serves it from memory with maximum efficiency.

- All requests to paths/queries are automatically redirected back to `/`.

**Server responsibilities:**
- Embed and process HTML template at compile time
- Minify and gzip compress the viewer
- Serve the viewer with proper headers
- Redirect all non-root paths to `/`

**Client responsibilities:**
- Fetch frames from a configurable API endpoint
- Manage multi-panel layouts with keyboard navigation
- Cache frames and deduplicate requests
- Provide state management for loaded frames

## What It Does

Pathless is a lightweight web server that:

1. **Embeds an HTML template** (`pathless.html`) at compile time using Go's `//go:embed` directive
2. **Processes the template** with environment-configurable values (title and API URL)
3. **Minifies the HTML** by removing comments, whitespace, and newlines
4. **Compresses with gzip** for optimal transfer size
5. **Serves from memory** - the processed HTML is stored in RAM, not read from disk

All processing happens **once** during initialization, making subsequent requests extremely fast.

### Server (Go)

```
┌─────────────────────────────────────┐
│ Compile Time                        │
├─────────────────────────────────────┤
│ //go:embed pathless.html            │
│ Template stored in binary           │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│ Runtime Initialization              │
├─────────────────────────────────────┤
│ 1. Process template (TITLE, API_URL)│
│ 2. Minify HTML                      │
│ 3. Gzip compress                    │
│ 4. Store in memory (once)           │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│ Request Handling                    │
├─────────────────────────────────────┤
│ GET / → Serve compressed HTML       │
│ GET * → 301 redirect to /           │
└─────────────────────────────────────┘
```

### Client (JavaScript)

```
┌──────────┐    ┌──────────┐    ┌──────────┐
│   One    │───▶│    Fx    │───▶│   Zero   │
│Controller│    │  Layout  │    │  Cache   │
└──────────┘    └──────────┘    └──────────┘
     │               │                │
     │               │                │
  Keyboard      Panel State      HTTP Fetch
  Events        Management       & Caching
```
## Caching Strategy

**Frame Cache:**
- Keyed by frame index
- Stores `{ data, headers }`
- Never expires (memory-bounded by frame count)

**Request Deduplication:**
- Pending requests tracked in `Map`
- Multiple calls to same key return same Promise
- Cleaned up on completion

**Content-Type Handling:**
- `image/*` → Blob URL (memory-efficient)
- `application/json` → Parsed object
- Default → Raw text

## Performance Characteristics

**Server:**
- HTML processed once at startup
- Gzip compression applied once
- Zero disk I/O per request
- Sub-millisecond response times

**Client:**
- Frames cached indefinitely
- Request deduplication
- Efficient DOM updates via `requestAnimationFrame`
- Script re-execution on panel update

# production

## Environment Variables

| Variable  | Default                   | Description                     |
| --------- | ------------------------- | ------------------------------- |
| `TITLE`   | `"hello_universe"`        | Page title shown in browser tab |
| `API_URL` | `"http://localhost:1002"` | Base URL for frame API endpoint |

# development
