package mpd

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Eyevinn/dash-mpd/xml"
	"github.com/barkimedes/go-deepcopy"
)

const (
	STATIC_TYPE  = "static"
	DYNAMIC_TYPE = "dynamic"
)

const (
	DASH_NAMESPACE                         = "urn:mpeg:dash:schema:mpd:2011"
	PROFILE_LIVE                           = "urn:mpeg:dash:profile:isoff-live:2011"
	PROFILE_ONDEMAND                       = "urn:mpeg:dash:profile:isoff-on-demand:2011"
	AUDIO_CHANNEL_CONFIGURATION_MPEG_DASH  = "urn:mpeg:dash:23003:3:audio_channel_configuration:2011"
	AUDIO_CHANNEL_CONFIGURATION_MPEG_DOLBY = "tag:dolby.com,2014:dash:audio_channel_configuration:2011"
	MIME_TYPE_VIDEO_MP4                    = "video/mp4"
	MIME_TYPE_AUDIO_MP4                    = "audio/mp4"
	MIME_TYPE_SUBTITLE_VTT                 = "text/vtt"
	MIME_TYPE_TTML                         = "application/ttml+xml"
)

// MPD is MPEG-DASH Media Presentation Description (MPD) as defined in ISO/IEC 23009-1 5'th edition.
//
// The tree of structs is generated from the corresponding XML Schema at https://github.com/MPEGGroup/DASHSchema
// but fine-tuned manually to handle default cases, listing enumerals, name space for xlink etc.
type MPD struct {
	XMLName                    xml.Name                   `xml:"MPD"`
	XMLNs                      string                     `xml:"xmlns,attr,omitempty"`
	SchemaLocation             string                     `xml:"http://www.w3.org/2001/XMLSchema-instance xsi:schemaLocation,attr,omitempty"`
	Id                         string                     `xml:"id,attr,omitempty"`
	Profiles                   ListOfProfilesType         `xml:"profiles,attr"`
	Type                       *string                    `xml:"type,attr,omitempty"` // Optional with default "static"
	AvailabilityStartTime      DateTime                   `xml:"availabilityStartTime,attr,omitempty"`
	AvailabilityEndTime        DateTime                   `xml:"availabilityEndTime,attr,omitempty"`
	PublishTime                DateTime                   `xml:"publishTime,attr,omitempty"`
	MediaPresentationDuration  *Duration                  `xml:"mediaPresentationDuration,attr"`
	MinimumUpdatePeriod        *Duration                  `xml:"minimumUpdatePeriod,attr"`
	MinBufferTime              *Duration                  `xml:"minBufferTime,attr"`
	TimeShiftBufferDepth       *Duration                  `xml:"timeShiftBufferDepth,attr"`
	SuggestedPresentationDelay *Duration                  `xml:"suggestedPresentationDelay,attr"`
	MaxSegmentDuration         *Duration                  `xml:"maxSegmentDuration,attr"`
	MaxSubsegmentDuration      *Duration                  `xml:"maxSubsegmentDuration,attr"`
	ProgramInformation         []*ProgramInformationType  `xml:"ProgramInformation"`
	BaseURL                    []*BaseURLType             `xml:"BaseURL"`
	Location                   []AnyURI                   `xml:"Location"`
	PatchLocation              []*PatchLocationType       `xml:"PatchLocation"`
	ServiceDescription         []*ServiceDescriptionType  `xml:"ServiceDescription"`
	InitializationSet          []*InitializationSetType   `xml:"InitializationSet"`
	InitializationGroup        []*UIntVWithIDType         `xml:"InitializationGroup"`
	InitializationPresentation []*UIntVWithIDType         `xml:"InitializationPresentation"`
	ContentProtection          []*ContentProtectionType   `xml:"ContentProtection"`
	Periods                    []*Period                  `xml:"Period"`
	Metrics                    []*MetricsType             `xml:"Metrics"`
	EssentialProperties        []*DescriptorType          `xml:"EssentialProperty"`
	SupplementalProperties     []*DescriptorType          `xml:"SupplementalProperty"`
	UTCTimings                 []*DescriptorType          `xml:"UTCTiming"`
	LeapSecondInformation      *LeapSecondInformationType `xml:"LeapSecondInformation"`
}

// GetType returns static or dynamic.
func (m *MPD) GetType() string {
	if m.Type == nil {
		return STATIC_TYPE
	}
	return *m.Type
}

// Clone creates a deep copy of mpd and sets parents.
func Clone(mpd *MPD) *MPD {
	cm := deepcopy.MustAnything(mpd).(*MPD)
	cm.SetParents()
	return cm
}

// PatchLocationType is Patch Location Type.
type PatchLocationType struct {
	XMLName xml.Name `xml:"PatchLocation"`
	Ttl     float64  `xml:"ttl,attr,omitempty"`
	Value   AnyURI   `xml:",chardata"`
}

