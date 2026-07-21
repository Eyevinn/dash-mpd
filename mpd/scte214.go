package mpd

import (
	"github.com/Eyevinn/dash-mpd/xml"
)

// SCTE-214 (ANSI/SCTE 214-1 2024) MPD extensions for DASH.
//
// This file maps the MPD extensions defined by SCTE 214-1 (MPEG DASH for
// IP-Based Cable Services Part 1: MPD Constraints and Extensions) onto Go
// types and constants:
//
//   - the supplementalProfiles and supplementalCodecs attributes extending
//     RepresentationBaseType (Sec. 11.1.6, Table 8); the fields live on
//     RepresentationBaseType in mpd.go and are therefore available on
//     AdaptationSet, Representation and SubRepresentation
//   - the ContentIdentifier element for UPID-based asset identification
//     (Sec. 9.2); the field lives on DescriptorType in mpd.go and is
//     therefore available on AssetIdentifier and SupplementalProperty
//   - the descriptor schemeIdUri values defined throughout the document
//     and listed in Annex A
//
// All extension elements and attributes belong to the SCTE214Namespace
// ("urn:scte:dash:scte214-extensions"). The scte214 prefix is used on
// marshal; on unmarshal any prefix bound to the namespace URI is matched.

// SCTE-214 XML namespace URIs (ANSI/SCTE 214-1 2024 Annex A).
const (
	// SCTE214Namespace is the XML namespace for SCTE-214 MPD extensions,
	// used by all editions of the specification after 2021.
	SCTE214Namespace = "urn:scte:dash:scte214-extensions"
	// SCTE214Namespace2021 is the XML namespace of the 2021 and prior
	// editions of the specification.
	SCTE214Namespace2021 = "urn:scte:dash:2021"
)

// Descriptor schemeIdUri values defined by ANSI/SCTE 214-1 2024.
const (
	// SCTE214SchemeIdAssetIdUpid identifies assets with SCTE-35 UPIDs.
	// AssetIdentifier and SupplementalProperty descriptors with this
	// scheme carry one or more ContentIdentifier elements (Sec. 9.2).
	SCTE214SchemeIdAssetIdUpid = "urn:scte:dash:asset-id:upid:2015"
	// SCTE214SchemeIdCEA608 signals CEA-608 closed caption services in
	// Accessibility descriptors (Sec. 8.2.3).
	SCTE214SchemeIdCEA608 = "urn:scte:dash:cc:cea-608:2015"
	// SCTE214SchemeIdCEA708 signals CEA-708 closed caption services in
	// Accessibility descriptors (Sec. 8.2.2).
	SCTE214SchemeIdCEA708 = "urn:scte:dash:cc:cea-708:2015"
	// SCTE214SchemeIdAssociatedService carries roles for non-accessibility
	// associated audio services in Role descriptors (Sec. 8.1).
	SCTE214SchemeIdAssociatedService = "urn:scte:dash:associated-service:2015"
	// SCTE214SchemeIdEssentialEvent marks an event scheme whose processing
	// is essential for the presentation, in EssentialProperty descriptors
	// (Annex A).
	SCTE214SchemeIdEssentialEvent = "urn:scte:dash:essential-event:2015"
	// SCTE214SchemeIdAssetEnd marks the last period of an asset, in Period
	// SupplementalProperty descriptors (Sec. 12.2).
	SCTE214SchemeIdAssetEnd = "urn:scte:dash:asset-end"
	// SCTE214SchemeIdAssetTime maps PeriodStart to an NPT or SMPTE
	// timestamp relative to the asset start (Sec. 12.2).
	SCTE214SchemeIdAssetTime = "urn:scte:dash:asset-time"
	// SCTE214SchemeIdUTCTime maps PeriodStart to the UTC capture time of
	// the first sample of the period (Sec. 12.2).
	SCTE214SchemeIdUTCTime = "urn:scte:dash:utc-time"
	// SCTE214SchemeIdPoweredBy identifies the software that generated the
	// MPD or the XLink response (Sec. 14.1).
	SCTE214SchemeIdPoweredBy = "urn:scte:dash:powered-by:2016"
	// SCTE214SchemeIdGenerationInfo carries the location and time of MPD
	// or XLink remote entity generation (Sec. 14.1).
	SCTE214SchemeIdGenerationInfo = "urn:scte:dash:generation-info:2016"
	// SCTE214SchemeIdGenerationRequest carries the HTTP request target
	// that produced the MPD or XLink response (Sec. 14.1).
	SCTE214SchemeIdGenerationRequest = "urn:scte:dash:generation-request:2016"
	// SCTE214SchemeIdDefaultLanguage signals the default media language in
	// SupplementalProperty descriptors. An optional URI fragment (#video,
	// #audio, #text or #commentary) narrows the content component type it
	// applies to (Sec. 8.1.2).
	SCTE214SchemeIdDefaultLanguage = "urn:scte:dash:default-language"
	// SCTE214SchemeIdIFrameTrack marks I-Frame track representations for
	// trick modes (Sec. 11.3.3).
	SCTE214SchemeIdIFrameTrack = "urn:scte:dash:i-frame-track:2021"
	// SCTE214SchemeIdSegmentModel signals conformance to the SCTE segment
	// model in SupplementalProperty descriptors (Sec. 10.1.1).
	SCTE214SchemeIdSegmentModel = "urn:scte:dash:segment-model:2024"
)

// SCTE214ContentIdentifierType is the scte214:ContentIdentifier element (XSD
// complexType UPID) defined in Sec. 9.2 of ANSI/SCTE 214-1 2024. Type is an
// SCTE-35 UPID type name and Value its textual representation, both as
// specified by table 9-7 of ANSI/SCTE 35. One or more ContentIdentifier
// elements appear inside an AssetIdentifier or SupplementalProperty
// descriptor with schemeIdUri SCTE214SchemeIdAssetIdUpid.
type SCTE214ContentIdentifierType struct {
	Type     string     `xml:"type,attr"`
	Value    string     `xml:"value,attr"`
	AnyAttrs []xml.Attr `xml:",any,attr"`
}

// NewSCTE214ContentIdentifier returns a ContentIdentifier with the given
// SCTE-35 UPID type name and value.
func NewSCTE214ContentIdentifier(idType, idValue string) *SCTE214ContentIdentifierType {
	return &SCTE214ContentIdentifierType{Type: idType, Value: idValue}
}
