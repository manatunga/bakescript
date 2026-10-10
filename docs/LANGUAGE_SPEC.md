# BakeScript Language Specification

**Status:** Draft design specification (pre-implementation)  
**File extension:** `.bs`  
**Implementation target:** Go  
**Initial execution engine:** Bytecode VM

> This document describes intended behavior, not features already implemented.
> **Decided** means explicitly chosen; **Provisional** means accepted as a working
> default but not yet fully specified; **Open** means further design is required.
> Examples are illustrative and are not yet executable.

## 1. Design philosophy [Decided]

BakeScript uses Python-style significant indentation, lowercase keywords, and
conversational English with selective cooking vocabulary. Avoid forcing a
cooking metaphor onto conventional programming concepts when that reduces clarity.

## 2. Lexical structure

- Significant indentation introduces suites after `:`. [Decided]
- Keywords are lowercase. [Decided]
- `#` begins a line comment. [Decided]
- Single and double quotes delimit strings: `'text'`, `"text"`. [Decided]
- Triple quotes delimit multiline strings; a triple-quoted string in a designated
  documentation position can act as a docstring. [Decided; precise positions provisional]
- Interpolated strings use `f"Hello {name}"` or `f'Hello {name}'`. [Decided]
- Identifier rules, tabs versus spaces, Unicode normalization, numeric literal
  formats, escapes, and indentation diagnostics remain to be specified. [Open]

## 3. Values and types

| Concept | Syntax | Status |
|---|---|---|
| Boolean true | `fresh` | Decided |
| Boolean false | `spoilt` | Decided |
| Null-like value | `empty` | Decided |
| Lists | `[1, 2, 3]` | Decided |
| Dictionaries | `{"name": "Alex"}` | Decided |
| Tuples | `(1, 2)` | Decided |
| Sets | `{"a", "b"}` | Decided |
| Empty set | `set()` | Decided |

`empty` is distinct from `0`, `""`, and `[]`. Dictionaries preserve insertion
order; sets are unordered. [Decided]

Conventional constructors/conversions use parentheses, e.g. `int("3")` and
`set(items)`, and are distinct from user-function invocation with `cook`. [Decided]

Precise numeric types, overflow, hashing/equality, tuple grouping versus tuple
construction, and object identity remain to be defined. [Open]

## 4. Variables and assignment [Decided]

```bakescript
prep score as 10        # mutable declaration
score as 20             # reassignment
prep later              # declared but uninitialized
later as 5
freeze maximum as 100   # deeply immutable
```

Reading an uninitialized variable is an error; this is different from a
variable explicitly holding `empty`. `freeze` prevents reassignment and nested
mutation, including mutation through aliases. Exact enforcement mechanics are
not yet specified.

Duplicate `prep` declarations in the same scope are disallowed.

## 5. Optional types and checking [Decided]

Both annotation forms have identical meaning:

```bakescript
prep age of type int as 20
prep score: int as 95

to multiply with a of type int, b of type int -> int:
    serve a * b

to add with a: int, b: int -> int:
    serve a + b
```

- `: Type` and `of type Type` are equivalent for variable declarations and
  function parameters.
- `-> Type` annotates a function's return type.
- `as` remains the initializer/default-value syntax.
- Annotations are optional. `bake check` performs gradual static checking.
- Diagnostics are nonblocking by default, and annotations do not automatically
  impose runtime type checks.
- Grouped annotations such as `(a, b): int` are **not** currently part of the design.

Type compatibility rules, generics, unions, and checker strictness are [Open].

## 6. Expressions and operators

| Category | Operators / forms | Status |
|---|---|---|
| Equality | `is`, `==` | Decided |
| Inequality | `is not`, `!=` | Decided |
| Ordering | `<`, `<=`, `>`, `>=` | Decided |
| Logical AND | `and`, `&&` | Decided |
| Logical OR | `or`, `||` | Decided |
| Logical NOT | `not`, `!` | Decided |
| Arithmetic | `+`, `-`, `*`, `/`, `//`, `%`, `**` | Decided |
| Membership | `in`, `not in` | Provisional |

Logical operators short-circuit. `/` performs floating-point division
(`4 / 2` yields `2.0`); `//` performs floor division (`-5 // 2` yields `-3`).
Conventional operator precedence is the provisional default; its exact table
still needs documenting.

## 7. Conditionals [Decided, with `instead` spelling provisional]

```bakescript
taste temperature > 180:
    pour "Hot"
instead temperature > 100:
    pour "Warm"
otherwise:
    pour "Cool"
```

`instead` is the currently preferred `elif` keyword, but its spelling was
previously marked provisional.

## 8. Iteration [Decided]

Collection iteration uses **`each`**, not `stir`:

```bakescript
each ingredient in ingredients:
    pour ingredient
```

Fixed-count repetition uses `stir`:

```bakescript
stir once:
    pour "One"

stir twice:
    pour "Two"

stir thrice:
    pour "Three"

stir 5 times:
    pour "Five"

stir (2 * 3) times:
    pour "Six"
```

Conditional repetition uses pre-test `until`: execute while its condition is
false, stop once true.

```bakescript
prep temperature as 100
until temperature >= 180:
    temperature as temperature + 10
```

