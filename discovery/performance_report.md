# Performance Audit: dash-mpd

Focus: reading (`xml.Unmarshal` / `mpd.ReadFromFile`) and writing
(`xml.MarshalIndent` / `mpd.(*MPD).Write`) of DASH MPDs.

Audit environment:
- macOS (darwin/arm64), Apple M4 Pro, Go 1.26.0
- Commit at audit: `a1042ec` (v0.14.1), clean tree
- Sample corpus: `mpd/testdata/schema-mpds/example_G15.mpd` (~60 lines,
  single period, DRM signalling and segment template). Representative of a
  small static MPD, not of large live MPDs with long SegmentTimelines.

---

## 1. Baseline numbers

Existing benchmarks (`mpd/mpd_test.go`), `go test -bench=. -benchmem -benchtime=5s`:

| Benchmark          | ns/op   | B/op   | allocs/op |
|--------------------|---------|--------|-----------|
| BenchmarkUnmarshal | ~57 600 | 37 080 | **703**   |
| BenchmarkMarshal   | ~26 500 | 19 320 | **108**   |

`BenchmarkClone` fails (stack overflow / `go-cmp` infinite recursion
following the `parent *MPD` back-pointer). That is a bug in the benchmark
set-up, covered in §5.

On a 60-line MPD we allocate **703 objects per parse and 108 per write**.
Every allocation costs both bytes and GC pressure; CPU profiles below
confirm that GC/scheduler overhead (`runtime.kevent`, `madvise`,
`pthread_cond_*`) dominates wall time for both operations.

### CPU/alloc profile highlights

Profiles were collected via `-cpuprofile` / `-memprofile` with
`-benchtime=5s`.

#### Unmarshal — `alloc_space`

```
  flat   flat%                                              cum%
38.79%  (*Decoder).rawToken           // Attr slice + string(data) copies
27.69%  reflect.unsafe_New            // *Period, *AdaptationSet, *string, *uint64 …
 8.56%  (*Decoder).Token
 3.84%  (*Decoder).name               // string(d.buf.Bytes())
 3.09%  (*Decoder).unmarshalAttr
 2.16%  copyValue                     // string(src) for chardata
 3.38%  mpd.ParseDuration / Duration.UnmarshalXMLAttr
 2.44%  regexp.(*Regexp).backtrack
```

#### Unmarshal — `cpu` (non-runtime only)

```
 4.42%  (*Decoder).getc               // byte-at-a-time read
 5.36%  (*fieldInfo).value + reflect.Value.Field
 4.73%  (*Decoder).text
 3.79%  runtime.mallocgcSmallScanNoHeader
 3.94%  isName  (unicode.Is on every read name)
```

Runtime counters (`kevent`, `madvise`, `pthread_cond_*`) together account
for >60 % of samples — symptom of the allocation volume above, not of
syscalls.

#### Marshal — `alloc_objects`

```
75.92%  (*printer).marshalAttr        // append(start.Attr, Attr{name, s})
 4.31%  strconv.FormatUint
 3.86%  (*printer).writeStart         // nsToPrefix / prefixToNS maps
 3.71%  mpd.(*Duration).String        // "PT" + string(buf[w:])
 3.59%  (*printer).createPrefix
 2.04%  bufio.NewWriterSize           // new 4 KiB buffer per Encode
 1.48%  xml.joinPrefixed              // prefix + ":" + name
 1.04%  strconv.FormatInt
```

Marshal CPU, like unmarshal, shows the reflect-based type resolution
(`reflect.implements`, `(*rtype).Implements`, `reflect.Value.Field`,
`(*fieldInfo).value`) as the dominant non-GC cost.

---

## 2. Architectural constraints that bound what we can do

1. **Patched `encoding/xml`.** The `xml/` package is a fork of Go's
   standard library with namespace-prefix support bolted on
   (`README.md:46-64`). Every optimisation lives in our tree, so it is
   available to us — but it also has to keep producing byte-identical
   output vs. the existing test corpus.
