# MPD test files

The files in this directory have been copied from [MPEGGroup/DASHSchema](https://github.com/MPEGGroup/DASHSchema).

Some of them have been slightly modified in order to make them compatible with the output from the mpd package.

The changes are:

* G7 - removed drm elements since not defined in DASH MPD
* G8 - removed comments inside representations and adaptationsets
* G20 - change availabilityTimeOffset from 7.500 to 7.5. Output default value marker="false" in Resync node

When introducing better duration parsing and generation, a lot of durations were then changed in half of the assets.

LICENSE for this content is specified as:

*Use of this repository and all contributions are subject to the ISO/IEC Directives including the ISO and JTC-1 Supplements: https://www.iso.org/directives-and-policies.html*