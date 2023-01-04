package mpd

// MPD is root
type MPD struct {
	XMLNs string `xml:"xmlns,attr,omitempty"`
	//XMLNsXSI                   string                     `xml:"xmlns:xsi,attr,omitempty"`
	SchemaLocation             string                     `xml:"http://www.w3.org/2001/XMLSchema-instance xsi:schemaLocation,attr"`
	Id                         string                     `xml:"id,attr,omitempty"`
	Profiles                   string                     `xml:"profiles,attr"`
	Type                       string                     `xml:"type,attr,omitempty"`
	AvailabilityStartTime      string                     `xml:"availabilityStartTime,attr,omitempty"`
	AvailabilityEndTime        string                     `xml:"availabilityEndTime,attr,omitempty"`
	PublishTime                string                     `xml:"publishTime,attr,omitempty"`
	MediaPresentationDuration  string                     `xml:"mediaPresentationDuration,attr,omitempty"`
	MinimumUpdatePeriod        string                     `xml:"minimumUpdatePeriod,attr,omitempty"`
	MinBufferTime              string                     `xml:"minBufferTime,attr"`
	TimeShiftBufferDepth       string                     `xml:"timeShiftBufferDepth,attr,omitempty"`
	SuggestedPresentationDelay string                     `xml:"suggestedPresentationDelay,attr,omitempty"`
	MaxSegmentDuration         string                     `xml:"maxSegmentDuration,attr,omitempty"`
	MaxSubsegmentDuration      string                     `xml:"maxSubsegmentDuration,attr,omitempty"`
	ProgramInformation         []*ProgramInformationType  `xml:"ProgramInformation"`
	BaseURL                    []*BaseURLType             `xml:"BaseURL"`
	Location                   []string                   `xml:"Location"`
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

// PatchLocationType is Patch Location Type
type PatchLocationType struct {
	Ttl   float64 `xml:"ttl,attr,omitempty"`
	Value string  `xml:",chardata"`
}

// PresentationType is Presentation Type enumeration
type PresentationType string

// PeriodType is Period
type PeriodType struct {
	XlinkHref            string                    `xml:"xlink:href,attr,omitempty"`
	XlinkActuate         string                    `xml:"xlink:actuate,attr,omitempty"`
	XlinkType            string                    `xml:"xlink:type,attr,omitempty"`
	XlinkShow            string                    `xml:"xlink:show,attr,omitempty"`
	Id                   string                    `xml:"id,attr,omitempty"`
	Start                string                    `xml:"start,attr,omitempty"`
	Duration             string                    `xml:"duration,attr,omitempty"`
	BitstreamSwitching   bool                      `xml:"bitstreamSwitching,attr,omitempty"`
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

// EventStreamType is Event Stream
type EventStreamType struct {
	XlinkHref              string       `xml:"xlink:href,attr,omitempty"`
	XlinkActuate           string       `xml:"xlink:actuate,attr,omitempty"`
	XlinkType              string       `xml:"xlink:type,attr,omitempty"`
	XlinkShow              string       `xml:"xlink:show,attr,omitempty"`
	SchemeIdUri            string       `xml:"schemeIdUri,attr"`
	Value                  string       `xml:"value,attr,omitempty"`
	Timescale              uint32       `xml:"timescale,attr,omitempty"`
	PresentationTimeOffset uint64       `xml:"presentationTimeOffset,attr,omitempty"`
	Event                  []*EventType `xml:"Event"`
}

// EventType is Event
type EventType struct {
	PresentationTime uint64 `xml:"presentationTime,attr,omitempty"`
	Duration         uint64 `xml:"duration,attr,omitempty"`
	Id               uint32 `xml:"id,attr,omitempty"`
	ContentEncoding  string `xml:"contentEncoding,attr,omitempty"`
	MessageData      string `xml:"messageData,attr,omitempty"`
}

// ContentEncodingType is Event Coding
type ContentEncodingType string

// InitializationSetType is Initialization Set
type InitializationSetType struct {
	XlinkHref      string            `xml:"xlink:href,attr,omitempty"`
	XlinkActuate   string            `xml:"xlink:actuate,attr,omitempty"`
	XlinkType      string            `xml:"xlink:type,attr,omitempty"`
	Id             uint32            `xml:"id,attr"`
	InAllPeriods   bool              `xml:"inAllPeriods,attr,omitempty"`
	ContentType    string            `xml:"contentType,attr,omitempty"`
	Par            string            `xml:"par,attr,omitempty"`
	MaxWidth       uint32            `xml:"maxWidth,attr,omitempty"`
	MaxHeight      uint32            `xml:"maxHeight,attr,omitempty"`
	MaxFrameRate   string            `xml:"maxFrameRate,attr,omitempty"`
	Initialization string            `xml:"initialization,attr,omitempty"`
	Accessibility  []*DescriptorType `xml:"Accessibility"`
	Role           []*DescriptorType `xml:"Role"`
	Rating         []*DescriptorType `xml:"Rating"`
	Viewpoint      []*DescriptorType `xml:"Viewpoint"`
	*RepresentationBaseType
}

// ServiceDescriptionType is Service Description
type ServiceDescriptionType struct {
	Id                 uint32                    `xml:"id,attr,omitempty"`
	Scope              []*DescriptorType         `xml:"Scope"`
	Latency            []*LatencyType            `xml:"Latency"`
	PlaybackRate       []*PlaybackRateType       `xml:"PlaybackRate"`
	OperatingQuality   []*OperatingQualityType   `xml:"OperatingQuality"`
	OperatingBandwidth []*OperatingBandwidthType `xml:"OperatingBandwidth"`
}

// LatencyType is Service Description Latency
type LatencyType struct {
	ReferenceId    uint32                 `xml:"referenceId,attr,omitempty"`
	Target         uint32                 `xml:"target,attr,omitempty"`
	Max            uint32                 `xml:"max,attr,omitempty"`
	Min            uint32                 `xml:"min,attr,omitempty"`
	QualityLatency []*UIntPairsWithIDType `xml:"QualityLatency"`
}

// PlaybackRateType is Service Description Playback Rate
type PlaybackRateType struct {
	Max float64 `xml:"max,attr,omitempty"`
	Min float64 `xml:"min,attr,omitempty"`
}

// OperatingQualityType is Service Description Operating Quality
type OperatingQualityType struct {
	MediaType     string `xml:"mediaType,attr,omitempty"`
	Min           uint32 `xml:"min,attr,omitempty"`
	Max           uint32 `xml:"max,attr,omitempty"`
	Target        uint32 `xml:"target,attr,omitempty"`
	Type          string `xml:"type,attr,omitempty"`
	MaxDifference uint32 `xml:"maxDifference,attr,omitempty"`
}

// OperatingBandwidthType is Service Description Operating Bandwidth
type OperatingBandwidthType struct {
	MediaType string `xml:"mediaType,attr,omitempty"`
	Min       uint32 `xml:"min,attr,omitempty"`
	Max       uint32 `xml:"max,attr,omitempty"`
	Target    uint32 `xml:"target,attr,omitempty"`
}

// UIntPairsWithIDType is UInt Pairs With ID
type UIntPairsWithIDType struct {
	Type string `xml:"type,attr,omitempty"`
	*UIntVectorType
}

// UIntVWithIDType is UInt Vector With ID
type UIntVWithIDType struct {
	Id          uint32 `xml:"id,attr"`
	Profiles    string `xml:"profiles,attr,omitempty"`
	ContentType string `xml:"contentType,attr,omitempty"`
	*UIntVectorType
}

// ListOfProfilesType is List of Profiles
type ListOfProfilesType string

// AdaptationSetType is Adaptation Set
type AdaptationSetType struct {
	XlinkHref               string                  `xml:"xlink:href,attr,omitempty"`
	XlinkActuate            string                  `xml:"xlink:actuate,attr,omitempty"`
	XlinkType               string                  `xml:"xlink:type,attr,omitempty"`
	XlinkShow               string                  `xml:"xlink:show,attr,omitempty"`
	Id                      uint32                  `xml:"id,attr,omitempty"`
	Group                   uint32                  `xml:"group,attr,omitempty"`
	Lang                    string                  `xml:"lang,attr,omitempty"`
	ContentType             string                  `xml:"contentType,attr,omitempty"`
	Par                     string                  `xml:"par,attr,omitempty"`
	MinBandwidth            uint32                  `xml:"minBandwidth,attr,omitempty"`
	MaxBandwidth            uint32                  `xml:"maxBandwidth,attr,omitempty"`
	MinWidth                uint32                  `xml:"minWidth,attr,omitempty"`
	MaxWidth                uint32                  `xml:"maxWidth,attr,omitempty"`
	MinHeight               uint32                  `xml:"minHeight,attr,omitempty"`
	MaxHeight               uint32                  `xml:"maxHeight,attr,omitempty"`
	MinFrameRate            string                  `xml:"minFrameRate,attr,omitempty"`
	MaxFrameRate            string                  `xml:"maxFrameRate,attr,omitempty"`
	SegmentAlignment        bool                    `xml:"segmentAlignment,attr,omitempty"`
	SubsegmentAlignment     bool                    `xml:"subsegmentAlignment,attr,omitempty"`
	SubsegmentStartsWithSAP uint32                  `xml:"subsegmentStartsWithSAP,attr,omitempty"`
	BitstreamSwitching      bool                    `xml:"bitstreamSwitching,attr,omitempty"`
	InitializationSetRef    *UIntVectorType         `xml:"initializationSetRef,attr,omitempty"`
	InitializationPrincipal string                  `xml:"initializationPrincipal,attr,omitempty"`
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

// RatioType is Ratio Type for sar and par
type RatioType string

// FrameRateType is Type for Frame Rate
type FrameRateType string

// RFC6838ContentTypeType is Type for RFC6838 Content Type
type RFC6838ContentTypeType string

// ContentComponentType is Content Component
type ContentComponentType struct {
	Id            uint32            `xml:"id,attr,omitempty"`
	Lang          string            `xml:"lang,attr,omitempty"`
	ContentType   string            `xml:"contentType,attr,omitempty"`
	Par           string            `xml:"par,attr,omitempty"`
	Tag           string            `xml:"tag,attr,omitempty"`
	Accessibility []*DescriptorType `xml:"Accessibility"`
	Role          []*DescriptorType `xml:"Role"`
	Rating        []*DescriptorType `xml:"Rating"`
	Viewpoint     []*DescriptorType `xml:"Viewpoint"`
}

// RepresentationType is Representation
type RepresentationType struct {
	Id                     string                   `xml:"id,attr"`
	Bandwidth              uint32                   `xml:"bandwidth,attr"`
	QualityRanking         uint32                   `xml:"qualityRanking,attr,omitempty"`
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
	Vbr       bool             `xml:"vbr,attr,omitempty"`
	ModelPair []*ModelPairType `xml:"ModelPair"`
}

// ModelPairType is Model Pair
type ModelPairType struct {
	BufferTime string `xml:"bufferTime,attr"`
	Bandwidth  uint32 `xml:"bandwidth,attr"`
}

// StringNoWhitespaceType is String without white spaces
type StringNoWhitespaceType string

// SubRepresentationType is SubRepresentation
type SubRepresentationType struct {
	Level            uint32            `xml:"level,attr,omitempty"`
	DependencyLevel  *UIntVectorType   `xml:"dependencyLevel,attr,omitempty"`
	Bandwidth        uint32            `xml:"bandwidth,attr,omitempty"`
	ContentComponent *StringVectorType `xml:"contentComponent,attr,omitempty"`
	*RepresentationBaseType
}

// RepresentationBaseType is Representation base (common attributes and elements)
type RepresentationBaseType struct {
	Profiles                  string                       `xml:"profiles,attr,omitempty"`
	Width                     uint32                       `xml:"width,attr,omitempty"`
	Height                    uint32                       `xml:"height,attr,omitempty"`
	Sar                       string                       `xml:"sar,attr,omitempty"`
	FrameRate                 string                       `xml:"frameRate,attr,omitempty"`
	AudioSamplingRate         *UIntVectorType              `xml:"audioSamplingRate,attr,omitempty"`
	MimeType                  string                       `xml:"mimeType,attr,omitempty"`
	SegmentProfiles           *ListOf4CCType               `xml:"segmentProfiles,attr,omitempty"`
	Codecs                    string                       `xml:"codecs,attr,omitempty"`
	ContainerProfiles         *ListOf4CCType               `xml:"containerProfiles,attr,omitempty"`
	MaximumSAPPeriod          float64                      `xml:"maximumSAPPeriod,attr,omitempty"`
	StartWithSAP              uint32                       `xml:"startWithSAP,attr,omitempty"`
	MaxPlayoutRate            float64                      `xml:"maxPlayoutRate,attr,omitempty"`
	CodingDependency          bool                         `xml:"codingDependency,attr,omitempty"`
	ScanType                  string                       `xml:"scanType,attr,omitempty"`
	SelectionPriority         uint32                       `xml:"selectionPriority,attr,omitempty"`
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

// ContentProtectionType is Content Protection
type ContentProtectionType struct {
	Robustness string `xml:"robustness,attr,omitempty"`
	RefId      string `xml:"refId,attr,omitempty"`
	Ref        string `xml:"ref,attr,omitempty"`
	*DescriptorType
}

// ResyncType is Resynchronization Point
type ResyncType struct {
	Type   uint32  `xml:"type,attr,omitempty"`
	DT     uint32  `xml:"dT,attr,omitempty"`
	DImax  float32 `xml:"dImax,attr,omitempty"`
	DImin  float32 `xml:"dImin,attr,omitempty"`
	Marker bool    `xml:"marker,attr,omitempty"`
}

// PR ...
type PR struct {
	PopularityRate uint32 `xml:"popularityRate,attr,omitempty"`
	Start          uint64 `xml:"start,attr,omitempty"`
	R              int    `xml:"r,attr,omitempty"`
}

// ContentPopularityRateType is Content Popularity Rate
type ContentPopularityRateType struct {
	Source            string `xml:"source,attr"`
	Sourcedescription string `xml:"source_description,attr,omitempty"`
	PR                []*PR  `xml:"PR"`
}

// LabelType is Label and Group Label
type LabelType struct {
	Id    uint32 `xml:"id,attr,omitempty"`
	Lang  string `xml:"lang,attr,omitempty"`
	Value string `xml:",chardata"`
}

// ProducerReferenceTimeType is Producer Reference time
type ProducerReferenceTimeType struct {
	Id                uint32          `xml:"id,attr"`
	Inband            bool            `xml:"inband,attr,omitempty"`
	Type              string          `xml:"type,attr,omitempty"`
	ApplicationScheme string          `xml:"applicationScheme,attr,omitempty"`
	WallClockTime     string          `xml:"wallClockTime,attr"`
	PresentationTime  uint64          `xml:"presentationTime,attr"`
	UTCTiming         *DescriptorType `xml:"UTCTiming"`
}

// ProducerReferenceTimeTypeType ...
type ProducerReferenceTimeTypeType string

// PreselectionType is Preselection
type PreselectionType struct {
	Id                     string            `xml:"id,attr,omitempty"`
	PreselectionComponents *StringVectorType `xml:"preselectionComponents,attr"`
	Lang                   string            `xml:"lang,attr,omitempty"`
	Order                  string            `xml:"order,attr,omitempty"`
	Accessibility          []*DescriptorType `xml:"Accessibility"`
	Role                   []*DescriptorType `xml:"Role"`
	Rating                 []*DescriptorType `xml:"Rating"`
	Viewpoint              []*DescriptorType `xml:"Viewpoint"`
	*RepresentationBaseType
}

// AudioSamplingRateType is Audio Sampling Rate
type AudioSamplingRateType *UIntVectorType

// SAPType is Stream Access Point type enumeration
type SAPType uint32

// VideoScanType is Video Scan type enumeration
type VideoScanType string

// TagType is Tag
type TagType string

// SubsetType is Subset
type SubsetType struct {
	Contains *UIntVectorType `xml:"contains,attr"`
	Id       string          `xml:"id,attr,omitempty"`
}

// SwitchingType is Switching
type SwitchingType struct {
	Interval uint32 `xml:"interval,attr"`
	Type     string `xml:"type,attr,omitempty"`
}

// SwitchingTypeType is Switching Type type enumeration
type SwitchingTypeType string

// RandomAccessType is Random Access
type RandomAccessType struct {
	Interval      uint32 `xml:"interval,attr"`
	Type          string `xml:"type,attr,omitempty"`
	MinBufferTime string `xml:"minBufferTime,attr,omitempty"`
	Bandwidth     uint32 `xml:"bandwidth,attr,omitempty"`
}

// RandomAccessTypeType is Random Access Type type enumeration
type RandomAccessTypeType string

// PreselectionOrderType is Preselection Order type
type PreselectionOrderType string

// SegmentBaseType is Segment information base
type SegmentBaseType struct {
	Timescale                uint32               `xml:"timescale,attr,omitempty"`
	EptDelta                 int                  `xml:"eptDelta,attr,omitempty"`
	PdDelta                  int                  `xml:"pdDelta,attr,omitempty"`
	PresentationTimeOffset   uint64               `xml:"presentationTimeOffset,attr,omitempty"`
	PresentationDuration     uint64               `xml:"presentationDuration,attr,omitempty"`
	TimeShiftBufferDepth     string               `xml:"timeShiftBufferDepth,attr,omitempty"`
	IndexRange               string               `xml:"indexRange,attr,omitempty"`
	IndexRangeExact          bool                 `xml:"indexRangeExact,attr,omitempty"`
	AvailabilityTimeOffset   float64              `xml:"availabilityTimeOffset,attr,omitempty"`
	AvailabilityTimeComplete bool                 `xml:"availabilityTimeComplete,attr,omitempty"`
	Initialization           *URLType             `xml:"Initialization"`
	RepresentationIndex      *URLType             `xml:"RepresentationIndex"`
	FailoverContent          *FailoverContentType `xml:"FailoverContent"`
}

// MultipleSegmentBaseType is Multiple Segment information base
type MultipleSegmentBaseType struct {
	Duration           uint32               `xml:"duration,attr,omitempty"`
	StartNumber        uint32               `xml:"startNumber,attr,omitempty"`
	EndNumber          uint32               `xml:"endNumber,attr,omitempty"`
	SegmentTimeline    *SegmentTimelineType `xml:"SegmentTimeline"`
	BitstreamSwitching *URLType             `xml:"BitstreamSwitching"`
	*SegmentBaseType
}

// URLType is Segment Info item URL/range
type URLType struct {
	SourceURL string `xml:"sourceURL,attr,omitempty"`
	Range     string `xml:"range,attr,omitempty"`
}

// SingleRFC7233RangeType ...
type SingleRFC7233RangeType string

// FCS ...
type FCS struct {
	T uint64 `xml:"t,attr"`
	D uint64 `xml:"d,attr,omitempty"`
}

// FailoverContentType is Failover Content
type FailoverContentType struct {
	Valid bool   `xml:"valid,attr,omitempty"`
	FCS   []*FCS `xml:"FCS"`
}

// SegmentListType is Segment List
type SegmentListType struct {
	XlinkHref    string            `xml:"xlink:href,attr,omitempty"`
	XlinkActuate string            `xml:"xlink:actuate,attr,omitempty"`
	XlinkType    string            `xml:"xlink:type,attr,omitempty"`
	XlinkShow    string            `xml:"xlink:show,attr,omitempty"`
	SegmentURL   []*SegmentURLType `xml:"SegmentURL"`
	*MultipleSegmentBaseType
}

// SegmentURLType is Segment URL
type SegmentURLType struct {
	Media      string `xml:"media,attr,omitempty"`
	MediaRange string `xml:"mediaRange,attr,omitempty"`
	Index      string `xml:"index,attr,omitempty"`
	IndexRange string `xml:"indexRange,attr,omitempty"`
}

// SegmentTemplateType is Segment Template
type SegmentTemplateType struct {
	Media              string `xml:"media,attr,omitempty"`
	Index              string `xml:"index,attr,omitempty"`
	Initialization     string `xml:"initialization,attr,omitempty"`
	BitstreamSwitching string `xml:"bitstreamSwitching,attr,omitempty"`
	*MultipleSegmentBaseType
}

// S ...
type S struct {
	T uint64 `xml:"t,attr,omitempty"`
	N uint64 `xml:"n,attr,omitempty"`
	D uint64 `xml:"d,attr"`
	R int    `xml:"r,attr,omitempty"`
	K uint64 `xml:"k,attr,omitempty"`
}

// SegmentTimelineType is Segment Timeline
type SegmentTimelineType struct {
	S []*S `xml:"S"`
}

// StringVectorType is Whitespace-separated list of strings
type StringVectorType []string

// ListOf4CCType is Whitespace separated list of 4CC
type ListOf4CCType []string

// FourCCType is 4CC as per latest 14496-12
type FourCCType string

// UIntVectorType is Whitespace-separated list of unsigned integers
type UIntVectorType []uint32

// BaseURLType is Base URL
type BaseURLType struct {
	ServiceLocation          string  `xml:"serviceLocation,attr,omitempty"`
	ByteRange                string  `xml:"byteRange,attr,omitempty"`
	AvailabilityTimeOffset   float64 `xml:"availabilityTimeOffset,attr,omitempty"`
	AvailabilityTimeComplete bool    `xml:"availabilityTimeComplete,attr,omitempty"`
	TimeShiftBufferDepth     string  `xml:"timeShiftBufferDepth,attr,omitempty"`
	RangeAccess              bool    `xml:"rangeAccess,attr,omitempty"`
	Value                    string  `xml:",chardata"`
}

// ProgramInformationType is Program Information
type ProgramInformationType struct {
	Lang               string `xml:"lang,attr,omitempty"`
	MoreInformationURL string `xml:"moreInformationURL,attr,omitempty"`
	Title              string `xml:"Title"`
	Source             string `xml:"Source"`
	Copyright          string `xml:"Copyright"`
}

// DescriptorType is Descriptor
type DescriptorType struct {
	SchemeIdUri string `xml:"schemeIdUri,attr"`
	Value       string `xml:"value,attr,omitempty"`
	Id          string `xml:"id,attr,omitempty"`
}

// MetricsType is Metrics
type MetricsType struct {
	Metrics   string            `xml:"metrics,attr"`
	Range     []*RangeType      `xml:"Range"`
	Reporting []*DescriptorType `xml:"Reporting"`
}

// RangeType is Metrics Range
type RangeType struct {
	Starttime string `xml:"starttime,attr,omitempty"`
	Duration  string `xml:"duration,attr,omitempty"`
}

// CodecsType is RFC6381 simp-list without enclosing double quotes
type CodecsType string

// LeapSecondInformationType is Leap Second Information
type LeapSecondInformationType struct {
	AvailabilityStartLeapOffset     int    `xml:"availabilityStartLeapOffset,attr"`
	NextAvailabilityStartLeapOffset int    `xml:"nextAvailabilityStartLeapOffset,attr,omitempty"`
	NextLeapChangeTime              string `xml:"nextLeapChangeTime,attr,omitempty"`
}
