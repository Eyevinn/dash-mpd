package mpd

// MPD is MPEG-DASH Media Presentation Description (MPD) as defined in ISO/IEC 23009-1 5'th edition.
//
// The tree of structs is generated from the corresponding XML Schema at https://github.com/MPEGGroup/DASHSchema
// but fine-tuned manually to handle default cases, listing enumerals, name space for xlink etc.
type MPD struct {
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
	Period                     []*PeriodType              `xml:"Period"`
	Metrics                    []*MetricsType             `xml:"Metrics"`
	EssentialProperty          []*DescriptorType          `xml:"EssentialProperty"`
	SupplementalProperty       []*DescriptorType          `xml:"SupplementalProperty"`
	UTCTiming                  []*DescriptorType          `xml:"UTCTiming"`
	LeapSecondInformation      *LeapSecondInformationType `xml:"LeapSecondInformation"`
}

// PatchLocationType is Patch Location Type.
type PatchLocationType struct {
	Ttl   float64 `xml:"ttl,attr,omitempty"`
	Value AnyURI  `xml:",chardata"`
}

// PeriodType is Period.
type PeriodType struct {
	XlinkHref            string                    `xml:"http://www.w3.org/1999/xlink xlink:href,attr,omitempty"`
	XlinkActuate         string                    `xml:"http://www.w3.org/1999/xlink xlink:actuate,attr,omitempty"` // default = "onRequest"
	XlinkType            string                    `xml:"http://www.w3.org/1999/xlink xlink:type,attr,omitempty"`    // fixed = "simple"
	XlinkShow            string                    `xml:"http://www.w3.org/1999/xlink xlink:show,attr,omitempty"`    // fixed = "embed"
	Id                   string                    `xml:"id,attr,omitempty"`
	Start                *Duration                 `xml:"start,attr"` // Mandatory for dynamic manifests. default = 0
	Duration             *Duration                 `xml:"duration,attr"`
	BitstreamSwitching   *bool                     `xml:"bitstreamSwitching,attr"`
	BaseURL              []*BaseURLType            `xml:"BaseURL"`
	SegmentBase          *SegmentBaseType          `xml:"SegmentBase"`
	SegmentList          *SegmentListType          `xml:"SegmentList"`
	SegmentTemplate      *SegmentTemplateType      `xml:"SegmentTemplate"`
	AssetIdentifier      *DescriptorType           `xml:"AssetIdentifier"`
	EventStream          []*EventStreamType        `xml:"EventStream"`
	ServiceDescription   []*ServiceDescriptionType `xml:"ServiceDescription"`
	ContentProtection    []*ContentProtectionType  `xml:"ContentProtection"`
	AdaptationSet        []*AdaptationSetType      `xml:"AdaptationSet"`
	Subset               []*SubsetType             `xml:"Subset"`
	SupplementalProperty []*DescriptorType         `xml:"SupplementalProperty"`
	EmptyAdaptationSet   []*AdaptationSetType      `xml:"EmptyAdaptationSet"`
	GroupLabel           []*LabelType              `xml:"GroupLabel"`
	Preselection         []*PreselectionType       `xml:"Preselection"`
}

// EventStreamType is Event Stream.
type EventStreamType struct {
	XlinkHref              string       `xml:"http://www.w3.org/1999/xlink xlink:href,attr,omitempty"`
	XlinkActuate           string       `xml:"http://www.w3.org/1999/xlink xlink:actuate,attr,omitempty"` // default = "onRequest"
	XlinkType              string       `xml:"http://www.w3.org/1999/xlink xlink:type,attr,omitempty"`    // fixed = "simple"
	XlinkShow              string       `xml:"http://www.w3.org/1999/xlink xlink:show,attr,omitempty"`    // fixed = "embed"
	SchemeIdUri            AnyURI       `xml:"schemeIdUri,attr"`
	Value                  string       `xml:"value,attr,omitempty"`
	Timescale              *uint32      `xml:"timescale,attr"`                        // default = 1
	PresentationTimeOffset uint64       `xml:"presentationTimeOffset,attr,omitempty"` // default is 0
	Event                  []*EventType `xml:"Event"`
}

