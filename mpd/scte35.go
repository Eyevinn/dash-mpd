package mpd

import (
	"github.com/Eyevinn/dash-mpd/xml"
)

// SCTE-35 (ANSI/SCTE 35 2023r2) splice information signalling for DASH MPD events.
//
// This file maps the SCTE-35 XML representation onto Go structs. The schema
// source is https://schemas.scte.org/35/scte_35_20230713.xsd (version 20230713,
// matching the ANSI/SCTE 35 2023r2 specification). The <Signal> wrapper element
// is defined by SCTE 214-1 (MPEG-DASH carriage); the SCTE-35 XSD itself defines
// it only as an xsd:group, with <SpliceInfoSection> and <Binary> as the two
// global elements that may appear directly inside a DASH <Event>.
//
// # Namespace handling
//
// Two namespace URIs are seen in the wild for the same schema:
//
//   - "http://www.scte.org/schemas/35"      (canonical, declared as
//     targetNamespace by all XSDs since 2020)
//   - "http://www.scte.org/schemas/35/2016" (legacy, still produced by a
//     significant fraction of packagers)
//
// The schema content is identical between the two URIs. Each SCTE-35 type in
// this file therefore carries its namespace URI in an untagged XMLName field.
// On unmarshal the URI is captured from the input document; on marshal it is
// replayed verbatim. As a result an unmarshal/marshal round trip preserves
// whichever URI the source manifest used.
//
// # Constructors
//
// A freshly built Signal/SpliceInfoSection/etc. value should be created via the
// New* helpers in this file. They populate XMLName with the canonical SCTE-35
// URI so that a newly constructed value marshals as a proper SCTE-35 element.
// A bare struct literal works too but the resulting element will be emitted
// without an xmlns="..." declaration, inheriting the surrounding DASH namespace
// — typically not what you want.
//
// # Two carriage forms inside <Event>
//
//	<Event ...>
//	  <Signal xmlns="http://www.scte.org/schemas/35">
//	    <SpliceInfoSection ...>...</SpliceInfoSection>
//	  </Signal>
//	</Event>
//
// is the form prescribed by SCTE 214-1 and produced by Unified Streaming,
// Broadpeak and similar packagers. AWS MediaTailor and a few others emit the
// SpliceInfoSection directly:
//
//	<Event ...>
//	  <scte35:SpliceInfoSection xmlns:scte35="http://www.scte.org/schemas/35" ...>...</scte35:SpliceInfoSection>
//	</Event>
//
// Both forms are valid: SCTE-35 declares SpliceInfoSection as a global element
// and the DASH MPD XSD's EventType has an <xs:any namespace="##other"
// processContents="lax"/> extension slot. EventType exposes them through two
// fields — Signal and SpliceInfoSection — see (*EventType).SpliceInfo for a
// uniform accessor.

// SCTE-35 namespace URIs.
const (
	// SCTE35Namespace is the canonical SCTE-35 XML namespace, declared as
	// targetNamespace by every XSD revision since 2020.
	SCTE35Namespace = "http://www.scte.org/schemas/35"
	// SCTE35Namespace2016 is the legacy URI still produced by a large
	// fraction of real-world packagers. Schema content is identical; the
	// types here capture whichever URI is present in the input.
	SCTE35Namespace2016 = "http://www.scte.org/schemas/35/2016"
)

// EventStream@schemeIdUri values for SCTE-35, defined by SCTE 214-1.
const (
	// SCTE35SchemeIdXML carries an XML representation of the SCTE-35
	// message: a <Signal> containing either <SpliceInfoSection> or
	// <Binary>, or one of those directly under <Event>.
	SCTE35SchemeIdXML = "urn:scte:scte35:2013:xml"
	// SCTE35SchemeIdXMLBin restricts the XML form to <Signal><Binary>...
	SCTE35SchemeIdXMLBin = "urn:scte:scte35:2014:xml+bin"
	// SCTE35SchemeIdBin is the raw binary form used for inband (emsg)
	// signalling, never carried in an MPD event.
	SCTE35SchemeIdBin = "urn:scte:scte35:2013:bin"
)

