# BakeScript Design Decisions

A lightweight decision log. **Decided** means an explicit current choice;
**Provisional** means a working assumption that may be revised.

| ID | Status | Decision | Reason / note |
|---|---|---|---|
| D-001 | Decided | Implement BakeScript in Go | Standalone tooling, native compiled runtime, room to optimize |
| D-002 | Decided | `.bs` extension and significant indentation | Concise, readable programs |
| D-003 | Decided | Conversational syntax with selective cooking vocabulary | Distinct identity without forcing metaphors |
| D-004 | Decided | `prep` mutable declaration, `freeze` deep immutability, `as` assignment/default | Readable data operations |
| D-005 | Decided | `fresh`, `spoilt`, `empty` | Language-specific value vocabulary |
| D-006 | Decided | `taste` / `otherwise` conditionals; `instead` currently provisional | Natural branching |
| D-007 | Decided | `each item in collection:` for iteration | More intuitive than `stir item in collection:` |
| D-008 | Decided | `stir N times:` and `stir once/twice/thrice:` for fixed repetition | `stir` suits counted repetition |
| D-009 | Decided | `until condition:` for conditional repetition | Clear stop condition |
| D-010 | Decided | `to` declaration, `cook` call, `serve` return | Consistent conversational function syntax |
| D-011 | Decided | Optional `of type Type` and `: Type` annotations | Expressive and compact forms are equivalent |
| D-012 | Decided | `-> Type` return annotation; `as` for defaults | Avoid type/default ambiguity |
| D-013 | Decided | Gradual static checking, nonblocking by default | Preserve dynamic execution |
| D-014 | Decided | `menu`, `order`, implicit `my` for objects | Cooking identity in OO constructs |
| D-015 | Decided | `attempt`/`rescue`/`always`, `burn` | Error handling vocabulary |
| D-016 | Decided | `grab` for module imports | Short, readable imports |
| D-017 | Decided | Initial Go bytecode compiler and VM | Good balance of control, performance potential, and complexity |
| D-018 | Decided | Execution backends remain replaceable/extensible | VM is implementation strategy, not permanent restriction |
| D-019 | Decided | SemVer-inspired versioning | Compatible minors/patches after 1.0; majors for breaking changes |
| D-020 | Decided | CPython parity or better is a long-term performance ambition | Benchmark-driven engineering target, not a guarantee |
| D-021 | Provisional | Python-like behavior for several edge cases | Needs explicit specification before compatibility promises |

## Superseded proposals

| Old proposal | Replacement | Reason |
|---|---|---|
| `stir item in collection:` | `each item in collection:` | `stir` feels ambiguous for iteration |
| Only `name of type Type` | Also allow `name: Type` | Reduce verbosity |
| `drop` for loop break | `stop` | Current preferred control-flow word |

## Not adopted for now

- Grouped parameter type shorthand `(a, b): int`.
- A mandatory second execution backend in the initial release.
- A comprehensive multi-backend IR before a second backend is justified.

## Maintaining this log

When a decision changes:
1. Update the language specification.
2. Record the old and new form here, with a brief rationale.
3. Update relevant examples and tests.
4. If public behavior is affected, assess version compatibility.

This log intentionally preserves the difference between decisions explicitly
settled in discussion and provisional defaults accepted to keep design moving.
