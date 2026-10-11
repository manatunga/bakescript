# GeeBreeze Technical Backlog

**Status:** Living engineering decision register. This file records deferred work so that "later" always has a review point. A **review commitment** is not a commitment to implement a particular technology.

## How to use this file

- **Committed:** intended product capability, subject to specification and feasibility.
- **Review required:** a concrete architectural decision that must be revisited at a milestone.
- **Exploratory:** a possible approach, not yet approved for implementation.
- **Resolved:** reviewed with a recorded decision, evidence, and follow-up.

At each milestone review, inspect all items whose trigger has been reached. Record the decision, benchmark/prototype evidence where applicable, and the next action. If deferred again, assign a **new explicit trigger** rather than leaving it open-ended.

## Tracked items

### TB-001 — Stack-based versus register-based VM

- **Status:** Review required
- **Current baseline:** Stack-based bytecode VM with local variable slots.
- **Reason for deferral:** Simpler initial compiler, instruction semantics, and debugging.
- **Review trigger:** First working VM and repeatable benchmark harness (target: post-v0.2.0 architecture review).
- **Review method:** Profile dispatch, operand-stack traffic, local access, bytecode size, and representative workloads. Evaluate specialized stack instructions and, if justified, prototype a small register-based execution path.
- **Decision criteria:** Measurable end-to-end improvement weighed against implementation complexity, maintainability, and correctness. No automatic migration.
- **Decision/evidence:** Pending.

### TB-002 — Goroutine-backed concurrency and multicore execution

- **Status:** Review required; capability under consideration
- **Current baseline:** Sequential execution; avoid unnecessary global mutable VM state.
- **Reason for deferral:** Shared-object safety, cancellation, exception propagation, and task lifetime semantics need design.
- **Review trigger:** Functions, call frames, closures, and runtime value ownership model are stable (target: after v0.5.0 feature milestone).
- **Review method:** Prototype isolated VM execution contexts on goroutines; test CPU-bound parallelism, I/O concurrency, race safety, cancellation, and shared immutable values.
- **Decision criteria:** Safe, predictable semantics and worthwhile performance/ergonomics; define user-facing syntax separately.
- **Guardrail:** Do not introduce a mandatory global VM execution lock without documenting the tradeoff and alternatives.
- **Decision/evidence:** Pending.

### TB-003 — Adaptive bytecode specialization

- **Status:** Exploratory
- **Current baseline:** Generic bytecode operations.
- **Reason for deferral:** Optimization without profiling risks complexity without benefit.
- **Review trigger:** Representative benchmarks and profiling of a correct baseline VM, after v0.6.0 or when hot operations are measurable.
- **Review method:** Prototype specialization for high-frequency numeric, call, or attribute operations; track deoptimization/correctness behavior.
- **Decision criteria:** Repeatable improvements across relevant workloads without unacceptable code complexity.
- **Decision/evidence:** Pending.

### TB-004 — JIT compilation

- **Status:** Exploratory
- **Current baseline:** Bytecode interpretation.
- **Reason for deferral:** Significant implementation and maintenance cost.
- **Review trigger:** Baseline VM optimized and remaining hotspots documented; review no later than the post-v1.0 performance-planning milestone.
- **Review method:** Evaluate available JIT strategies, compilation overhead, portability, warmup, and representative workloads.
- **Decision criteria:** Benefits justify complexity relative to improving the interpreter.
- **Decision/evidence:** Pending.

### TB-005 — AOT / native backend

- **Status:** Exploratory
- **Current baseline:** Go-hosted bytecode VM.
- **Reason for deferral:** Stable semantics and real use cases should precede a second backend.
- **Review trigger:** Core semantics stabilize near v1.0 and native compilation requirements emerge; formal review at v1.0 architecture planning.
- **Review method:** Compare native compilation, Go transpilation, and existing bytecode packaging against startup, deployment, and performance needs.
- **Decision criteria:** Clear practical benefit and manageable compatibility burden.
- **Decision/evidence:** Pending.

### TB-006 — WebAssembly target