// Selected segmentation_type_id values from Table 23 of ANSI/SCTE 35 2023r2.
// Values not listed here (custom or rare) are still valid; segmentationTypeId
// on SegmentationDescriptorType accepts any uint8.
const (
	SegTypeNotIndicated                                uint8 = 0x00
	SegTypeContentIdentification                       uint8 = 0x01
	SegTypeProgramStart                                uint8 = 0x10
	SegTypeProgramEnd                                  uint8 = 0x11
	SegTypeProgramEarlyTermination                     uint8 = 0x12
	SegTypeProgramBreakaway                            uint8 = 0x13
	SegTypeProgramResumption                           uint8 = 0x14
	SegTypeProgramRunoverPlanned                       uint8 = 0x15
	SegTypeProgramRunoverUnplanned                     uint8 = 0x16
	SegTypeProgramOverlapStart                         uint8 = 0x17
	SegTypeProgramBlackoutOverride                     uint8 = 0x18
	SegTypeProgramJoin                                 uint8 = 0x19
	SegTypeChapterStart                                uint8 = 0x20
	SegTypeChapterEnd                                  uint8 = 0x21
	SegTypeBreakStart                                  uint8 = 0x22
	SegTypeBreakEnd                                    uint8 = 0x23
	SegTypeOpeningCreditStart                          uint8 = 0x24
	SegTypeOpeningCreditEnd                            uint8 = 0x25
	SegTypeClosingCreditStart                          uint8 = 0x26
	SegTypeClosingCreditEnd                            uint8 = 0x27
	SegTypeProviderAdvertisementStart                  uint8 = 0x30
	SegTypeProviderAdvertisementEnd                    uint8 = 0x31
	SegTypeDistributorAdvertisementStart               uint8 = 0x32
	SegTypeDistributorAdvertisementEnd                 uint8 = 0x33
	SegTypeProviderPlacementOpportunityStart           uint8 = 0x34
	SegTypeProviderPlacementOpportunityEnd             uint8 = 0x35
	SegTypeDistributorPlacementOpportunityStart        uint8 = 0x36
	SegTypeDistributorPlacementOpportunityEnd          uint8 = 0x37
	SegTypeProviderOverlayPlacementOpportunityStart    uint8 = 0x38
	SegTypeProviderOverlayPlacementOpportunityEnd      uint8 = 0x39
	SegTypeDistributorOverlayPlacementOpportunityStart uint8 = 0x3A
	SegTypeDistributorOverlayPlacementOpportunityEnd   uint8 = 0x3B
	SegTypeProviderPromoStart                          uint8 = 0x3C
	SegTypeProviderPromoEnd                            uint8 = 0x3D
	SegTypeDistributorPromoStart                       uint8 = 0x3E
	SegTypeDistributorPromoEnd                         uint8 = 0x3F
	SegTypeUnscheduledEventStart                       uint8 = 0x40
	SegTypeUnscheduledEventEnd                         uint8 = 0x41
	SegTypeAlternateContentOpportunityStart            uint8 = 0x42
	SegTypeAlternateContentOpportunityEnd              uint8 = 0x43
	SegTypeProviderAdBlockStart                        uint8 = 0x44
	SegTypeProviderAdBlockEnd                          uint8 = 0x45
	SegTypeDistributorAdBlockStart                     uint8 = 0x46
	SegTypeDistributorAdBlockEnd                       uint8 = 0x47
	SegTypeNetworkStart                                uint8 = 0x50
	SegTypeNetworkEnd                                  uint8 = 0x51
)

