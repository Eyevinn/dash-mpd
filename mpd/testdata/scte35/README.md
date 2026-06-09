# SCTE-35 test fixtures

DASH MPDs carrying SCTE-35 splice information sections, used by
`TestDecodeEncodeMPDs` to validate the unmarshal/marshal round-trip
of the types in `mpd/scte35.go`.

The schema source is [`scte_35_20230713.xsd`][xsd] (ANSI/SCTE 35
2023r2), retrieved from <https://schemas.scte.org/35/>. The carriage
in DASH events (the `<Signal>` wrapper element and the
`urn:scte:scte35:2013:xml` / `urn:scte:scte35:2014:xml+bin` schemeIdUri
values) is defined by [SCTE 214-1][214-1].

Each fixture is hand-written and verified to round-trip
unchanged through the package after a single marshal pass.
Durations are written in their canonical form
(e.g. `PT48H10M2.036S`, not `PT173402.036S`) and DASH attributes use
the order produced by the marshaller; both are existing constraints
of the round-trip test that the SCTE-35 fixtures inherit.

| File | What it covers |
|---|---|
| `time_signal_segmentation.mpd` | `<Signal>` + `TimeSignal` + `SegmentationDescriptor` + `SegmentationUpid` (the canonical ad-marker form) |
| `splice_insert.mpd` | `SpliceInsert` with `Program`/`SpliceTime` and `BreakDuration` |
| `splice_insert_immediate.mpd` | `SpliceInsert` with `spliceImmediateFlag="true"` (no `SpliceTime`) |
| `signal_binary_xmlbin.mpd` | `<Signal>` + `Binary` (`urn:scte:scte35:2014:xml+bin` form) using the **legacy** `http://www.scte.org/schemas/35/2016` namespace, demonstrating round-trip pass-through of the namespace URI |
| `direct_splice_info_section.mpd` | `SpliceInfoSection` placed directly inside `<Event>` without a `<Signal>` wrapper (AWS MediaTailor style) |
| `delivery_restrictions.mpd` | `SegmentationDescriptor` with `DeliveryRestrictions` and a hex-binary `SegmentationUpid` |
| `splice_null_heartbeat.mpd` | `SpliceNull` heartbeat |
| `splice_schedule.mpd` | `SpliceSchedule` with a single scheduled `Event` |
| `private_command.mpd` | `PrivateCommand` with `PrivateBytes` payload |
| `time_descriptor.mpd` | `TimeSignal` accompanied by a `TimeDescriptor` (TAI) |
| `audio_descriptor.mpd` | `AudioDescriptor` with two `AudioChannel` children |
| `avail_dtmf_descriptors.mpd` | `SpliceInsert` accompanied by `AvailDescriptor` and `DTMFDescriptor` |

[xsd]: https://schemas.scte.org/35/scte_35_20230713.xsd
[214-1]: https://account.scte.org/standards/library/catalog/scte-214-1-mpeg-dash-for-ip-based-cable-services-part1-mpd-constraints-and-extensions/
