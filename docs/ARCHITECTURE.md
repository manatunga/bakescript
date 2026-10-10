# BakeScript Architecture

**Status:** Proposed architecture; no production implementation is implied.

## Initial execution pipeline

```text
Source (.bs)
    |
    v
Lexer / indentation handling
    |
    v
Parser ---> AST
    |
    v
Name resolution / semantic analysis (incremental)
    |
    v
Bytecode compiler
    |
    v
Bytecode + constants
    |
    v
Go-based stack VM
    |
    v
Program execution
```

A static checker (`bake check`) may share the frontend but must not be required
for normal execution by default.

## Responsibilities

| Component | Primary responsibility |
|---|---|
| Lexer | Tokens, source locations, indentation |
| Parser | Syntax tree and syntax errors |
| Resolver | Scope and name binding |
| Checker | Optional gradual type diagnostics |
| Compiler | Convert validated AST to VM instructions |
| VM | Execute instructions, manage frames and stack |
| Runtime | Values, collections, functions, objects, errors, built-ins |
| CLI | `bake run`, `bake check`, `bake fmt`, REPL |

## VM direction

Start with a straightforward stack-based bytecode VM in Go.

Likely early building blocks:
- Instruction/opcode representation.
- Constant pool.
- Operand/value stack.
- Indexed local variable slots.
- Call frames and return addresses.
- Branch instructions and loop compilation.
- Runtime values and error propagation.
- Closure/upvalue support once basic functions work.

Do not prematurely commit to a specific value representation, bytecode encoding,
GC scheme, or object layout. Prototype, test, and benchmark alternatives.

Go's garbage collector can manage Go-allocated runtime objects, though language
semantics such as `freeze` and object identity still require explicit design.

## Performance engineering

Measure first. Potential later optimizations:
- Faster dispatch and compact instructions.
- Specialized numeric operations.
- Reduced allocations and boxing.
- Local-slot and call-frame improvements.
- Inline caches for property/method lookup.
- Adaptive specialization.

Performance must be compared on equivalent, correctness-tested workloads.
No particular speedup over CPython is assumed.

## Backend independence

**The bytecode VM is the initial backend, not BakeScript's identity.**

Possible future architecture:

```text
Source
  |
Frontend / semantic model
  |
  +-------------------+-------------------+
  |                   |                   |
Bytecode VM      Native/AOT backend   Wasm backend
  |
Optional JIT / adaptive specialization
```

This is conceptual. A shared IR may become useful if multiple backends are
actually pursued; **do not build an elaborate multi-backend IR for v0.1**.

Keep parsing and language semantics separate from VM-specific machinery so
future backends can be added without redefining BakeScript's behavior.

## Compatibility boundaries

- Source-level semantics become a compatibility commitment at 1.0.
- Internal bytecode layout is unstable unless explicitly standardized later.
- Backend changes that preserve public behavior need not trigger a major version.
- A new optional backend can be a compatible minor feature.
- Breaking public language behavior requires a major release after 1.0.

## Decisions still open

- Runtime value representation.
- Exact bytecode instruction set and encoding.
- Exception unwinding mechanics.
- Module caching and initialization.
- Deep immutability enforcement.
- Debug information and source maps.
- Whether/when to introduce a backend-neutral IR.