type Period struct {
	XMLName                xml.Name                  `xml:"Period"`
	XlinkHref              string                    `xml:"http://www.w3.org/1999/xlink xlink:href,attr,omitempty"`
	XlinkActuate           string                    `xml:"http://www.w3.org/1999/xlink xlink:actuate,attr,omitempty"` // default = "onRequest"
	XlinkType              string                    `xml:"http://www.w3.org/1999/xlink xlink:type,attr,omitempty"`    // fixed = "simple"
	XlinkShow              string                    `xml:"http://www.w3.org/1999/xlink xlink:show,attr,omitempty"`    // fixed = "embed"
	Id                     string                    `xml:"id,attr,omitempty"`
	Start                  *Duration                 `xml:"start,attr"` // Mandatory for dynamic manifests. default = 0
	Duration               *Duration                 `xml:"duration,attr"`
	BitstreamSwitching     *bool                     `xml:"bitstreamSwitching,attr"`
	BaseURLs               []*BaseURLType            `xml:"BaseURL"`
	SegmentBase            *SegmentBaseType          `xml:"SegmentBase"`
	SegmentList            *SegmentListType          `xml:"SegmentList"`
	SegmentTemplate        *SegmentTemplateType      `xml:"SegmentTemplate"`
	AssetIdentifier        *DescriptorType           `xml:"AssetIdentifier"`
	EventStreams           []*EventStreamType        `xml:"EventStream"`
	ServiceDescriptions    []*ServiceDescriptionType `xml:"ServiceDescription"`
	ContentProtections     []*ContentProtectionType  `xml:"ContentProtection"`
	AdaptationSets         []*AdaptationSetType      `xml:"AdaptationSet"`
	Subsets                []*SubsetType             `xml:"Subset"`
	SupplementalProperties []*DescriptorType         `xml:"SupplementalProperty"`
	EmptyAdaptationSets    []*AdaptationSetType      `xml:"EmptyAdaptationSet"`
	GroupLabels            []*LabelType              `xml:"GroupLabel"`
	Preselections          []*PreselectionType       `xml:"Preselection"`
	parent                 *MPD                      `xml:"-"`
}

func (p *Period) SetParent(m *MPD) {
	p.parent = m
}

func (p *Period) Parent() *MPD {
	return p.parent
}

// EventStreamType is EventStream or InbandEventStream.
// It therefore has no XMLName.
type EventStreamType struct {
	XlinkHref              string       `xml:"http://www.w3.org/1999/xlink xlink:href,attr,omitempty"`
	XlinkActuate           string       `xml:"http://www.w3.org/1999/xlink xlink:actuate,attr,omitempty"` // default = "onRequest"
	XlinkType              string       `xml:"http://www.w3.org/1999/xlink xlink:type,attr,omitempty"`    // fixed = "simple"
	XlinkShow              string       `xml:"http://www.w3.org/1999/xlink xlink:show,attr,omitempty"`    // fixed = "embed"
	SchemeIdUri            AnyURI       `xml:"schemeIdUri,attr"`
	Value                  string       `xml:"value,attr,omitempty"`
	Timescale              *uint32      `xml:"timescale,attr"`                        // default = 1
	PresentationTimeOffset uint64       `xml:"presentationTimeOffset,attr,omitempty"` // default is 0
	Events                 []*EventType `xml:"Event"`
}

// EventType is Event.
type EventType struct {
	XMLName          xml.Name            `xml:"Event"`
	PresentationTime uint64              `xml:"presentationTime,attr,omitempty"` // default is 0
	Duration         uint64              `xml:"duration,attr,omitempty"`
	Id               uint32              `xml:"id,attr"`
	ContentEncoding  ContentEncodingType `xml:"contentEncoding,attr,omitempty"`
	MessageData      string              `xml:"messageData,attr,omitempty"`
}

// InitializationSetType is Initialization Set.
type InitializationSetType struct {
	XMLName         xml.Name               `xml:"InitializationSet"`
	XlinkHref       string                 `xml:"http://www.w3.org/1999/xlink xlink:href,attr,omitempty"`
	XlinkActuate    string                 `xml:"http://www.w3.org/1999/xlink xlink:actuate,attr,omitempty"` // default = "onRequest"
	XlinkType       string                 `xml:"http://www.w3.org/1999/xlink xlink:type,attr,omitempty"`    // fixed = "simple"
	Id              uint32                 `xml:"id,attr"`
	InAllPeriods    *bool                  `xml:"inAllPeriods,attr"` // default is true
	ContentType     RFC6838ContentTypeType `xml:"contentType,attr,omitempty"`
	Par             RatioType              `xml:"par,attr,omitempty"`
	MaxWidth        uint32                 `xml:"maxWidth,attr,omitempty"`
	MaxHeight       uint32                 `xml:"maxHeight,attr,omitempty"`
	MaxFrameRate    string                 `xml:"maxFrameRate,attr,omitempty"`
	Initialization  AnyURI                 `xml:"initialization,attr,omitempty"`
	Accessibilities []*DescriptorType      `xml:"Accessibility"`
	Roles           []*DescriptorType      `xml:"Role"`
	Ratings         []*DescriptorType      `xml:"Rating"`
	Viewpoints      []*DescriptorType      `xml:"Viewpoint"`
	RepresentationBaseType
}

// ServiceDescriptionType is Service Description.
type ServiceDescriptionType struct {
	XMLName             xml.Name                  `xml:"ServiceDescription"`
	Id                  uint32                    `xml:"id,attr"`
	Scopes              []*DescriptorType         `xml:"Scope"`
	Latencies           []*LatencyType            `xml:"Latency"`
	PlaybackRates       []*PlaybackRateType       `xml:"PlaybackRate"`
	OperatingQualities  []*OperatingQualityType   `xml:"OperatingQuality"`
	OperatingBandwidths []*OperatingBandwidthType `xml:"OperatingBandwidth"`
}

// LatencyType is Service Description Latency (Annex K.4.2.2).
type LatencyType struct {
	XMLName          xml.Name               `xml:"Latency"`
	ReferenceId      uint32                 `xml:"referenceId,attr"`
	Target           *uint32                `xml:"target,attr"`
	Max              *uint32                `xml:"max,attr"`
	Min              *uint32                `xml:"min,attr"`
	QualityLatencies []*UIntPairsWithIDType `xml:"QualityLatency"`
}

// PlaybackRateType is Service Description Playback Rate.
type PlaybackRateType struct {
	XMLName xml.Name `xml:"PlaybackRate"`
	Max     float64  `xml:"max,attr,omitempty"`
	Min     float64  `xml:"min,attr,omitempty"`
}

