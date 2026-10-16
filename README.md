<p align="center">
  <img src=".github/branding/readme-banner.png" alt="Eloque Programming Language Banner" width="100%">
</p>

# The Eloque Programming Language

**Eloque** is a proposed general-purpose programming language with structured, conversational-English-inspired syntax and precise, statically checked semantics. It is intended to be approachable for beginners while providing predictable diagnostics and tooling for AI-assisted development.

> **Status:** Early development — Milestone 1 (Frontend foundations) in progress. The lexer, initial AST structures, parser infrastructure, Pratt arithmetic expression parser, and initial statement and program parsers are implemented and tested. The parser supports `let` declarations, `change` assignments, and `say` statements containing the currently supported arithmetic expressions, identifiers, and literals. It constructs complete program ASTs from prepared token streams, with source spans and fail-fast syntax diagnostics. Lexer–parser integration and full frontend conformance remain pending. Eloque cannot yet execute source programs. The working name, `.elq` extension, and CLI name remain provisional. The examples below are **illustrative, not executable**.

## A glimpse of Eloque

```text
let greeting be "Hello"
let name be "Alex"

to greet with (person as String) giving String:
    give "{greeting}, {person}!"

say greet with (name)

let numbers be [1, 2, 3]
add 4 to numbers

for each number in numbers:
    say number

define Person:
    name as String
    age as Int

    to introduce with ():
        say "I'm {self.name} and I'm {self.age}."

let person be Person with (name be "Alex", age be 20)
person.introduce with ()
```

This example illustrates the intended language, including features that will arrive at different implementation milestones. It is **not** a statement that these features currently run.

## Design principles

- **Readable, not ambiguous:** conversational keywords follow a deterministic grammar; Eloque does not attempt to interpret unrestricted English.
- **Static safety:** types are checked before execution, with inference for eligible local bindings and explicit function signatures.
- **Predictable behavior:** strict booleans, checked integer arithmetic, explicit mutation and defined value semantics.
- **Useful diagnostics:** errors should be clear to beginners and available in structured form for editors and AI agents.
- **Concurrency by design:** safe native concurrency is a longer-term goal, not a feature of the initial implementation.
- **Measure before optimizing:** the first runtime is a Go-based stack bytecode VM; future optimization and backend choices depend on evidence.

## Documentation

- [Language specification](docs/LANGUAGE_SPEC.md) — intended syntax and semantics
- [Design decisions](docs/DESIGN_DECISIONS.md) — current agreements and superseded syntax
- [Architecture](docs/ARCHITECTURE.md) — proposed compiler, checker, runtime and VM
- [Roadmap](docs/ROADMAP.md) — incremental milestones and provisional releases
- [Technical backlog](docs/TECHNICAL_BACKLOG.md) — deferred engineering decisions and review triggers

## Planned implementation

Eloque is planned to be implemented in **Go**, with the following initial pipeline:

```text
.elq source
    ↓
Lexer → Parser → AST
    ↓
Name resolution → Static type checker
    ↓
Bytecode compiler → Stack-based VM
    ↓
Program execution
```

The first implementation milestones will establish the Go project, lexer, parser and source diagnostics, followed by a **small statically checked executable subset**. More advanced features such as collections, custom types, modules, exceptions and concurrency will be introduced incrementally.

No installation instructions are provided yet because there is no working Eloque interpreter or compiler to install.

## Project status

Eloque is in active frontend development under Milestone 1.

The lexer and foundational AST/parser infrastructure are implemented and tested. The Pratt expression parser supports integer and string literals, identifiers, grouping, unary arithmetic, and binary arithmetic with the documented precedence and associativity rules.

The initial statement parser supports `let`, `change`, and `say`. Complete-program parsing supports newline-separated statements, blank lines, and EOF termination. Parsed statements and programs retain source spans, while malformed syntax produces structured, fail-fast diagnostics.

Expression and program parsing have undergone correctness, robustness, and performance testing. Initial program-parsing benchmarks demonstrate approximately linear scaling across workloads containing 3, 300, and 3,000 statements. These benchmarks measure parsing from prepared tokens, excluding lexing.

The M1-P3 implementation has been completed after a final checkpoint review. Lexer–parser integration, broader frontend conformance, mandatory static type checking, bytecode compilation, and virtual-machine execution remain future work.

The language specification describes the intended language, not the currently executable feature set. Syntax, release boundaries, and technical details may evolve before version 1.0.