// EventType is Event.
type EventType struct {
	PresentationTime uint64              `xml:"presentationTime,attr,omitempty"` // default is 0
	Duration         uint64              `xml:"duration,attr,omitempty"`
	Id               uint32              `xml:"id,attr"`
	ContentEncoding  ContentEncodingType `xml:"contentEncoding,attr,omitempty"`
	MessageData      string              `xml:"messageData,attr,omitempty"`
}

// InitializationSetType is Initialization Set.
type InitializationSetType struct {
	XlinkHref      string                 `xml:"http://www.w3.org/1999/xlink xlink:href,attr,omitempty"`
	XlinkActuate   string                 `xml:"http://www.w3.org/1999/xlink xlink:actuate,attr,omitempty"` // default = "onRequest"
	XlinkType      string                 `xml:"http://www.w3.org/1999/xlink xlink:type,attr,omitempty"`    // fixed = "simple"
	Id             uint32                 `xml:"id,attr"`
	InAllPeriods   *bool                  `xml:"inAllPeriods,attr"` // default is true
	ContentType    RFC6838ContentTypeType `xml:"contentType,attr,omitempty"`
	Par            RatioType              `xml:"par,attr,omitempty"`
	MaxWidth       uint32                 `xml:"maxWidth,attr,omitempty"`
	MaxHeight      uint32                 `xml:"maxHeight,attr,omitempty"`
	MaxFrameRate   string                 `xml:"maxFrameRate,attr,omitempty"`
	Initialization AnyURI                 `xml:"initialization,attr,omitempty"`
	Accessibility  []*DescriptorType      `xml:"Accessibility"`
	Role           []*DescriptorType      `xml:"Role"`
	Rating         []*DescriptorType      `xml:"Rating"`
	Viewpoint      []*DescriptorType      `xml:"Viewpoint"`
	*RepresentationBaseType
}

// ServiceDescriptionType is Service Description.
type ServiceDescriptionType struct {
	Id                 uint32                    `xml:"id,attr"`
	Scope              []*DescriptorType         `xml:"Scope"`
	Latency            []*LatencyType            `xml:"Latency"`
	PlaybackRate       []*PlaybackRateType       `xml:"PlaybackRate"`
	OperatingQuality   []*OperatingQualityType   `xml:"OperatingQuality"`
	OperatingBandwidth []*OperatingBandwidthType `xml:"OperatingBandwidth"`
}

// LatencyType is Service Description Latency (Annex K.4.2.2).
type LatencyType struct {
	ReferenceId    uint32                 `xml:"referenceId,attr"`
	Target         *uint32                `xml:"target,attr"`
	Max            *uint32                `xml:"max,attr"`
	Min            *uint32                `xml:"min,attr"`
	QualityLatency []*UIntPairsWithIDType `xml:"QualityLatency"`
}

// PlaybackRateType is Service Description Playback Rate.
type PlaybackRateType struct {
	Max float64 `xml:"max,attr,omitempty"`
	Min float64 `xml:"min,attr,omitempty"`
}

// OperatingQualityType is Service Description Operating Quality.
type OperatingQualityType struct {
	MediaType     string `xml:"mediaType,attr,omitempty"` // default is "any"
	Min           uint32 `xml:"min,attr,omitempty"`
	Max           uint32 `xml:"max,attr,omitempty"`
	Target        uint32 `xml:"target,attr,omitempty"`
	Type          AnyURI `xml:"type,attr,omitempty"`
	MaxDifference uint32 `xml:"maxDifference,attr,omitempty"`
}

// OperatingBandwidthType is Service Description Operating Bandwidth.
type OperatingBandwidthType struct {
	MediaType string `xml:"mediaType,attr,omitempty"` // default is "all"
	Min       uint32 `xml:"min,attr,omitempty"`
	Max       uint32 `xml:"max,attr,omitempty"`
	Target    uint32 `xml:"target,attr,omitempty"`
}

// UIntPairsWithIDType is UInt Pairs With ID.
type UIntPairsWithIDType struct {
	Type AnyURI `xml:"type,attr,omitempty"`
	*UIntVectorType
}

// UIntVWithIDType is UInt Vector With ID
type UIntVWithIDType struct {
	Id          uint32                 `xml:"id,attr"`
	Profiles    ListOfProfilesType     `xml:"profiles,attr,omitempty"`
	ContentType RFC6838ContentTypeType `xml:"contentType,attr,omitempty"`
	*UIntVectorType
}

