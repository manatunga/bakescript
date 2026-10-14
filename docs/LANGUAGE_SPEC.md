# Eloque Language Specification

**Status:** Active language design specification, revised October 2026. Examples are illustrative, not executable.
**Working name / CLI / extension:** Eloque / `eloque` / `.elq` (naming checks remain pending).  
**Implementation direction:** Go frontend, static checker, bytecode compiler, stack-based VM.

**Status labels:** **Decided** means an agreed language rule, not an implemented feature; **Provisional** means the broad direction is accepted but edge cases remain; **Deferred** means not part of the initial executable subset. This document describes the *target language*, not a promise that every feature ships in v0.1.0.

## 1. Principles [Decided]

Eloque is a general-purpose programming language with structured conversational syntax, not unrestricted natural-language parsing. It prioritizes novice readability, unambiguous static semantics and machine-readable diagnostics for AI-assisted development, and a foundation for safe native concurrency. Prefer concise, regular grammar over excessive English verbosity. Execution is deterministic subject to explicit external inputs and eventual concurrency semantics.

## 2. Lexical rules [Decided]


- Significant indentation uses spaces, with widths aligned to multiples of four; tabs are invalid. An increase to any greater four-space-aligned width introduces one indentation level (`INDENT`), regardless of the size of the increase. A decrease must return to a previously active indentation width, emitting one `DEDENT` for each exited level. A `:` introduces an indented block. Blank and comment-only lines do not affect indentation.
- One statement per logical line. Newlines and indentation are insignificant inside matched parentheses `(...)`, allowing multiline expressions and inline comments. Implicit continuation inside `[...]` and `{...}` is planned for their respective syntax implementations.
- `#` starts a line comment outside a string.
- Identifiers use ASCII letters, digits, and `_`, cannot start with a digit, and are case-sensitive. Reserved keywords cannot be identifiers. Snake_case is the recommended naming convention.
- Strings may use single (`'`) or double (`"`) quotes. The matching quote terminates the literal. The lexer recognizes the escapes `\\`, `\"`, `\'`, `\n`, `\t`, and `\r`, preserving their original source spelling in token lexemes. String interpolation using `{expression}` is planned but is not implemented by the initial lexer. Precise interpolation and Unicode grapheme rules remain to be specified.
- Multiword operators are parsed from individual keyword tokens, e.g. `is greater than`.
- Decimal integer literals consist of ASCII digits. Leading zeros are forbidden except for the literal `0` itself. A digit sequence immediately followed by an ASCII letter or underscore is a malformed numeric literal and is rejected during lexical analysis (e.g. `123abc`, `42_name`). When both an invalid identifier suffix and leading zeros occur (e.g. `007abc`), the invalid-suffix diagnostic takes precedence.

## 3. Types and literals [Decided unless noted]

| Type | Meaning | Examples |
|---|---|---|
| `Int` | Signed 64-bit integer; checked overflow | `42`, `-7` |
| `Float` | IEEE 754 binary64 | `3.14`, `0.5` |
| `String` | Immutable Unicode text, stored as UTF-8 | `"Hello"` |
| `Bool` | Exactly `true` or `false` | `true` |
| Optional `T?` | Value of type `T` or absence | `String?` |

`none` is the absence literal, valid in optional contexts; it is not a general substitute for a non-optional value. `List[T]`, `Map[K, V]`, and user-defined types are statically checked. Literal examples: `[1, 2, 3]`, `{"name": "Alex"}`. Tuple and set types are **deferred**, not established core syntax.

- Integer overflow raises `IntegerOverflow`; integer division by zero raises an error.
- Float arithmetic is approximate; use the shortest round-trippable decimal when displaying a Float. Rounding is explicit.
- Implicit `Int` to `Float` promotion is permitted only when exactness can be established. Precise rules for value-dependent conversions and Float exceptional values remain provisional.
- No implicit truthiness; conditions require `Bool`.

## 4. Bindings, mutability, scope [Decided]

```elq
let age be 20
change age to 21
let fixed max_attempts as Int be 3
let frozen names as List[String] be ["Alice", "Bob"]
let nickname as String? be none
```