// SignalType is the <Signal> wrapper element specified by SCTE 214-1 for
// carrying an SCTE-35 message inside a DASH <Event>. It contains exactly one
// of SpliceInfoSection (parsed XML form) or Binary (base64 of the binary
// splice_info_section); see ANSI/SCTE 35 2023r2 section 12.
//
// XMLName preserves the namespace URI as seen in the input document. Use
// NewSignal to construct a fresh value with the canonical URI.
type SignalType struct {
	XMLName           xml.Name
	SpliceInfoSection *SpliceInfoSectionType `xml:"SpliceInfoSection,omitempty"`
	Binary            *BinaryType            `xml:"Binary,omitempty"`
}

// SpliceInfoSectionType maps the SpliceInfoSection element of clause 9.6 of
// ANSI/SCTE 35 2023r2. The XSD declares an xsd:choice between SpliceNull,
// SpliceSchedule, SpliceInsert, TimeSignal, BandwidthReservation and
// PrivateCommand — exactly one of those fields should be non-nil in a valid
// message. The descriptor slices implement the trailing xsd:choice
// maxOccurs="unbounded"; ordering between descriptor kinds is not preserved on
// re-marshal (see limitation #3 of the package README).
type SpliceInfoSectionType struct {
	XMLName                 xml.Name
	SapType                 *uint8                        `xml:"sapType,attr,omitempty"` // default 3
	PreRollMilliSeconds     uint32                        `xml:"preRollMilliSeconds,attr,omitempty"`
	PtsAdjustment           uint64                        `xml:"ptsAdjustment,attr,omitempty"`
	ProtocolVersion         *uint8                        `xml:"protocolVersion,attr,omitempty"` // fixed 0
	Tier                    *uint16                       `xml:"tier,attr,omitempty"`
	AnyAttrs                []xml.Attr                    `xml:",any,attr"`
	Ext                     *ExtType                      `xml:"Ext,omitempty"`
	EncryptedPacket         *EncryptedPacketType          `xml:"EncryptedPacket,omitempty"`
	SpliceNull              *SpliceNullType               `xml:"SpliceNull,omitempty"`
	SpliceSchedule          *SpliceScheduleType           `xml:"SpliceSchedule,omitempty"`
	SpliceInsert            *SpliceInsertType             `xml:"SpliceInsert,omitempty"`
	TimeSignal              *TimeSignalType               `xml:"TimeSignal,omitempty"`
	BandwidthReservation    *BandwidthReservationType     `xml:"BandwidthReservation,omitempty"`
	PrivateCommand          *PrivateCommandType           `xml:"PrivateCommand,omitempty"`
	AvailDescriptors        []*AvailDescriptorType        `xml:"AvailDescriptor"`
	DTMFDescriptors         []*DTMFDescriptorType         `xml:"DTMFDescriptor"`
	SegmentationDescriptors []*SegmentationDescriptorType `xml:"SegmentationDescriptor"`
	TimeDescriptors         []*TimeDescriptorType         `xml:"TimeDescriptor"`
	AudioDescriptors        []*AudioDescriptorType        `xml:"AudioDescriptor"`
	PrivateDescriptors      []*PrivateDescriptorType      `xml:"PrivateDescriptor"`
}

// BinaryType is the <Binary> element used by the urn:scte:scte35:2014:xml+bin
// scheme and as one branch of the SCTE-35 Signal group. Value is the base64
// encoding of the raw splice_info_section bytes.
type BinaryType struct {
	XMLName    xml.Name
	SignalType string `xml:"signalType,attr,omitempty"` // default "SpliceInfoSection"; "private:..." also allowed
	Value      string `xml:",chardata"`
}

// EncryptedPacketType is the EncryptedPacket child of SpliceInfoSection. The
// SCTE-35 XSD declares both attributes as required.
type EncryptedPacketType struct {
	XMLName             xml.Name
	EncryptionAlgorithm uint8      `xml:"encryptionAlgorithm,attr"`
	CwIndex             uint8      `xml:"cwIndex,attr"`
	AnyAttrs            []xml.Attr `xml:",any,attr"`
	Ext                 *ExtType   `xml:"Ext,omitempty"`
}

// SpliceNullType is the splice_null() command. Carries no data of its own.
type SpliceNullType struct {
	XMLName xml.Name
	Ext     *ExtType `xml:"Ext,omitempty"`
}