// AdaptationSetType is Adaptation Set.
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
	Accessibility           []*DescriptorType       `xml:"Accessibility"`
	Role                    []*DescriptorType       `xml:"Role"`
	Rating                  []*DescriptorType       `xml:"Rating"`
	Viewpoint               []*DescriptorType       `xml:"Viewpoint"`
	ContentComponent        []*ContentComponentType `xml:"ContentComponent"`
	BaseURL                 []*BaseURLType          `xml:"BaseURL"`
	SegmentBase             *SegmentBaseType        `xml:"SegmentBase"`
	SegmentList             *SegmentListType        `xml:"SegmentList"`
	SegmentTemplate         *SegmentTemplateType    `xml:"SegmentTemplate"`
	Representation          []*RepresentationType   `xml:"Representation"`
	*RepresentationBaseType
}

// ContentComponentType is Content Component.
type ContentComponentType struct {
	Id            *uint32                `xml:"id,attr"`
	Lang          string                 `xml:"lang,attr,omitempty"`
	ContentType   RFC6838ContentTypeType `xml:"contentType,attr,omitempty"`
	Par           RatioType              `xml:"par,attr,omitempty"`
	Tag           string                 `xml:"tag,attr,omitempty"`
	Accessibility []*DescriptorType      `xml:"Accessibility"`
	Role          []*DescriptorType      `xml:"Role"`
	Rating        []*DescriptorType      `xml:"Rating"`
	Viewpoint     []*DescriptorType      `xml:"Viewpoint"`
}

// RepresentationType is Representation.
type RepresentationType struct {
	Id                     string                   `xml:"id,attr"`
	Bandwidth              uint32                   `xml:"bandwidth,attr"`
	QualityRanking         *uint32                  `xml:"qualityRanking,attr,omitempty"`
	DependencyId           *StringVectorType        `xml:"dependencyId,attr,omitempty"`
	AssociationId          *StringVectorType        `xml:"associationId,attr,omitempty"`
	AssociationType        *ListOf4CCType           `xml:"associationType,attr,omitempty"`
	MediaStreamStructureId *StringVectorType        `xml:"mediaStreamStructureId,attr,omitempty"`
	BaseURL                []*BaseURLType           `xml:"BaseURL"`
	ExtendedBandwidth      []*ExtendedBandwidthType `xml:"ExtendedBandwidth"`
	SubRepresentation      []*SubRepresentationType `xml:"SubRepresentation"`
	SegmentBase            *SegmentBaseType         `xml:"SegmentBase"`
	SegmentList            *SegmentListType         `xml:"SegmentList"`
	SegmentTemplate        *SegmentTemplateType     `xml:"SegmentTemplate"`
	*RepresentationBaseType
}

// ExtendedBandwidthType is Extended Bandwidth Model
type ExtendedBandwidthType struct {
	Vbr       bool             `xml:"vbr,attr,omitempty"` // default is false
	ModelPair []*ModelPairType `xml:"ModelPair"`
}

// ModelPairType is Model Pair
type ModelPairType struct {
	BufferTime *Duration `xml:"bufferTime,attr"`
	Bandwidth  uint32    `xml:"bandwidth,attr"`
}

// SubRepresentationType is SubRepresentation
type SubRepresentationType struct {
	Level            *uint32           `xml:"level,attr,omitempty"`
	DependencyLevel  *UIntVectorType   `xml:"dependencyLevel,attr,omitempty"`
	Bandwidth        uint32            `xml:"bandwidth,attr,omitempty"`
	ContentComponent *StringVectorType `xml:"contentComponent,attr,omitempty"`
	*RepresentationBaseType
}