2. **Schema-generated struct tree with many pointer fields.** `mpd/mpd.go`
   has 71 struct types and 439 `xml:"…"` tags. Optional scalars are
   modelled as pointers (`*bool`, `*uint32`, `*Duration`, `*string`) to
   distinguish "absent" from "zero". Every populated optional scalar is
   therefore one heap allocation. Slices such as
   `Periods []*Period` are slices of pointers. This is unavoidable
   without a wider schema redesign, and it explains the 27.7 % share of
   `reflect.unsafe_New` in alloc_space.
3. **Linear attribute matching.** `README.md:90-95` already notes this.
   `xml.unmarshal` walks `tinfo.fields` linearly for every incoming
   attribute (`xml/read.go:458-487`). With ~30-field types like
   `AdaptationSetType` and `RepresentationType`, each attribute does a
   linear scan.

---

## 3. Findings, ranked by expected impact

### H1 — `mpd.ParseDuration` runs the same regex twice *(high impact, low risk)*

**Location:** `mpd/duration.go:181-231`

```go
if !xmlDurationRegex.Match([]byte(str)) { … }         // first scan, + []byte alloc
var parts = xmlDurationRegex.FindStringSubmatch(str)   // second scan + allocs
```

- Two full regex traversals per duration.
- `[]byte(str)` is a wasted copy.
- Every call allocates `parts` (a `[]string` of 5 elements) plus an
  internal `bitState` (profile shows 1.52 % of all alloc_space in
  `regexp.(*bitState).reset`).
- Then `strings.TrimRight` + `strconv.Atoi` allocates per segment.

Impact on the MPD corpus: `ParseDuration` alone accounts for **3.4 %** of
alloc_space and 1.8 % of alloc_objects on a single-period small MPD; on
live MPDs with many `SegmentTimeline S@d` or `Event@duration` attributes
the proportion scales roughly with period/segment count. Segments are
usually typed as raw `uint64` though, so the biggest wins come from MPDs
with many `availabilityTimeOffset`, `minBufferTime`, etc.

**Recommendation:** drop regex. Hand-write a scanner for the
`P[nD][T[nH][nM][nS]]` grammar — one pass, no heap traffic in the happy
path. Fallback to an explicit error message on malformed input.
Estimated: >10× speed-up on `ParseDuration`, eliminates all of its
allocations, shrinks unmarshal allocs by ~2-3 %.

### H2 — `xml.rawToken` allocates a fresh `Attr` slice and string per attribute *(high impact)*

**Location:** `xml/xml.go:796-867`

```go
attr = []Attr{}
for { …
    a.Name, _ = d.nsname()      // string(d.buf.Bytes()) per name
    a.Value  = string(data)     // string copy out of d.buf
    attr = append(attr, a)      // slice grow
}
return StartElement{name, attr}, nil
```

Per-attribute cost: up to **three allocations** (name string, value
string, slice extension). Profile: 38.8 % of unmarshal alloc_space sits
in `rawToken`, split roughly:

- 12.1 % — `attr = append(attr, a)` (slice backing array reallocations)
- 6.8 % — `a.Value = string(data)`
- 6.8 % — `a.Name = d.nsname()`
- 4.3 % — `return StartElement{name, attr}`

**Recommendations** (ordered by effort):

1. Pre-size the slice: `attr = make([]Attr, 0, 8)` avoids most of the
   growth allocations for the typical ≤8 attributes/element in DASH.
   One-liner, measurable win.
2. Intern XML names. In an MPD the same ~60 tag names (`Period`,
   `AdaptationSet`, `Representation`, `SegmentTemplate`, …) and the same
   ~100 attribute names (`id`, `bandwidth`, `mimeType`, `codecs`,
   `timescale`, …) appear thousands of times. Keep a per-`Decoder`
   `map[string]string` that turns the `d.buf.Bytes()` slice into a
   canonical `string` on first occurrence and returns the interned one
   thereafter. Costs one map lookup per name, saves the `string(b)`
   allocation for every repeat. The canonicaliser needs to be on the
   Decoder (not global) to avoid unbounded growth across long-running
   services.
3. Pool `[]Attr` and `StartElement` backing storage. The token returned
   by `rawToken` is consumed almost immediately by
   `Decoder.Token`/`unmarshal`. A scoped `sync.Pool` for the Attr slice
   is feasible but invasive — the Attr slice currently outlives the
   call. Defer this until (1) and (2) have been measured.

