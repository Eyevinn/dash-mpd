# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- Nothing yet

## [0.11.1] - 2024-01-17

### Fixed

- Update the DASH-IF ClearKey definitions to follow DASH-IF IOP v5.0

## [0.11.0] - 2024-01-04

### Added

- GetContentProtections, GetMimeType, GetCodecs, and GetSegmentTemplate methods for AdaptationSet and Representation
- ContentProtection elements and name spaces for Marlin DRM and DASH-IF ClearKey

### Fixed

- ContentProtection and other parts of RepresentationBaseType moved before other elements in AdaptationSet and Representation

## [0.10.0] - 2023-05-26

### Changed

- `NewMPD` function now also sets DASH namespace
- `mpd.Write` now has two parameters to set indentation and an optional XML header
- renamed constants StaticMpdType and DynamicMpdType to `STATIC_TYPE` and `DYNAMIC_TYPE`

### Added

Lots of convenience functions to create MPDs

- `mpd.WriteToString` function to return a string
- constants for many common values like audio-channel-configuration
- `rep.SetSegmentBase` method
- `NewBaseURL` function for simple BaseURL cases
- `NewDescriptor` function for creating a DescriptorType instance
- `NewRole` function for creating a new DescriptorType for a role
- `listOfTypes.AddProfile` method to add a profile
- `NewAdaptationSetWithParams` function for generating AdaptationSet
- `NewRepresentationWithID` function for generating a representation with a few parameters
- `NewAudioRepresentation` function for audio representations
- `NewVideoRepresentation` function for video representation
- `Seconds2DurPtrFloat64` to generate a pointer to duration specified as seconds using float64

## [0.9.1] - 2023-05-17

- Same as 0.9.0, but new version number due to mistake in release process

## [0.9.0] - 2023-05-17

### Added

- Parents of AdaptationSets, Representations, and SubRepresentations
- Methods AppendAdaptationSet to Period, and so on
- New Clone methods

## [0.8.0] - 2023-04-10

### Changed

- PeriodType renamed Period
- PeriodType is now PTRegular, PTEarlyAvaiable etc
- GetRepInit changed to GetInit() method on Representation
- GetRepMedia changed to GetMedia() method on Representation
- NewMPD() now takes mpdType as argument

### Added

- New Duration method Seconds()
- New DateTime method ConvertToSeconds()
- New Period methods AbsoluteStart(), GetStart(), GetType(), GetIndex()
- Parent/SetParent methods for Period, AdaptationSet, Representataion, and SubRepresentation
- New Append methods the above that appends and sets parent
- New type PeriodType
- Declared errors: ErrPeriodNotFound and similar
- New function MPDFromBytes() to unmarshal MPD

## [0.7.0] - 2023-03-24

### Changed

- ConvertToDateTime() crops fraction of second to milliseconds instead of full seconds

### Fixed

- Infinite value of availabilityTimeOffset is now marshalled/unmarshalled as "INF"

### Added

- ConvertToDateTimeMS is new function

## [0.6.1] - 2023-03-07

### Changed

- Updated dependencies

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

[Unreleased]: https://github.com/Eyevinn/dash-mpd/compare/v0.11.1...HEAD
[0.11.1]: https://github.com/Eyevinn/dash-mpd/compare/v0.11.0...v0.11.1
[0.11.0]: https://github.com/Eyevinn/dash-mpd/compare/v0.10.0...v0.11.0
[0.10.0]: https://github.com/Eyevinn/dash-mpd/compare/v0.9.1...v0.10.0
[0.9.1]: https://github.com/Eyevinn/dash-mpd/compare/v0.9.0...v0.9.1
[0.9.0]: https://github.com/Eyevinn/dash-mpd/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/Eyevinn/dash-mpd/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/Eyevinn/dash-mpd/compare/v0.6.1...v0.7.0
[0.6.1]: https://github.com/Eyevinn/dash-mpd/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/Eyevinn/dash-mpd/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/Eyevinn/dash-mpd/releases/tag/v0.5.0
