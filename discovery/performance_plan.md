# Performance Improvement Plan: dash-mpd

This plan turns the findings of `discovery/performance_report.md` into a
concrete, ordered set of implementable tasks. It is written for a junior
developer: each task lists what to change, which files and line ranges
to touch, how to verify correctness, and how to measure the impact.

The work is split into 11 independent steps. Land them as separate PRs
(or at least separate commits) in the order given, so each change can be
measured against the previous one. **Never land more than one
performance-impacting change per commit** — otherwise you can't tell
which one helped.

---

## 0. Preflight: set up a repeatable measurement workflow

Before touching any code, make sure you can reliably measure the
benchmarks. Every step below is gated on "did the numbers get better?".

1. Create a working branch:

   ```sh
   git switch -c perf-h1-through-hN   # pick a name per step
   ```

2. Record the baseline on `main` (or current `perf-audit`) **once**:

   ```sh
   git switch main
   go test ./mpd -run=^$ -bench=Unmarshal\|Marshal \
       -benchmem -benchtime=5s -count=10 \
       > /tmp/dash-mpd-baseline.txt
   ```

   The `-count=10` gives `benchstat` enough samples to compute a mean
   and confidence interval. `-run=^$` disables unit tests during the
   bench run. We exclude `BenchmarkClone` from `-bench=` because it
   currently crashes (see step 1).

3. Install `benchstat` if you don't already have it:

   ```sh
   go install golang.org/x/perf/cmd/benchstat@latest
   ```

4. For every subsequent change, generate a new result file and compare:

   ```sh
   go test ./mpd -run=^$ -bench=Unmarshal\|Marshal \
       -benchmem -benchtime=5s -count=10 \
       > /tmp/dash-mpd-after.txt
   benchstat /tmp/dash-mpd-baseline.txt /tmp/dash-mpd-after.txt
   ```

5. Correctness gate on every step: `go test ./...` must pass. In
   particular `TestDecodeEncodeMPDs` (in `mpd/mpd_test.go:22`) is the
   round-trip regression test — it parses every MPD in `mpd/testdata/`
   and re-serialises it, checking that the result is semantically
   identical. It will catch almost any behaviour regression in the
   `xml/` package.

Keep the baseline file around for the whole project; the "cumulative
effect" at the end is measured against it, not against the previous
step.

---

## 1. Fix `BenchmarkClone` so the full bench suite runs (H10)

**Why first:** without this, `go test -bench=.` aborts with a stack
overflow, and you can't run all three benchmarks in the same command.
Not a performance win by itself, but it unblocks everything else.

**File:** `mpd/mpd_test.go:118-129`

**Current code:**

```go
func BenchmarkClone(b *testing.B) {
    data, err := os.ReadFile("testdata/schema-mpds/example_G15.mpd")
    require.NoError(b, err)
    mpd := m.MPD{}
    err = xml.Unmarshal(data, &mpd)
    require.NoError(b, err)
    mpdCopy := m.Clone(&mpd)
    cmp.Equal(&mpd, mpdCopy)          // <-- crashes: infinite recursion via parent back-pointer
    for i := 0; i < b.N; i++ {
        _ = m.Clone(&mpd)
    }
}
```

**Steps:**

1. Delete the `mpdCopy := m.Clone(&mpd)` and `cmp.Equal(&mpd, mpdCopy)`
   lines. They are outside the timed loop and do nothing useful. They
   are the direct cause of the stack overflow: `cmp.Equal` walks the
   whole struct tree and follows `parent *MPD` pointers set by
   `SetParents()`, which creates infinite recursion.
2. The benchmark does **not** call `mpd.SetParents()`, which means
   `m.Clone` runs against a tree without parent back-pointers. That's
   fine for measuring clone cost. Leave it off.
3. If the `go-cmp` import becomes unused, remove it from the import
   block.

**Result after:**

```go
func BenchmarkClone(b *testing.B) {
    data, err := os.ReadFile("testdata/schema-mpds/example_G15.mpd")
    require.NoError(b, err)
    mpd := m.MPD{}
    err = xml.Unmarshal(data, &mpd)
    require.NoError(b, err)
    for i := 0; i < b.N; i++ {
        _ = m.Clone(&mpd)
    }
}
```

**Verify:**

```sh
go test ./mpd -run=^$ -bench=Clone -benchmem -benchtime=1s
```

Should report a number instead of crashing. Commit.

---

## 2. Rewrite `mpd.ParseDuration` to remove regex (H1)

**Why:** The current implementation runs the same regex twice per call
(`Match` then `FindStringSubmatch`), copies the input to `[]byte`,
allocates a 5-element `[]string` plus internal regex state, and then
makes another round of allocations via `strings.TrimRight` +
`strconv.Atoi`. Profile share: ~3.4 % of unmarshal alloc_space on a
small MPD; scales with attribute count on live MPDs.

