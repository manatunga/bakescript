# Eloque Architecture

**Status:** Proposed architecture, pre-implementation. The working language name, `.elq` extension and CLI command remain provisional. Nothing in this document implies an executable compiler exists.

## Goals and boundaries

Eloque is a general-purpose, structured conversational language designed for novice readability, precise static analysis and eventual safe concurrency. The implementation is written in Go, beginning with a **stack-based bytecode virtual machine**. Language semantics must not depend on the bytecode representation.

This document describes the target architecture and the order in which it will be built. The [language specification](LANGUAGE_SPEC.md) governs syntax and observable behavior; [design decisions](DESIGN_DECISIONS.md) records settled and superseded choices.

## Compilation and execution pipeline

```text
Source file (.elq)
       |
       v
Lexer (tokens, indentation, source spans)
       |
       v
Parser (AST and syntax diagnostics)
       |
       v
Name resolution and semantic analysis
       |
       v
Static type checker (mandatory for normal compilation)
       |
       v
Bytecode compiler (constants, instructions, debug metadata)
       |
       v
Stack-based Go VM
       |
       v
Program output or typed runtime error
```

**Static checking is not optional or gradual.** Eloque permits local type inference, but programs with unresolved type errors must not proceed to normal bytecode execution. A future `eloque check` command should run the frontend without executing user code. The frontend should collect useful diagnostics where safe instead of aborting at the first error.

The implementation may initially recognize only a small subset of the full design. Unsupported syntax should be rejected clearly, never silently accepted with guessed semantics.

## Component responsibilities

| Component | Responsibility |
|---|---|
| CLI | Route `run`, future `check`, and eventual formatting/debugging commands |
| Source manager | Track file names, line/column positions and source spans |
| Lexer | Recognize identifiers, keywords, literals, operators, newlines, indentation and dedentation |
| Parser | Construct AST nodes and report grammatical errors, including malformed conversational phrases |
| Resolver | Bind names, enforce scope rules and detect undeclared names and prohibited shadowing |
| Type checker | Infer eligible local types, validate explicit signatures, optional types, calls, operators and return paths |
| Bytecode compiler | Lower validated AST into executable instructions and constant pools |
| VM | Manage value stack, frames, instructions, calls, control flow and exception propagation |
| Runtime | Represent values, implement built-ins and collection semantics, enforce checked operations and immutability |
| Diagnostics | Render human-readable errors and provide a machine-readable structured form |

A single implementation package may initially handle more than one responsibility. Package boundaries should emerge from real dependencies rather than an elaborate empty directory tree.

## Frontend constraints

### Lexical and syntactic rules

- Four spaces per indentation level; tabs prohibited. Blank/comment-only lines do not change indentation.
- A colon introduces an indented block; newlines ordinarily terminate statements.
- Parentheses, brackets and braces permit implicit line continuation.
- Keywords and conversational operator phrases have precise grammar; this is **not** natural-language parsing.
- Prefer tokenizing individual words and recognizing multiword constructs in the parser.
- Preserve source spans on tokens and AST nodes for precise diagnostics.

### Name resolution and typing

- `let x be value` declares a mutable binding; `change x to value` updates it.
- `let fixed` prohibits rebinding; `let frozen` requires deeply immutable values.
- Names cannot shadow visible bindings in enclosing scopes.
- Function parameters require `as Type`; functions returning values require `giving Type`.
- Initialized local bindings can infer types. `none` is permitted only where an optional type is expected.
- Type checking is a required compilation stage, including for CLI execution.
- Any flow-sensitive narrowing must be invalidated where mutation or concurrency makes it unsafe.

### Evaluation

- Evaluate operands and function arguments left to right; `and`/`or` short-circuit.
- Preserve lexical scope and deterministic sequential execution by default.
- Top-level statements execute in source order.
- Imports initialize modules once per execution; reject circular imports in the initial design.
- No automatic `main` invocation.

## Runtime value model

The observable semantics of lists, maps and structured custom values are **value semantics**: assigning one to another produces logically independent values. The runtime may use copies, copy-on-write or structural sharing if mutation remains observationally independent.

- `Int`: signed 64-bit with checked overflow.
- `Float`: IEEE 754 binary64; no implicit lossy conversion.
- `String`: immutable Unicode text.
- `Bool`: strict booleans, no truthiness conversion.
- Optional values: `T?` can hold `none`.
- Lists/maps: typed operations, negative indexing for supported sequence types, end-exclusive slices, safe `find` access.

