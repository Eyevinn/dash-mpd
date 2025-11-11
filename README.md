![Test](https://github.com/Eyevinn/dash-mpd/workflows/Go/badge.svg)
[![golangci-lint](https://github.com/Eyevinn/dash-mpd/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/Eyevinn/dash-mpd/actions/workflows/golangci-lint.yml)
[![GoDoc](https://godoc.org/github.com/Eyevinn/dash-mpd?status.svg)](http://godoc.org/github.com/Eyevinn/dash-mpd)
[![Go Report Card](https://goreportcard.com/badge/github.com/Eyevinn/dash-mpd)](https://goreportcard.com/report/github.com/Eyevinn/dash-mpd)
[![license](https://img.shields.io/github/license/Eyevinn/dash-mpd.svg)](https://github.com/Eyevinn/dash-mpd/blob/master/LICENSE)

# DASH-MPD - A complete MPEG-DASH MPD parser/writer

This MPEG-DASH MPD implementation is meant to include all elements from
the MPEG DASH specification (ISO/IEC 23009-1 5'th edition) by starting from the
MPD XML schema and auto-generating all data structures.
It should also handle namespaces and schemaLocation properly.

It has been enhanced with extra structures from the Common Encryption specification
ISO/IEC 23001-7 and proprietary structures and name spaces for some DRM systems.

## XML Schemas

The XML Schema for MPEG DASH MPD being used is the 5'th edition fetched from the
[DASHSchema repo](https://github.com/MPEGGroup/DASHSchema) commit `169fbd3`.

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
* Add ContentProtection elements and corresponding name spaces

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

1. Only elements and attributes specified in the data structures are kept. Unknown elements and
   attributes are silently discarded.
2. All XML comments are removed
3. The output order of XML attributes is given by the order in the structure, which is in turn
   comes from the XML schema. The order is therefore often different from the input document.
4. Mapping to float numbers may not preserve the exact value of the input.
5. Addition of extra name spaces, such as specific DRM systems must be done explicitly.
6. In the output, the name-spaces are added to the level where they are used and not at the top MPD level.
   The output names are also fixed, and may differ from the input names.
7. Durations are mapped to nanoseconds and back. This may change the duration slightly. All trailing zeros
   are also removed, as are the minutes and seconds counts if they are zero and a bigger unit is present.

## Tests

The MPD marshaling/unmarshaling is tested by by using many MPDs in the `testdata/schema-mpds`,
`testdata/go-dash-fixtures`, and `testdata/livesim` directories.
See the README.md files in these directories for their
origins and the small tweaks needed to make durations and some other values consistent after a
unmarshaling/marshaling process.

## Performance

Including all elements, including rarely used ones, has some performance penalty.
One such contribution comes from matching attributes which is done in a linear fashion.
Some benchmark tests are included in the `mpd/mpd_test.go` file.

## Commits, ChangeLog

This project aims to follow Semantic Versioning and
[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).
There is a manual [ChangeLog](CHANGELOG.md) that should be updated with
each commit.

With making a local git commit, this project uses [pre-commit hooks][precmt].
You therefore need Python3 installed and run `make prepare` to create a Python
virtual envinvorment at `.venv` and then run `source .venv/bin/activate` in your
shell in order to activate that virtual environment to get the [pre-commit`][precmt] mechanism to work.

## License

MIT, see [LICENSE](LICENSE).

## Support

Join our [community on Slack](http://slack.streamingtech.se) where you can post any questions regarding any of our open source projects. Eyevinn's consulting business can also offer you:

* Further development of this component
* Customization and integration of this component into your platform
* Support and maintenance agreement

Contact [sales@eyevinn.se](mailto:sales@eyevinn.se) if you are interested.

## About Eyevinn Technology

[Eyevinn Technology](https://www.eyevinntechnology.se) is an independent consultant firm specialized in video and streaming. Independent in a way that we are not commercially tied to any platform or technology vendor. As our way to innovate and push the industry forward we develop proof-of-concepts and tools. The things we learn and the code we write we share with the industry in [blogs](https://dev.to/video) and by open sourcing the code we have written.

Want to know more about Eyevinn and how it is to work here. Contact us at work@eyevinn.se!

[precmt]: https://pre-commit.com