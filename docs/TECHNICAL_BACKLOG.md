# Eloque Technical Backlog

**Status:** Living engineering decision register; pre-implementation. A review commitment is not a commitment to implement a specific technology.

See [LANGUAGE_SPEC.md](LANGUAGE_SPEC.md) for intended language behavior, [DESIGN_DECISIONS.md](DESIGN_DECISIONS.md) for settled decisions, and [ROADMAP.md](ROADMAP.md) for implementation milestones.

## How to use this register

- **Committed:** intended capability, with details or delivery timing still subject to engineering validation.
- **Review required:** specific engineering decision with an explicit review trigger.
- **Exploratory:** possible technology or optimization, not approved for implementation.
- **Resolved:** reviewed, with decision and evidence recorded.

At each milestone, inspect items whose triggers have arrived. Record the decision, evidence and next action. If deferred, assign a new explicit trigger. Language rules already settled in the specification must not be silently reopened as runtime implementation choices.

## VM and performance

### TB-001 — Stack-based versus register-based VM
- **Status:** Review required.
- **Baseline:** Go stack-based bytecode VM, with local slots and call frames.
- **Trigger:** After first working VM and repeatable benchmark harness (approximately M2).
- **Method:** Profile dispatch, stack traffic, bytecode size and representative programs. Compare simpler stack optimizations before prototyping a register VM.
- **Criterion:** Measurable end-to-end gains that justify additional compiler/runtime complexity.
- **Evidence:** Pending.

### TB-002 — Safe native concurrency
- **Status:** Committed long-term goal; design review required.
- **Baseline:** Sequential VM; avoid unnecessary global mutable state.
- **Trigger:** Once functions, value ownership and mutation semantics are stable, before concurrency implementation (M9 or later).
- **Method:** Design task lifetime, cancellation, isolation, error propagation, and explicit shared-state rules. Prototype goroutine-backed execution contexts and test race safety.
- **Criterion:** Safe, understandable semantics with useful concurrency; Go goroutines are an implementation mechanism, not the language contract.
- **Evidence:** Pending.

### TB-003 — Adaptive bytecode specialization
- **Status:** Exploratory.
- **Baseline:** Correct generic bytecode operations.
- **Trigger:** After meaningful workloads and profiling expose hot operations.
- **Method:** Prototype numeric/call/property specialization; measure allocation, dispatch and regression behavior.
- **Criterion:** Repeatable benefits without unacceptable complexity.
- **Evidence:** Pending.

### TB-004 — JIT compilation
- **Status:** Exploratory.
- **Baseline:** Bytecode interpretation.
- **Trigger:** Baseline VM optimized and remaining bottlenecks documented; formal review after initial stabilization.
- **Method:** Assess compilation overhead, portability, warmup, runtime integration and maintenance.
- **Criterion:** Clear benefit beyond improving the interpreter.
- **Evidence:** Pending.

### TB-005 — AOT/native backend
- **Status:** Exploratory.
- **Baseline:** Go-hosted bytecode VM.
- **Trigger:** Stable semantics and a concrete deployment/performance use case, no later than 1.0 architecture planning.
- **Method:** Compare native backend, transpilation and bytecode packaging.
- **Criterion:** Practical advantage and manageable compatibility burden.
- **Evidence:** Pending.

### TB-006 — WebAssembly
- **Status:** Exploratory.
- **Baseline:** Native Go runtime.
- **Trigger:** Concrete browser, portable sandbox or embedded execution requirement; otherwise 1.0 ecosystem planning.
- **Method:** Compare compiling Go runtime to Wasm against a dedicated backend.
- **Criterion:** Clear use case with sustainable support costs.
- **Evidence:** Pending.

### TB-007 — CPython comparative performance
- **Status:** Committed measurement program; parity is an aspiration, not a guarantee.
- **Baseline:** No measurements.
- **Trigger:** First usable VM (M2), then major runtime milestones.
- **Method:** Pin versions, hardware, warmup and inputs. Compare equivalent, correctness-tested arithmetic, loops, calls, strings, collections and object workloads as implemented.
- **Criterion:** Reproducible findings that guide optimization, not marketing claims based on isolated microbenchmarks.
- **Evidence:** Pending.

