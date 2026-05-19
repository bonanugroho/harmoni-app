---
layout: doc
date: 2026-05-19
---

# Technology Stack

## Core language & runtime
- **Language**: JavaScript (CommonJS modules)
- **Runtime**: Node.js (inferred from `package.json` type)

## Package management
- **Package manager**: npm (no lockfile such as `package-lock.json` or `yarn.lock` present)
- **Scripts**: No custom npm scripts are defined in `package.json`.

## Primary production dependencies
- `@opencode-ai/plugin` @ 1.15.4 – core OpenCode plugin used by the GSD framework.

## Development dependencies (absent)
- No `devDependencies` were listed. When the application code is added, consider adding tools such as:
  - `eslint` / `prettier` for linting & formatting
  - `jest` or `vitest` for testing
  - `typescript` for static typing (optional)

## Ecosystem summary
- No additional language ecosystems (Python, Go, Rust, etc.) were detected.
- No containerization (`Dockerfile`) or CI configuration files (`.github/workflows/`) are present yet.

*This document will need to be updated as the project adopts new libraries, frameworks, or runtime environments.*