// SpliceScheduleType is the splice_schedule() command. It contains 0..255
// scheduled events, each with UTC splice times (xsd:dateTime).
type SpliceScheduleType struct {
	XMLName xml.Name
	Ext     *ExtType                   `xml:"Ext,omitempty"`
	Events  []*SpliceScheduleEventType `xml:"Event"`
}

// SpliceScheduleEventType is a single <Event> inside <SpliceSchedule>. The
// Program / Components choice mirrors the XSD; valid messages set exactly one.
type SpliceScheduleEventType struct {
	XMLName                    xml.Name
	SpliceEventId              *uint32                        `xml:"spliceEventId,attr,omitempty"`
	EventIdComplianceFlag      *bool                          `xml:"eventIdComplianceFlag,attr,omitempty"`
	SpliceEventCancelIndicator *bool                          `xml:"spliceEventCancelIndicator,attr,omitempty"` // default 0
	OutOfNetworkIndicator      *bool                          `xml:"outOfNetworkIndicator,attr,omitempty"`
	UniqueProgramId            *uint16                        `xml:"uniqueProgramId,attr,omitempty"`
	AvailNum                   *uint8                         `xml:"availNum,attr,omitempty"`
	AvailsExpected             *uint8                         `xml:"availsExpected,attr,omitempty"`
	AnyAttrs                   []xml.Attr                     `xml:",any,attr"`
	Ext                        *ExtType                       `xml:"Ext,omitempty"`
	Program                    *SpliceScheduleProgramType     `xml:"Program,omitempty"`
	Components                 []*SpliceScheduleComponentType `xml:"Component"`
	BreakDuration              *BreakDurationType             `xml:"BreakDuration,omitempty"`
}

// SpliceScheduleProgramType is the program-mode <Program> inside a scheduled
// event. The UTC splice time is required.
type SpliceScheduleProgramType struct {
	XMLName       xml.Name
	UtcSpliceTime DateTime   `xml:"utcSpliceTime,attr"`
	AnyAttrs      []xml.Attr `xml:",any,attr"`
	Ext           *ExtType   `xml:"Ext,omitempty"`
}

// SpliceScheduleComponentType is the component-mode <Component> inside a
// scheduled event. Both componentTag and utcSpliceTime are required.
type SpliceScheduleComponentType struct {
	XMLName       xml.Name
	ComponentTag  uint8      `xml:"componentTag,attr"`
	UtcSpliceTime DateTime   `xml:"utcSpliceTime,attr"`
	AnyAttrs      []xml.Attr `xml:",any,attr"`
	Ext           *ExtType   `xml:"Ext,omitempty"`
}

// SpliceInsertType is the splice_insert() command. The Program / Components
// choice mirrors the XSD; valid messages set exactly one.
type SpliceInsertType struct {
	XMLName                    xml.Name
	SpliceEventId              *uint32                      `xml:"spliceEventId,attr,omitempty"`
	SpliceEventCancelIndicator *bool                        `xml:"spliceEventCancelIndicator,attr,omitempty"` // default 0
	OutOfNetworkIndicator      *bool                        `xml:"outOfNetworkIndicator,attr,omitempty"`
	SpliceImmediateFlag        *bool                        `xml:"spliceImmediateFlag,attr,omitempty"`
	EventIdComplianceFlag      *bool                        `xml:"eventIdComplianceFlag,attr,omitempty"`
	UniqueProgramId            *uint16                      `xml:"uniqueProgramId,attr,omitempty"`
	AvailNum                   *uint8                       `xml:"availNum,attr,omitempty"`
	AvailsExpected             *uint8                       `xml:"availsExpected,attr,omitempty"`
	Ext                        *ExtType                     `xml:"Ext,omitempty"`
	Program                    *SpliceInsertProgramType     `xml:"Program,omitempty"`
	Components                 []*SpliceInsertComponentType `xml:"Component"`
	BreakDuration              *BreakDurationType           `xml:"BreakDuration,omitempty"`
}