**File:** `mpd/duration.go:180-231`

**Grammar we must parse** (XSD duration subset that this library
accepts):

```
P[nD][T[nH][nM][n(.n)?S]]
```

- Must start with `P`.
- Optionally `nD` (days).
- If any of hours/minutes/seconds appear, they must be after `T`.
- `H`, `M`, `S` appear in that order, each optional, each preceded by
  an integer. Seconds may be fractional.
- `-` is rejected (the existing code does this; keep that behaviour).

**Implementation steps:**

1. Delete the regex constants (`rStart`, `rDays`, `rTime`, `rHours`,
   `rMinutes`, `rSeconds`, `rEnd`) and the `xmlDurationRegex` variable
   at `mpd/duration.go:36-46`. Also remove the `regexp` import.
2. Replace the body of `ParseDuration` with a single-pass scanner. The
   scanner walks the string, reading runs of digits (and optional
   `.` for seconds), then expects the corresponding designator letter.
   Each designator letter must only appear once, and they must appear
   in the order `D`, `H`, `M`, `S`.

3. Suggested structure:

   ```go
   func ParseDuration(str string) (time.Duration, error) {
       if len(str) < 3 {
           return 0, newParseDurationError("at least one number and designator are required")
       }
       if str[0] != 'P' {
           return 0, newParseDurationError("duration must be in the format: P[nD][T[nH][nM][nS]]")
       }

       var (
           total     time.Duration
           i         = 1         // cursor, positioned just after 'P'
           afterT    = false     // have we seen 'T' yet?
           seenAny   = false     // for the "at least one number" check
           lastUnit  byte        // tracks order; must strictly increase
       )

       // Unit ordering, used to enforce "D before H before M before S".
       rank := func(c byte) int {
           switch c {
           case 'D': return 1
           case 'H': return 2
           case 'M': return 3
           case 'S': return 4
           }
           return 0
       }

       for i < len(str) {
           c := str[i]
           if c == 'T' {
               if afterT {
                   return 0, newParseDurationError("duplicate T in duration")
               }
               afterT = true
               i++
               continue
           }
           if c == '-' {
               return 0, newParseDurationError("duration cannot be negative")
           }
           // Must now read digits (and maybe '.', only if looking for S).
           start := i
           sawDot := false
           for i < len(str) {
               b := str[i]
               if b >= '0' && b <= '9' {
                   i++
                   continue
               }
               if b == '.' && !sawDot {
                   sawDot = true
                   i++
                   continue
               }
               break
           }
           if i == start {
               return 0, newParseDurationError("expected digit in duration")
           }
           if i == len(str) {
               return 0, newParseDurationError("number without unit in duration")
           }
           unit := str[i]
           i++
           if rank(unit) == 0 {
               return 0, newParseDurationError("unexpected character in duration")
           }
           if rank(unit) <= rank(lastUnit) {
               return 0, newParseDurationError("duration units out of order")
           }
           if unit != 'D' && !afterT {
               return 0, newParseDurationError("H/M/S must follow T")
           }
           if sawDot && unit != 'S' {
               return 0, newParseDurationError("decimal only allowed on S")
           }
           lastUnit = unit
           seenAny = true

           numStr := str[start:i-1] // the slice of digits (and maybe '.') without the unit letter

           switch unit {
           case 'D':
               n, err := strconv.ParseUint(numStr, 10, 64)
               if err != nil {
                   return 0, newParseDurationError("error parsing Days")
               }
               total += time.Duration(n) * 24 * time.Hour
           case 'H':
               n, err := strconv.ParseUint(numStr, 10, 64)
               if err != nil {
                   return 0, newParseDurationError("error parsing Hours")
               }
               total += time.Duration(n) * time.Hour
           case 'M':
               n, err := strconv.ParseUint(numStr, 10, 64)
               if err != nil {
                   return 0, newParseDurationError("error parsing Minutes")
               }
               total += time.Duration(n) * time.Minute
           case 'S':
               f, err := strconv.ParseFloat(numStr, 64)
               if err != nil {
                   return 0, newParseDurationError("error parsing Seconds")
               }
               total += time.Duration(f * float64(time.Second))
           }
       }
       if !seenAny {
           return 0, newParseDurationError("duration must contain at least one value")
       }
       return total, nil
   }
   ```

   Notes:
   - `strconv.ParseUint` / `ParseFloat` accept a `string`, not
     `[]byte`, but `str[start:i-1]` is a substring of the input; no
     allocation.
   - No `strings.TrimRight` calls; we slice the digits out directly.
   - The error messages match the previous ones as closely as
     possible — tests may assert against them. If
     `go test ./...` fails because a test asserts a specific error
     text, copy the old text.

