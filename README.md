# Dagger 🗡

The DAG-powered, product management tool built for AI.

## Docs Augmented Generation + Dagger

[Docs Augmented Generation](https://github.com/mbaquerizo/docs-augmented-generation), or DAG, is an approach to optimizing agentic code generation using well-structured and complete documentation. All decisions made, from idea to ticket, published and injected as context to your agent.

Dagger is where those documents go, and where your agent reads them from. Simple and elegant UI for humans, plus a blazing fast text-based interface for agents.

DAG provides the structure, backed by [dag-research](https://github.com/mbaquerizo/dag-research), dagger provides the interface for agents to do what they need to do, quickly and correctly.

## How it works

You tell your agent, "Work on DG-42". DG-42 is just a ticket. To a human, it might look like any other ticket in any other project management app.

But when the agent retrieves it, it gets every relevant piece of information from idea, to architectural decisions, to acceptance criteria, curated from the rich collection of documents Dagger preserved and polished during the lifecycle of the feature.

## Try it out

__Note: BYOK__ First, provide the key for your favorite model.

The entry point of Dagger is the DAG-powered chat interface.

**Have a feature idea  or just had a meeting about it?** Say, "I want to plan a new feature" or paste in your meeting transcript. Dagger will identify the stakeholder, product manager, and engineering voices and ask any questions it needs to create and publish a higher-level feature planning doc. It'll give you the planning doc ID and url when it's done.

**Want to get into the finer technical details?** Say, "Help me architect X" or paste the planning doc link. Dagger will guide you through a deep architectural decision making session, and then publish the ADR. You and your team could then decide

**Ready to create tickets?** Pass it the ADR or say "Create tickets for X". Dagger will create a thorough specification using your past discussions as context and output a DAG-ready ticket or collection of PR-sized tickets.

No missing acceptance criteria or scenarios. Your work agent will be equipped to create the most accurate code possibe, using the discussions and decisions you made throughout the lifecylce of that feature.

## Human UI

Humans will never be replaced. Your users need a product or service created with empathy, so it's important to us that your team of visionaries, architects, testers, and decision makers has a simple yet elegant user interface tuned for building that product. The chat interface is familiar and clean. The docs interface is well organized, and the ticket interface is focused.

## Agent interface

Agents interface via the Dagger AI secure API. It gets a focused version of the ticket or doc, complete with all the context it needs to be up to speed, served as raw Markdown.

## Development

Prerequisites: Go 1.26+

Setup:
```
git clone <repo>
cd dagger
go mod tidy
```

Run migrations:
```
make migrate
```

Run:
```
make run
# or: go run ./cmd/api
```

Health check:
```
curl http://localhost:8080/healthz
# → OK
```