# What this file is for

Working with AI on a codebase creates a problem that doesn't exist when
humans work alone: the code and the *intent* behind it drift apart. One
prompt produces a structure the AI then forgets. A second prompt, written
two weeks later, quietly reshapes it. Public interfaces get renamed,
packages absorb responsibilities that weren't theirs, helper functions
proliferate. Nothing is obviously wrong — every individual change was
reasonable in isolation — but the architecture erodes until no one, human
or AI, has a clear picture of what the codebase actually is.

`spore.yaml` is an attempt to fix this by making the architecture an
explicit artifact. Instead of hoping the code communicates the design, the
design lives in a file that both humans and agents read and edit directly.
An agent that needs to change the program changes the spec. A separate step
generates or updates code from the spec. The code becomes the *derivation*
of intent, not its only record.

# The spec is the source of truth

The load-bearing idea is that `spore.yaml` is authoritative. When the spec
and the code disagree, the spec is correct and the code is wrong — by
definition. That inversion is the whole point. In a codebase-first
workflow, the code is the only durable record of decisions and the spec
(if there is one) is perpetually stale documentation. In a spec-first
workflow, the code is regenerable and the spec is the thing you protect.

This only works if two conditions hold. First, the spec has to capture
enough about the system that regenerating the code from it produces
something functionally equivalent to what you had. Second, the spec has to
capture enough that *two independent regenerations* produce essentially
the same code — because if they don't, every regeneration introduces
churn, and churn destroys the trust that makes the whole loop viable.

# What the spec should contain

The rule for deciding whether something belongs in the spec is simple:
**would a human reviewer notice or care if this changed silently between
two regenerations of the same spec?** If yes, it belongs in the spec. If
no, it belongs in the code.

A package's name, its public functions and their signatures, the types
that cross package boundaries, the interface contracts between components,
the error values callers can depend on, the dependency graph, the
lifecycle of long-running routines, the external channels the app exposes
— all of these produce visible, breaking differences when they change.
They belong in the spec.

The body of a function, the names of local variables, the organization of
helper code within a package, the exact wording of an error message,
whether a loop is a `for range` or a `for i := 0` — none of these are
architectural. Pinning them in the spec turns the spec into pseudocode and
defeats the purpose; the point isn't to specify *how* the code works but
*what contracts it upholds*. The code-writer AI is allowed, and expected,
to make reasonable local decisions.

This principle matters more than any specific list of fields. The template
that accompanies this document shows one way to organize the ideas —
packages, exports, routines, channels, error strategy — but the fields
themselves are not sacred. A project without goroutines has no use for a
`routines` block. A library with no external surface has no use for
`channels`. When you hit a case the current shape doesn't cover, add a
field; when a section goes consistently unused, remove it. The fields
exist in service of the principle, not the other way around.

# Verification closes the loop

Declaring the spec authoritative doesn't make it true. A spec nobody
checks is just a wish. The other half of the system is a set of
verification scripts that read the spec and inspect the code, confirming
that the code actually matches what the spec says.

Every package claimed in the spec must exist on disk. Every exported
function must have the declared signature. Every type listed must be
present with the declared fields. Every interface claimed must be
satisfied by the types that claim to implement it. Every dependency
declared must match the actual import graph — and no undeclared
dependencies may exist. Every sentinel error named must be a package-level
variable at the given path.

When verification fails, one of two things is true: either the code was
changed without updating the spec (fix the code or update the spec), or
the spec was changed without updating the code (run the code-writer, or
do it by hand). Either way, the divergence is surfaced immediately rather
than accumulating invisibly. This is what keeps the spec honest — not
good intentions, but a mechanical check that runs often and fails loudly.

The checks are language-specific because "what counts as a public
function" or "what counts as an import cycle" depends on the language.
That's why the structure includes per-language verify packages: each one
knows how to inspect its language's source code and compare it against
the generic claims the spec makes.

# Evolving the spec itself

The shape of `spore.yaml` is not a fixed standard. It is a working
attempt at an open problem — describing a codebase's architecture in a
form precise enough to regenerate from, readable enough for humans to
maintain, and structured enough for scripts to verify against.

Anyone improving this file should hold the principles tightly and the
format loosely. If a new kind of architectural fact matters for the
systems you're building — event flows, state machines, schedules,
configuration schemas, migration contracts — add a section for it. If a
section has only ever been filled with placeholder values, remove it. If
the template's assumptions about packages or channels don't translate
cleanly to your language, change the assumptions.

The test is always the same: does the spec, when read, let an AI agent
(or a new human contributor) understand the system well enough to change
it safely? Does a regenerated codebase from this spec produce essentially
the code that already exists? Does verification catch the cases where
they diverge? If yes, the spec is doing its job. If no, the spec needs to
grow, shrink, or reshape — regardless of what any previous version looked
like.