// OperatingQualityType is Service Description Operating Quality.
type OperatingQualityType struct {
	XMLName       xml.Name `xml:"OperatingQuality"`
	MediaType     string   `xml:"mediaType,attr,omitempty"` // default is "any"
	Min           uint32   `xml:"min,attr,omitempty"`
	Max           uint32   `xml:"max,attr,omitempty"`
	Target        uint32   `xml:"target,attr,omitempty"`
	Type          AnyURI   `xml:"type,attr,omitempty"`
	MaxDifference uint32   `xml:"maxDifference,attr,omitempty"`
}

// OperatingBandwidthType is Service Description Operating Bandwidth.
type OperatingBandwidthType struct {
	XMLName   xml.Name `xml:"OperatingBandwidth"`
	MediaType string   `xml:"mediaType,attr,omitempty"` // default is "all"
	Min       uint32   `xml:"min,attr,omitempty"`
	Max       uint32   `xml:"max,attr,omitempty"`
	Target    uint32   `xml:"target,attr,omitempty"`
}

// UIntPairsWithIDType is UInt Pairs With ID.
type UIntPairsWithIDType struct {
	XMLName xml.Name `xml:"UIntPairsWithID"`
	Type    AnyURI   `xml:"type,attr,omitempty"`
	*UIntVectorType
}

// UIntVWithIDType is UInt Vector With ID
type UIntVWithIDType struct {
	XMLName     xml.Name               `xml:"UIntVWithID"`
	Id          uint32                 `xml:"id,attr"`
	Profiles    ListOfProfilesType     `xml:"profiles,attr,omitempty"`
	ContentType RFC6838ContentTypeType `xml:"contentType,attr,omitempty"`
	UIntVectorType
}

// AdaptationSetType is AdaptationSet or EmptyAdaptationSet.
// Note that XMLName is not set, since the same structure is used also for EmptyAdaptationSet.
type AdaptationSetType struct {
	XlinkHref               string                  `xml:"xlink:href,attr,omitempty"`
	XlinkActuate            string                  `xml:"xlink:actuate,attr,omitempty"` // default is "onRequest"
	XlinkType               string                  `xml:"xlink:type,attr,omitempty"`    // fixed "simple"
	XlinkShow               string                  `xml:"xlink:show,attr,omitempty"`    // fixed "embed"
	Id                      *uint32                 `xml:"id,attr"`
	Group                   uint32                  `xml:"group,attr,omitempty"`
	Lang                    string                  `xml:"lang,attr,omitempty"`
	ContentType             RFC6838ContentTypeType  `xml:"contentType,attr,omitempty"`
	Par                     RatioType               `xml:"par,attr,omitempty"`
	MinBandwidth            uint32                  `xml:"minBandwidth,attr,omitempty"`
	MaxBandwidth            uint32                  `xml:"maxBandwidth,attr,omitempty"`
	MinWidth                uint32                  `xml:"minWidth,attr,omitempty"`
	MaxWidth                uint32                  `xml:"maxWidth,attr,omitempty"`
	MinHeight               uint32                  `xml:"minHeight,attr,omitempty"`
	MaxHeight               uint32                  `xml:"maxHeight,attr,omitempty"`
	MinFrameRate            string                  `xml:"minFrameRate,attr,omitempty"`
	MaxFrameRate            string                  `xml:"maxFrameRate,attr,omitempty"`
	SegmentAlignment        bool                    `xml:"segmentAlignment,attr,omitempty"`        // default = false
	SubsegmentAlignment     bool                    `xml:"subsegmentAlignment,attr,omitempty"`     // default = false
	SubsegmentStartsWithSAP uint32                  `xml:"subsegmentStartsWithSAP,attr,omitempty"` // default = 0
	BitstreamSwitching      *bool                   `xml:"bitstreamSwitching,attr"`
	InitializationSetRef    *UIntVectorType         `xml:"initializationSetRef,attr,omitempty"`
	InitializationPrincipal AnyURI                  `xml:"initializationPrincipal,attr,omitempty"`
	Accessibilities         []*DescriptorType       `xml:"Accessibility"`
	Roles                   []*DescriptorType       `xml:"Role"`
	Ratings                 []*DescriptorType       `xml:"Rating"`
	Viewpoints              []*DescriptorType       `xml:"Viewpoint"`
	ContentComponents       []*ContentComponentType `xml:"ContentComponent"`
	BaseURLs                []*BaseURLType          `xml:"BaseURL"`
	SegmentBase             *SegmentBaseType        `xml:"SegmentBase"`
	SegmentList             *SegmentListType        `xml:"SegmentList"`
	SegmentTemplate         *SegmentTemplateType    `xml:"SegmentTemplate"`
	Representations         []*RepresentationType   `xml:"Representation"`
	parent                  *Period                 `xml:"-"`
	RepresentationBaseType
}

func (a *AdaptationSetType) SetParent(p *Period) {
	a.parent = p
}

func (a *AdaptationSetType) Parent() *Period {
	return a.parent
}

// Clone returns a deep copy of the AdaptationSet with parent links set.
func (a *AdaptationSetType) Clone() *AdaptationSetType {
	ac := deepcopy.MustAnything(a).(*AdaptationSetType)
	ac.SetParents()
	return ac
}

// GetSegmentTemplate returns the segment template of the AdaptationSet or nil if not set.
func (a *AdaptationSetType) GetSegmentTemplate() *SegmentTemplateType {
	return a.SegmentTemplate
}

// GetMimeType returns the mime type.
func (a *AdaptationSetType) GetMimeType() string {
	return a.MimeType
}

// GetCodecs returns codecs string of the adaptation set.
func (a *AdaptationSetType) GetCodecs() string {
	return a.Codecs
}

