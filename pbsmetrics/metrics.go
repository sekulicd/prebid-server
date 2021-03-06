package pbsmetrics

import (
	"time"

	"github.com/prebid/prebid-server/openrtb_ext"
)

// Labels defines the labels that can be attached to the metrics.
type Labels struct {
	Source        DemandSource
	RType         RequestType
	PubID         string // exchange specific ID, so we cannot compile in values
	Browser       Browser
	CookieFlag    CookieFlag
	RequestStatus RequestStatus
}

// AdapterLabels defines the labels that can be attached to the adapter metrics.
type AdapterLabels struct {
	Source        DemandSource
	RType         RequestType
	Adapter       openrtb_ext.BidderName
	PubID         string // exchange specific ID, so we cannot compile in values
	Browser       Browser
	CookieFlag    CookieFlag
	AdapterBids   AdapterBid
	AdapterErrors map[AdapterError]struct{}
}

// ImpLabels defines metric labels describing the impression type.
type ImpLabels struct {
	BannerImps bool
	VideoImps  bool
	AudioImps  bool
	NativeImps bool
}

// RequestLabels defines metric labels describing the result of a network request.
type RequestLabels struct {
	RequestStatus RequestStatus
}

// Label typecasting. Se below the type definitions for possible values

// DemandSource : Demand source enumeration
type DemandSource string

// ImpMediaType : Media type described in the "imp" JSON object  TODO is this still needed?
type ImpMediaType string

// RequestType : Request type enumeration
type RequestType string

// Browser type enumeration
type Browser string

// CookieFlag : User ID cookie exists flag
type CookieFlag string

// RequestStatus : The request return status
type RequestStatus string

// AdapterBid : Whether or not the adapter returned bids
type AdapterBid string

// AdapterError : Errors which may have occurred during the adapter's execution
type AdapterError string

// CacheResult : Cache hit/miss
type CacheResult string

// PublisherUnknown : Default value for Labels.PubID
const PublisherUnknown = "unknown"

// The demand sources
const (
	DemandWeb     DemandSource = "web"
	DemandApp     DemandSource = "app"
	DemandUnknown DemandSource = "unknown"
)

func DemandTypes() []DemandSource {
	return []DemandSource{
		DemandWeb,
		DemandApp,
		DemandUnknown,
	}
}

// The request types (endpoints)
const (
	ReqTypeLegacy   RequestType = "legacy"
	ReqTypeORTB2Web RequestType = "openrtb2-web"
	ReqTypeORTB2App RequestType = "openrtb2-app"
	ReqTypeAMP      RequestType = "amp"
	ReqTypeVideo    RequestType = "video"
)

// The media types described in the "imp" json objects
const (
	ImpTypeBanner ImpMediaType = "banner"
	ImpTypeVideo  ImpMediaType = "video"
	ImpTypeAudio  ImpMediaType = "audio"
	ImpTypeNative ImpMediaType = "native"
)

func RequestTypes() []RequestType {
	return []RequestType{
		ReqTypeLegacy,
		ReqTypeORTB2Web,
		ReqTypeORTB2App,
		ReqTypeAMP,
		ReqTypeVideo,
	}
}

func ImpTypes() []ImpMediaType {
	return []ImpMediaType{
		ImpTypeBanner,
		ImpTypeVideo,
		ImpTypeAudio,
		ImpTypeNative,
	}
}

// Browser flag; at this point we only care about identifying Safari
const (
	BrowserSafari Browser = "safari"
	BrowserOther  Browser = "other"
)

func BrowserTypes() []Browser {
	return []Browser{
		BrowserSafari,
		BrowserOther,
	}
}

// Cookie flag
const (
	CookieFlagYes     CookieFlag = "exists"
	CookieFlagNo      CookieFlag = "no"
	CookieFlagUnknown CookieFlag = "unknown"
)

func CookieTypes() []CookieFlag {
	return []CookieFlag{
		CookieFlagYes,
		CookieFlagNo,
		CookieFlagUnknown,
	}
}

// Request/return status
const (
	RequestStatusOK          RequestStatus = "ok"
	RequestStatusBadInput    RequestStatus = "badinput"
	RequestStatusErr         RequestStatus = "err"
	RequestStatusNetworkErr  RequestStatus = "networkerr"
	RequestStatusBlacklisted RequestStatus = "blacklistedacctorapp"
)

