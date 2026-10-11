# GeeBreeze

**GeeBreeze** is a proposed general-purpose programming language with conversational, intuitive syntax and precise semantics. The language is planned to be implemented from scratch in **Go**, starting with a bytecode compiler and VM.

> **Status:** Design / pre-implementation. GeeBreeze is a **provisional name**, pending availability and naming-conflict checks. The `.gb` extension and `geebreeze` command are also provisional. None of the examples below are executable yet.

```text
bring square from maths

fix greeting at "Hello"
let name be ask "What's your name?"

to greet with person as "Guest":
    say f"{greeting}, {person}!"

call greet with name

let numbers be [1, 2, 3]
add 4 to numbers

each number in numbers:
    say number

define Person:
    to setup with name:
        my.name = name

let person be create Person with "Alex"
```

## Principles

- Natural, readable instructions without pretending to parse unrestricted English.
- Concise expression, explicit semantics, and useful tooling.
- Optional type annotations and gradual static checking.
- A serious runtime roadmap, benchmarked rather than assumed.
- Core syntax is **provisionally settled**, not yet a formal grammar.

## Documents

- [Language specification](docs/LANGUAGE_SPEC.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Roadmap](docs/ROADMAP.md)
- [Design decisions and superseded syntax](docs/DESIGN_DECISIONS.md)
- [Technical backlog and review triggers](docs/TECHNICAL_BACKLOG.md)

## Implementation approach

The initial Go implementation will use a lexer, parser, semantic analysis, bytecode compiler, and stack-based VM. The stack-versus-register decision must be revisited after a functional VM and representative benchmarks. Performance, concurrency, JIT, AOT and WebAssembly are tracked as explicit future investigations rather than promises of completed features.
