---
layout: doc
date: 2026-05-19
---

# Code Conventions & Style Guide

The repository currently contains only the OpenCode tooling scripts, which follow the style conventions of the Get‑Shit‑Done framework:

- **JavaScript style** – CommonJS modules, `eslint` configuration is not present; the tooling uses standard Node.js idioms.
- **Naming** – Files are snake‑case or kebab‑case; functions and variables use `camelCase`.
- **Error handling** – Synchronous code throws exceptions; asynchronous code (if added) should use `async/await` with proper `try/catch` blocks.
- **Documentation** – Inline comments are used sparingly; future code should include JSDoc‑style comments for exported functions.

## Missing conventions

No linting configuration (`.eslintrc`, `prettier.config.js`) or formatting rules exist. When application code is introduced, consider adding:

- **ESLint** with the Airbnb or Standard style guide.
- **Prettier** for consistent formatting.
- **TypeScript** for static typing (optional but recommended).
- **Commit linting** to enforce conventional commit messages.

## Recommendations

1. Add an ESLint configuration file (`.eslintrc.json`).
2. Add a Prettier configuration (`.prettierrc`).
3. Adopt a coding standard (e.g., Airbnb) and enforce it with a pre‑commit hook.
4. Document any domain‑specific patterns (e.g., error‑code enums) in this file as they emerge.
