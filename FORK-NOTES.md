# Fork notes

Transistor's fork of `Notifuse/liquidgo`, forked at upstream `c0b2d83`.

## Why this fork exists

Transistor renders podcast websites from Liquid themes. Today Ruby renders
them; we're moving rendering to Go, first in `receiver` (the theme development
tool, rendering on the developer's machine) and later in a Go website server.
Go output has to be **byte-identical** to the Ruby Shopify Liquid gem 5.13.0 so
a site looks the same whichever engine renders it.

liquidgo was the closest Go port, but it diverged from Ruby in ways real themes
hit (truthiness, Integer vs Float math, `{% liquid %}` blocks, date parsing and
more), and it had no way for custom filters to read the render context. This
fork fixes those. The plan is to upstream each change as its own PR (below) and
return to upstream liquidgo once they land.

Every fix has a self-contained regression test in `parity/`, each naming the
Ruby semantics it pins down. Expectations are always derived by running the real
gem, never by reasoning about what Ruby "should" do; that has been wrong more
than once.

## Fixes

| Area | Behaviour corrected |
|---|---|
| `split` | `"" \| split: ","` returned `[""]`; trailing empties were kept |
| `{% liquid %}` | tag regex lacked `(?s)`, so multi-line tags were unknown tags |
| `{% liquid %}` | a **blank line silently aborted the rest of the block** — no error |
| filters | keyword args were parsed then never passed to the filter |
| `{% case %}` | `Blank()` read the tag's own nodelist, stripping enclosing whitespace |
| `date` | could not parse Ruby `Time#to_s` (`2024-03-01 09:07:00 -0500`) |
| `strftime` | rewritten as a scanner: `%-d`, `%l`, `%s`, `%:z`, `%U`/`%W` now work, and literal text can no longer collide with a Go layout token |
| `contains` | substring-matched string array elements; Ruby's `Array#include?` is exact |
| `==` | never compared across Integer/Float, so `1 == 1.0` was false |
| arithmetic | `ToNumber` normalises to float64, so Integer op Integer produced a Float; also `sum` and `abs` |
| `divided_by` | Integer/Integer must be integer division, flooring |
| truthiness | `""` was falsy in 4 places; in Ruby only `nil` and `false` are falsy |
| `truncate` | counted bytes, mangling any non-ASCII text or multi-byte ellipsis |
| `escape_once` | re-escaped existing entities (`&amp;` → `&amp;amp;`) |
| `slice` | negative offsets returned empty; Ruby counts back from the end |
| output | arrays render as joined elements (`123`), not `[1, 2, 3]` |
| output | floats render Ruby-style: `9.0`, `1.0e+15`, threshold at exp ≥ 15 or ≤ −5 |
| `sum`, `abs` | returned Floats for Integer input, so `{{ posts \| sum: 'n' }}` rendered `43.0` |
| strict parse | accepted `\| truncate:` and trailing commas; a colon promises an argument |
| `{% for %}` | **a map rendered nothing at all** — Ruby iterates a Hash as `[key, value]` pairs |

## Additions

Engine features Ruby Liquid has that liquidgo lacked. Tests live in
`liquid/strainer_template_test.go`.

| Addition | Why |
|---|---|
| `ContextualFilter` | Ruby filters read `@context` (registers, e.g. the locale for `t`). liquidgo shared one registered filter instance across all renders with no context. A registered filter implementing `WithContext(*liquid.Context) interface{}` now gets a per-render copy bound to that render's context. |
| Filter errors raise | A filter returning a non-nil `error` had it silently dropped. It now raises, rendering `Liquid error: <message>` for Liquid errors and `Liquid error: internal` otherwise, as Ruby does. Standard filters that return errors (`divided_by: 0`, bad base64) now surface them too. |

## Upstream PR plan

One PR per row. The first group leaves upstream's suite green; the second
changes behaviour upstream tests assert, so each needs maintainer agreement
that Ruby is the reference. Each PR brings its test and updates any upstream
test it contradicts.

**Green on upstream's suite**

