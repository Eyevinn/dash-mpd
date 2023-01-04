# DASH-MPD - A complete MPEG DASH MPD parser/writer

It is the goal of this MPEG-DASH MPD implementation to include all elements from
the DASH specification by starting from the XML schema and auto-generating
all data structures, as well as handling namespaces and schemaLocation.

A first use case for it is a new DASH-IF live-source simulator in Go,
to replace the original one in [Python](https://github.com/dash-Industry-Forum/dash-live-source-simulator).

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

Mapping of the full MPD Schema to Go structures was a good start, but it had some
limitations and issues, so all structures have been scrutinized, and modified where
needed. The main modifications made were:

* Change of top-level MPD type
* Addition of name spaces including xlink name space
* Change of some attributes to remove `omitempty` or become pointers,
  such as`availabilityTimeComplete`. It is a bool, but it should either have the
  value `false` or be absent.
* Add type comments and document enum values for certain types
* Change names to plural for all subelement slices

## XML handling

The handling of XML name spaces in the Go standard library `encoding/xml` is incomplete,
and has a number of quirks and limitations which makes it impossible to generate
namespaces and namespace prefixes in the standard way we see in many places including
XML schemas and DASH MPDs.

There are a number of pull requests to improve the situation, and in particular
[PR #48641 - encoding/xml: support xmlns prefixes](https://github.com/golang/go/pull/48641)
includes an extension to the XML struct tags that make it possible to specify name
space prefixes.

Since this functionality has not made its way into the standard library after more
than a year, the full patched version of the `encoding/xml` package is included here
in the `xml` directory.

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
