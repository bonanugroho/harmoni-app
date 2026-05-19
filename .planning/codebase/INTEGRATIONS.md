---
layout: doc
date: 2026-05-19
---

# External Integrations

The repository does not contain explicit configuration files for external services such as databases, third‑party APIs, authentication providers, or message queues.

- No `docker-compose.yml` or Kubernetes manifests were found.
- No `.env` files or other environment variable templates are present.
- No Prisma, TypeORM, or other ORM configuration files were detected.

## Potential integration points (future)

| Integration | Typical config file | Reason for inclusion |
|-------------|--------------------|----------------------|
| **Database** | `prisma/schema.prisma`, `ormconfig.json` | Persistent storage for user data.
| **Authentication** | `.auth0rc`, `next-auth.config.js` | OAuth, SSO, magic‑link support.
| **Payments** | `stripe.config.js` | Billing and subscription handling.
| **External APIs** | `src/api/*.ts` with OpenAPI specs | Connect to third‑party services (e.g., Slack, Twilio).
| **Message Queue** | `bullmq.config.js`, `kafka.yaml` | Async job processing.
| **Search** | `elastic.config.js`, `meilisearch.yaml` | Full‑text search capabilities.

*These entries are placeholders; when the application codebase grows and such services are added, the corresponding configuration files should be added and this document updated accordingly.*

## Action items

- When adding a new service, create a dedicated config file and list it here.
- Consider adding a `.env.example` to document required environment variables.
- Review this document after each major feature addition to keep integration coverage up‑to‑date.