| PR | Issue | Fix | Tests |
|---|---|---|---|
| `{% liquid %}` tag parsing | Multi-line tags inside `{% liquid %}` were unknown tags, and a blank line silently ended the block | `(?s)` on the full-token regex; skip blank lines instead of treating them as end of input | `parity/liquid_tag_test.go`, `parity/liquid_tag_blank_line_test.go` |
| `case` whitespace | `Blank()` read the tag's own nodelist, stripping whitespace around `{% case %}` | Evaluate blankness from the branches | `parity/case_blank_test.go` |
| Filter keyword arguments | `t: title: x` parsed the keyword args then dropped them | Pass the keyword hash as the last argument | `parity/filter_kwargs_test.go` |
| Strict filter arguments | Strict mode accepted `\| truncate:` and trailing commas | Require an argument after a colon | `parity/strict_filter_args_test.go` |
| `date` input and `strftime` | Couldn't parse Ruby `Time#to_s`; strftime missed `%-d`, `%l`, `%s`, `%:z`, `%U`/`%W` and could collide with Go layout tokens | Parse Ruby's format; rewrite strftime as a directive scanner | `parity/date_parsing_test.go`, `parity/strftime_test.go` |
| String filters | `split` kept trailing empties; `truncate` counted bytes; `escape_once` double-escaped entities | Ruby's `split` semantics; count runes; skip existing entities | `parity/split_test.go`, `parity/string_filters_test.go` |
| `contains` on arrays | String elements substring-matched | Exact match for strings (Ruby `Array#include?`) | `parity/contains_test.go` |
| Value output | Arrays rendered as `[1, 2, 3]`; floats not Ruby-formatted | Join arrays; Ruby float formatting (`9.0`, `1.0e+15`) | `parity/value_rendering_test.go` |
| `for` over a map | A map rendered nothing | Iterate `[key, value]` pairs (sorted; `liquid.OrderedKeys` for insertion order) | `parity/hash_iteration_test.go` |
| `ContextualFilter` | Custom filters couldn't see the render context | Per-render copy via `WithContext` | `TestContextualFilterSeesRenderRegisters` |
| Filter errors | Errors returned by filters were dropped | Raise them; render as `Liquid error: …` | `TestFilterErrorsRenderAsLiquidErrors` |

**Changes behaviour upstream tests assert**

| PR | Issue | Fix | Upstream tests to update |
|---|---|---|---|
| Truthiness | `""` was falsy in 4 places; Ruby treats only `nil` and `false` as falsy | `IsTruthy` used everywhere | `TestIfTagRenderToOutputBufferEmptyStringCondition`, `TestUnlessTagWithEmptyString` |
| `slice` negative offsets | Negative offsets returned empty; Ruby counts back from the end | Ruby slice bounds, by character | `TestStandardFiltersSlice`, `TestSliceCollectionEdgeCases` |
| Integer arithmetic | Integer op Integer returned a Float (`43.0`), `divided_by` didn't floor, `1 == 1.0` was false | Track integer-ness from the inputs; floor division; numeric equality across types | `TestPlus`, `TestMinus`, `TestTimes`, `TestModulo`, `TestDividedBy`, `TestSum`, `TestAbs` |

Upstream's own suite is inconsistent on arithmetic: `TestSum` asserts
`float64(3)`, while `integration/comprehensive_test.go` asserts that
`{{ posts | sum: 'comments_count' }}` renders `"43"`, an Integer. The
integration test is the Ruby-correct one, which is a good argument for the PR.

`performance/` fails to build without the `reference-liquid` submodule; that is
pre-existing and unrelated.

## Known gaps

- **Float arithmetic precision.** Ruby's math filters run through `BigDecimal`,
  so `0.1 | plus: 0.2` is exactly `0.3`; Go does binary float arithmetic. Needs
  decimal arithmetic throughout the math filters. No Transistor theme does float
  arithmetic. Pinned by `TestFloatArithmeticPrecisionIsAKnownGap`.
- **Hash iteration order.** `{% for %}` over a map now yields `[key, value]`
  pairs, but sorted by key: Go maps have no order to recover, while Ruby hashes
  iterate in insertion order. A value implementing `liquid.OrderedKeys` can
  supply Ruby's order, though a plain Go map cannot — that needs a struct, which
  then loses property access. No theme iterates a hash today.
- **Coverage beyond Transistor themes.** Parity is proven for what Transistor
  themes use (10 tags, 19 standard filters). `render`, `capture`, `raw`, `echo`,
  `cycle` and the array/math filters our themes don't use aren't in a golden set
  yet.
