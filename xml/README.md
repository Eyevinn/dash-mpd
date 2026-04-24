# Patched XML

The handling of XML name spaces in the Go standard library `encoding/xml` is incomplete,
and has a number of quirks and limitations which makes it impossible to generate
namespaces and namespace prefixes in the standard way we see in many places including
XML schemas and DASH MPDs.

There are a number of pull requests to improve the situation, and in particular
[PR #48641 - encoding/xml: support xmlns prefixes](https://github.com/golang/go/pull/48641)
includes an extension to the XML struct tags that make it possible to specify name
space prefixes.

This package is a copy of PR#48641. More specifically, commit
`a79fce7705e3102ca642bb981bb3829b3fa0824f` from the branch `xmlns-fixes` in
[ydnar/go](https://github.com/ydnar/go), which is rebased on Go master
as of 2024-10-28.

Hopefully, this code will be merged into the Go standard library at some point.

Escaping tab and newline in `CharData` byte content led to problems with some
DASHSchema test content, so this is disabled: `EscapeText` passes
`escapeNewline=false`, and the `'\t'` case is removed from `escapeText`. Tab and
newline are still escaped in `(*printer).EscapeString`, which is used for
attribute values.