4. **Unit-test the new parser aggressively.** The existing tests cover
   the happy path. Add at least:

   - `"PT0S"` → 0
   - `"P1D"` → 24h
   - `"PT1H"`, `"PT1M"`, `"PT1S"` → each unit
   - `"PT0.5S"` → 500ms
   - `"P1DT2H3M4.5S"` → full combo
   - Error cases: `""`, `"P"`, `"PT"`, `"P1M"` (M before T),
     `"PT1.5M"` (decimal not on S), `"P1H"` (H before T), `"P-1D"`,
     `"P1D2H"` (H before T), `"P1Y"` (unsupported unit), `"P1S1H"`
     (order), `"P1D1D"` (duplicate).

**Verify:**

```sh
go test ./mpd -run=Duration      # all duration tests
go test ./...                    # full suite, especially TestDecodeEncodeMPDs
benchstat /tmp/dash-mpd-baseline.txt /tmp/dash-mpd-after.txt
```

Expected: ~2-3 % drop in `BenchmarkUnmarshal` allocs, bigger win on
MPDs with many `xs:duration` attributes.

---

## 3. Pre-size the `[]Attr` slice in the parser (H2.1)

**Why:** `xml.rawToken` starts each element with `attr = []Attr{}` (a
zero-length slice with zero capacity) and grows it via `append`. Each
growth copies the backing array. Profile share of the grow path: ~12 %
of unmarshal alloc_space.

**File:** `xml/xml.go:796-867` (start-element parsing)

**Steps:**

1. Find line 811: `attr = []Attr{}`
2. Change to: `attr = make([]Attr, 0, 8)`

   `8` is chosen because most DASH elements have ≤8 attributes.
   Elements with more will still grow correctly — this only changes
   the initial capacity.

That's the whole change. One line.

**Verify:**

```sh
go test ./...
benchstat ...
```

Expected: small drop in alloc_objects / alloc_space for Unmarshal.
Probably 5–10 % allocs.

---

## 4. Pre-size `start.Attr` in the writer (H3.2)

**Why:** Symmetrical to step 3, for marshal. `marshalAttr` appends
into `start.Attr` one attribute at a time, causing the same grow
churn.

**File:** `xml/marshal.go` — the function that builds the
`StartElement` before calling `marshalAttr` for each field. Look at
`marshalValue` (roughly `xml/marshal.go:400-500`) and `defaultStart`
(`xml/marshal.go:692` onwards).

**Steps:**

1. Open `xml/typeinfo.go`. Add a method (or inline compute) that
   returns the count of `fAttr` fields for a given `*typeInfo`:

   ```go
   func (ti *typeInfo) numAttrFields() int {
       n := 0
       for i := range ti.fields {
           if ti.fields[i].flags&fMode == fAttr {
               n++
           }
       }
       return n
   }
   ```

   Cache it on the `typeInfo` struct so it's computed once per type:

   ```go
   type typeInfo struct {
       xmlname  *fieldInfo
       fields   []fieldInfo
       nAttrs   int    // NEW: populated at the end of getTypeInfo
   }
   ```

   In `getTypeInfo`, after the fields slice is fully built, walk it
   once and set `tinfo.nAttrs`.

2. In `xml/marshal.go`, find every place `start := StartElement{…}` or
   `start.Attr = append(start.Attr, startTemplate.Attr...)` runs
   before the attribute loop. The call site you want is in
   `marshalValue`, just before the `for i := range tinfo.fields` loop
   that dispatches `fAttr`. Pre-size:

   ```go
   if cap(start.Attr) < tinfo.nAttrs {
       start.Attr = make([]Attr, 0, tinfo.nAttrs)
   }
   ```

   Do **not** unconditionally overwrite — `defaultStart` may already
   have appended template attrs. Grow, don't clobber.

**Verify:**

```sh
go test ./...
benchstat ...
```

Expected: 5–10 % drop in Marshal alloc_objects.

---

## 5. Map-based attribute lookup on `*typeInfo` (H8)

**Why:** In `xml/read.go:458-487` each incoming attribute is matched
against `tinfo.fields` by linear scan. Types like
`AdaptationSetType` and `RepresentationType` have ~30 fields. Each
attribute is therefore O(fields). Profile: ~5 % CPU. This moves up
in relative cost once we fix the allocation-heavy paths above.

**File:** `xml/typeinfo.go` + `xml/read.go:458-487`

**Steps:**

1. Extend `typeInfo`:

   ```go
   type typeInfo struct {
       xmlname *fieldInfo
       fields  []fieldInfo
       nAttrs  int                      // (from step 4)
       attrs   map[attrKey]*fieldInfo   // NEW
       anyAttr *fieldInfo               // NEW: first fAny|fAttr field, nil if none
   }

   type attrKey struct {
       space string
       local string
   }
   ```

