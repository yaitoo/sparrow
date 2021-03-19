package auth

import (
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// Metadata is a convenience wrapper definiting extra functions on the metadata.
type Metadata metadata.MD

func (m Metadata) AsMD() *metadata.MD {
	md := metadata.MD(m)
	return &md
}

// FromIncoming extracts an inbound metadata from the server-side context.
//
// This function always returns a Metadata wrapper of the metadata.MD, in case the context doesn't have metadata it returns
// a new empty Metadata.
func FromIncoming(ctx context.Context) Metadata {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return Metadata(metadata.Pairs())
	}
	return Metadata(md)
}

// FromOutgoing extracts an outbound metadata from the client-side context.
//
// This function always returns a Metadata wrapper of the metadata.MD, in case the context doesn't have metadata it returns
// a new empty Metadata.
func FromOutgoing(ctx context.Context) Metadata {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return Metadata(metadata.Pairs())
	}
	return Metadata(md)
}

// ToOutgoing sets the given Metadata as a client-side context for dispatching.
func (m Metadata) ToOutgoing(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.MD(m))
}

// ToIncoming sets the given Metadata as a server-side context for dispatching.
//
// This is mostly useful in ServerInterceptors..
func (m Metadata) ToIncoming(ctx context.Context) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.MD(m))
}

// FirstValue retrieves a single value from the metadata.
//
// It works analogously to http.Header.Get, returning the first value if there are many set. If the value is not set,
// an empty string is returned.
func (m Metadata) FirstValue(key string, defalutValue string) string {

	vv, ok := m[strings.ToLower(key)]
	if !ok {
		return defalutValue
	}
	return vv[0]
}

// LastValue retrieves a single value from the metadata.
func (m Metadata) LastValue(key string, defalutValue string) string {

	vv, ok := m[strings.ToLower(key)]
	if !ok {
		return defalutValue
	}
	return vv[len(vv)-1]
}

// Values retrieves all values from the metadata.
func (m Metadata) Values(key string) []string {

	vv, ok := m[strings.ToLower(key)]
	if !ok {
		return []string{}
	}
	return vv
}

// Del retrieves a single value from the metadata.
//
// It works analogously to http.Header.Del, deleting all values if they exist.
func (m Metadata) Del(key string) Metadata {

	delete(m, strings.ToLower(key))
	return m
}

// Set sets the given value in a metadata.
//
// It works analogously to http.Header.Set, overwriting all previous metadata values.
func (m Metadata) Set(key string, value string) Metadata {
	m[strings.ToLower(key)] = []string{value}
	return m
}

// Add retrieves a single value from the metadata.
//
// It works analogously to http.Header.Add, as it appends to any existing values associated with key.
func (m Metadata) Add(key string, value string) Metadata {
	k := strings.ToLower(key)
	m[k] = append(m[k], value)
	return m
}