`let` declares a mutable binding; `change ... to ...` reassigns. `fixed` prevents rebinding but does not necessarily freeze the referenced value; `frozen` enforces deep immutability. Initializers are required in the current design; uninitialized bindings are not supported. Local binding types may be inferred statically from their initializers. Block scopes are lexical. **Shadowing is prohibited** when an enclosing declaration of the same name is visible; unrelated scopes may reuse names.

Lists, maps and structured values use **value semantics**: assignment creates an independent *logical* value. The runtime may optimize with copy-on-write or structural sharing without observable aliasing. Explicit references/shared mutable state are deferred to the concurrency design.

## 5. Expressions [Decided unless noted]

Arithmetic: `+`, `-`, `*`, `/` (real division), `//` (floor division), `%` (floor-compatible remainder). For example, `-5 // 2` is `-3` and `-7 % 3` is `2`. Division by zero raises an exception. Exact rules for mixed numeric operands require conformance tests.

Comparisons: `is`, `is not`, `is greater than`, `is less than`, `is at least`, `is at most`. Structural equality is defined for compatible lists, maps, and custom types where member equality is supported. Chained comparisons are **not supported in v0.1**.

Boolean operations: `not`, `and`, `or`, with short-circuit evaluation. `not age is at least 18` means `not (age is at least 18)`.

**Precedence, highest to lowest:** grouping/call/access; unary numeric minus; multiplication/division/modulo; addition/subtraction; ordering; equality; `not` over a comparison; `and`; `or`. Parentheses override precedence. Same-precedence arithmetic operators associate left-to-right. Operands and call arguments evaluate left-to-right, except short-circuiting operands that are skipped. The formal Pratt binding-power table will refine interactions between `not` and other expressions before parser implementation.

**Exponentiation:** The lexer recognizes `**` as a distinct operator. Its precedence relative to unary minus and its associativity remain provisional and must be finalized before arithmetic expression parsing.

## 6. Conditions and loops [Decided]

```elq
if age is at least 18:
    say "Adult"
otherwise if age is at least 13:
    say "Teen"
otherwise:
    say "Child"

for each name in names:
    say name

repeat 3 times:
    say "Again"

repeat while attempts is less than 3:
    change attempts to attempts + 1

repeat until ready:
    stop
```

In this example, `attempts` and `ready` are assumed to have been declared earlier. `repeat while` and `repeat until` test their condition before each iteration. `stop` exits the innermost loop; `skip` continues to its next iteration. Collection mutation during iteration and repeat-count edge cases remain provisional.

## 7. Functions [Decided]

```elq
to calculate_total with (price as Float, quantity as Float) giving Float:
    give price * quantity

to greet with (name as String, excited as Bool be false):
    say "Hello, {name}!"

let total be calculate_total with (20.0, 3.0)
greet with (name be "Alex")
greet with ("Alex", excited be true)
```

`to` defines a function, `give` returns a value. Parameters **require explicit type annotations**. Functions returning a value require `giving Type`; omitting it defines a function with no value return, where `give value` is invalid. Calls **always** use `with (...)`, including zero arguments (`greet with ()` where applicable). Named arguments use `be`, as do parameter defaults. Evaluation order is left-to-right. Function names conventionally start with verbs, but this is a style guideline.

First-class functions, closures, variadics and generics require further specification; older `call`, `outer`, and anonymous `with x give ...` forms are **not** grandfathered into the new grammar.

## 8. Collections and text [Decided unless noted]

```elq
let names be ["Alice", "Bob", "Charlie"]
say item at 0 in names
say item at -1 in names
let missing_name as String? be find item at 99 in names
say slice of names from 0 to 2
say size of names
add "Dana" to names
remove "Bob" from names
remove all "Alice" from names

let data be {"name": "Alex"}
say value for "name" in data
let missing_value as String? be find value for "missing" in data
add key "city" with value "Colombo" to data
change value for "city" in data to "Kandy"
remove key "city" from data
```

