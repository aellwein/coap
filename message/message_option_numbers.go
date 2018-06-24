package message

type OptionFormat uint8

const (
	Empty OptionFormat = iota
	Opaque
	Uint
	String
)

type OptionDefinition struct {
	C       bool
	U       bool
	N       bool
	R       bool
	Name    string
	Format  OptionFormat
	Default interface{}
}

const (
	IfMatch       OptionNumberType = iota + 1
	UriHost                        = 3
	ETag                           = 4
	IfNoneMatch                    = 5
	UriPort                        = 7
	LocationPath                   = 8
	UriPath                        = 11
	ContentFormat                  = 12
	MaxAge                         = 14
	UriQuery                       = 15
	Accept                         = 17
	LocationQuery                  = 20
	ProxyUri                       = 35
	ProxyScheme                    = 39
	Size1                          = 60
)

// Lookup table for possible options.
var OptionLookupTable = map[OptionNumberType]OptionDefinition{
	IfMatch: {
		C:      true,
		U:      false,
		N:      false,
		R:      true,
		Name:   "If-Match",
		Format: Opaque,
	},
	UriHost: {
		C:      true,
		U:      true,
		N:      false,
		R:      false,
		Name:   "Uri-Host",
		Format: String,
	},
	ETag: {
		C:      false,
		U:      false,
		N:      false,
		R:      true,
		Name:   "ETag",
		Format: Opaque,
	},
	IfNoneMatch: {
		C:      true,
		U:      false,
		N:      false,
		R:      false,
		Name:   "If-None-Match",
		Format: Empty,
	},
	UriPort: {
		C:      true,
		U:      true,
		N:      false,
		R:      false,
		Name:   "Uri-Port",
		Format: Uint,
	},
	LocationPath: {
		C:      false,
		U:      false,
		N:      false,
		R:      true,
		Name:   "Location-Path",
		Format: String,
	},
	UriPath: {
		C:      true,
		U:      true,
		N:      false,
		R:      true,
		Name:   "Uri-Path",
		Format: String,
	},
	ContentFormat: {
		C:      false,
		U:      false,
		N:      false,
		R:      false,
		Name:   "Content-Format",
		Format: Uint,
	},
	MaxAge: {
		C:       false,
		U:       true,
		N:       false,
		R:       false,
		Name:    "Max-Age",
		Format:  Uint,
		Default: 60,
	},
	UriQuery: {
		C:      true,
		U:      true,
		N:      false,
		R:      true,
		Name:   "Uri-Query",
		Format: String,
	},
	Accept: {
		C:      true,
		U:      false,
		N:      false,
		R:      false,
		Name:   "Accept",
		Format: Uint,
	},
	LocationQuery: {
		C:      false,
		U:      false,
		N:      false,
		R:      true,
		Name:   "Location-Query",
		Format: Uint,
	},
	ProxyUri: {
		C:      true,
		U:      true,
		N:      false,
		R:      false,
		Name:   "Proxy-Uri",
		Format: String,
	},

	ProxyScheme: {
		C:      true,
		U:      true,
		N:      false,
		R:      false,
		Name:   "Proxy-Scheme",
		Format: Uint,
	},

	Size1: {
		C:      false,
		U:      false,
		N:      true,
		R:      false,
		Name:   "Size1",
		Format: Uint,
	},
}

// option number is in uint16 range
type OptionNumberType uint16
