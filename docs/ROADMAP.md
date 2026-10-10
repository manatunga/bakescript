# BakeScript Roadmap

**Status:** Aspirational planning document, not a release commitment.

## Release philosophy

BakeScript uses Semantic Versioning as a guideline:

- `MAJOR`: incompatible public language/API changes after 1.0.
- `MINOR`: backward-compatible new capabilities after 1.0.
- `PATCH`: backward-compatible fixes.
- Before `1.0`, the language contract is still evolving and `0.x` releases
  may contain breaking changes.

A performance breakthrough does **not** automatically require a major release.
An internal bytecode format is not currently part of the public compatibility contract.

## Initial 0.x series (proposed)

| Release | Theme | Demonstrable milestone |
|---|---|---|
| `0.1.0` | Frontend foundations | Lexer, parser, AST, syntax diagnostics, basic CLI |
| `0.2.0` | First execution | Bytecode compiler and VM; values, arithmetic, variables, `pour` |
| `0.3.0` | Control flow | `taste`/`instead`/`otherwise`, `each`, `stir`, `until`, `skip`, `stop` |
| `0.4.0` | Functions | `to`, `cook`, `serve`, recursion, local scope |
| `0.5.0` | Advanced functions | Closures, defaults, named arguments, varargs |
| `0.6.0` | Collections | Lists, maps, tuples, sets, indexing, slicing |
| `0.7.0` | Objects | `menu`, `order`, `my`, methods, inheritance |
| `0.8.0` | Integration | Modules, imports, exceptions, essential built-ins |
| `0.9.0` | Developer experience | REPL, formatter, gradual type checker, standard library, stabilization |
| `1.0.0` | Stable language | Documented semantics, tests, usable tooling, compatibility policy |

Versions may be combined, split, or reordered as implementation reveals
dependencies. An internal milestone does not need to become a public release.

## Approximate development phases

1. **Go readiness (1–2 weeks):** structs, slices/maps, pointers, interfaces,
   error handling, tests. Existing Go basics need not be relearned.
2. **Lexer and parser (2–4 weeks):** source locations, indentation, grammar, AST.
3. **Compiler and first VM (3–5 weeks):** instructions, constants, stack,
   variables, branches.
4. **Functions and control flow (1–2 months):** calls, frames, closures, loops.
5. **Language breadth (2–3 months):** collections, classes, exceptions, modules.
6. **Tooling and release (2–3 months):** CLI, REPL, formatter, checker, docs.
7. **Performance (ongoing, with a focused 3–6-month phase):** benchmark,
   profile, optimize, and compare against CPython.

These are **rough part-time estimates**, not deadlines. A complete v1.0
may take considerably longer than the initial prototype.

## Longer-term 1.x direction

Possible compatible releases include:
- Standard library expansion.
- VM specialization, inline caches, and allocation reduction.
- Better diagnostics, debugger support, and formatting.
- Package tooling and editor integration.
- Experimental optional execution backends.

`2.0` or `3.0` should occur only when breaking compatibility is justified,
not simply because a feature or runtime improvement is important.

## Performance ambition

**Long-term goal:** approach or surpass CPython performance on a representative
set of pure-language workloads, with correctness and reproducibility.

Measurement plan:
- Establish a baseline against a pinned CPython version.
- Compare equivalent algorithms, not superficially similar source code.
- Benchmark arithmetic, loops, calls/recursion, strings, collections, and
  object access separately.
- Report environment, iterations, warmup policy, memory use, and variance.
- Profile before optimizing; avoid presenting illustrative speed multipliers
  as measured predictions.
- Retain regression tests alongside benchmarks.

CPython parity is a **goal**, not a release promise or guaranteed outcome.

## Other long-term opportunities

- Adaptive bytecode specialization.
- Optional JIT or AOT/native backend if profiling supports the investment.
- WebAssembly backend if portability use cases justify it.
- Language server, editor extensions, package ecosystem, and community docs.

## Near-term next actions

- [ ] Finalize minimum Go readiness topics.
- [ ] Establish repository skeleton and Go module.
- [ ] Turn decided syntax into lexer fixtures.
- [ ] Write token and indentation tests.
- [ ] Implement scanner, then parser, then bytecode compiler/VM.
- [ ] Maintain the specification and design log alongside code changes.