// SpliceInsertProgramType is the program-mode <Program> inside a splice
// insert. SpliceTime is optional and absent when spliceImmediateFlag="true".
type SpliceInsertProgramType struct {
	XMLName    xml.Name
	AnyAttrs   []xml.Attr      `xml:",any,attr"`
	Ext        *ExtType        `xml:"Ext,omitempty"`
	SpliceTime *SpliceTimeType `xml:"SpliceTime,omitempty"`
}

// SpliceInsertComponentType is the component-mode <Component> inside a splice
// insert. componentTag is required; SpliceTime is optional.
type SpliceInsertComponentType struct {
	XMLName      xml.Name
	ComponentTag uint8           `xml:"componentTag,attr"`
	AnyAttrs     []xml.Attr      `xml:",any,attr"`
	Ext          *ExtType        `xml:"Ext,omitempty"`
	SpliceTime   *SpliceTimeType `xml:"SpliceTime,omitempty"`
}

// TimeSignalType is the time_signal() command. SpliceTime is required.
type TimeSignalType struct {
	XMLName    xml.Name
	Ext        *ExtType        `xml:"Ext,omitempty"`
	SpliceTime *SpliceTimeType `xml:"SpliceTime"`
}

// BandwidthReservationType is the bandwidth_reservation() command. It carries
// no data of its own.
type BandwidthReservationType struct {
	XMLName xml.Name
	Ext     *ExtType `xml:"Ext,omitempty"`
}

// PrivateCommandType is the private_command() command. identifier is the
// 32-bit owner ID; PrivateBytes carries the owner-defined payload as hex.
type PrivateCommandType struct {
	XMLName      xml.Name
	Identifier   uint32   `xml:"identifier,attr"`
	Ext          *ExtType `xml:"Ext,omitempty"`
	PrivateBytes string   `xml:"PrivateBytes,omitempty"` // xsd:hexBinary
}

// SpliceTimeType is the splice_time() structure carrying a PTS-domain time.
// ptsTime is optional (absent indicates a time-specified-flag of 0 in the
// binary form).
type SpliceTimeType struct {
	XMLName  xml.Name
	PtsTime  *uint64    `xml:"ptsTime,attr,omitempty"`
	AnyAttrs []xml.Attr `xml:",any,attr"`
	Ext      *ExtType   `xml:"Ext,omitempty"`
}

// BreakDurationType is the break_duration() structure shared by SpliceInsert
// and SpliceScheduleEvent. Both attributes are required.
type BreakDurationType struct {
	XMLName    xml.Name
	AutoReturn bool       `xml:"autoReturn,attr"`
	Duration   uint64     `xml:"duration,attr"`
	AnyAttrs   []xml.Attr `xml:",any,attr"`
	Ext        *ExtType   `xml:"Ext,omitempty"`
}

// AvailDescriptorType is the avail_descriptor() splice descriptor.
type AvailDescriptorType struct {
	XMLName         xml.Name
	ProviderAvailId uint32     `xml:"providerAvailId,attr"`
	AnyAttrs        []xml.Attr `xml:",any,attr"`
	Ext             *ExtType   `xml:"Ext,omitempty"`
}

// DTMFDescriptorType is the DTMF_descriptor() splice descriptor.
type DTMFDescriptorType struct {
	XMLName  xml.Name
	Preroll  *uint8     `xml:"preroll,attr,omitempty"`
	Chars    string     `xml:"chars,attr,omitempty"` // restricted to [0-9#*]+
	AnyAttrs []xml.Attr `xml:",any,attr"`
	Ext      *ExtType   `xml:"Ext,omitempty"`
}

