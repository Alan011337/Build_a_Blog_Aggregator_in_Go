# Gator — Go + PostgreSQL Blog Aggregator

A command-line RSS/blog aggregator built as a backend learning project in **Go + PostgreSQL**.

The project is useful as implementation evidence for backend-oriented technical fluency: CLI command design, persistence, SQL-backed data access, user/feed state, periodic aggregation, and turning external feed data into stored application data.

This is a **learning project**, not a production feed platform.

## What this project demonstrates

- Structuring a multi-command Go CLI application
- Persisting users, feeds, follows, and posts in PostgreSQL
- Working with SQL-backed data-access code
- Maintaining application state across commands
- Fetching and parsing RSS/feed data
- Running periodic aggregation workflows
- Separating command orchestration, configuration, database access, and domain behavior

## System flow

```text
CLI command
   ↓
Command handler / application logic
   ↓
PostgreSQL queries + persisted state
   ↑
RSS feed fetch / aggregation loop
```

The central product behavior is simple: users register or log in, add/follow feeds, run aggregation, and browse posts that have been fetched and persisted.

## Data and state model

A useful way to reason about the system is to separate several kinds of state:

- **Current-user state** — which user the CLI is acting as.
- **Product state** — users, feeds, follows, and posts persisted in PostgreSQL.
- **External state** — RSS feeds that can change independently of the application.
- **Aggregation state** — periodic work that fetches external content and converts it into local product data.

This separation matters because a feature that looks like “show me posts” depends on identity, external-network behavior, persistence, ingestion logic, and timing.

## Technical context

- **Language:** Go
- **Database:** PostgreSQL
- **Data access:** SQL / PostgreSQL driver
- **External data:** RSS feeds
- **Interface:** command-line application
- **State:** local configuration + persisted database records

## Representative commands

The nested project README contains the full setup and command reference. Representative flows include:

- `register <name>` — create a user
- `login <name>` — switch the current user
- `addfeed <name> <url>` — add a feed
- `agg <duration>` — periodically fetch feeds
- `browse [limit]` — inspect stored posts

## Repository structure

The implementation lives in:

[`Build_a_Blog_Aggregator_in_Go/`](./Build_a_Blog_Aggregator_in_Go)

That directory contains the Go source, configuration logic, database integration, and detailed run instructions.

## Why this matters for product / technical work

I am building technical fluency so I can make better product decisions and communicate with engineering teams with more precision. Gator helps me reason concretely about:

- how product state is represented in a database;
- where application behavior belongs relative to persistence;
- how periodic/background work differs from request-driven flows;
- how external data is ingested, transformed, and stored;
- what constraints appear when a seemingly simple product feature crosses CLI, network, database, and scheduling boundaries;
- how stale or failed external data should be distinguished from internal application state.

For a Product / TPM context, this is particularly useful for discussing ingestion pipelines, persisted state, background jobs, data freshness, idempotency, and failure boundaries without pretending that a learning project is a production distributed system.

## Reviewer guide

A useful discussion of this repository should be able to answer:

1. Which data should be persisted, and which state can remain local or transient?
2. What is different about an aggregation loop versus a request-driven API flow?
3. What can go wrong when external RSS data is fetched repeatedly?
4. How would you prevent duplicate or stale data from confusing product behavior?
5. If aggregation stops working, how would you separate network, parsing, database, and scheduling failures?
6. What would need to change before this could support multiple concurrent users or production reliability expectations?

## Evidence boundary

This repository supports claims about Go application structure, PostgreSQL persistence, SQL-backed state, RSS ingestion, background/periodic workflows, and backend-system reasoning. It does **not** by itself establish production feed-platform ownership, distributed-systems expertise, high-scale concurrency, or production SRE/reliability experience.

The intended signal is **backend-system literacy and implementation practice**, not senior backend-engineering expertise.

## Related evidence

- [`chirpy`](https://github.com/Alan011337/chirpy) — Go HTTP API + PostgreSQL + authentication + webhooks
- [`build_an_ai_agent`](https://github.com/Alan011337/build_an_ai_agent) — Python tool-using AI agent + function calling
- [Haven Product Portfolio](https://somber-tamarillo-df3.notion.site/Haven-AI-native-Product-Portfolio-3d8ad9856018811b96ffe8e33a7d48ef) — my flagship AI-native product work

For a curated map of my public technical projects, see:
[`Alan-Zeng-Git-Hub`](https://github.com/Alan011337/Alan-Zeng-Git-Hub)
