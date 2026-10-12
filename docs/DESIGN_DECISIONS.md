# Eloque Design Decisions

**Status:** Pre-implementation design log. **Decided** means agreed direction or syntax, not implemented; **Provisional** requires additional specification.

| ID | Status | Current decision |
|---|---|---|
| D-001 | Decided | Implement in Go; initially lexer → parser → bytecode compiler → VM |
| D-002 | Provisional | Working name Eloque, CLI `eloque`, extension `.elq`; check conflicts |
| D-003 | Decided | Conversational, concise, rigorous, general-purpose; culinary theme retired |
| D-004 | Decided | Significant indentation, `:`, lowercase keywords, `#` comments |
| D-005 | Decided | Mutable `let x be value`; optional uninitialized `let x`; assignment `=` |
| D-006 | Decided | Deep immutable `fix x at value`; nested/alias mutation prohibited |
| D-007 | Decided | `true`, `false`, `nothing`; conventional collection literals |
| D-008 | Decided | Optional `: Type` / `of type Type`; gradual checking |
| D-009 | Decided | Return annotations `giving Type` / `-> Type` |
| D-010 | Decided | `if` / `else if` / `otherwise`; block scope |
| D-011 | Decided | `each`, `repeat`, `while`, `until`; `skip` / `stop` |
| D-012 | Decided | `to` declaration, `call` invocation, `give` return, defaults with `as` |
| D-013 | Decided | `...items`, `named ...options`, lexical closures, `outer` |
| D-014 | Decided | Anonymous expression function `with x give x * 2` |
| D-015 | Decided | Parentheses for grouping and nested calls |
| D-016 | Decided | `define` class, `from` inheritance, `setup`, `create`, `my`, `parent` |
| D-017 | Decided | `raise` / `try` / `catch` / `always` |
| D-018 | Decided | `bring` for imports; aliases with `as` |
| D-019 | Decided | `say` output; `ask` input |
| D-020 | Decided | `add`, `remove`, `remove every` collection operations |
| D-021 | Decided | `is`/`==` value equality; `is not`/`!=`; logical word/symbol aliases |
| D-022 | Decided | Stack VM first; mandatory register-VM comparison after working VM/benchmarks |
| D-023 | Provisional | CPython parity or better is a long-term benchmark target, not a guarantee |
| D-024 | Provisional | Backend extensibility, SemVer-inspired versioning, standard library |
| D-025 | Provisional | Preserve prior indexing/slicing alternatives and public-by-default members |
| D-026 | Open | Built-in `set()` / `int()` invocation versus mandatory `call` |
| D-027 | Open | Precise grammar, evaluation order, typing, errors, modules, class MRO |

## Superseded syntax (historical record)

| Original BakeScript | Eloque replacement |
|---|---|
| `prep x as 1`, `x as 2` | `let x be 1`, `x = 2` |
| `freeze x as 1` | `fix x at 1` |
| `fresh`, `spoilt`, `empty` | `true`, `false`, `nothing` |
| `taste`, `instead`, `otherwise` | `if`, `else if`, `otherwise` |
| `stir 3 times` | `repeat 3 times` |
| `pour`, `sip` | `say`, `ask` |
| `cook greet`, `serve x` | `call greet`, `give x` |
| `*args`, `**kwargs` | `...args`, `named ...kwargs` |
| `lambda x: x * 2` | `with x give x * 2` |
| `menu Person`, `menu Child based on Parent` | `define Person`, `define Child from Parent` |
| `order Person with ...` | `create Person with ...` |
| `my.x as value` | `my.x = value` |
| `burn`, `attempt`, `rescue`, `always` | `raise`, `try`, `catch`, `always` |
| `grab maths` | `bring maths` |
| `toss`, `pluck`, `dispose` | `add`, `remove`, `remove every` |
| `bake` CLI, `.bs` | Provisional `eloque` CLI, `.elq` |

**Unchanged unless a conflict arises:** `each`, `until`, `skip`, `stop`, `size of`, `at` indexing, bracket indexing/slicing, `slice ... from ... to ...`, `setup`, `my`, `parent`, dotted members, named/default arguments, lexical closures, first-class functions, optional annotations, and the technical backlog.

## Change process

When a decision changes, update the specification, this log, examples, and relevant implementation tests. Do not treat provisional semantics as settled merely because an example uses them.