Indexing is zero-based; negative indices count backward from the end (`-1` is last). Slices include the start and exclude the end. Direct out-of-range indexing raises `IndexOutOfRange`; missing map keys raise `KeyNotFound`. `find` variants return an optional result (`none` when absent). List removal removes the first matching element; `remove all` removes all matches. Map equality is independent of entry order. String concatenation uses `+`; `size of` and slicing apply to strings as well. Unicode character versus grapheme indexing details remain provisional.

## 9. Custom types and interfaces [Decided design; implementation deferred]

```elq
define Person:
    name as String
    age as Int

    to greet with ():
        say "Hello, I'm {self.name}!"

let alice be Person with (name be "Alice", age be 25)
alice.greet with ()
change alice.age to 26

define interface Speaker:
    to speak with () giving String
```

`self` is implicit inside methods. Interfaces are satisfied **structurally and implicitly** by compatible method signatures. Traditional class inheritance, `setup`, `create`, `my`, and `parent` from the older draft are superseded. Field initialization, visibility diagnostics, and mutation rules for frozen instances require additional tests.

## 10. Pattern matching [Decided design; implementation deferred]

```elq
match status:
    when 200:
        say "OK"
    when 404:
        say "Missing"
    otherwise:
        say "Other"
```

Only the first matching branch executes; there is no fallthrough. `when value named name:` binds the matched value. Statement matches may omit `otherwise`; value-producing matches must be exhaustive. The expression form and exhaustiveness algorithm remain to be specified.

## 11. Exceptions [Decided]

```elq
try:
    raise DivisionError with ("Cannot divide by zero")
catch DivisionError named exc:
    say "Invalid calculation"
catch error named exc:
    say "Unexpected error"
finally:
    say "Cleanup"
```

At least one `catch` or `finally` is required. Uncaught exceptions propagate after `finally`; `finally` also runs on normal exits and `give`. Exceptions are **unchecked**: functions need not declare their possible exceptions. Static tooling may warn about potentially unhandled failures. Detailed exception hierarchy and precedence when `finally` itself fails remain provisional.

## 12. Modules and entry point [Decided]

```elq
import math_utils
import square from math_utils
import square, cube from math_utils
import math_utils named math
import square named squared from math_utils
```

A `.elq` file can execute top-level statements in source order; no automatic `main` invocation is required. An imported module initializes once per program execution, with dependency initialization in defined order. Circular imports are rejected in the initial implementation. Top-level declarations are public by default, except names beginning `_`, which are **compiler-enforced module-private**. Precise module lookup paths and packaging are deferred.

## 13. Execution, safety and diagnostics [Decided direction]

Normal execution is sequential by default. Assignment evaluates its right-hand side before changing the target. Runtime failures are typed and include source locations. Human-readable diagnostics are the default; structured machine-readable diagnostics (e.g. JSON) are a design requirement for agent tooling. The runtime should support configurable instruction/time/memory limits for sandboxed execution; normal CLI execution has no arbitrary universal instruction cap. Native safe concurrency is a language goal, but syntax, ownership and shared-state rules are **deferred**.

**Lexical diagnostics:** The initial lexer uses fail-fast error handling: scanning stops at the first detected lexical error. Diagnostics identify the source line and column and describe the invalid input. Lexical errors are compiler diagnostics, not runtime exceptions. Error recovery and reporting multiple lexical errors in a single pass may be introduced in a later implementation.

## 14. Tooling and backend [Provisional roadmap]

Planned commands include `eloque run`, `eloque check`, `eloque fmt` and eventually a REPL. The initial Go implementation uses a lexer, parser, AST, mandatory static type checker, bytecode compiler and stack-based VM. The first runnable subset will be narrower than this specification. Bytecode format is internal and unstable. Alternative VM designs, JIT, AOT, Wasm, and CPython performance comparisons are tracked separately; no performance parity is promised.

## 15. Remaining specification work

Before or during implementation, formalize: token/escape tables; Pratt precedence details; complete AST grammar; numeric conversion and Float edge cases; diagnostics schema; collection iteration and Unicode indexing; type compatibility, narrowing and generics; exception hierarchy; module resolution; interfaces; concurrency and ownership. These are **implementation-time refinements**, not reasons to postpone a first lexer.
