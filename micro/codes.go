package micro

import (
	"google.golang.org/grpc/codes"
)

const (
	// ErrContinue           codes.Code = 100 // RFC 7231, 6.2.1
	// ErrSwitchingProtocols codes.Code = 101 // RFC 7231, 6.2.2
	// ErrProcessing         codes.Code = 102 // RFC 2518, 10.1

	// ErrOK                   codes.Code = 200 // RFC 7231, 6.3.1
	// ErrCreated              codes.Code = 201 // RFC 7231, 6.3.2
	// ErrAccepted             codes.Code = 202 // RFC 7231, 6.3.3
	// ErrNonAuthoritativeInfo codes.Code = 203 // RFC 7231, 6.3.4
	// ErrNoContent            codes.Code = 204 // RFC 7231, 6.3.5
	// ErrResetContent         codes.Code = 205 // RFC 7231, 6.3.6
	// ErrPartialContent       codes.Code = 206 // RFC 7233, 4.1
	// ErrMultiErr          codes.Code = 207 // RFC 4918, 11.1
	// ErrAlreadyReported      codes.Code = 208 // RFC 5842, 7.1
	// ErrIMUsed               codes.Code = 226 // RFC 3229, 10.4.1

	// ErrMultipleChoices   codes.Code = 300 // RFC 7231, 6.4.1
	// ErrMovedPermanently  codes.Code = 301 // RFC 7231, 6.4.2
	// ErrFound             codes.Code = 302 // RFC 7231, 6.4.3
	// ErrSeeOther          codes.Code = 303 // RFC 7231, 6.4.4
	// ErrNotModified       codes.Code = 304 // RFC 7232, 4.1
	// ErrUseProxy          codes.Code = 305 // RFC 7231, 6.4.5
	// _                       codes.Code = 306 // RFC 7231, 6.4.6 (Unused)
	// ErrTemporaryRedirect codes.Code = 307 // RFC 7231, 6.4.7
	// ErrPermanentRedirect codes.Code = 308 // RFC 7538, 3

	// ErrBadRequest                   codes.Code = 400 // RFC 7231, 6.5.1
	//ErrUnauthorized                 codes.Code = 401 // RFC 7235, 3.1
	ErrPaymentRequired codes.Code = 402 // RFC 7231, 6.5.2
	//ErrForbidden                    codes.Code = 403 // RFC 7231, 6.5.3
	//ErrNotFound                     codes.Code = 404 // RFC 7231, 6.5.4
	ErrMethodNotAllowed             codes.Code = 405 // RFC 7231, 6.5.5
	ErrNotAcceptable                codes.Code = 406 // RFC 7231, 6.5.6
	ErrProxyAuthRequired            codes.Code = 407 // RFC 7235, 3.2
	ErrRequestTimeout               codes.Code = 408 // RFC 7231, 6.5.7
	ErrConflict                     codes.Code = 409 // RFC 7231, 6.5.8
	ErrGone                         codes.Code = 410 // RFC 7231, 6.5.9
	ErrLengthRequired               codes.Code = 411 // RFC 7231, 6.5.10
	ErrPreconditionFailed           codes.Code = 412 // RFC 7232, 4.2
	ErrRequestEntityTooLarge        codes.Code = 413 // RFC 7231, 6.5.11
	ErrRequestURITooLong            codes.Code = 414 // RFC 7231, 6.5.12
	ErrUnsupportedMediaType         codes.Code = 415 // RFC 7231, 6.5.13
	ErrRequestedRangeNotSatisfiable codes.Code = 416 // RFC 7233, 4.4
	ErrExpectationFailed            codes.Code = 417 // RFC 7231, 6.5.14
	ErrTeapot                       codes.Code = 418 // RFC 7168, 2.3.3
	ErrUnprocessableEntity          codes.Code = 422 // RFC 4918, 11.2
	ErrLocked                       codes.Code = 423 // RFC 4918, 11.3
	ErrFailedDependency             codes.Code = 424 // RFC 4918, 11.4
	ErrUpgradeRequired              codes.Code = 426 // RFC 7231, 6.5.15
	ErrPreconditionRequired         codes.Code = 428 // RFC 6585, 3
	ErrTooManyRequests              codes.Code = 429 // RFC 6585, 4
	ErrRequestHeaderFieldsTooLarge  codes.Code = 431 // RFC 6585, 5
	ErrUnavailableForLegalReasons   codes.Code = 451 // RFC 7725, 3

	//ErrInternalServerError           codes.Code = 500 // RFC 7231, 6.6.1
	ErrNotImplemented                codes.Code = 501 // RFC 7231, 6.6.2
	ErrBadGateway                    codes.Code = 502 // RFC 7231, 6.6.3
	ErrServiceUnavailable            codes.Code = 503 // RFC 7231, 6.6.4
	ErrGatewayTimeout                codes.Code = 504 // RFC 7231, 6.6.5
	ErrHTTPVersionNotSupported       codes.Code = 505 // RFC 7231, 6.6.6
	ErrVariantAlsoNegotiates         codes.Code = 506 // RFC 2295, 8.1
	ErrInsufficientStorage           codes.Code = 507 // RFC 4918, 11.5
	ErrLoopDetected                  codes.Code = 508 // RFC 5842, 7.2
	ErrNotExtended                   codes.Code = 510 // RFC 2774, 7
	ErrNetworkAuthenticationRequired codes.Code = 511 // RFC 6585, 6
)

const (
	//ErrUnknown unkown error code
	//ErrUnknown = codes.Unknown

	//ErrPasswordMismatch login and password are mismatched
	ErrPasswordMismatch codes.Code = 1000

	//ErrObjectUnavailable object is unavaiable/disabled
	ErrObjectUnavailable codes.Code = 1001
)
