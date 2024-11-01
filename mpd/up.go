package mpd

const (
	UrlParameterNamespace2014 = "urn:mpeg:dash:schema:urlparam:2014"
	UrlParameterNamespace2016 = "urn:mpeg:dash:schema:urlparam:2016"
)

// UrlQueryInfoType is defined in Annex I of ISO/IEC 23009-1.
// Its namespace is UrlParameterNamespace2014.
type UrlQueryInfoType struct {
	QueryTemplate  string `xml:"queryTemplate,attr,omitempty"`
	UseMPDUrlQuery bool   `xml:"useMPDUrlQuery,attr,omitempty"`
	QueryString    string `xml:"queryString,attr,omitempty"`
	XlinkHref      string `xml:"http://www.w3.org/1999/xlink xlink:href,attr,omitempty"`
	XlinkActuate   string `xml:"http://www.w3.org/1999/xlink xlink:actuate,attr,omitempty"` // default = "onRequest"
	XlinkType      string `xml:"http://www.w3.org/1999/xlink xlink:type,attr,omitempty"`    // fixed = "simple"
	XlinkShow      string `xml:"http://www.w3.org/1999/xlink xlink:show,attr,omitempty"`    // fixed = "embed"
}

// ExtendedUrlInfoType is defined in Annex I of ISO/IEC 23009-1.
// Its namespace is UrlParameterNamespace2016.
type ExtendedUrlInfoType struct {
	UrlQueryInfoType
	IncludeInRequests string `xml:"includeInRequests,attr,omitempty"` // default = "segment"
	HeaderParamSource string `xml:"headerParamSource,attr,omitempty"` // default = "segment"
	SameOriginPolicy  bool   `xml:"sameOriginPolicy,attr,omitempty"`  // default = "false"
}
