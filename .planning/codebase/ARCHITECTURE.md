---
layout: doc
date: 2026-05-19
---

# Architecture Overview

The repository currently contains only the OpenCode tooling infrastructure (`.opencode/` directory) and no application source code. Consequently, there are no defined architectural layers, entry points, or domain‑specific components.

## High‑level structure

- **`.opencode/`** – Contains the Get‑Shit‑Done framework, custom agents, skills, and workflow hooks. It drives the GSD orchestration but is not part of the target product architecture.
- **Project root** – Holds configuration files (`package.json`, `AGENTS.md`, PRD docs) and the `.planning/` folder for planning artefacts.

## Expected future layers (placeholders)

When the actual application code is added, the typical architecture will likely include:

1. **Entry point** – e.g., `src/index.ts` or `main.js` that bootstraps the server/client.
2. **Domain layer** – business logic, services, and use‑case implementations.
3. **Infrastructure layer** – data access, external API clients, persistence.
4. **Presentation layer** – UI components (React, Vue, etc.) or CLI commands.

These layers should be reflected in the `ARCHITECTURE.md` once source files are present.
