# GeeBreeze Language Specification

**Status:** Draft, pre-implementation. Examples are illustrative, not executable yet.  
**Working name:** GeeBreeze (availability and naming conflicts not fully verified)  
**Provisional extension:** `.gbz` (formerly `.bs`)  
**Implementation:** Go; initially a bytecode compiler and VM

**Labels:** **Decided** = explicitly agreed syntax or intent; **Provisional** = accepted direction with details to formalize; **Open** = unresolved behavior.

## 1. Design philosophy [Decided]

GeeBreeze is a general-purpose programming language with conversational, intuitive syntax and rigorous semantics. It is **not** unrestricted natural language. Favor natural, concise instructions without sacrificing precision or introducing forced metaphors. Significant indentation and `:` delimit blocks; keywords are lowercase.

## 2. Lexical structure

- `#` starts a comment. [Decided]
- Strings use single or double quotes; triple quotes support multiline strings and docstrings. [Decided; docstring placement provisional]
- Interpolation uses `f"Hello, {name}!"`. [Decided]
- Identifier rules, Unicode handling, indentation errors, tabs, escapes, numeric literal syntax, and semicolon policy need formalization. [Open]

## 3. Values and collections [Decided]

| Concept | Form |
|---|---|
| True / false / null-like | `true`, `false`, `nothing` |
| List | `[1, 2, 3]` |
| Insertion-ordered dictionary | `{"name": "Alex"}` |
| Tuple | `(1, 2)` |
| Set | `{"a", "b"}` |
| Empty set | `set()` |

`nothing` differs from an uninitialized variable. Set/dictionary literal disambiguation and empty `{}` behavior require formal grammar. Numeric model, equality/hashing and identity remain open. Built-in constructors/conversions were previously written as `int("3")` and `set(items)`; their invocation relationship to mandatory `call` **requires review**, not silent grandfathering.

## 4. Declarations, assignment, scope [Decided]

```text
let score be 10
score = 20
score += 5
let later
let count: int
let amount of type float
later = 5
fix maximum at 100
fix pi: float at 3.14
fix threshold of type int at 50
```

- `let` declares a mutable binding; `be` supplies its initializer.
- `fix` declares a deeply immutable binding; `at` supplies its **required** initializer.
- `=` reassigns; compound assignments such as `+=`, `-=`, `*=`, `/=` are supported.
- `let` may be uninitialized; reading it before assignment is an error, distinct from holding `nothing`.
- Static checking should detect possibly uninitialized reads when feasible; runtime checks actual invalid reads.
- Function and conditional/loop blocks create lexical scopes. Assignment can update a binding in an enclosing scope; declaration creates a new binding in its scope.
- Deep immutability prevents nested mutation through aliases; enforcement mechanics remain open.

## 5. Optional types [Decided]

`name: Type` and `name of type Type` are equivalent. Types are optional; gradual static checking is planned. Both `giving Type` and `-> Type` annotate return types.

```text
to multiply with a: int, b of type int giving int:
    give a * b

to square with x: int -> int:
    give x * x
```

Checker strictness, inference, unions, generics, runtime enforcement and compatibility rules are open. Static diagnostics are intended to be nonblocking by default.

## 6. Expressions and operators

| Category | Forms |
|---|---|
| Equality / inequality | `is`, `==`, `is not`, `!=` |
| Ordering | `<`, `<=`, `>`, `>=` |
| Logic | `and`, `&&`, `or`, `||`, `not`, `!` |
| Arithmetic | `+`, `-`, `*`, `/`, `//`, `%`, `**` |
| Membership | `in`, `not in` |
| Grouping | `(expression)` |

`is` means **value equality**, not object identity. Logical operators short-circuit. `/` is floating-point division; `//` floors (`-5 // 2` is `-3`). Exact precedence and operator behavior for custom types need specification.

## 7. Branching [Decided]

```text
if temperature > 180:
    say "Hot"
else if temperature > 100:
    say "Warm"
otherwise:
    say "Cool"
```

## 8. Loops [Decided]

```text
each fruit in fruits:
    say fruit

repeat 5 times:
    say "Again"

repeat once:
    say "Once"

repeat twice:
    say "Twice"

repeat thrice:
    say "Thrice"

while temperature < 180:
    temperature += 10

until ready:
    say "Waiting"
```

`while` runs while true; `until` runs while false. `skip` continues and `stop` breaks the innermost loop. Repeat-count evaluation, mutation during iteration, iterator protocol, and zero/negative counts need precise rules.

## 9. Functions, calls and closures [Decided unless noted]

```text
to greet with name as "Guest":
    say f"Hello, {name}!"

to add with a: int, b: int giving int:
    give a + b

call greet
call greet with "Alex"
call greet with name as "Alex"
let result be call add with 2, 3
let fn be add

let double be with x give x * 2
```