// RepresentationBaseType is Representation base (common attributes and elements).
type RepresentationBaseType struct {
	Profiles                  ListOfProfilesType           `xml:"profiles,attr,omitempty"`
	Width                     uint32                       `xml:"width,attr,omitempty"`
	Height                    uint32                       `xml:"height,attr,omitempty"`
	Sar                       RatioType                    `xml:"sar,attr,omitempty"`
	FrameRate                 FrameRateType                `xml:"frameRate,attr,omitempty"`
	AudioSamplingRate         *UIntVectorType              `xml:"audioSamplingRate,attr,omitempty"`
	MimeType                  string                       `xml:"mimeType,attr,omitempty"`
	SegmentProfiles           *ListOf4CCType               `xml:"segmentProfiles,attr,omitempty"`
	Codecs                    string                       `xml:"codecs,attr,omitempty"`
	ContainerProfiles         *ListOf4CCType               `xml:"containerProfiles,attr,omitempty"`
	MaximumSAPPeriod          float64                      `xml:"maximumSAPPeriod,attr,omitempty"`
	StartWithSAP              uint32                       `xml:"startWithSAP,attr,omitempty"`
	MaxPlayoutRate            float64                      `xml:"maxPlayoutRate,attr,omitempty"`
	CodingDependency          *bool                        `xml:"codingDependency,attr,omitempty"`
	ScanType                  VideoScanType                `xml:"scanType,attr,omitempty"`
	SelectionPriority         *uint32                      `xml:"selectionPriority,attr"` // default = 1
	Tag                       string                       `xml:"tag,attr,omitempty"`
	FramePacking              []*DescriptorType            `xml:"FramePacking"`
	AudioChannelConfiguration []*DescriptorType            `xml:"AudioChannelConfiguration"`
	ContentProtection         []*ContentProtectionType     `xml:"ContentProtection"`
	OutputProtection          *DescriptorType              `xml:"OutputProtection"`
	EssentialProperty         []*DescriptorType            `xml:"EssentialProperty"`
	SupplementalProperty      []*DescriptorType            `xml:"SupplementalProperty"`
	InbandEventStream         []*EventStreamType           `xml:"InbandEventStream"`
	Switching                 []*SwitchingType             `xml:"Switching"`
	RandomAccess              []*RandomAccessType          `xml:"RandomAccess"`
	GroupLabel                []*LabelType                 `xml:"GroupLabel"`
	Label                     []*LabelType                 `xml:"Label"`
	ProducerReferenceTime     []*ProducerReferenceTimeType `xml:"ProducerReferenceTime"`
	ContentPopularityRate     []*ContentPopularityRateType `xml:"ContentPopularityRate"`
	Resync                    []*ResyncType                `xml:"Resync"`
}

// ContentProtectionType is Content Protection.
type ContentProtectionType struct {
	Robustness string `xml:"robustness,attr,omitempty"`
	RefId      string `xml:"refId,attr,omitempty"`
	Ref        string `xml:"ref,attr,omitempty"`
	DefaultKID string `xml:"urn:mpeg:cenc:2013 cenc:default_KID,attr,omitempty"`
	// Pssh is PSSH Box with namespace "urn:mpeg:cenc:2013" and prefix "cenc".
	Pssh *PsshType `xml:"urn:mpeg:cenc:2013 cenc:pssh,omitempty"`
	// MSPro is Microsoft PlayReady provisioning data with namespace "urn:microsoft:playready and "prefix "mspr".
	MSPro *MSProType `xml:"urn:microsoft:playready mspr:pro,omitempty"`
	*DescriptorType
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
	Type   uint32   `xml:"type,attr"` // default = 0
	DT     *uint32  `xml:"dT,attr"`
	DImax  *float32 `xml:"dImax,attr"`
	DImin  float32  `xml:"dImin,attr"`  // default = 0
	Marker bool     `xml:"marker,attr"` // default = false
}

// PR is PR element defined in Table 47.
type PR struct {
	PopularityRate uint32  `xml:"popularityRate,attr"`
	Start          *uint64 `xml:"start,attr"`
	R              int     `xml:"r,attr,omitempty"` // default = 0
}

// ContentPopularityRateType is Content Popularity Rate.
type ContentPopularityRateType struct {
	Source            string `xml:"source,attr"`
	Sourcedescription string `xml:"source_description,attr,omitempty"`
	PR                []*PR  `xml:"PR"`
}

// LabelType is Label and Group Label.
type LabelType struct {
	Id    uint32 `xml:"id,attr,omitempty"` // default = 0
	Lang  string `xml:"lang,attr,omitempty"`
	Value string `xml:",chardata"`
}

