![Test](https://github.com/Eyevinn/mp4ff/workflows/Go/badge.svg)
![golangci-lint](https://github.com/Eyevinn/mp4ff/workflows/golangci-lint/badge.svg?branch=master)
[![GoDoc](https://godoc.org/github.com/Eyevinn/mp4ff?status.svg)](http://godoc.org/github.com/Eyevinn/mp4ff)
[![Go Report Card](https://goreportcard.com/badge/github.com/Eyevinn/mp4ff)](https://goreportcard.com/report/github.com/Eyevinn/mp4ff)
[![license](https://img.shields.io/github/license/Eyevinn/mp4ff.svg)](https://github.com/Eyevinn/mp4ff/blob/master/LICENSE)

# DASH-MPD - A complete MPEG-DASH MPD parser/writer

This MPEG-DASH MPD implementation is meant to include all elements from
the MPEG DASH specification (ISO/IEC 23009-1 5'th edition) by starting from the
MPD XML schema and auto-generating all data structures.
It should also handle namespaces and schemaLocation properly.

It has been enhanced with extra structures from the Common Encryption specification
ISO/IEC 23001-7 and proprietary structures and name spaces for some DRM systems.

## XML Schemas

The XML Schema for MPEG DASH MPD being used is the 5'th edition fetched from the
[DASHSchema repo](https://github.com/MPEGGroup/DASHSchema) commit `993cb92`.

It was fed to a forked modified version of [xgen](https://github.com/xuri/xgen).
The fork modified `xgen` so that it could deal with DOCTYPE,
and also removed the suffix `Attr` from all attribute names. The fork can be found at
[tobbee xgen](https://github.com/tobbee/xgen/tree/shorten-attr).

The process of generating this start MPD structures is done in the Github repo
[tobbee/mpdgen](https://github.com/tobbee/mpdgen). The output provided the starting
point for the MPD type of this project.

## Adjustments of the MPD structures

Mapping of the full MPD Schema to Go structures provides a good start, but it has some
limitations and issues. Therefore, all structures have been scrutinized, and modified where
needed. The main modifications made were:

* Change of top-level MPD type
* Add xlink name space
* Change some attributes to remove `omitempty` or become pointers.
  An example is `availabilityTimeComplete`. It is boolean, but it should either have the
  value `false` or be absent
* Add type comments and document enum values for certain types
* Change names to plural for all subelement slices, e.g. Periods instead of Period
* Add `cenc:pssh` and `mspr:pro` DRM elements

## XML handling

The handling of XML name spaces in the Go standard library `encoding/xml` is incomplete,
and has a number of quirks and limitations which makes it impossible to generate
namespaces and namespace prefixes in the standard way used in many places including
XML schemas and DASH MPDs.

There are a number of pull requests to improve the situation, and in particular
[PR #48641 - encoding/xml: support xmlns prefixes](https://github.com/golang/go/pull/48641)
includes an extension to the XML struct tags that make it possible to specify name
space prefixes.

Since this functionality has not yet (Go 1.19) made its way into the standard library, after more
than a year, the full patched version of the `encoding/xml` package is included here
in the `xml` directory. Therefore, the `github.com/Eyevinn/dash-mpd/xml` package should
be used together with the XML structure in the package, rather than the standard library version.
The package functions `mpd.ReadFromFile()`, `mpd.ReadFromString()`, and `mpd.Write()` do use
that XML library.

## Limitations to MPD manipulations

Since the MPD is mapped to Go structures, there are some limitations to what can be processed
and how the output looks:

1. Only elements and attributes specified in the data structures are kept. Unkonwn elements and
   attributes are silently discarded.
2. All XML comments are removed
3. The output order of XML attributes is given by the order in the structure, which is in turn
   comes from the XML schema. The order is therefore often different from the input document.
4. Mapping to float numbers may not preserve the exact value of the input.
5. Addition of extra name spaces, such as specific DRM systems must be done explicitly.
6. Durations are mapped to nanoseconds and back. This may change the duration slightly. All trailing zeros
   are also removed, as are the minutes and seconds counts if they are zero and a bigger unit is present.

## Tests

The MPD marshaling/unmarshaling is tested by by using many MPDs in the `testdata/schema-mpds` and
`testdata/go-dash-fixtures` directories. See the README.md files in these directorie for their
origins and the small tweaks needed to make durations and some other values consistent after a
unmarshaling/marshaling process.

## Performance

Including all elements, including rarely used ones, has some performance penalty.
One such contribution comes from matching attributes which is done in a linear fashion.
Some benchmark tests are included in the `mpd/mpd_test.go` file.
