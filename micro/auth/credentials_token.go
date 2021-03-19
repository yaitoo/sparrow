package auth

import (
	"context"
)

// TokenCredential  authorize via token
type TokenCredential struct {
	RequireTLS bool
	Token      string
}

//GetRequestMetadata build token metadata
func (c TokenCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"x-token": c.Token,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (c TokenCredential) RequireTransportSecurity() bool {
	if c.RequireTLS {
		return true
	}

	return false
}

//TokenFromContext get token from context
func TokenFromContext(ctx context.Context) string {
	md := FromIncoming(ctx)
	return md.LastValue("x-token", "")
}