### TB-008 — Runtime value representation
- **Status:** Review required.
- **Baseline:** Undecided; Go GC can manage Go-allocated objects.
- **Trigger:** Before implementing the first VM value type (M2).
- **Method:** Compare Go interfaces, tagged representations and pointer-based values against checked `Int`, binary64 `Float`, `Bool`, `String`, optional values and future collections.
- **Criterion:** Correct semantics, maintainability, acceptable allocation/dispatch costs.
- **Evidence:** Pending.

### TB-009 — Intermediate representation and backend separation
- **Status:** Review required.
- **Baseline:** AST directly to bytecode; frontend semantics independent of VM.
- **Trigger:** Second backend approved or optimization passes become unwieldy.
- **Method:** Identify real requirements and prototype the smallest useful IR.
- **Criterion:** Lower overall complexity, not speculative extensibility.
- **Evidence:** Pending.

## Frontend and language semantics

### TB-010 — Conversational grammar and call parsing
- **Status:** Review required.
- **Baseline:** `name with (...)` for calls, including zero arguments; multiword operators recognized by the parser.
- **Trigger:** Before expression parsing in M1.
- **Method:** Write grammar and conformance fixtures for nested calls, named/default arguments, precedence, grouping, indexing and method calls.
- **Criterion:** Deterministic, unambiguous parsing and precise diagnostics.
- **Note:** The old `call ... with ...` and bare built-in invocation debate is superseded; do not treat it as active syntax.
- **Evidence:** Pending.

### TB-011 — Scope, capture and immutability
- **Status:** Review required.
- **Baseline:** Block scopes, no shadowing of enclosing visible names, `fixed` binding immutability and `frozen` deep immutability.
- **Trigger:** Scope resolution in M2; revisit alias/capture behavior before mutable collections and closures.
- **Method:** Test nested scopes, assignment, recursion, frozen values and aliasing. Define capture mechanics only when closures enter scope.
- **Criterion:** No accidental mutation through alternate references; diagnostics reflect declared rules.
- **Evidence:** Pending.

### TB-012 — Static type inference and checking
- **Status:** Committed; implementation review required.
- **Baseline:** Strict static types, explicit parameter annotations, local inference and `T?`; type errors block normal execution.
- **Trigger:** Before first executable subset in M2.
- **Method:** Specify expression typing, assignment compatibility, return-path validation, optional narrowing and invalidation, and error recovery. Add positive/negative checker fixtures for every implemented feature.
- **Criterion:** Predictable soundness boundaries and actionable diagnostics; never silently downgrade type errors to warnings.
- **Evidence:** Pending.

### TB-013 — Numeric semantics and conversions
- **Status:** Review required.
- **Baseline:** Checked signed 64-bit `Int`; IEEE 754 binary64 `Float`; floor division/remainder; exact-only implicit Int-to-Float promotion.
- **Trigger:** Before arithmetic checker and VM implementation (M2).
- **Method:** Test overflow, division by zero, negative operands, mixed types, representability, NaN/infinity policy and shortest round-trip formatting.
- **Criterion:** Compiler/runtime agreement and reproducible edge-case behavior.
- **Evidence:** Pending.

### TB-014 — Collection value semantics and mutation
- **Status:** Review required.
- **Baseline:** Lists, maps and structured values are logically independent when assigned; negative sequence indexing and optional safe lookup are designed.
- **Trigger:** Before mutable collections in M5.
- **Method:** Compare eager copying, copy-on-write and structural sharing; test nested mutation, `fixed` versus `frozen`, slicing and map-key behavior.
- **Criterion:** Preserve observable value independence without excessive implementation complexity.
- **Evidence:** Pending.

### TB-015 — Unicode and string indexing
- **Status:** Review required.
- **Baseline:** Immutable UTF-8-backed Unicode strings with character-oriented access, not raw byte indexing.
- **Trigger:** Before string indexing/slicing in M5.
- **Method:** Decide Unicode code points versus grapheme clusters; test combining marks, emoji sequences and malformed source input.
- **Criterion:** Well-documented predictable behavior with viable implementation cost.
- **Evidence:** Pending.

