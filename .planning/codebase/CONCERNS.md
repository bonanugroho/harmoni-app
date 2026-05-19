---
layout: doc
date: 2026-05-19
---

# Technical Concerns & Debt

The repository primarily consists of the Get‑Shit‑Done framework. Scanning the codebase for common markers (`TODO`, `FIXME`, `BUG`, `NOTE`) yielded **301** occurrences across the tooling files. These are expected in a framework that evolves rapidly.

## Top‑level concern categories

| Category | Typical locations | Example markers |
|----------|-------------------|-----------------|
| **Framework maintenance** | `.opencode/` hooks and plugins | `// TODO: improve performance of guard` |
| **Documentation gaps** | Hook files, skill templates | `// NOTE: add more examples` |
| **Legacy compatibility** | Compatibility shims for older runtimes | `// FIXME: remove once Node 20 is mandatory` |
| **Security reviews** | Guard scripts that exec user code | `// BUG: sanitize inputs` |
| **Testing placeholders** | No test suites yet | `// TODO: add unit tests for guard` |

## Action guidance

- **Prioritize security‑related concerns** (e.g., input sanitisation in guard hooks) before exposing the framework to user projects.
- **Address documentation notes** to keep the developer experience smooth.
- **Schedule a dedicated cleanup sprint** to reduce the TODO count once the core product stabilises.

*The full list of line numbers can be extracted with `grep -R -n -E "TODO|FIXME|BUG|NOTE" .` if deeper investigation is required.*