// SegmentationDescriptorType is the segmentation_descriptor() splice
// descriptor — the workhorse for ad insertion signalling. See clause 10.3.3
// of ANSI/SCTE 35 2023r2 and Table 23 for segmentationTypeId values
// (constants SegType* are defined in this package).
type SegmentationDescriptorType struct {
	XMLName                          xml.Name
	SegmentationEventId              *uint32 `xml:"segmentationEventId,attr,omitempty"`
	SegmentationEventCancelIndicator *bool   `xml:"segmentationEventCancelIndicator,attr,omitempty"` // default 0
	// SegmentationEventIdComplianceIndicator captures the corresponding XSD
	// attribute. The 2020+ XSDs declare it with the literal name
	// "segmentationEventIdComplianceIndicatorbute" — a known SCTE-35 schema
	// typo. Real-world manifests use the corrected name without the trailing
	// "bute" and the XSD typo is treated as errata; we match the corrected
	// name for interop and document the discrepancy here.
	SegmentationEventIdComplianceIndicator *bool                                  `xml:"segmentationEventIdComplianceIndicator,attr,omitempty"`
	SegmentationDuration                   *uint64                                `xml:"segmentationDuration,attr,omitempty"` // < 2^40
	SegmentationTypeId                     *uint8                                 `xml:"segmentationTypeId,attr,omitempty"`
	SegmentNum                             *uint8                                 `xml:"segmentNum,attr,omitempty"`
	SegmentsExpected                       *uint8                                 `xml:"segmentsExpected,attr,omitempty"`
	SubSegmentNum                          *uint8                                 `xml:"subSegmentNum,attr,omitempty"`
	SubSegmentsExpected                    *uint8                                 `xml:"subSegmentsExpected,attr,omitempty"`
	AnyAttrs                               []xml.Attr                             `xml:",any,attr"`
	Ext                                    *ExtType                               `xml:"Ext,omitempty"`
	DeliveryRestrictions                   *DeliveryRestrictionsType              `xml:"DeliveryRestrictions,omitempty"`
	SegmentationUpids                      []*SegmentationUpidType                `xml:"SegmentationUpid"`
	Components                             []*SegmentationDescriptorComponentType `xml:"Component"`
}

// DeliveryRestrictionsType is the <DeliveryRestrictions> child element of
// SegmentationDescriptor. All four attributes are XSD-required.
type DeliveryRestrictionsType struct {
	XMLName                xml.Name
	WebDeliveryAllowedFlag bool       `xml:"webDeliveryAllowedFlag,attr"`
	NoRegionalBlackoutFlag bool       `xml:"noRegionalBlackoutFlag,attr"`
	ArchiveAllowedFlag     bool       `xml:"archiveAllowedFlag,attr"`
	DeviceRestrictions     uint8      `xml:"deviceRestrictions,attr"` // 0..3, see Table 21
	AnyAttrs               []xml.Attr `xml:",any,attr"`
	Ext                    *ExtType   `xml:"Ext,omitempty"`
}

// SegmentationUpidType is the <SegmentationUpid> child of
// SegmentationDescriptor. The element body is an xsd:token whose meaning is
// dictated by the SegmentationUpidType attribute and the optional Format
// attribute.
//
// Real-world manifests sometimes carry a "segmentationUpidLength" attribute
// that was removed from the 2023 XSD; AnyAttrs preserves it through round
// trips.
type SegmentationUpidType struct {
	XMLName                xml.Name
	SegmentationUpidType   *uint8     `xml:"segmentationUpidType,attr,omitempty"`
	FormatIdentifier       *uint32    `xml:"formatIdentifier,attr,omitempty"`
	SegmentationUpidFormat string     `xml:"segmentationUpidFormat,attr,omitempty"` // text|hexbinary|base-64|private:...
	AnyAttrs               []xml.Attr `xml:",any,attr"`
	Value                  string     `xml:",chardata"`
}

// SegmentationDescriptorComponentType is the <Component> child of
// SegmentationDescriptor (distinct from the splice-time Component variants).
type SegmentationDescriptorComponentType struct {
	XMLName      xml.Name
	ComponentTag uint8      `xml:"componentTag,attr"`
	PtsOffset    uint64     `xml:"ptsOffset,attr"` // PTSType
	AnyAttrs     []xml.Attr `xml:",any,attr"`
	Ext          *ExtType   `xml:"Ext,omitempty"`
}

