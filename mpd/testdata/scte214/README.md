# SCTE-214 test fixtures

DASH MPDs carrying SCTE-214 MPD extensions, used by
`TestDecodeEncodeMPDs` to validate the unmarshal/marshal round-trip
of the types and attributes in `mpd/scte214.go`.

The extensions are defined by [ANSI/SCTE 214-1 2024][214-1] (MPEG DASH
for IP-Based Cable Services Part 1: MPD Constraints and Extensions):
the `supplementalProfiles` and `supplementalCodecs` attributes on
RepresentationBaseType (Sec. 11.1.6, Table 8) and the
`ContentIdentifier` element for UPID-based asset identification
(Sec. 9.2). All of them belong to the
`urn:scte:dash:scte214-extensions` namespace (Annex A).

Each fixture is hand-written and verified to round-trip unchanged
through the package. DASH attributes use the order produced by the
marshaller and the `xmlns:scte214` declaration sits on the carrying
element, both existing constraints of the round-trip test.
Attribute values follow those seen in production manifests
(Dolby Vision 8.1 over HEVC, HDR10+) and the signaling examples of
Sec. 11.1.7.

| File | What it covers |
|---|---|
| `supplemental_codecs_representation.mpd` | `supplementalCodecs`/`supplementalProfiles` on `Representation`, with different codec strings per representation (Dolby Vision 8.1 over HEVC) |
| `supplemental_codecs_adaptation_set.mpd` | the attributes on `AdaptationSet` with a space-separated multi-value `supplementalProfiles` list (`db1p cdm4`, Dolby Vision 8.1 + HDR10+) |
| `asset_id_upid.mpd` | `ContentIdentifier` elements inside `AssetIdentifier` and `SupplementalProperty` descriptors with the `urn:scte:dash:asset-id:upid:2015` scheme |

[214-1]: https://account.scte.org/standards/library/catalog/scte-214-1-mpeg-dash-for-ip-based-cable-services-part1-mpd-constraints-and-extensions/
