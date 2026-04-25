# Sporecaster

> The Sporecaster defines the shape of your codebase and verifies it holds — giving AI agents the context to keep it growing without becoming messy and unmaintainable. Your spore is the single source of truth. The code follows.

---

## What it is

The Sporecaster is a spec-driven AI workflow for software development. Instead of pointing an AI at your code and hoping for the best, you define the architecture first — in a file called a **spore** — and everything flows from that.

The spore describes your application: its packages, how they connect, what flows between them. The codebase is what grew from it. Sporecast keeps the two in sync.

---

## How it works

**1. Define your spore**
A `spore.yaml` file at the root of your project defines the intended architecture — packages, channels, dependencies, and structure. This is the truth your code must reflect.

**2. Verify it holds**
Sporecast checks that the actual codebase matches the spore. Package declarations, channel ownership, import structure — if something drifts from the spec, it gets flagged.

**3. Grow with agents**
A set of AI agent personalities use the spore as their primary context. They understand the intended architecture before touching a single file.

| Agent | Role |
|---|---|
| **Architect** | Evolves the spore, proposes structural changes |
| **Builder** | Implements features within the defined structure |
| **Debugger** | Diagnoses issues with full architectural context |

---

## Status

Early and experimental. Currently targeting Go projects. The goal is to expand to multiple languages — starting with Ruby on Rails.

A web visualiser to explore spores and their relationships is also planned.

---

## Roadmap



---

*Cast the architecture. Grow the codebase.*