2. Populate `attrs` and `anyAttr` at the end of `getTypeInfo`, once the
   `fields` slice is final (so the `*fieldInfo` pointers we store are
   stable — **important: store pointers into the final slice, not into
   an in-progress one**):

   ```go
   tinfo.attrs = make(map[attrKey]*fieldInfo, tinfo.nAttrs)
   for i := range tinfo.fields {
       fi := &tinfo.fields[i]
       switch fi.flags & fMode {
       case fAttr:
           // Some fields apply to any namespace (fi.xmlns == "");
           // we register them under ("", fi.name). The lookup in
           // read.go must first try (a.Name.Space, fi.name), then
           // fall back to ("", fi.name).
           tinfo.attrs[attrKey{fi.xmlns, fi.name}] = fi
       case fAny | fAttr:
           if tinfo.anyAttr == nil {
               tinfo.anyAttr = fi
           }
       }
   }
   ```

3. Rewrite the inner loop in `xml/read.go:458-487`:

   ```go
   for _, a := range start.Attr {
       fi, ok := tinfo.attrs[attrKey{a.Name.Space, a.Name.Local}]
       if !ok {
           // A field registered with empty namespace matches any namespace.
           fi, ok = tinfo.attrs[attrKey{"", a.Name.Local}]
       }
       if ok {
           strv := fi.value(sv, initNilPointers)
           if err := d.unmarshalAttr(strv, a); err != nil {
               return err
           }
           continue
       }
       if tinfo.anyAttr != nil {
           strv := tinfo.anyAttr.value(sv, initNilPointers)
           if err := d.unmarshalAttr(strv, a); err != nil {
               return err
           }
       }
   }
   ```

4. **Critical pitfall:** the original code allowed
   `finfo.xmlns == "" || finfo.xmlns == a.Name.Space` (two cases).
   Preserve both — the two-step lookup above covers them.

**Verify:** `TestDecodeEncodeMPDs` is the backstop. If any MPD in
`mpd/testdata/` deserialises differently, the lookup logic is wrong.

Expected: 3–5 % drop in Unmarshal CPU, no allocation change.

---

## 6. Cache `Implements` checks per type (H9)

