# BakeScript

BakeScript is an experimental, cooking-inspired programming language designed around human readability, conversational syntax, and expressive simplicity. Built from scratch in Go, it combines Python-style indentation with intuitive keywords and familiar programming constructs.

BakeScript's initial implementation targets a bytecode compiler and virtual machine, with room for future execution backends and runtime optimizations.

> **Development status:** BakeScript is currently in its early design and development stages. Its syntax, semantics, and implementation are subject to change.

```bakescript
prep ingredients as ["Flour", "Sugar", "Eggs"]

each ingredient in ingredients:
    pour ingredient

stir 3 times:
    pour "Mixing..."

to double with value: int -> int:
    serve value * 2

pour cook double with 5
```

## Project documents

- [Language specification](docs/LANGUAGE_SPEC.md) — intended syntax and semantics; distinguishes decided rules from provisional ones.
- [Roadmap](docs/ROADMAP.md) — proposed releases and long-term objectives.
- [Architecture](docs/ARCHITECTURE.md) — planned compiler and VM design, without locking out future backends.
- [Design decisions](docs/DESIGN_DECISIONS.md) — rationale and superseded choices.
- [Technical backlog](docs/TECHINCAL_BACKLOG.md) - tracked technical objectives that have been deferred.

## Principles

- Readability and approachable English-like syntax, without excessive metaphors.
- Optional type annotations and gradual static checking.
- A practical, general-purpose language rather than a novelty-only esolang.
- Correctness first, then measurable performance optimization.
- An initial Go bytecode VM, with execution-model flexibility in the future.

There is no production runtime yet. Examples describe the *intended* language.