func RequestStatuses() []RequestStatus {
	return []RequestStatus{
		RequestStatusOK,
		RequestStatusBadInput,
		RequestStatusErr,
		RequestStatusNetworkErr,
		RequestStatusBlacklisted,
	}
}

// Adapter bid response status.
const (
	AdapterBidPresent AdapterBid = "bid"
	AdapterBidNone    AdapterBid = "nobid"
)

func AdapterBids() []AdapterBid {
	return []AdapterBid{
		AdapterBidPresent,
		AdapterBidNone,
	}
}

// Adapter execution status
const (
	AdapterErrorBadInput            AdapterError = "badinput"
	AdapterErrorBadServerResponse   AdapterError = "badserverresponse"
	AdapterErrorTimeout             AdapterError = "timeout"
	AdapterErrorFailedToRequestBids AdapterError = "failedtorequestbid"
	AdapterErrorUnknown             AdapterError = "unknown_error"
)

func AdapterErrors() []AdapterError {
	return []AdapterError{
		AdapterErrorBadInput,
		AdapterErrorBadServerResponse,
		AdapterErrorTimeout,
		AdapterErrorFailedToRequestBids,
		AdapterErrorUnknown,
	}
}

const (
	// CacheHit represents a cache hit i.e the key was found in cache
	CacheHit CacheResult = "hit"
	// CacheMiss represents a cache miss i.e that key wasn't found in cache
	// and had to be fetched from the backend
	CacheMiss CacheResult = "miss"
)

// CacheResults returns possible cache results i.e. cache hit or miss
func CacheResults() []CacheResult {
	return []CacheResult{
		CacheHit,
		CacheMiss,
	}
}

// UserLabels : Labels for /setuid endpoint
type UserLabels struct {
	Action RequestAction
	Bidder openrtb_ext.BidderName
}

// RequestAction : The setuid request result
type RequestAction string

// /setuid action labels
const (
	RequestActionSet    RequestAction = "set"
	RequestActionOptOut RequestAction = "opt_out"
	RequestActionGDPR   RequestAction = "gdpr"
	RequestActionErr    RequestAction = "err"
)

// RequestActions returns possible setuid action labels
func RequestActions() []RequestAction {
	return []RequestAction{
		RequestActionSet,
		RequestActionOptOut,
		RequestActionGDPR,
		RequestActionErr,
	}
}

// MetricsEngine is a generic interface to record PBS metrics into the desired backend
// The first three metrics function fire off once per incoming request, so total metrics
// will equal the total numer of incoming requests. The remaining 5 fire off per outgoing
// request to a bidder adapter, so will record a number of hits per incoming request. The
// two groups should be consistent within themselves, but comparing numbers between groups
// is generally not useful.
type MetricsEngine interface {
	RecordConnectionAccept(success bool)
	RecordConnectionClose(success bool)
	RecordRequest(labels Labels)                           // ignores adapter. only statusOk and statusErr fom status
	RecordImps(labels ImpLabels)                           // RecordImps across openRTB2 engines that support the 'Native' Imp Type
	RecordLegacyImps(labels Labels, numImps int)           // RecordImps for the legacy engine
	RecordRequestTime(labels Labels, length time.Duration) // ignores adapter. only statusOk and statusErr fom status
	RecordAdapterRequest(labels AdapterLabels)
	RecordAdapterPanic(labels AdapterLabels)
	// This records whether or not a bid of a particular type uses `adm` or `nurl`.
	// Since the legacy endpoints don't have a bid type, it can only count bids from OpenRTB and AMP.
	RecordAdapterBidReceived(labels AdapterLabels, bidType openrtb_ext.BidType, hasAdm bool)
	RecordAdapterPrice(labels AdapterLabels, cpm float64)
	RecordAdapterTime(labels AdapterLabels, length time.Duration)
	RecordCookieSync()
	RecordAdapterCookieSync(adapter openrtb_ext.BidderName, gdprBlocked bool)
	RecordUserIDSet(userLabels UserLabels) // Function should verify bidder values
	RecordStoredReqCacheResult(cacheResult CacheResult, inc int)
	RecordStoredImpCacheResult(cacheResult CacheResult, inc int)
	RecordPrebidCacheRequestTime(success bool, length time.Duration)
}