### H3 — `xml.marshalAttr` allocates an `Attr` and a number string per attribute *(high impact)*

**Location:** `xml/marshal.go:607-688`, `marshalSimple` at lines 883-916.

Each struct attribute does:

1. Recurse into `marshalSimple`, which calls
   `strconv.FormatInt/FormatUint/FormatFloat/FormatBool` returning a
   *string* (`strconv.FormatUint = 4.31 %` of marshal allocs).
2. Build `Attr{name, s}` and `append` it onto `start.Attr`.

75.92 % of marshal allocations come from this path. Two complementary
fixes:

1. **Skip the Attr struct entirely for simple fields.** `writeStart`
   (`xml/marshal.go:746-846`) is the only consumer of `start.Attr`, and
   it only reads each Attr once. Instead of materialising a `[]Attr`
   then iterating it, `marshalValue` could call a `writeAttr(prefix,
   name, value)` method that writes directly into the `bufio.Writer`
   after `writeStart` has emitted the element's prefix logic. Trade-off:
   custom `MarshalXMLAttr` implementations (e.g.
   `mpd/duration.go:48-53`) return an `Attr`, so we still need the
   intermediate representation for those paths — but numeric/string
   attributes are the vast majority.
2. **Pre-size the attr slice** analogously to H2: `start.Attr =
   make([]Attr, 0, N)` where N is the struct's attribute count from
   `tinfo`. Cheap, reduces growslice churn until (1) lands.
3. **Append ints directly.** `marshalSimple` already uses
   `strconv.AppendInt` in the `fCharData` branch (lines 984-999). The
   attribute branch should do the same with a scratch `[64]byte` and
   return `[]byte` so we never materialise a Go string for transient
   numeric values.

### H4 — `mpd.(*Duration).String` allocates the final string via concatenation *(medium)*

**Location:** `mpd/duration.go:69-137`

```go
return "PT" + string(buf[w:])
```

3.71 % of marshal allocations. The `buf` has 32 bytes and the scratch
already lives on the stack — we only pay to turn it into a `string`.

**Recommendation:** pre-write `'P','T'` into `buf` when constructing the
output so the function returns `string(buf[w:])` with no concatenation,
and then have `Duration.MarshalXMLAttr` return an attribute whose value
shares a pooled buffer where possible. A simpler intermediate step: keep
the prefix in `buf` and remove the `+` operator (saves one allocation
per call).

### H5 — `xml.Encoder` allocates a 4 KiB `bufio.Writer` per call *(medium)*

**Location:** `xml/marshal.go:141-145`

```go
func NewEncoder(w io.Writer) *Encoder {
    e := &Encoder{printer{Writer: bufio.NewWriter(w)}}
    …
}
```

`bufio.NewWriterSize` accounts for 2.04 % of marshal allocs per call. On
`mpd.(*MPD).WriteToString`/`Write` this is per-invocation overhead. For
services that write many MPDs per second (live origin / packager) this
adds up.

**Recommendations:**

1. Expose a `sync.Pool` of `*Encoder` and reuse them between calls. The
   `printer` struct holds transient state (`elements`, `depth`,
   `putNewline`) that `Encode` leaves behind at zero values, so a
   simple reset before returning to the pool is enough.
2. Alternatively, if the caller already supplies a buffer
   (`mpd.(*MPD).Write` writes to a `bytes.Buffer`/`http.ResponseWriter`
   that is itself buffered), expose an API that bypasses the internal
   bufio — `NewEncoderUnbuffered(io.Writer)`.

### H6 — Double buffering between `Marshal` and the caller *(medium)*

**Location:** `mpd/io.go:43-71`

```go
func (m *MPD) Write(w io.Writer, indent string, withHeader bool) (int, error) {
    data, err := xml.MarshalIndent(m, "", indent)           // full copy in memory
    …
    n, err = w.Write(data)                                  // then copy to w
}
```

`MarshalIndent` returns a `[]byte` that is copied out. For large live
MPDs (multi-period, long timelines, tens to hundreds of KB) this doubles
peak memory and adds a `memmove`. `WriteToString` adds a *third* copy
(`string(data)`).