**Why:** `xml/marshal.go:473-493` asks the reflect system four
questions per value ("does this implement Marshaler? does its
pointer? does it implement TextMarshaler? does its pointer?"). These
are deterministic given the type. Profile: ~7 % of marshal CPU.

**File:** `xml/typeinfo.go` + `xml/marshal.go:473-493`

**Steps:**

1. Extend `typeInfo`:

   ```go
   type typeInfo struct {
       // ... existing ...
       implementsMarshaler        bool
       implementsAddrMarshaler    bool
       implementsMarshalerAttr    bool
       implementsAddrMarshalerAttr bool
       implementsTextMarshaler    bool
       implementsAddrTextMarshaler bool
   }
   ```

2. Populate at the end of `getTypeInfo`:

   ```go
   tinfo.implementsMarshaler         = typ.Implements(marshalerType)
   tinfo.implementsAddrMarshaler     = reflect.PointerTo(typ).Implements(marshalerType)
   tinfo.implementsMarshalerAttr     = typ.Implements(marshalerAttrType)
   tinfo.implementsAddrMarshalerAttr = reflect.PointerTo(typ).Implements(marshalerAttrType)
   tinfo.implementsTextMarshaler     = typ.Implements(textMarshalerType)
   tinfo.implementsAddrTextMarshaler = reflect.PointerTo(typ).Implements(textMarshalerType)
   ```

   Note: `getTypeInfo` currently only runs for struct types. You need
   to ensure these booleans are available for **every** type the
   marshaller inspects, including non-struct types used as attribute
   values (e.g. `Duration`). The simplest route is a separate, thinner
   cache keyed on `reflect.Type`:

   ```go
   type typeImpl struct {
       marshaler, addrMarshaler,
       marshalerAttr, addrMarshalerAttr,
       textMarshaler, addrTextMarshaler bool
   }
   var typeImplMap sync.Map // map[reflect.Type]typeImpl

   func getTypeImpl(t reflect.Type) typeImpl {
       if v, ok := typeImplMap.Load(t); ok {
           return v.(typeImpl)
       }
       pt := reflect.PointerTo(t)
       v := typeImpl{
           marshaler:        t.Implements(marshalerType),
           addrMarshaler:    pt.Implements(marshalerType),
           marshalerAttr:    t.Implements(marshalerAttrType),
           addrMarshalerAttr: pt.Implements(marshalerAttrType),
           textMarshaler:    t.Implements(textMarshalerType),
           addrTextMarshaler: pt.Implements(textMarshalerType),
       }
       typeImplMap.Store(t, v)
       return v
   }
   ```

3. In `marshalValue` (around `xml/marshal.go:473-493`) and
   `marshalAttr` (around `xml/marshal.go:608-652`), replace
   `val.Type().Implements(marshalerType)` with
   `getTypeImpl(val.Type()).marshaler`, and similarly for the other
   three checks. The `val.CanInterface()` / `val.CanAddr()` guards
   stay as they are — those are fast and depend on runtime state, not
   type identity.

**Verify:**

```sh
go test ./...
benchstat ...
```

Expected: 3–5 % drop in Marshal CPU, no allocation change.

---

## 7. Bypass `Attr` construction for simple attributes (H3.1 + H3.3)

**Why:** `marshalAttr` materialises a whole `Attr{name, string}` and
appends it to `start.Attr` even for plain numeric fields. For the
numeric branch it also calls `strconv.FormatInt/FormatUint/
FormatFloat`, each of which allocates a Go string. Together this is
~75 % of marshal allocations.

**This step is the single biggest marshal win, and also the most
invasive.** Do it after steps 3–6 so you're landing it on a stable
base.

**File:** `xml/marshal.go` (`marshalAttr` at lines 606-688,
`marshalSimple` at lines 883-916, `writeStart` at lines 746-846)

**Steps:**

1. Read `writeStart` end-to-end. Understand how it currently iterates
   `start.Attr` and writes each as `name="value"` to the buffered
   writer. The goal of this step is to let `marshalAttr` write
   directly to that buffer for numeric/string attributes, while still
   falling back to the `start.Attr` list for custom
   `MarshalerAttr`/`TextMarshaler` paths.

2. Split `marshalValue`'s attribute loop into two passes:

   - **Pass A** (existing): iterate fields, but for each attribute
     decide whether it has a custom marshaller. If it does, call it
     and append the returned `Attr` to `start.Attr`. If it does *not*,
     skip — we'll handle it in Pass B.
   - In between, call `writeStart(start)` once, which writes the
     element opening plus the accumulated `start.Attr` (custom ones
     only).
   - **Pass B** (new): iterate fields again, this time for every
     "simple" attribute, write `name="value"` directly to
     `p.Writer` using a new helper:

     ```go
     func (p *printer) writeSimpleAttr(prefix, local string, val reflect.Value) error {
         p.Writer.WriteByte(' ')
         if prefix != "" {
             p.WriteString(prefix)
             p.WriteByte(':')
         }
         p.WriteString(local)
         p.WriteString(`="`)
         if err := p.writeSimpleAttrValue(val); err != nil {
             return err
         }
         p.WriteByte('"')
         return nil
     }
     ```

   The element's `>` must be written after both passes complete.

3. Implement `writeSimpleAttrValue`. This is essentially a copy of
   `marshalSimple` but it writes directly to the buffered writer
   instead of returning `(string, []byte, error)`. For integers and
   floats, use the `strconv.AppendInt` / `AppendUint` /
   `AppendFloat` / `AppendBool` family with a stack-local scratch
   buffer:

   ```go
   var scratch [64]byte
   switch val.Kind() {
   case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
       b := strconv.AppendInt(scratch[:0], val.Int(), 10)
       _, err := p.Writer.Write(b)
       return err
   case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
       b := strconv.AppendUint(scratch[:0], val.Uint(), 10)
       _, err := p.Writer.Write(b)
       return err
   case reflect.Float32, reflect.Float64:
       b := strconv.AppendFloat(scratch[:0], val.Float(), 'g', -1, val.Type().Bits())
       _, err := p.Writer.Write(b)
       return err
   case reflect.Bool:
       b := strconv.AppendBool(scratch[:0], val.Bool())
       _, err := p.Writer.Write(b)
       return err
   case reflect.String:
       return EscapeText(p.Writer, []byte(val.String()))
       // or a direct string-escape helper if one exists
   // ... other kinds as needed ...
   }
   ```

   `EscapeText` already exists in the package for XML-escaping.

4. **XML escaping for attribute values:** the current code runs string
   attribute values through `EscapeText` or its equivalent before
   writing. Do **not** skip this — an attribute containing `<`, `&`,
   `"` etc. must be escaped. Re-use the existing escape routine.

5. Keep the old `start.Attr`-based path alive for:
   - Any type implementing `MarshalerAttr` (e.g. `mpd.Duration`).
   - Any type implementing `encoding.TextMarshaler`.
   - Fields tagged `,any,attr`.

   These emit attributes through their custom method, and the result
   can't be written directly without the `Attr` intermediate (the
   method returns one). Use the `typeImpl` cache from step 6 to decide
   cheaply.