// GetRepresentations returns the ContentProtections of the adaptation set.
func (a *AdaptationSetType) GetContentProtections() []*ContentProtectionType {
	return a.ContentProtections
}

// ContentComponentType is Content Component.
type ContentComponentType struct {
	XMLName         xml.Name               `xml:"ContentComponent"`
	Id              *uint32                `xml:"id,attr"`
	Lang            string                 `xml:"lang,attr,omitempty"`
	ContentType     RFC6838ContentTypeType `xml:"contentType,attr,omitempty"`
	Par             RatioType              `xml:"par,attr,omitempty"`
	Tag             string                 `xml:"tag,attr,omitempty"`
	Accessibilities []*DescriptorType      `xml:"Accessibility"`
	Roles           []*DescriptorType      `xml:"Role"`
	Ratings         []*DescriptorType      `xml:"Rating"`
	Viewpoints      []*DescriptorType      `xml:"Viewpoint"`
}

// RepresentationType is Representation.
type RepresentationType struct {
	XMLName                xml.Name                 `xml:"Representation"`
	Id                     string                   `xml:"id,attr"`
	Bandwidth              uint32                   `xml:"bandwidth,attr"`
	QualityRanking         *uint32                  `xml:"qualityRanking,attr,omitempty"`
	DependencyId           *StringVectorType        `xml:"dependencyId,attr,omitempty"`
	AssociationId          *StringVectorType        `xml:"associationId,attr,omitempty"`
	AssociationType        *ListOf4CCType           `xml:"associationType,attr,omitempty"`
	MediaStreamStructureId *StringVectorType        `xml:"mediaStreamStructureId,attr,omitempty"`
	BaseURLs               []*BaseURLType           `xml:"BaseURL"`
	ExtendedBandwidths     []*ExtendedBandwidthType `xml:"ExtendedBandwidth"`
	SubRepresentations     []*SubRepresentationType `xml:"SubRepresentation"`
	SegmentBase            *SegmentBaseType         `xml:"SegmentBase"`
	SegmentList            *SegmentListType         `xml:"SegmentList"`
	SegmentTemplate        *SegmentTemplateType     `xml:"SegmentTemplate"`
	parent                 *AdaptationSetType       `xml:"-"` // adaptation set
	RepresentationBaseType
}

func (r *RepresentationType) SetParent(p *AdaptationSetType) {
	r.parent = p
}

func (r *RepresentationType) Parent() *AdaptationSetType {
	return r.parent
}

// SetSegmentBase sets SegmentBaseType for RepresentationType.
func (r *RepresentationType) SetSegmentBase(initSize, sidxSize uint32, indexRangeExact bool) {
	initRange := fmt.Sprintf("0-%d", initSize-1)
	indexRange := fmt.Sprintf("%d-%d", initSize, initSize+sidxSize-1)
	r.SegmentBase = &SegmentBaseType{
		IndexRange: indexRange,
		Initialization: &URLType{
			Range: initRange,
		},
	}
	if indexRangeExact {
		r.SegmentBase.IndexRangeExact = true
	}
}

// ExtendedBandwidthType is Extended Bandwidth Model
type ExtendedBandwidthType struct {
	XMLName    xml.Name         `xml:"ExtendedBandwidth"`
	Vbr        bool             `xml:"vbr,attr,omitempty"` // default is false
	ModelPairs []*ModelPairType `xml:"ModelPair"`
}

// ModelPairType is Model Pair
type ModelPairType struct {
	XMLName    xml.Name  `xml:"ModelPair"`
	BufferTime *Duration `xml:"bufferTime,attr"`
	Bandwidth  uint32    `xml:"bandwidth,attr"`
}

// SubRepresentationType is SubRepresentation
type SubRepresentationType struct {
	XMLName          xml.Name            `xml:"SubRepresentation"`
	Level            *uint32             `xml:"level,attr,omitempty"`
	DependencyLevel  *UIntVectorType     `xml:"dependencyLevel,attr,omitempty"`
	Bandwidth        uint32              `xml:"bandwidth,attr,omitempty"`
	ContentComponent *StringVectorType   `xml:"contentComponent,attr,omitempty"`
	parent           *RepresentationType `xml:"-"`
	RepresentationBaseType
}

func (s *SubRepresentationType) SetParent(p *RepresentationType) {
	s.parent = p
}

func (s *SubRepresentationType) Parent() *RepresentationType {
	return s.parent
}

