# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] - 2023-03-03

### API change

- Changed composition of some MPD elements to not use pointers

### Added

- Functions to create new elements
- Example code to generate new MPD from scratch
- Added XMLNames to structure with unique names in MPD
- Methods SetTimescale and GetTimescale for SegmentTemplate

## [0.5.0] - 2023-01-20

### Added

- First release with complete support of all XMLSchema elements
- Tests with well-known MPDs
- Tweaked XML library to support namespaces

[0.6.0]: https://github.com/Eyevinn/dash-mpd/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/Eyevinn/dash-mpd/releases/tag/v0.5.0