- `to` declares; `call` invokes; `give` returns. Bare `give` returns without a value; a function falling through returns `nothing` [Provisional].
- `with` introduces parameters or arguments; omitted for zero-argument functions.
- Default parameters use `as`; named arguments use `name as value`.
- Positional arguments precede named arguments.
- Functions are first-class; referencing `add` does **not** call it.
- Lexical closures capture outer bindings. `outer count` marks intention to mutate a captured variable; reads need no `outer`.
- Variadic positional parameter `...items` collects a list; `...numbers: int` constrains each element.
- Variadic named parameter `named ...options` collects a dictionary.
- Unpacking: `call sum with ...values` and `call configure with named ...options`.
- Parameter order: required, defaulted, positional variadic, named variadic.
- Recursion supported.

### Grouped nested calls [Decided]

```text
call display with (call calculate with 10, 20), 30
let result be (10 + 5) * 2
```

Parentheses disambiguate nested calls and expression grouping. The exact parser boundary for unparenthesized calls, anonymous-function precedence, and builtin call forms remains **Open**.

## 10. Collections, indexing and mutation

```text
let first be items[0]
let also_first be items at 0
let middle be items[1:3]
let other_middle be slice items from 1 to 3
say size of items

add "orange" to fruits
remove "banana" from fruits
remove every "apple" from fruits
```

- Zero-based indexing; slices end-exclusive. `[start:stop:step]` retained provisionally.
- `size of` returns length.
- `add value to collection` adds an element; initially specified for lists as append.
- `remove value from collection` removes the first matching occurrence.
- `remove every value from collection` removes all matching occurrences.
- Previous behavior that removal with no match errors is retained **provisionally**, pending collection-specific error rules.
- Other collection types' `add`/`remove` semantics remain open.
- `slice items from 1 to 3` and other conversational indexing forms are retained from the original draft.

## 11. Classes and inheritance [Decided syntax]

```text
define Person:
    to setup with name:
        my.name = name

    to greet:
        say f"Hello, I'm {my.name}!"

define Student from Person:
    to setup with name, school:
        call parent.setup with name
        my.school = school

let student be create Student with "Alex", "IIT"
call student.greet
say student.name
```

`define` declares a class; `define Child from Parent:` expresses inheritance. `create` constructs an instance and calls `setup`. `my` denotes the current instance; `parent` accesses inherited behavior. Dot member access is supported. Previously accepted `name of object` access is retained **provisionally** pending a grammar audit. Public-by-default members and overriding are provisional; MRO, class attributes, initialization and multiple inheritance are open.

## 12. Errors [Decided syntax]

```text
try:
    raise "Something went wrong"
catch error:
    say error
always:
    say "Cleanup"
```

`raise` raises an error; `try` / `catch` / `always` handle errors and cleanup. `always` runs on normal and exceptional exits, including `give`, `skip`, and `stop` (subject to runtime termination rules). Typed catches, hierarchy, matching, rethrow and propagation remain open.

## 13. Modules [Decided syntax]

```text
bring maths
bring square, cube from maths
bring mathematics as maths
bring square as sq from mathematics
bring recipes.cakes
```

Each `.gb` source file is intended as a module. Import aliases use `as`; inheritance and imports both use `from` in unambiguous contexts. Relative/circular imports and module initialization remain provisional pending precise resolution rules.

## 14. Input/output and standard library

```text
say "Hello"
let name be ask "Your name: "
let answer be ask
```

`say` outputs; `ask` reads input with optional prompt. Exact newline, formatting and EOF behavior are open. Math, files, dates and JSON modules are planned.

## 15. Tooling [Provisional]

```text
geebreeze run hello.gb
geebreeze check hello.gb
geebreeze fmt hello.gb
geebreeze build hello.gb
geebreeze
```

`run`, `check`, `fmt` and REPL are planned. `build` is a **future-facing placeholder**, not a promise that native compilation is present in the initial VM release. CLI behavior and packaging are open.

## 16. Implementation independence and compatibility

The initial backend is a Go-based stack bytecode VM, but syntax/semantics should not depend on that implementation. A register VM review is required after the first working VM and benchmark harness (target around v0.2.0). JIT, AOT, Wasm and concurrency research remain tracked separately in `TECHNICAL_BACKLOG.md`. CPython-level performance or better is an **aspiration, not a claim**. Internal bytecode is not a stable public contract.

## 17. Open grammar and semantics work

1. Formal EBNF/Pratt grammar, precedence, nested calls and constructors.
2. Block-scoped declarations, assignment resolution and `outer` edge cases.
3. Collection literals, identity, equality, hashing, numeric model and deep immutability.
4. Typed/default/variadic argument binding and evaluation order.
5. Iterators, slice variants and collection mutation/error rules.
6. Class initialization, inherited member lookup, `of` member syntax and method resolution.
7. Exception taxonomy, typed matching and cleanup control-flow precedence.
8. Module initialization, relative imports, circular imports and packaging.
9. String interpolation/escaping, Unicode and diagnostics.
10. CLI, stdlib, formatter and version stability.