// RepresentationBaseType is Representation base (common attributes and elements).
type RepresentationBaseType struct {
	Profiles                   ListOfProfilesType           `xml:"profiles,attr,omitempty"`
	Width                      uint32                       `xml:"width,attr,omitempty"`
	Height                     uint32                       `xml:"height,attr,omitempty"`
	Sar                        RatioType                    `xml:"sar,attr,omitempty"`
	FrameRate                  FrameRateType                `xml:"frameRate,attr,omitempty"`
	AudioSamplingRate          *UIntVectorType              `xml:"audioSamplingRate,attr,omitempty"`
	MimeType                   string                       `xml:"mimeType,attr,omitempty"`
	SegmentProfiles            *ListOf4CCType               `xml:"segmentProfiles,attr,omitempty"`
	Codecs                     string                       `xml:"codecs,attr,omitempty"`
	ContainerProfiles          *ListOf4CCType               `xml:"containerProfiles,attr,omitempty"`
	MaximumSAPPeriod           float64                      `xml:"maximumSAPPeriod,attr,omitempty"`
	StartWithSAP               uint32                       `xml:"startWithSAP,attr,omitempty"`
	MaxPlayoutRate             float64                      `xml:"maxPlayoutRate,attr,omitempty"`
	CodingDependency           *bool                        `xml:"codingDependency,attr,omitempty"`
	ScanType                   VideoScanType                `xml:"scanType,attr,omitempty"`
	SelectionPriority          *uint32                      `xml:"selectionPriority,attr"` // default = 1
	Tag                        string                       `xml:"tag,attr,omitempty"`
	FramePackings              []*DescriptorType            `xml:"FramePacking"`
	AudioChannelConfigurations []*DescriptorType            `xml:"AudioChannelConfiguration"`
	ContentProtections         []*ContentProtectionType     `xml:"ContentProtection"`
	OutputProtection           *DescriptorType              `xml:"OutputProtection"`
	EssentialProperties        []*DescriptorType            `xml:"EssentialProperty"`
	SupplementalProperties     []*DescriptorType            `xml:"SupplementalProperty"`
	InbandEventStreams         []*EventStreamType           `xml:"InbandEventStream"`
	Switchings                 []*SwitchingType             `xml:"Switching"`
	RandomAccesses             []*RandomAccessType          `xml:"RandomAccess"`
	GroupLabels                []*LabelType                 `xml:"GroupLabel"`
	Labels                     []*LabelType                 `xml:"Label"`
	ProducerReferenceTimes     []*ProducerReferenceTimeType `xml:"ProducerReferenceTime"`
	ContentPopularityRates     []*ContentPopularityRateType `xml:"ContentPopularityRate"`
	Resyncs                    []*ResyncType                `xml:"Resync"`
}

func (r *RepresentationType) GetSegmentTemplate() *SegmentTemplateType {
	if r.SegmentTemplate == nil {
		return r.parent.GetSegmentTemplate()
	}
	return r.SegmentTemplate
}

// GetInit returns the representation's initialization URI with replaced identifiers.
//
// TODO: Apply BaseURLs
func (r *RepresentationType) GetInit() (string, error) {
	a := r.parent
	if a == nil {
		return "", ErrParentNotSet
	}
	st := r.GetSegmentTemplate()
	if st == nil {
		return "", ErrSegmentTemplateNotSet
	}
	initialization := st.Initialization
	initialization = strings.ReplaceAll(initialization, "$RepresentationID$", r.Id)
	initialization = strings.ReplaceAll(initialization, "$Bandwidth$", strconv.Itoa(int(r.Bandwidth)))
	return initialization, nil
}

// GetRepMedia returns the representation's media path with replaced ID and bandwidth identifiers.
//
// TODO: Apply BaseURLs.
func (r *RepresentationType) GetMedia() (string, error) {
	a := r.parent
	if a == nil || r == nil {
		return "", ErrParentNotSet
	}
	st := r.GetSegmentTemplate()
	if st == nil {
		return "", ErrSegmentTemplateNotSet
	}
	media := st.Media
	media = strings.ReplaceAll(media, "$RepresentationID$", r.Id)
	media = strings.ReplaceAll(media, "$Bandwidth$", strconv.Itoa(int(r.Bandwidth)))

	return media, nil
}

// GetMimeType returns the representation's or its parent's mime type.
func (r *RepresentationType) GetMimeType() string {
	if r.MimeType == "" {
		return r.parent.GetMimeType()
	}
	return r.MimeType
}

// GetCodecs returns the representation's or its parent's codecs string.
func (r *RepresentationType) GetCodecs() string {
	if r.Codecs == "" {
		return r.parent.GetCodecs()
	}
	return r.Codecs
}

// GetContentProtections returns the representation's or its parent's content protections.
func (r *RepresentationType) GetContentProtections() []*ContentProtectionType {
	if len(r.ContentProtections) == 0 {
		return r.parent.GetContentProtections()
	}
	return r.ContentProtections
}

// ContentProtectionType is Content Protection.
type ContentProtectionType struct {
	XMLName    xml.Name `xml:"ContentProtection"`
	Robustness string   `xml:"robustness,attr,omitempty"`
	RefId      string   `xml:"refId,attr,omitempty"`
	Ref        string   `xml:"ref,attr,omitempty"`
	DefaultKID string   `xml:"urn:mpeg:cenc:2013 cenc:default_KID,attr,omitempty"`
	// Pssh is PSSH Box with namespace "urn:mpeg:cenc:2013" and prefix "cenc".
	Pssh *PsshType `xml:"urn:mpeg:cenc:2013 cenc:pssh,omitempty"`
	// MSPro is Microsoft PlayReady provisioning data with namespace "urn:microsoft:playready and "prefix "mspr".
	MSPro *MSProType `xml:"urn:microsoft:playready mspr:pro,omitempty"`
	DescriptorType
}

// PsshType is general PSSH box as defined in ISO/IEC 23001-7 (Common Encryption Format).
type PsshType struct {
	Value string `xml:",chardata"`
}

// MSProType is Microsoft PlayReady provisioning data.
type MSProType struct {
	Value string `xml:",chardata"`
}

// ResyncType is Resynchronization Point.
type ResyncType struct {
	XMLName xml.Name `xml:"Resync"`
	Type    uint32   `xml:"type,attr"` // default = 0
	DT      *uint32  `xml:"dT,attr"`
	DImax   *float32 `xml:"dImax,attr"`
	DImin   float32  `xml:"dImin,attr"`  // default = 0
	Marker  bool     `xml:"marker,attr"` // default = false
}

// PR is PR element defined in Table 47.
type PR struct {
	XMLName        xml.Name `xml:"PR"`
	PopularityRate uint32   `xml:"popularityRate,attr"`
	Start          *uint64  `xml:"start,attr"`
	R              int      `xml:"r,attr,omitempty"` // default = 0
}

// ContentPopularityRateType is Content Popularity Rate.
type ContentPopularityRateType struct {
	XMLName           xml.Name `xml:"ContentPopularityRate"`
	Source            string   `xml:"source,attr"`
	Sourcedescription string   `xml:"source_description,attr,omitempty"`
	PR                []*PR    `xml:"PR"`
}