// ProducerReferenceTimeType is Producer Reference time.
type ProducerReferenceTimeType struct {
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
	Id                     string                `xml:"id,attr,omitempty"` // default = "1"
	PreselectionComponents *StringVectorType     `xml:"preselectionComponents,attr"`
	Lang                   string                `xml:"lang,attr,omitempty"`
	Order                  PreselectionOrderType `xml:"order,attr,omitempty"`
	Accessibility          []*DescriptorType     `xml:"Accessibility"`
	Role                   []*DescriptorType     `xml:"Role"`
	Rating                 []*DescriptorType     `xml:"Rating"`
	Viewpoint              []*DescriptorType     `xml:"Viewpoint"`
	*RepresentationBaseType
}

// AudioSamplingRateType is UIntVectorType with 1 or 2 components.
type AudioSamplingRateType *UIntVectorType

// SubsetType is Subset.
type SubsetType struct {
	Contains *UIntVectorType `xml:"contains,attr"`
	Id       string          `xml:"id,attr,omitempty"`
}

// SwitchingType is Switching
type SwitchingType struct {
	Interval uint32            `xml:"interval"`
	Type     SwitchingTypeType `xml:"type,attr,omitempty"` // default = "media"
}

// SwitchingTypeType is enumeration "media", "bitstream".
type SwitchingTypeType string

// RandomAccessType is Random Access
type RandomAccessType struct {
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
	AvailabilityTimeOffset   float64              `xml:"availabilityTimeOffset,attr,omitempty"`
	AvailabilityTimeComplete *bool                `xml:"availabilityTimeComplete,attr"`
	Initialization           *URLType             `xml:"Initialization"`
	RepresentationIndex      *URLType             `xml:"RepresentationIndex"`
	FailoverContent          *FailoverContentType `xml:"FailoverContent"`
}

// MultipleSegmentBaseType is Multiple Segment information base.
type MultipleSegmentBaseType struct {
	Duration           *uint32              `xml:"duration,attr"`
	StartNumber        *uint32              `xml:"startNumber,attr"`
	EndNumber          *uint32              `xml:"endNumber,attr"`
	SegmentTimeline    *SegmentTimelineType `xml:"SegmentTimeline"`
	BitstreamSwitching *URLType             `xml:"BitstreamSwitching"`
	*SegmentBaseType
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
	*MultipleSegmentBaseType
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
	*MultipleSegmentBaseType
}

// S is the S element of SegmentTimeline. All time units in media timescale.
type S struct {
	// T is start time of first Segment in the the series relative to presentation time offset.
	T *uint64 `xml:"t,attr"`
	// N is the Segment number of the first Segment in the series.
	N *uint64 `xml:"n,attr"`
	// D is the Segment duration or the duration of a Segment sequence.
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
	ServiceLocation          string    `xml:"serviceLocation,attr,omitempty"`
	ByteRange                string    `xml:"byteRange,attr,omitempty"`
	AvailabilityTimeOffset   *float64  `xml:"availabilityTimeOffset,attr"`
	AvailabilityTimeComplete *bool     `xml:"availabilityTimeComplete,attr"`
	TimeShiftBufferDepth     *Duration `xml:"timeShiftBufferDepth,attr"`
	RangeAccess              bool      `xml:"rangeAccess,attr,omitempty"` // default = false
	Value                    AnyURI    `xml:",chardata"`
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

// MetricsType is Metrics.
type MetricsType struct {
	Metrics   string            `xml:"metrics,attr"`
	Range     []*RangeType      `xml:"Range"`
	Reporting []*DescriptorType `xml:"Reporting"`
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

// ProducerReferenceTimeTypeType is enumaration "encoder", "captured", "application".
type ProducerReferenceTimeTypeType string

// PreselectionOrderType is enumeration "undefined", "time-ordered", "fully-ordered".
type PreselectionOrderType string

// SingleRFC7233RangeType is range defined in RFC7233 ([0-9]*)(\-([0-9]*))?).
type SingleRFC7233RangeType string

// AnyURI is xsd:anyURI http://www.datypic.com/sc/xsd/t-xsd_anyURI.html.
type AnyURI string

// DateTime is xs:dateTime https://www.w3.org/TR/xmlschema-2/#dateTime (almost ISO 8601).
type DateTime string
