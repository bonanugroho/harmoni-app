---
layout: doc
date: 2026-05-19
---

# Testing Strategy

At present the repository contains only the GSD tooling and no application code, therefore no tests exist.

## Existing test infrastructure

- No test runner (Jest, Mocha, Vitest, etc.) is configured.
- No `__tests__` or `spec/` directories are present.

## Recommended testing setup for future development

1. **Choose a framework** – For a JavaScript/TypeScript project, Jest or Vitest are common choices.
2. **Add configuration** – Create `jest.config.js` or `vitest.config.ts` with sensible defaults.
3. **Write unit tests** – Place them alongside source files (`src/**/*.test.ts`) or in a dedicated `tests/` folder.
4. **Integrate with CI** – Add a script in `package.json` (e.g., `"test": "jest"`) and ensure the CI pipeline runs `npm test`.
5. **Coverage** – Enable coverage reporting (`--coverage`) and enforce a minimum threshold.

When the first source files are added, populate this document with concrete test plans, framework choices, and coverage goals.