// LabelType is Label and Group Label.
type LabelType struct {
	XMLName xml.Name `xml:"Label"`
	Id      uint32   `xml:"id,attr,omitempty"` // default = 0
	Lang    string   `xml:"lang,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

// ProducerReferenceTimeType is Producer Reference time.
type ProducerReferenceTimeType struct {
	XMLName           xml.Name                      `xml:"ProducerReferenceTime"`
	Id                uint32                        `xml:"id,attr"`
	Inband            bool                          `xml:"inband,attr,omitempty"` // default = false
	Type              ProducerReferenceTimeTypeType `xml:"type,attr,omitempty"`   // default = encoder
	ApplicationScheme string                        `xml:"applicationScheme,attr,omitempty"`
	WallClockTime     string                        `xml:"wallClockTime,attr"`
	PresentationTime  uint64                        `xml:"presentationTime,attr"`
	UTCTiming         *DescriptorType               `xml:"UTCTiming"`
}

// PreselectionType is Preselection.
type PreselectionType struct {
	XMLName                xml.Name              `xml:"Preselection"`
	Id                     string                `xml:"id,attr,omitempty"` // default = "1"
	PreselectionComponents *StringVectorType     `xml:"preselectionComponents,attr"`
	Lang                   string                `xml:"lang,attr,omitempty"`
	Order                  PreselectionOrderType `xml:"order,attr,omitempty"`
	Accessibilities        []*DescriptorType     `xml:"Accessibility"`
	Roles                  []*DescriptorType     `xml:"Role"`
	Ratings                []*DescriptorType     `xml:"Rating"`
	Viewpoints             []*DescriptorType     `xml:"Viewpoint"`
	RepresentationBaseType
}

// AudioSamplingRateType is UIntVectorType with 1 or 2 components.
type AudioSamplingRateType *UIntVectorType

// SubsetType is Subset.
type SubsetType struct {
	XMLName  xml.Name        `xml:"Subset"`
	Contains *UIntVectorType `xml:"contains,attr"`
	Id       string          `xml:"id,attr,omitempty"`
}

// SwitchingType is Switching
type SwitchingType struct {
	XMLName  xml.Name          `xml:"Switching"`
	Interval uint32            `xml:"interval"`
	Type     SwitchingTypeType `xml:"type,attr,omitempty"` // default = "media"
}

// SwitchingTypeType is enumeration "media", "bitstream".
type SwitchingTypeType string

// RandomAccessType is Random Access
type RandomAccessType struct {
	XMLName       xml.Name             `xml:"RandomAccess"`
	Interval      uint32               `xml:"interval,attr"`
	Type          RandomAccessTypeType `xml:"type,attr,omitempty"` // default = "closed"
	MinBufferTime *Duration            `xml:"minBufferTime,attr"`
	Bandwidth     uint32               `xml:"bandwidth,attr,omitempty"`
}

// RandomAccessTypeType is enumeration "closed", "open", "gradual".
type RandomAccessTypeType string

// SegmentBaseType is Segment information base.
type SegmentBaseType struct {
	Timescale                *uint32              `xml:"timescale,attr"` // default = 1
	EptDelta                 *int                 `xml:"eptDelta,attr"`
	PdDelta                  *int                 `xml:"pdDelta,attr"`
	PresentationTimeOffset   *uint64              `xml:"presentationTimeOffset,attr,omitempty"`
	PresentationDuration     *uint64              `xml:"presentationDuration,attr,omitempty"`
	TimeShiftBufferDepth     Duration             `xml:"timeShiftBufferDepth,attr,omitempty"`
	IndexRange               string               `xml:"indexRange,attr,omitempty"`
	IndexRangeExact          bool                 `xml:"indexRangeExact,attr,omitempty"`
	AvailabilityTimeOffset   FloatInf64           `xml:"availabilityTimeOffset,attr,omitempty"`
	AvailabilityTimeComplete *bool                `xml:"availabilityTimeComplete,attr"`
	Initialization           *URLType             `xml:"Initialization"`
	RepresentationIndex      *URLType             `xml:"RepresentationIndex"`
	FailoverContent          *FailoverContentType `xml:"FailoverContent"`
}

// GetTimescale returns timescale including default value 1 if not set.
func (s SegmentBaseType) GetTimescale() uint32 {
	if s.Timescale == nil {
		return 1
	}
	return *s.Timescale
}

// SetTimescale sets timescale and uses default value 1.
func (s *SegmentBaseType) SetTimescale(timescale uint32) {
	if timescale == 1 {
		s.Timescale = nil
	}
	s.Timescale = &timescale
}

// MultipleSegmentBaseType is Multiple Segment information base.
type MultipleSegmentBaseType struct {
	Duration           *uint32              `xml:"duration,attr"`
	StartNumber        *uint32              `xml:"startNumber,attr"`
	EndNumber          *uint32              `xml:"endNumber,attr"`
	SegmentTimeline    *SegmentTimelineType `xml:"SegmentTimeline"`
	BitstreamSwitching *URLType             `xml:"BitstreamSwitching"`
	SegmentBaseType
}

// URLType is Segment Info item URL/range.
type URLType struct {
	SourceURL AnyURI `xml:"sourceURL,attr,omitempty"`
	Range     string `xml:"range,attr,omitempty"`
}

// FCS is Failover Content Section.
type FCS struct {
	T uint64 `xml:"t,attr"`
	D uint64 `xml:"d,attr,omitempty"`
}

// FailoverContentType is Failover Content.
type FailoverContentType struct {
	Valid *bool  `xml:"valid,attr"` // default = true
	FCS   []*FCS `xml:"FCS"`
}

// SegmentListType is Segment List.
type SegmentListType struct {
	XlinkHref    string            `xml:"http://www.w3.org/1999/xlink xlink:href,attr,omitempty"`
	XlinkActuate string            `xml:"http://www.w3.org/1999/xlink xlink:actuate,attr,omitempty"` // default = "onRequest"
	XlinkType    string            `xml:"http://www.w3.org/1999/xlink xlink:type,attr,omitempty"`    // fixed = "simple"
	XlinkShow    string            `xml:"xlink:show,attr,omitempty"`                                 // fixed "embed"
	SegmentURL   []*SegmentURLType `xml:"SegmentURL"`
	MultipleSegmentBaseType
}

// SegmentURLType is Segment URL.
type SegmentURLType struct {
	Media      AnyURI                 `xml:"media,attr,omitempty"`
	MediaRange SingleRFC7233RangeType `xml:"mediaRange,attr,omitempty"`
	Index      AnyURI                 `xml:"index,attr,omitempty"`
	IndexRange SingleRFC7233RangeType `xml:"indexRange,attr,omitempty"`
}

// SegmentTemplateType is Segment Template
type SegmentTemplateType struct {
	Media              string `xml:"media,attr,omitempty"`
	Index              string `xml:"index,attr,omitempty"`
	Initialization     string `xml:"initialization,attr,omitempty"`
	BitstreamSwitching string `xml:"bitstreamSwitching,attr,omitempty"`
	MultipleSegmentBaseType
}

// S is the S element of SegmentTimeline. All time units in media timescale.
// Defined in ISO/IEC 23009-1 Section 5.3.9.6
type S struct {
	// T is presentation time of first Segment in sequence relative to presentationTimeOffset.
	T *uint64 `xml:"t,attr"`
	// N is is first Segment number in Segment sequence relative startNumber
	N *uint64 `xml:"n,attr"`
	// D is the Segment duration.
	D uint64 `xml:"d,attr"`
	// R is repeat count (how many times to repeat. -1 is unlimited)
	R int `xml:"r,attr,omitempty"` // default = 0
	// K is the number of Segments that are included in a Segment Sequence.
	K *uint64 `xml:"k,attr"` // default = 1
}

// SegmentTimelineType is Segment Timeline.
type SegmentTimelineType struct {
	S []*S `xml:"S"`
}

// BaseURLType is Base URL.
type BaseURLType struct {
	ServiceLocation          string     `xml:"serviceLocation,attr,omitempty"`
	ByteRange                string     `xml:"byteRange,attr,omitempty"`
	AvailabilityTimeOffset   FloatInf64 `xml:"availabilityTimeOffset,attr,omitempty"`
	AvailabilityTimeComplete *bool      `xml:"availabilityTimeComplete,attr"`
	TimeShiftBufferDepth     *Duration  `xml:"timeShiftBufferDepth,attr"`
	RangeAccess              bool       `xml:"rangeAccess,attr,omitempty"` // default = false
	Value                    AnyURI     `xml:",chardata"`
}

// NewBaseURL returns a new BaseURLType with Value set.
func NewBaseURL(value string) *BaseURLType {
	return &BaseURLType{
		Value: AnyURI(value),
	}
}

// ProgramInformationType is Program Information.
type ProgramInformationType struct {
	Lang               string `xml:"lang,attr,omitempty"`
	MoreInformationURL AnyURI `xml:"moreInformationURL,attr,omitempty"`
	Title              string `xml:"Title,omitempty"`
	Source             string `xml:"Source,omitempty"`
	Copyright          string `xml:"Copyright,omitempty"`
}

// DescriptorType is Descriptor.
type DescriptorType struct {
	SchemeIdUri AnyURI `xml:"schemeIdUri,attr,omitempty"`
	Value       string `xml:"value,attr,omitempty"`
	Id          string `xml:"id,attr,omitempty"`
}

// NewDescriptor returns a new DescriptorType.
func NewDescriptor(schemeIdURI, value, id string) *DescriptorType {
	return &DescriptorType{
		SchemeIdUri: AnyURI(schemeIdURI),
		Value:       value,
		Id:          id,
	}
}

func NewRole(value string) *DescriptorType {
	return NewDescriptor("urn:mpeg:dash:role:2011", value, "")
}

// MetricsType is Metrics.
type MetricsType struct {
	Metrics    string            `xml:"metrics,attr"`
	Ranges     []*RangeType      `xml:"Range"`
	Reportings []*DescriptorType `xml:"Reporting"`
}

// RangeType is Metrics Range
type RangeType struct {
	Starttime *Duration `xml:"starttime,attr"`
	Duration  *Duration `xml:"duration,attr"`
}

// LeapSecondInformationType is Leap Second Information
type LeapSecondInformationType struct {
	AvailabilityStartLeapOffset     int      `xml:"availabilityStartLeapOffset,attr"`
	NextAvailabilityStartLeapOffset int      `xml:"nextAvailabilityStartLeapOffset,attr,omitempty"`
	NextLeapChangeTime              DateTime `xml:"nextLeapChangeTime,attr,omitempty"`
}

// ListOfProfilesType is comma-separated list of profiles.
type ListOfProfilesType string

// AddProfile adds a profile to the comma-separated list of profiles.
func (l ListOfProfilesType) AddProfile(profile string) ListOfProfilesType {
	if l == "" {
		return ListOfProfilesType(profile)
	}
	return ListOfProfilesType(string(l) + "," + profile)
}

// StringVectorType is Whitespace-separated list of strings.
type StringVectorType string

// ListOf4CCType is Whitespace separated list of 4CC.
type ListOf4CCType string

// UIntVectorType is Whitespace-separated list of unsigned integers.
type UIntVectorType string

// ContentEncodingType is an enum with single value "base64".
type ContentEncodingType string

// RatioType is Ratio Type for sar and par ([0-9]*:[0-9]*)
type RatioType string

// FrameRateType is Type for Frame Rate ([0-9]+(/[1-9][0-9]*)?).
type FrameRateType string

// RFC6838ContentTypeType is Type for RFC6838 Content Type.
type RFC6838ContentTypeType string

// StringNoWhitespaceType is String without white spaces.
type StringNoWhitespaceType string

// VideoScanType is enumeration "progressive", "interlaced", "unknown".
type VideoScanType string

// ProducerReferenceTimeTypeType is enumeration "encoder", "captured", "application".
type ProducerReferenceTimeTypeType string

// PreselectionOrderType is enumeration "undefined", "time-ordered", "fully-ordered".
type PreselectionOrderType string

// SingleRFC7233RangeType is range defined in RFC7233 ([0-9]*)(\-([0-9]*))?).
type SingleRFC7233RangeType string

// AnyURI is xsd:anyURI http://www.datypic.com/sc/xsd/t-xsd_anyURI.html.
type AnyURI string

// DateTime is xs:dateTime https://www.w3.org/TR/xmlschema-2/#dateTime (almost ISO 8601).
type DateTime string

// NewMPD returns a new empty MPD with the right type.
func NewMPD(mpdType string) *MPD {
	return &MPD{
		Type:  &mpdType,
		XMLNs: DASH_NAMESPACE,
	}
}

// NewPeriod returns a new empty Period.
func NewPeriod() *Period {
	return &Period{}
}

// NewAdaptationSet returns a new empty AdaptationSet.
func NewAdaptationSet() *AdaptationSetType {
	return &AdaptationSetType{}
}

func NewAdaptationSetWithParams(contentType, mimeType string, segmentAlignment bool, startsWithSAP uint32) *AdaptationSetType {
	return &AdaptationSetType{
		RepresentationBaseType: RepresentationBaseType{
			MimeType:     mimeType,
			StartWithSAP: startsWithSAP,
		},
		SegmentAlignment: segmentAlignment,
		ContentType:      RFC6838ContentTypeType(contentType),
	}
}

// NewRepresentation returns a new empty Representation.
func NewRepresentation() *RepresentationType {
	return &RepresentationType{}
}

// NewRepresentationWithID returns a new empty Representation with the given ID and parameters.
func NewRepresentationWithID(id, codec, mimeType string, bandwidth int) *RepresentationType {
	r := &RepresentationType{
		RepresentationBaseType: RepresentationBaseType{
			Codecs:   codec,
			MimeType: mimeType,
		},
		Id:        id,
		Bandwidth: uint32(bandwidth),
	}
	return r
}

// NewAudioRepresentation returns a new audio representation with parameters.
func NewAudioRepresentation(id, codec, mimeType string, bandwidth, audioSamplingRate int) *RepresentationType {
	r := NewRepresentationWithID(id, codec, mimeType, bandwidth)
	r.AudioSamplingRate = Ptr(UIntVectorType(fmt.Sprintf("%d", audioSamplingRate)))
	return r
}

// NewVideoRepresentation returns a new  video representation with parameters.
func NewVideoRepresentation(id, codec, mimeType, frameRate string, bandwidth, width, height int) *RepresentationType {
	r := NewRepresentationWithID(id, codec, mimeType, bandwidth)
	r.FrameRate = FrameRateType(frameRate)
	r.Width = uint32(width)
	r.Height = uint32(height)
	return r
}

// NewSubRepresentation returns a new empty SubRepresentation.
func NewSubRepresentation() *SubRepresentationType {
	return &SubRepresentationType{}
}

// NewSegmentTemplate returns a new empty SegmentTemplate.
func NewSegmentTemplate() *SegmentTemplateType {
	return &SegmentTemplateType{}
}

// NewSegmentList returns a new empty SegmentList.
func NewSegmentList() *SegmentListType {
	return &SegmentListType{}
}

// NewSegmentTemplate returns a new empty SegmentTimeline.
func NewSegmentTimeline() *SegmentTimelineType {
	return &SegmentTimelineType{}
}

// NewInitializationSet returns a new empty InitializationSet.
func NewInitializationSet() *InitializationSetType {
	return &InitializationSetType{}
}

// NewPreselection returns a new empty Preselection.
func NewPreselection() *PreselectionType {
	return &PreselectionType{}
}

// NewContentProtection returns a new empty ContentProtection.
func NewContentProtection() *ContentProtectionType {
	return &ContentProtectionType{}
}

// NewProducerReferenceTime returns a new empty ProducerReferenceTime.
func NewProducerReferenceTime() *ProducerReferenceTimeType {
	return &ProducerReferenceTimeType{}
}

// NewUIntVWithID returns a new empty UIntVWithID.
func NewUIntVWithID() *UIntVWithIDType {
	return &UIntVWithIDType{}
}

// FloatInf64 is a float64 which renders as "INF" for +Inf
type FloatInf64 float64

func (f *FloatInf64) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if f == nil {
		return xml.Attr{}, nil
	}
	fl := float64(*f)
	if fl == math.Inf(+1) {
		return xml.Attr{Name: name, Value: "INF"}, nil
	}
	val := strconv.FormatFloat(fl, 'f', -1, 64)
	return xml.Attr{Name: name, Value: val}, nil
}

func (f *FloatInf64) UnmarshalXMLAttr(attr xml.Attr) error {
	if attr.Value == "INF" {
		*f = FloatInf64(math.Inf(+1))
		return nil
	}
	fl, err := strconv.ParseFloat(attr.Value, 64)
	if err != nil {
		return err
	}
	*f = FloatInf64(fl)
	return nil
}