The VM's actual value representation (Go interfaces, tagged union, pointers, etc.) remains **open** until implementation. Go's garbage collector can manage Go-allocated objects, but it does not automatically implement Eloque's value semantics, deep freezing or safe shared-state policy.

## VM direction

Start with a minimal stack machine that can execute validated constants, arithmetic, `say` and bindings. Add branches, loops, calls, collections and exceptions as the frontend supports them.

Likely building blocks:

- Opcode/instruction representation and constant pool.
- Operand stack and local variable slots.
- Call frames and return addresses.
- Branch and loop instructions.
- Source-location/debug metadata.
- Typed runtime failures and exception unwinding.

Avoid premature commitments to bytecode encoding, object layout, garbage-collection strategies beyond Go-managed allocation, or a multi-backend intermediate representation.

A small tree evaluator can be used as a temporary test aid if useful, but **the primary planned runtime is bytecode plus VM**.

## Errors and diagnostics

Diagnostics are a first-class interface for both human programmers and AI agents.

- Compile-time errors should include an error category/code, message, file, line, column and relevant source span.
- Runtime failures should carry typed exceptions and source locations when available.
- Human-readable CLI diagnostics are the default; structured machine-readable output should be supported by the tooling design.
- Invalid direct indexing, missing map keys, integer overflow and division by zero must not silently return misleading values.
- Exceptions are unchecked; static analysis may warn about potentially unhandled errors without requiring exception declarations.

Exact diagnostic codes, JSON schema and exception hierarchy are deferred until their respective implementation milestones.

## Modules

Each `.elq` file is intended to be a module. `import` supports whole-module and selected-name imports, with aliases introduced by `named`. Top-level names beginning with `_` are private and access must be rejected by the compiler.

Module initialization runs at most once per program execution. Initialization order follows a defined dependency order; circular imports are rejected initially. Path resolution, caching representation and failure recovery need further specification.

## Concurrency and sandboxing

The initial runtime is **sequential**. Eventual native concurrency is a language goal, not a feature of the first VM.

Explicit references, ownership rules, shared mutable state, cancellation, task lifetime and exception propagation should be designed together. Go goroutines may be a runtime mechanism, but Eloque's user-facing safety guarantees must be specified independently.

Normal CLI execution has no arbitrary universal instruction budget. Sandboxed/agent execution should eventually support configurable instruction, time and memory limits. Hard guarantees depend on the implementation and hosting environment; not every infinite loop can be proven statically.

## Performance and future backends

Measure correctness first, then benchmark representative workloads. A stack-based VM is the baseline. Compare a register-based alternative only after the working VM and a repeatable benchmark harness exist.

Potential later work includes opcode specialization, reduced allocations, inline caches, optional JIT, AOT/native compilation and WebAssembly. None is a committed v0.1 feature. CPython parity is an aspiration, **not** an established result.

Keep language semantics separate from VM internals so future backends can preserve observable behavior without rewriting the language contract. Do not build a general-purpose IR until a real second-backend or optimization need justifies it.

## Compatibility

- Before 1.0, syntax and semantics may change, with documentation and conformance tests updated together.
- At 1.0, published source-level semantics become the compatibility commitment.
- Internal bytecode is not a stable public interchange format unless standardized later.
- Backend changes that preserve observable behavior need not require a major version.

## Open architecture decisions and review gates

| Decision | Review trigger |
|---|---|
| Runtime value representation and efficient value semantics | Before first VM value implementation |
| Exact instruction encoding and source-map design | First bytecode compiler milestone |
| Float edge cases, conversion rules and formatting | Numeric runtime and checker implementation |
| Collection freezing and safe mutation | Before mutable collection implementation |
| Exception taxonomy and unwinding precedence | Exception milestone |
| Module resolution and caching | Import milestone |
| Closure representation and advanced functions | When closures enter scope |
| Explicit references, ownership and concurrency | Dedicated concurrency design milestone |
| Structured diagnostic schema and sandbox limits | CLI tooling / agent execution milestone |
| Backend-neutral IR | Only when a second backend or real optimization need arises |

Review triggers are tracked in [TECHNICAL_BACKLOG.md](TECHNICAL_BACKLOG.md). Avoid interpreting a deferred implementation detail as an unsettled language principle.