`skip` continues the loop; `stop` breaks it. Fixed-repeat evaluation timing,
iteration mutation rules, and iterator protocol are [Open].

**Superseded:** `stir item in collection:` is retired. Use `each item in collection:`.

## 9. Functions, calls, and scope

```bakescript
to greet with name as "Guest":
    pour f"Hello, {name}!"

to square with x:
    serve x * x

cook greet
cook greet with "Alex"
cook greet with name as "Alex"

prep result as cook square with 5
prep fn as square
```

- Declare with `to`, call with `cook`, return with `serve`. [Decided]
- No-argument declarations/calls omit `with`. [Decided]
- Positional arguments precede named arguments. [Decided]
- Default parameters use `as`. [Decided]
- `*values` captures positional varargs as a tuple; `**options` captures
  named arguments as a dictionary. Unpacking in calls is supported. [Decided]
- Functions are first-class and support lexical closures. [Decided]
- `lambda x: x * 2` is a provisional anonymous-function form.
- Functions have lexical scopes; conditional/loop blocks do not introduce
  separate scopes. `prep` inside a function shadows an outer variable.
  `outer name` enables assignment to an enclosing variable. [Decided]
- Recursion is supported; a function without `serve` implicitly returns
  `empty`. [Provisional]

Call-expression precedence and complex combinations of named/default/varargs
arguments require formal grammar work. [Open]

## 10. Collections, indexing, and mutation

Both conventional and conversational indexing are supported:

```bakescript
prep first as items[0]
prep also_first as items at 0
prep middle as items[1:3]
prep other_middle as slice items from 1 to 3
prep beginning as slice items to 2
prep ending as slice items from 2
```

Indexing is zero-based; slice ends are exclusive. `[start:stop:step]` is a
provisional extension.

```bakescript
pour size of items
toss "Butter" into items
pluck "Sugar" from items
dispose "Flour" from items
```

- `size of` returns length.
- `toss` appends to a list.
- `pluck` removes the first matching value.
- `dispose` removes all matching values.
- Both removal forms error if no match exists; list order is preserved.

Other conventional collection/string operations are provisionally accepted
but not individually specified.

## 11. Classes and inheritance [Decided]

```bakescript
menu Recipe:
    to setup with name:
        my.name as name

    to describe:
        pour my.name

menu Cake based on Recipe:
    to describe:
        cook parent.describe
        pour "Cake!"

prep cake as order Cake with "Chocolate"
cook cake.describe
cook describe of cake
pour cake.name
pour name of cake
```

- `menu` declares a class; `order` instantiates it.
- `setup` is automatically called on construction.
- `my` is implicit instance access; no explicit `self` parameter.
- `menu Child based on Parent:` declares inheritance.
- `cook parent.describe` and `cook describe from parent` are parent-call forms.
- Dot and conversational `of` access are both supported.
- Members are public by default; overriding/inheritance details remain
  provisional.

Method resolution, multiple inheritance, class versus instance attributes,
and initializer chaining remain [Open].

## 12. Errors [Decided; some details provisional]

```bakescript
attempt:
    burn "Something went wrong"
rescue error:
    pour error
always:
    pour "Cleanup"
```

`burn` raises an error. `attempt`/`rescue`/`always` provide error handling.
`always` runs even when control exits via `serve`, `stop`, or `skip`.
Typed exceptions and propagation are provisional; exception hierarchy and
matching grammar remain [Open].

## 13. Modules [Decided; resolution details provisional]

```bakescript
grab maths
grab square, cube from maths
grab mathematics as maths
grab square as sq from mathematics
grab recipes.cakes
```

Every `.bs` file is a module. Relative imports and circular import behavior
are provisionally Python-like, but exact resolution rules are [Open].

## 14. Input, output, and built-ins [Decided]

```bakescript
pour "Hello"
prep name as sip "Your name: "
prep answer as sip
```

`pour` outputs a value. `sip` is an input expression with an optional prompt.
Precise formatting, newline behavior, and EOF handling are [Open].

Basic math, files, dates, and JSON standard-library modules are planned
[Provisional].

## 15. Tooling [Provisional]

- `bake run` — execute programs.
- `bake check` — static analysis and type diagnostics.
- `bake fmt` — format source.
- `bake` REPL — interactive execution.

CLI flags, exit codes, and formatter style remain [Open].

## 16. Compatibility and implementation independence

This specification defines BakeScript semantics, not a particular backend.
The first implementation targets a Go-based bytecode VM; future interpreters,
JITs, AOT compilers, or other backends may implement the same semantics.
Internal bytecode is **not** promised as a stable public format at this stage.

## 17. Outstanding design questions

1. Formal grammar, precedence table, and indentation/tokenization rules.
2. Exact numerical model, integer overflow, and floating-point behavior.
3. Value identity, equality, hashability, and deep-freeze aliasing.
4. Function argument binding and call-expression grammar.
5. Iterator protocol, collection mutation during iteration, and loop scoping.
6. Class inheritance and method-resolution semantics.
7. Exception taxonomy and rescue matching.
8. Module import resolution and initialization.
9. Standard library APIs and user-facing diagnostics.
10. Stability policy for syntax before and after 1.0.
