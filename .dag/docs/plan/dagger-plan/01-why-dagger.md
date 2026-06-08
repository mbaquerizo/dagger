# Why Dagger

## The Problem

AI code generation quality is bottlenecked by context quality. Agents are given vague tickets, undocumented assumptions, and no visibility into why decisions were made. The result: hallucinated implementations, missed edge cases, and multiple correction rounds.

The context needed for good code generation exists — it's just not documented. It's in conversations, in meeting transcripts, in architectural decisions that were made but never written down.

## The Thesis

Curated, structured documentation beats RAG and ad-hoc prompt engineering for AI code generation. A ticket linked to its motivating ADR, the code exploration that mapped the relevant codebase, and the pitch that started it all gives an agent everything it needs to generate correct code on the first try.

This is Docs Augmented Generation (DAG).

## What Dagger Is

Dagger is the platform where those documents live and connect. It's a product management tool built from the ground up for the AI era:

- **Graph of connected documents** — every pitch, ADR, code exploration, and ticket is linked by explicit relationships. Agents traverse the graph to get full decision context.
- **Agent-first API** — the primary consumer is an AI agent. Tickets are served as optimized markdown with all linked context assembled into a single document. No UI chrome, no parsing overhead.
- **Human UI for the humans** — clean web interface for product managers, architects, and decision-makers to create, browse, and manage the document graph.

## What Differentiates It

| | Jira / Linear | Dagger |
|---|---|---|
| Primary consumer | Humans | Agents |
| Doc format | Free-form text | Structured docs with typed relationships |
| Context resolution | Manual search | Automatic — traverse decision graph |
| AI integration | Bolt-on feature | Built-in from day 1 |
| Token optimization | None | Structured for minimal token waste |

This is not "Jira with AI." This is an AI context server with a product management UI.

## Target User

Engineering orgs where developers use AI coding tools (Claude Code, Cursor, Copilot) daily. The product manager or engineering manager plans features in Dagger. Developers tell their agent "Work on DG-42." The agent fetches full context from Dagger's API and generates correct code informed by every decision made along the way.

## Why Now

1. AI coding tool adoption is accelerating — every major dev tool has an AI agent now
2. The bottleneck has shifted from "can the model code?" to "does the model have the right context?"
3. DAG conventions already exist and are proven in the DAG plugin
4. dag-research will produce evidence that shapes optimal documentation — Dagger is where those insights get applied