### TB-016 — Exceptions and cleanup
- **Status:** Review required.
- **Baseline:** Typed unchecked exceptions, `try`/`catch`/`finally`, and source-aware failures.
- **Trigger:** Before exception implementation (M7).
- **Method:** Define hierarchy, matching, propagation, return/loop exits through `finally`, and precedence when cleanup also fails.
- **Criterion:** Deterministic unwinding and clear diagnostics.
- **Evidence:** Pending.

### TB-017 — Modules and visibility
- **Status:** Review required.
- **Baseline:** `import`, aliases with `named`, public-by-default top-level symbols, underscore-private names, once-only initialization and rejected circular imports.
- **Trigger:** Before module implementation (M7).
- **Method:** Define path resolution, dependency order, initialization failure behavior and caching.
- **Criterion:** Deterministic behavior and compile-time privacy enforcement.
- **Evidence:** Pending.

### TB-018 — Structured types, interfaces and matching
- **Status:** Review required.
- **Baseline:** `define` structured values, implicit `self`, structural interfaces and `match`/`when` without fallthrough; no class inheritance.
- **Trigger:** Before M6.
- **Method:** Specify construction, methods, structural conformance, recursive values, exhaustiveness and value semantics.
- **Criterion:** Coherent static checks and understandable errors.
- **Evidence:** Pending.

## Tooling, agents and execution safety

### TB-019 — Structured diagnostics
- **Status:** Committed; schema review required.
- **Baseline:** Human-readable diagnostics plus machine-readable output for agents and editor tools.
- **Trigger:** Begin source spans in M1; finalize initial schema before M8 tooling.
- **Method:** Define stable error codes, severity, file/line/column/spans, related locations and optional suggested fixes.
- **Criterion:** Diagnostics are actionable, deterministic and easy to parse without requiring prose interpretation.
- **Evidence:** Pending.

### TB-020 — Sandboxed execution budgets
- **Status:** Committed long-term capability; design review required.
- **Baseline:** Normal CLI has no arbitrary global instruction limit. Sandboxed runs should support configurable instruction, time and memory budgets.
- **Trigger:** Before agent-oriented execution or untrusted-code hosting (M9 or earlier if required).
- **Method:** Define metering points, cancellation, resource isolation, host guarantees and failure categories.
- **Criterion:** Explicit, testable limits with no misleading security guarantees.
- **Evidence:** Pending.

### TB-021 — References and ownership
- **Status:** Review required.
- **Baseline:** Value semantics by default; explicit shared mutable references are deferred.
- **Trigger:** Before public concurrency/shared-state syntax.
- **Method:** Evaluate ownership, borrowing or other explicit-sharing models against ease of learning and safety.
- **Criterion:** Preserve predictable mutation while enabling useful concurrent programs.
- **Evidence:** Pending.

## Milestone review checklist

- [ ] **M1:** TB-010 grammar; begin TB-019 source spans.
- [ ] **Before M2:** TB-008 values, TB-012 checker, TB-013 numeric behavior, TB-011 scope.
- [ ] **After working M2:** TB-001 stack/register review; TB-007 benchmark baseline.
- [ ] **Before M5:** TB-014 collections/immutability; TB-015 Unicode.
- [ ] **Before M6:** TB-018 custom types/interfaces/matching.
- [ ] **Before M7:** TB-016 exceptions; TB-017 modules.
- [ ] **M8:** TB-019 diagnostics/tooling schema.
- [ ] **Before M9 or sandbox deployment:** TB-002 concurrency, TB-020 budgets, TB-021 ownership.
- [ ] **When measured evidence or concrete use cases arise:** TB-003/004/005/006/009.

## Review record template

```text
Date / milestone:
Item ID:
Decision: implement / retain baseline / defer / reject
Evidence (tests, benchmarks, prototypes, links):
Reasoning:
Owner and next action:
Next review trigger (if deferred):
```

Do not mark an item resolved merely because a target release has passed. Review it or explicitly reschedule with a rationale.