// TimeDescriptorType is the time_descriptor() splice descriptor (TAI offset).
type TimeDescriptorType struct {
	XMLName    xml.Name
	TaiSeconds *uint64    `xml:"taiSeconds,attr,omitempty"` // < 2^48
	TaiNs      *uint32    `xml:"taiNs,attr,omitempty"`
	UtcOffset  *uint16    `xml:"utcOffset,attr,omitempty"`
	AnyAttrs   []xml.Attr `xml:",any,attr"`
	Ext        *ExtType   `xml:"Ext,omitempty"`
}

// AudioDescriptorType is the audio_descriptor() splice descriptor. It carries
// 1..15 AudioChannel children describing the audio service streams. Note that
// the SCTE-35 XSD spells these attributes with mixed case (BitStreamMode,
// NumChannels, FullSrvcAudio).
type AudioDescriptorType struct {
	XMLName       xml.Name
	AnyAttrs      []xml.Attr          `xml:",any,attr"`
	Ext           *ExtType            `xml:"Ext,omitempty"`
	AudioChannels []*AudioChannelType `xml:"AudioChannel"`
}

// AudioChannelType is a single <AudioChannel> entry inside an
// AudioDescriptor. ISOCode, BitStreamMode, NumChannels and FullSrvcAudio are
// XSD-required; componentTag is optional (0xFF means "use PMT order").
type AudioChannelType struct {
	XMLName       xml.Name
	ISOCode       string     `xml:"ISOCode,attr"` // xsd:language
	BitStreamMode uint8      `xml:"BitStreamMode,attr"`
	NumChannels   uint8      `xml:"NumChannels,attr"`
	FullSrvcAudio uint8      `xml:"FullSrvcAudio,attr"` // 0 or 1
	ComponentTag  *uint8     `xml:"componentTag,attr,omitempty"`
	AnyAttrs      []xml.Attr `xml:",any,attr"`
	Ext           *ExtType   `xml:"Ext,omitempty"`
}

// PrivateDescriptorType is the private_descriptor() catch-all owned by the
// identifier value (per SMPTE Registration Authority).
type PrivateDescriptorType struct {
	XMLName       xml.Name
	DescriptorTag uint8      `xml:"descriptorTag,attr"`
	Identifier    uint32     `xml:"identifier,attr"`
	AnyAttrs      []xml.Attr `xml:",any,attr"`
	Ext           *ExtType   `xml:"Ext,omitempty"`
	PrivateBytes  string     `xml:"PrivateBytes,omitempty"` // xsd:hexBinary
}

// ExtType is the SCTE-35 extensibility element. It can appear inside almost
// every SCTE-35 type and accepts arbitrary attributes and child elements
// (xs:any processContents="lax"). InnerXML captures the raw nested XML so
// extension content survives an unmarshal/marshal round trip; AnyAttrs
// captures arbitrary attributes likewise.
type ExtType struct {
	XMLName  xml.Name
	AnyAttrs []xml.Attr `xml:",any,attr"`
	InnerXML string     `xml:",innerxml"`
}

// canonicalSCTE35Name returns an xml.Name with the canonical SCTE-35 namespace
// and the given local name. Used by the New* constructors below.
func canonicalSCTE35Name(local string) xml.Name {
	return xml.Name{Space: SCTE35Namespace, Local: local}
}

// NewSignal returns a Signal with XMLName preset to the canonical SCTE-35
// namespace. Populate SpliceInfoSection or Binary (not both) before marshal.
func NewSignal() *SignalType {
	return &SignalType{XMLName: canonicalSCTE35Name("Signal")}
}

// NewSpliceInfoSection returns a SpliceInfoSection with XMLName preset to the
// canonical SCTE-35 namespace.
func NewSpliceInfoSection() *SpliceInfoSectionType {
	return &SpliceInfoSectionType{XMLName: canonicalSCTE35Name("SpliceInfoSection")}
}

