# pathless

Zero-dependency viewport allocator. 

## Overview

**pathless** delivers rich, interactive experiences from a single domain (e.g., `timefactory.io`), without exposing internal content paths. Requests can only be served from `/`. All requests to paths/queries are automatically redirected back to `/`, ensuring the peer facing URL remains clean. Content is dynamically allocated to responsive panels, providing a distraction-free interface with no scrollbars. sit behind a reverse proxy, delivers a full 

- **Responsive panels** - Frames auto-scale with CSS `object-fit: contain`
- **Zero scrollbars** - Clean, distraction-free interface

frames are a finite pool of simulataneously observable html content, cached after first fetch. state is managed panel -> frame -> state. 

## `window.pathless`

The `window.pathless` object provides the API facilitating interaction between `panels` (viewport) and `frames` (html/js/css).    

#### `pathless.context()`
Returns the DOM element of the focused panel, its current frame, and its panel specific state.

#### `pathless.update(key, value)`

Automatically detects the focused panel and current frame, uses key-value pair's to save and restore state across navigation and layout changes.

#### `pathless.fetch(url, opts)`
Fetches data from `url` (string) with standard fetch `opts` (object, optional). Returns a promise resolving to `{ data, headers }`. 

- `data`: Parsed response (blob URL for images, JSON for JSON responses, text otherwise)
- `headers`: Response headers object

Caching and request deduplication available via `opts.key`. 

#### `pathless.onKey(handler)`

Event handler used to register keybinds in frames, automatically scoped to the focused panel.