---
layout: doc
date: 2026-05-19
---

# Repository Structure

```
/ (project root)
├─ .git/                     # Git repository
├─ .opencode/                # GSD framework, agents, skills, hooks
│   ├─ get‑shit‑done/
│   ├─ hooks/
│   └─ plugins/
├─ .planning/                # Planning artefacts (codebase map, project docs)
│   └─ codebase/
│       ├─ ARCHITECTURE.md
│       ├─ CONCERNS.md
│       ├─ CONVENTIONS.md
│       ├─ INTEGRATIONS.md
│       ├─ STACK.md
│       ├─ STRUCTURE.md
│       └─ TESTING.md
├─ AGENTS.md                 # Agent definitions for GSD
├─ PRD‑requirements‑en.md   # Product requirement document (English)
├─ PRD‑specifications‑en.md # Specification document (English)
└─ package.json              # Minimal npm manifest for tooling
```

*Note*: No `src/`, `lib/`, or other application directories exist yet. They will be added as development progresses.
