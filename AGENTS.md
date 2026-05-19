## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:
- ALWAYS read graphify-out/GRAPH_REPORT.md before reading any source files, running grep/glob searches, or answering codebase questions. The graph is your primary map of the codebase.
- IF graphify-out/wiki/index.md EXISTS, navigate it instead of reading raw files
- For cross-module "how does X relate to Y" questions, prefer `graphify query "<question>"`, `graphify path "<A>" "<B>"`, or `graphify explain "<concept>"` over grep — these traverse the graph's EXTRACTED + INFERRED edges instead of scanning files
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).

<!-- GSD:project-start source:PROJECT.md -->
## Project

**Harmoni**

Harmoni is a community financial‑management web application for neighborhood‑scale administrations (Rukun Tetangga/RT and Rukun Warga/RW). It provides transparent, accountable income and expenditure reporting for Residents, RT Officers, and RW Officers, accessible via mobile‑first browsers.

**Core Value:** Transparency and accountability of community finances – if the reporting layer fails, the whole system loses trust.

### Constraints

- **Security**: Must use PASETO tokens and enforce Casbin policies per territory.
- **Data Isolation**: RT 01 officers cannot view RT 02 data.
- **Responsiveness**: UI must work on low‑end mobile browsers.
<!-- GSD:project-end -->

<!-- GSD:stack-start source:codebase/STACK.md -->
## Technology Stack

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
## Ecosystem summary
- No additional language ecosystems (Python, Go, Rust, etc.) were detected.
- No containerization (`Dockerfile`) or CI configuration files (`.github/workflows/`) are present yet.
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

- **JavaScript style** – CommonJS modules, `eslint` configuration is not present; the tooling uses standard Node.js idioms.
- **Naming** – Files are snake‑case or kebab‑case; functions and variables use `camelCase`.
- **Error handling** – Synchronous code throws exceptions; asynchronous code (if added) should use `async/await` with proper `try/catch` blocks.
- **Documentation** – Inline comments are used sparingly; future code should include JSDoc‑style comments for exported functions.
## Missing conventions
- **ESLint** with the Airbnb or Standard style guide.
- **Prettier** for consistent formatting.
- **TypeScript** for static typing (optional but recommended).
- **Commit linting** to enforce conventional commit messages.
## Recommendations
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

## High‑level structure
- **`.opencode/`** – Contains the Get‑Shit‑Done framework, custom agents, skills, and workflow hooks. It drives the GSD orchestration but is not part of the target product architecture.
- **Project root** – Holds configuration files (`package.json`, `AGENTS.md`, PRD docs) and the `.planning/` folder for planning artefacts.
## Expected future layers (placeholders)
<!-- GSD:architecture-end -->

<!-- GSD:skills-start source:skills/ -->
## Project Skills

No project skills found. Add skills to any of: `.claude/skills/`, `.agents/skills/`, `.cursor/skills/`, `.github/skills/`, or `.codex/skills/` with a `SKILL.md` index file.
<!-- GSD:skills-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->

<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