**Recommendation:** stream via `xml.NewEncoder(w)` + `enc.Indent(…)` +
`enc.Encode(m)`. Signature of `Write` stays the same; remove the
intermediate `[]byte`. Combine with H5's pool for a bounded-memory fast
path.

### H7 — `ReadFromFile` buffers the whole file, then re-reads via `bytes.Reader` *(low-medium)*

**Location:** `mpd/io.go:11-24`

```go
data, err := os.ReadFile(path)
…
err = xml.Unmarshal(data, &mpd)   // bytes.NewReader internally
```

Two alternatives:

- Keep as is — for local file I/O this is rarely the bottleneck and the
  `bytes.Reader` path is already zero-copy.
- Switch to `xml.NewDecoder(bufio.NewReader(f)).Decode(&mpd)` for
  `io.Reader` sources (HTTP body, stdin). This is the more important
  case: today a caller who has an `io.Reader` must `io.ReadAll` first.

**Recommendation:** add `mpd.ReadFromReader(io.Reader) (*MPD, error)`
that streams through a bufio reader. Leaves `ReadFromFile` semantics
alone.

### H8 — Reflection-driven attribute lookup is O(fields × attrs) per element *(medium, already known)*

**Location:** `xml/read.go:458-487`

```go
for _, a := range start.Attr {
    …
    for i := range tinfo.fields {
        finfo := &tinfo.fields[i]
        switch finfo.flags & fMode {
        case fAttr:
            if a.Name.Local == finfo.name && … {
```

The `README` already calls this out. In CPU profiles today it sits at
~5.4 % (`(*fieldInfo).value` + `reflect.Value.Field`), eclipsed by
allocation pressure — so attacking H1-H3 first will shift this up in
relative cost.

**Recommendation:** extend the cached `*typeInfo` with an
`attrs map[string]*fieldInfo` (and, if needed, an `attrsByNS` map). The
existing `tinfoMap sync.Map` already memoises per-type info, so we only
pay map construction once per type. Swap the inner loop for a single
lookup.

Side benefit: removes the need for the `any`-attr scan rewriting in
`read.go:472-486`. Keep a small parallel list of `fAny|fAttr` fields.

### H9 — `reflect.implements` checks done 4× per marshalled value *(low-medium)*

**Location:** `xml/marshal.go:473-493`

Each value asks: does `typ` implement `Marshaler`? If not, does `*typ`?
Then the same pair for `TextMarshaler`. Profile shows
`reflect.implements` + `(*rtype).Implements` ≈ 7 % of marshal CPU.

**Recommendation:** cache per-type booleans
(`implementsMarshaler`, `implementsAddrMarshaler`, … ) on the
`*typeInfo`. That pushes the lookup into the `sync.Map` we already have
and replaces four `Implements` calls with four bool loads.

### H10 — `BenchmarkClone` is broken *(blocker for cloning work)*

**Location:** `mpd/mpd_test.go:118-129`

```go
cmp.Equal(&mpd, mpdCopy)       // not what you want
for i := 0; i < b.N; i++ {
    _ = m.Clone(&mpd)
}
```

`go-cmp` traverses the struct and follows the `parent` back-pointer that
`SetParents()` sets on `Period`/`AdaptationSetType`/… creating an
infinite recursion; it stack-overflows on entry. This is what makes
`go test -bench=.` abort rather than reporting Clone numbers.

Regardless of the bench fix, `Clone` uses
`github.com/barkimedes/go-deepcopy` (`mpd/mpd.go:101-105`,
`mpd/period.go:178-182`), a reflection-based deep copier that is ~10×
slower than a generated or marshal-roundtrip clone for similar trees.

**Recommendations:**

1. Fix the benchmark: remove the errant `cmp.Equal` call (it does
   nothing useful there) and move the expected parent re-wiring into
   the benchmarked path only once, outside the loop.
