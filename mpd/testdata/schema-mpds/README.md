# MPD test files

The files in this directory are copied from
[MPEGGroup/DASHSchema](https://github.com/MPEGGroup/DASHSchema), branch
`6th-Ed`, commit `a855144` (2026-04-22), and re-saved through this
package's marshaller so that decode→encode is byte-stable.

The re-save normalises:

* Attribute order on each element (matches the order of the Go struct fields).
* Duration values to their canonical short form
  (e.g. `PT3256S` → `PT54M16S`, `PT1.500000S` → `PT1.5S`).
* `Resync` always emits `marker="false"` when not set explicitly.

A few examples were further tweaked compared to upstream:

* G7 — drm elements removed (not defined in DASH MPD).
* G8 — comments inside Representation/AdaptationSet removed.
* G20 — `availabilityTimeOffset` is `7.5` (not `7.500`).

LICENSE for this content is specified as:

*Use of this repository and all contributions are subject to the ISO/IEC
Directives including the ISO and JTC-1 Supplements:
https://www.iso.org/directives-and-policies.html*
