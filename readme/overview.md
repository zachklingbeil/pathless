# pathless

Zero-dependency viewport allocator. 

## Overview

**pathless** delivers rich, interactive experiences from a single domain (e.g., `timefactory.io`), without exposing internal content paths. Requests can only be served from `/`. All requests to paths/queries are automatically redirected back to `/`, ensuring the peer facing URL remains clean. Content is dynamically allocated to responsive panels, providing a distraction-free interface with no scrollbars. sit behind a reverse proxy, delivers a full 

- **Responsive panels** - Frames auto-scale with CSS `object-fit: contain`
- **Zero scrollbars** - Clean, distraction-free interface

frames are a finite pool of simulataneously observable html content, cached after first fetch. state is managed panel -> frame -> state. 