// NewBinary returns a Binary element (the urn:scte:scte35:2014:xml+bin form)
// with XMLName preset to the canonical SCTE-35 namespace and Value set to v.
func NewBinary(v string) *BinaryType {
	return &BinaryType{XMLName: canonicalSCTE35Name("Binary"), Value: v}
}

// NewSpliceNull returns a SpliceNull command in the canonical namespace.
func NewSpliceNull() *SpliceNullType {
	return &SpliceNullType{XMLName: canonicalSCTE35Name("SpliceNull")}
}

// NewSpliceSchedule returns a SpliceSchedule command in the canonical namespace.
func NewSpliceSchedule() *SpliceScheduleType {
	return &SpliceScheduleType{XMLName: canonicalSCTE35Name("SpliceSchedule")}
}

// NewSpliceInsert returns a SpliceInsert command in the canonical namespace.
func NewSpliceInsert() *SpliceInsertType {
	return &SpliceInsertType{XMLName: canonicalSCTE35Name("SpliceInsert")}
}

// NewTimeSignal returns a TimeSignal command in the canonical namespace.
func NewTimeSignal() *TimeSignalType {
	return &TimeSignalType{XMLName: canonicalSCTE35Name("TimeSignal")}
}

// NewBandwidthReservation returns a BandwidthReservation command in the
// canonical namespace.
func NewBandwidthReservation() *BandwidthReservationType {
	return &BandwidthReservationType{XMLName: canonicalSCTE35Name("BandwidthReservation")}
}

// NewPrivateCommand returns a PrivateCommand with the given owner identifier
// and XMLName preset to the canonical SCTE-35 namespace.
func NewPrivateCommand(identifier uint32) *PrivateCommandType {
	return &PrivateCommandType{XMLName: canonicalSCTE35Name("PrivateCommand"), Identifier: identifier}
}

// NewSpliceTime returns a SpliceTime element in the canonical namespace.
func NewSpliceTime() *SpliceTimeType {
	return &SpliceTimeType{XMLName: canonicalSCTE35Name("SpliceTime")}
}

// NewBreakDuration returns a BreakDuration element in the canonical namespace
// with the given autoReturn and duration values.
func NewBreakDuration(autoReturn bool, duration uint64) *BreakDurationType {
	return &BreakDurationType{
		XMLName:    canonicalSCTE35Name("BreakDuration"),
		AutoReturn: autoReturn,
		Duration:   duration,
	}
}

// NewSegmentationDescriptor returns a SegmentationDescriptor in the canonical
// namespace.
func NewSegmentationDescriptor() *SegmentationDescriptorType {
	return &SegmentationDescriptorType{XMLName: canonicalSCTE35Name("SegmentationDescriptor")}
}

// NewSegmentationUpid returns a SegmentationUpid in the canonical namespace
// with the given value.
func NewSegmentationUpid(value string) *SegmentationUpidType {
	return &SegmentationUpidType{XMLName: canonicalSCTE35Name("SegmentationUpid"), Value: value}
}

// IsSCTE35 reports whether es carries SCTE-35 events (i.e. its SchemeIdUri
// is one of the URIs defined in SCTE 214-1 for SCTE-35 carriage in DASH).
func (es *EventStreamType) IsSCTE35() bool {
	switch string(es.SchemeIdUri) {
	case SCTE35SchemeIdXML, SCTE35SchemeIdXMLBin, SCTE35SchemeIdBin:
		return true
	}
	return false
}

// SpliceInfo returns the SCTE-35 SpliceInfoSection of the event, whether it
// is wrapped in a Signal element (per SCTE 214-1) or attached directly to
// Event (per AWS MediaTailor and other tooling). Returns nil if neither is
// present.
func (e *EventType) SpliceInfo() *SpliceInfoSectionType {
	if e.Signal != nil && e.Signal.SpliceInfoSection != nil {
		return e.Signal.SpliceInfoSection
	}
	return e.SpliceInfoSection
}