2. Quantify `Clone` latency. If it's hot in real workloads (it isn't
   used inside read/write paths, so audit callers first), replace
   `go-deepcopy` with either:
   - A generated `Clone()` method on `*MPD` and its substructs —
     straightforward from the existing XML schema, zero reflection at
     runtime.
   - A marshal-then-unmarshal round-trip, which reuses all the other
     improvements on this list and is only ~2× the cost of marshal+parse.

Cloning is out of scope for the requested read/write focus, but the
broken benchmark blocks measuring any of the above against the other
benchmarks in the same `go test -bench=.` run — fix it for workflow
reasons.

### H11 — Per-byte reads from `bytes.Reader` *(low)*

**Location:** `xml/xml.go:924-945` (`getc`)

CPU profile puts `getc` at 4.42 %, `bytes.(*Reader).ReadByte` at 1.42 %.
`Unmarshal` wraps the input in a `bytes.Reader` (native `ByteReader`);
each byte pays one interface call. A custom hot loop that reads a slice
into a ring would be faster but significantly more invasive in the
forked parser.

**Recommendation:** leave alone until H1–H4 are done — the allocation
reductions will also reduce the per-byte CPU overhead indirectly via
lower GC work.

---

## 4. Suggested ordering

Measured against `BenchmarkUnmarshal`/`BenchmarkMarshal`:

1. **H1** (`ParseDuration` rewrite) — self-contained, high impact,
   risk-free. Ship first.
2. **H10** fix the `Clone` benchmark — trivial, but required so we can
   see regressions in the same run.
3. **H2.1** pre-size attr slice + **H3.2** pre-size attr slice in
   marshal — two one-line changes, easy to validate via
   `-benchmem`.
4. **H8** attribute-lookup map on `*typeInfo` + **H9** cached
   `implements` booleans — small change to `xml/typeinfo.go`, removes
   the most visible non-GC CPU time.
5. **H3.1 / H3.3** bypass `Attr` for simple attributes — biggest
   marshal win but touches the hot-path API in `xml/marshal.go`. Needs
   round-trip tests; the existing `TestDecodeEncodeMPDs` over all
   testdata MPDs is a good regression net.
6. **H2.2** name interning on the Decoder — requires a small design
   call about bounding the intern table; defer until the earlier wins
   are measured.
7. **H5 / H6** encoder pool and streaming `Write` — do together; easy
   follow-up once H3 has landed.
8. **H11 / H7** low-priority cleanups.

A rough target after H1–H5 is:

- Unmarshal: ~700 allocs → ~350 allocs, ~57 µs → ~40 µs
- Marshal: ~108 allocs → ~30 allocs, ~26 µs → ~15 µs

Conservative, based on the profile shares above; will need
re-benchmarking at each step.

---

## 5. Benchmark coverage gaps

The current benchmarks only cover a **single 60-line MPD**. They do not
exercise the realistic hot paths:

- Large live MPD with a long `SegmentTimeline` (hundreds to thousands
  of `S` elements per Representation). This is the hot case for DASH
  packagers / manifest manipulators.
- Multi-period MPD (SCTE-35 splice points), ~dozens of periods.
- On-demand profile MPD with many Representations (HDR/SDR/bitrate
  ladders).

Adding a parametrised benchmark that iterates all files in
`mpd/testdata/` (or the subset that is semantically realistic) would
catch regressions that the current benchmarks miss, and would give a
more honest picture of where optimisations matter.

Additionally:
- `BenchmarkClone` needs the fix in H10.
- A `BenchmarkReadFromFile` would let us distinguish I/O from parse
  cost, relevant once H7's streaming reader lands.
- A `BenchmarkWrite` that writes to `io.Discard` (rather than
  marshalling to a buffer) would validate H6.

---

## 6. Summary

The parser/writer are correct and well-tested, but their cost is
dominated by small, per-element allocations: attribute slices and
strings in the parser, `Attr` structs and numeric strings in the
writer, plus a redundant regex pass in `ParseDuration`. None of the
findings require touching the generated struct tree — all live in
`xml/*.go`, `mpd/duration.go`, and `mpd/io.go`. Delivering H1 through
H5 should roughly halve parse allocations and cut marshal allocations
by ~3-4×, with no behaviour change (the existing round-trip test suite
over the `testdata` MPDs guards correctness).