- **Status:** Exploratory
- **Current baseline:** Native Go runtime.
- **Reason for deferral:** Browser/portable runtime use cases are not yet established.
- **Review trigger:** Concrete browser, sandbox, or portable deployment requirement; otherwise review at v1.0 ecosystem planning.
- **Review method:** Evaluate compiling the Go runtime to Wasm versus a dedicated backend, including I/O and runtime limitations.
- **Decision criteria:** Demonstrable use case and sustainable implementation.
- **Decision/evidence:** Pending.

### TB-007 — CPython comparative performance

- **Status:** Committed measurement program; performance parity is an ambition, not a guarantee.
- **Current baseline:** No GeeBreeze measurements yet.
- **Reason for staging:** Correctness and comparable workloads must precede performance claims.
- **First trigger:** First usable VM (target: v0.2.0), then repeat at each major runtime milestone.
- **Review method:** Pin Go/GeeBreeze and CPython versions, hardware, inputs, warmup and iteration policy; compare equivalent implementations of loops, arithmetic, functions, strings, collections, and objects; report wall time, variance, and memory.
- **Decision criteria:** Use findings to prioritize optimization; do not infer overall parity from a single microbenchmark.
- **Decision/evidence:** Pending.

### TB-008 — Runtime value representation and memory management

- **Status:** Review required
- **Current baseline:** Not selected; Go GC is available for Go-managed objects.
- **Reason for deferral:** Premature boxing/object-layout decisions can constrain the VM.
- **Review trigger:** Before implementing the core VM value type (v0.2.0 design gate).
- **Review method:** Compare tagged values, interfaces, pointers, allocation cost, object identity, and `freeze` semantics; benchmark prototypes if practical.
- **Decision criteria:** Correct semantics, manageable implementation, acceptable allocation and dispatch costs.
- **Decision/evidence:** Pending.

### TB-009 — Intermediate representation and backend separation

- **Status:** Review required
- **Current baseline:** AST to bytecode compiler; keep language semantics distinct from VM details.
- **Reason for deferral:** A general-purpose IR is premature with one backend.
- **Review trigger:** When a second backend is approved or major optimization passes become difficult to express.
- **Review method:** Identify actual cross-backend/optimization requirements; prototype the smallest useful IR.
- **Decision criteria:** Reduced complexity overall, not speculative extensibility.
- **Decision/evidence:** Pending.

## Milestone review checklist

- [ ] **Before v0.2.0 VM implementation:** TB-008 value representation.
- [ ] **After first working VM / v0.2.0:** TB-001 stack vs register review; TB-007 first CPython benchmark baseline.
- [ ] **After v0.5.0 function milestone:** TB-002 goroutine concurrency design review.
- [ ] **After v0.6.0 collections milestone:** TB-003 specialization feasibility; update TB-007 benchmarks.
- [ ] **At v1.0 planning:** TB-005 native/AOT and TB-006 Wasm reviews; TB-007 performance status.
- [ ] **After baseline VM optimization / post-v1.0 planning:** TB-004 JIT review.
- [ ] **When a second backend is approved:** TB-009 IR review.

## Review record template

```text
Date / milestone:
Item ID:
Evidence (benchmarks, tests, prototypes, links):
Decision: implement / retain baseline / defer / reject
Reasoning:
Owner and next action:
Next review trigger (if deferred):
```

Do not close an item solely because its target release number has passed; either review it or explicitly reschedule it with rationale.

## TB-009 — Conversational-call grammar and built-in invocation

- **Status:** Review required.
- **Current baseline:** `call` is mandatory for user-defined invocation; older draft examples used `set()` and `int()` directly.
- **Review trigger:** Before freezing the parser grammar / implementing expression parsing.
- **Review method:** Test nested calls, built-in conversions, constructors, parentheses, variadics, method calls and operator precedence against a formal grammar and small corpus.
- **Decision required:** Whether built-ins use `call`, have a narrow syntactic exception, or require a separate design.

## TB-010 — Block scope, closures and deep immutability

- **Status:** Review required.
- **Current baseline:** Conditional/loop blocks introduce scopes; `outer` permits writes to captured variables; `fix` freezes nested values across aliases.
- **Review trigger:** Before resolver, closure representation and mutable collection runtime implementation.
- **Review method:** Write scoping and aliasing test matrices for shadowing, loop capture, early exit, uninitialized reads and mutation through references.
- **Decision required:** Scope lifetime, captured-cell model and frozen object/value representation.