6. **This change requires a thorough round-trip test run.** Run
   `TestDecodeEncodeMPDs` (which parses every `.mpd` under
   `mpd/testdata/` and checks semantic equality after re-encoding)
   after each sub-change. If output byte-order for attributes changes
   (we now emit them in two passes), `xmltree.Equal` still treats the
   XML as equal — attribute ordering is not semantically significant.

**Verify:**

```sh
go test ./...
benchstat ...
```

Expected: ~50 % drop in Marshal alloc_objects and ~50 % in ns/op.

**If anything breaks:** revert and land it as a smaller series of
changes, e.g. first just move numeric attrs to the new path, then
strings, then booleans.

---

## 8. Intern attribute / element names on the Decoder (H2.2)

**Why:** In any given MPD, the strings `"Period"`, `"AdaptationSet"`,
`"Representation"`, `"id"`, `"bandwidth"`, `"mimeType"`, etc. appear
hundreds or thousands of times. Each appearance currently allocates a
fresh Go string via `string(d.buf.Bytes())` in `Decoder.nsname` and
related helpers. A per-Decoder intern map collapses those into one
allocation per unique name.

**File:** `xml/xml.go` (`Decoder` struct, `nsname`, `name`, and the
places where `string(d.buf.Bytes())` appears)

**Steps:**

1. Add a field to `Decoder`:

   ```go
   nameIntern map[string]string
   ```

   Initialise lazily (first use creates it). Reset it in
   `Decoder.Reset()` if such a function exists, or add a comment that
   the intern map lives for the Decoder's lifetime.

2. Add a helper:

   ```go
   func (d *Decoder) internBytes(b []byte) string {
       if d.nameIntern == nil {
           d.nameIntern = make(map[string]string, 128)
       }
       if s, ok := d.nameIntern[string(b)]; ok { // map lookup with []byte key — Go special-case, no alloc
           return s
       }
       s := string(b)
       d.nameIntern[s] = s
       return s
   }
   ```

   The `d.nameIntern[string(b)]` form is specifically optimised by the
   Go compiler since 1.5: the `string(b)` conversion in a map key
   doesn't actually allocate. Look it up to confirm for the Go version
   in use.

3. Replace every occurrence of `string(d.buf.Bytes())` that is
   returning a **name** (tag name, attribute name, namespace prefix
   alias) with `d.internBytes(d.buf.Bytes())`. Do **not** intern
   attribute *values* or character data — those have unbounded
   cardinality and would leak.

   Specifically look at:
   - `d.nsname()` → name comes from `d.buf.Bytes()`.
   - `d.name()` → same.
   - The namespace URI → technically also bounded in an MPD; can be
     interned.

4. **Bound the intern map.** Because this is per-Decoder and a
   Decoder is normally short-lived (one `Unmarshal` call), unbounded
   growth is not a concern. If a caller holds onto a Decoder and feeds
   it a stream of unrelated documents, the map grows. Add a simple
   cap:

   ```go
   const maxInternEntries = 1024
   ...
   if len(d.nameIntern) < maxInternEntries {
       d.nameIntern[s] = s
   }
   return s
   ```

   If over the cap, return the fresh string without storing — graceful
   degradation.

**Verify:**

```sh
go test ./...
benchstat ...
```

Expected: 10–15 % drop in Unmarshal alloc_objects on any real MPD (the
savings scale with element/attribute repetition).

---

## 9. Pool and reuse `xml.Encoder` instances (H5)

**Why:** `xml.NewEncoder` calls `bufio.NewWriter` which allocates a
4 KiB buffer every call. For a service marshaling many MPDs per
second that's 4 KiB of garbage per call. Profile: ~2 % of Marshal
allocs *per call* — bigger in absolute terms when you do many.

**File:** `xml/marshal.go:141-145`, plus a new pool helper

**Steps:**

1. Add a `sync.Pool` of `*Encoder` at package scope. Pool elements
   should have their `bufio.Writer` reset to a scratch
   `*bytes.Buffer` on release, then re-targeted at the caller's
   `io.Writer` on acquire:

   ```go
   var encoderPool = sync.Pool{
       New: func() any {
           return &Encoder{printer{Writer: bufio.NewWriter(io.Discard)}}
       },
   }
   ```

2. Expose a new API: `AcquireEncoder(io.Writer) *Encoder` and
   `ReleaseEncoder(*Encoder)`. The acquirer calls
   `enc.Writer.Reset(w)` before returning the encoder; the releaser
   resets internal state (`elements`, `depth`, `prefix`, `putNewline`,
   any maps) to zero and puts the encoder back. Do **not** expose the
   pool — only the acquire/release pair — so callers can't accidentally
   use an encoder twice.

3. **Do not change `NewEncoder`'s signature or semantics.** Keep it
   allocating a fresh one; there are callers outside this repo.

4. In `mpd/io.go`, update `Write` / `WriteToString` to use the pooled
   path. See step 10 for the recommended shape, since these two
   changes combine naturally.

**Verify:**

```sh
go test ./...
benchstat ...
```

Expected: noticeable drop in Marshal allocs/op (the 4 KiB buffer is
the single largest allocation per Marshal call).

---

## 10. Stream `mpd.(*MPD).Write` directly, no intermediate buffer (H6)

**Why:** `Write` currently calls `xml.MarshalIndent` which returns a
`[]byte`, then copies that `[]byte` to the caller's `io.Writer`.
`WriteToString` adds a third copy. For a 500 KB live MPD we hold a
500 KB buffer we don't need.

**File:** `mpd/io.go:42-71`

**Steps:**

1. Rewrite `Write` to stream:

   ```go
   func (m *MPD) Write(w io.Writer, indent string, withHeader bool) (int, error) {
       cw := &countingWriter{w: w}
       if withHeader {
           if _, err := io.WriteString(cw, xml.Header); err != nil {
               return cw.n, err
           }
       }
       enc := xml.AcquireEncoder(cw)
       defer xml.ReleaseEncoder(enc)
       enc.Indent("", indent)
       if err := enc.Encode(m); err != nil {
           return cw.n, err
       }
       if err := enc.Flush(); err != nil {   // make sure bufio flushes into cw
           return cw.n, err
       }
       return cw.n, nil
   }
   ```

   Add a small `countingWriter` (wraps an `io.Writer` and counts bytes
   written) in the same file.

2. Rewrite `WriteToString` to use the streaming path:

   ```go
   func (m *MPD) WriteToString(indent string, withHeader bool) (string, error) {
       var buf bytes.Buffer
       buf.Grow(4096)
       if _, err := m.Write(&buf, indent, withHeader); err != nil {
           return "", err
       }
       return buf.String(), nil
   }
   ```

   Even better: use `strings.Builder` — `buf.String()` on
   `bytes.Buffer` copies, `strings.Builder.String()` does not. Check
   whether `strings.Builder` satisfies the `io.Writer` requirements
   used by `Write` (it does).

3. Make sure the test `TestDecodeEncodeMPDs` still passes — it calls
   `xml.MarshalIndent` directly, not `Write`, so it's unaffected. But
   add a new test that explicitly round-trips via `Write` /
   `ReadFromBytes` on a big testdata MPD.

**Verify:**

```sh
go test ./...
```

For a big-MPD benchmark, see step 12.

---

## 11. Tidy `Duration.String` concatenation (H4)

**Why:** `mpd/duration.go:69-137` returns
`"PT" + string(buf[w:])`. The `+` creates an extra string. This is
only ~3.7 % of marshal allocs but trivial to fix.

**File:** `mpd/duration.go:69-137`

**Steps:**

1. Pre-seed `'P'` and `'T'` at the very beginning of `buf`, or leave
   two bytes of head-room in `buf` so the final `w--` can write them.
   Easiest: at the very end, instead of `return "PT" + string(buf[w:])`,
   do:

   ```go
   w--
   buf[w] = 'T'
   w--
   buf[w] = 'P'
   return string(buf[w:])
   ```

   Check that the 32-byte buffer still has room. The max output is
   `PT2540400H10M10.000S` = 20 chars; we have plenty.

2. Handle the negative case: `-` goes before `P`, so move the
   `if neg { buf[w] = '-' }` after the `P`/`T` prefix:

   ```go
   w--
   buf[w] = 'T'
   w--
   buf[w] = 'P'
   if neg {
       w--
       buf[w] = '-'
   }
   return string(buf[w:])
   ```

**Verify:** This function has direct unit tests in the `mpd` package.
Run them:

```sh
go test ./mpd -run=Duration
```

Expected: 1 fewer allocation per Marshal of a `Duration` attribute
(small absolute win, but it's a one-liner).

---

## 12. Fill benchmark coverage gaps

**Why:** The existing benchmarks only exercise a 60-line MPD. Real
hot paths (large live manifests with long `SegmentTimeline`s, on-demand
ladders, multi-period content) behave very differently and are not
currently measured. Without these, regressions on real workloads will
slip through.

**File:** `mpd/mpd_test.go`

**Steps:**

1. Add a parametrised benchmark that iterates a representative subset
   of `mpd/testdata/`. Pattern:

   ```go
   func BenchmarkUnmarshalCorpus(b *testing.B) {
       files, err := filepath.Glob("testdata/schema-mpds/*.mpd")
       require.NoError(b, err)
       for _, f := range files {
           data, err := os.ReadFile(f)
           require.NoError(b, err)
           b.Run(filepath.Base(f), func(b *testing.B) {
               b.ReportAllocs()
               for i := 0; i < b.N; i++ {
                   mpd := m.MPD{}
                   _ = xml.Unmarshal(data, &mpd)
               }
           })
       }
   }
   ```

   Mirror this for `BenchmarkMarshalCorpus`.

2. Add a `BenchmarkWriteToDiscard` that uses `(*MPD).Write` into
   `io.Discard` to validate the streaming path from step 10:

   ```go
   func BenchmarkWriteToDiscard(b *testing.B) {
       mpd, err := m.ReadFromFile("testdata/schema-mpds/example_G15.mpd")
       require.NoError(b, err)
       b.ReportAllocs()
       for i := 0; i < b.N; i++ {
           _, _ = mpd.Write(io.Discard, "  ", true)
       }
   }
   ```

3. **Add a large-MPD fixture** if one is not already present. Look
   through `mpd/testdata/livesim/` for a long-`SegmentTimeline`
   example. If none is big enough, synthesise one with many `<S>`
   children (a few hundred) and check it in as
   `testdata/other/long-timeline.mpd`. Add an unmarshal and marshal
   benchmark specifically for it — that's where most real-world cost
   lives.

**Verify:** `go test ./mpd -bench=. -benchmem` runs all benchmarks
cleanly, including the new ones.

---

## 13. (Optional, low priority) `mpd.ReadFromReader` for `io.Reader` sources (H7)

**Why:** Today a caller with an `io.Reader` (HTTP body, stdin) must
`io.ReadAll` before calling `MPDFromBytes`. A streaming API avoids
the intermediate buffer.

**File:** `mpd/io.go`

**Steps:**

1. Add:

   ```go
   func ReadFromReader(r io.Reader) (*MPD, error) {
       br := bufio.NewReader(r)
       mpd := MPD{}
       if err := xml.NewDecoder(br).Decode(&mpd); err != nil {
           return nil, err
       }
       mpd.SetParents()
       return &mpd, nil
   }
   ```

2. Leave `ReadFromFile` as it is — `os.ReadFile` into `Unmarshal` is
   already fine for local files and avoids double-buffering through
   bufio.

**Verify:** add a unit test that calls `ReadFromReader` on a
`bytes.Reader` wrapping a testdata MPD, and checks it round-trips to
the same output as `ReadFromFile`.

---

## 14. (Optional, defer) H11 — per-byte reader

Skip. The report's recommendation is "leave alone until H1–H4 are
done, then re-profile". After landing steps 2–8 you should re-run
`go test -cpuprofile` and see whether `getc` still shows up. If it
does and Unmarshal is still a concern, revisit.

---

## Expected cumulative effect

Measured from baseline to end of step 11:

| Benchmark          | Baseline       | Target          | How we get there        |
|--------------------|----------------|-----------------|-------------------------|
| BenchmarkUnmarshal | ~57 µs, 703 allocs | ~40 µs, ~350 allocs | Steps 2, 3, 5, 8         |
| BenchmarkMarshal   | ~26 µs, 108 allocs | ~15 µs, ~30 allocs  | Steps 4, 6, 7, 9, 10, 11 |

These are estimates from the profile shares in the audit report. The
actual numbers per step must be **measured**, not assumed — if a step
doesn't produce the expected improvement, stop and investigate before
layering more changes on top.

## Correctness invariants you must preserve

- **Round-trip equality:** `TestDecodeEncodeMPDs` (`mpd/mpd_test.go:22`)
  parses every file under `mpd/testdata/` and re-serialises it,
  verifying semantic equality via `aqwari.net/xml/xmltree`. Run this
  after every single step. If it fails, stop and debug before moving on.
- **Public API:** do not change any exported function's signature in
  the `mpd` package without a deprecation cycle. The `xml` package is
  a fork used internally — it's safer to change, but still, don't
  remove exported names.
- **No behaviour changes on invalid input:** `ParseDuration` must
  still return an error on every malformed input that it currently
  rejects. Add explicit test cases for each class of malformed input
  before rewriting (step 2).
- **Concurrency:** `sync.Pool` (step 9) and `sync.Map` (steps 5, 6)
  are both safe for concurrent use. The per-`Decoder` intern map
  (step 8) is **not** shared across decoders, so no locking needed.
  Don't be tempted to share it globally — it would need a mutex and
  would grow unboundedly.

## Things explicitly out of scope

- Redesigning the schema structs to avoid pointer-valued optional
  fields (`*bool`, `*uint32`, …). That would eliminate a large
  fraction of the remaining `reflect.unsafe_New` allocations but is a
  breaking API change and a much larger project.
- Generating `MarshalXML` / `UnmarshalXML` methods per type. Also a
  big change; revisit only if the reflection-based hot path is still
  dominant after all of the above.
- Replacing `go-deepcopy` in `Clone`. Out of scope for the read/write
  focus (cloning isn't used by parse or write paths).
