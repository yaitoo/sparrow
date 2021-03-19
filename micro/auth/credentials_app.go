package auth

import (
	"context"
)

//AppCredential authorize via appkey and appsecret
type AppCredential struct {
	RequireTLS bool
	Key        string
	Secret     string
}

//GetRequestMetadata build app metadata
func (c AppCredential) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {

	return map[string]string{
		"x-app-key":    c.Key,
		"x-app-secret": c.Secret,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (c AppCredential) RequireTransportSecurity() bool {
	return c.RequireTLS
}

//AppFromContext get app from context
func AppFromContext(ctx context.Context) (string, string) {
	md := FromIncoming(ctx)
	return md.LastValue("x-app-key", ""), md.LastValue("x-app-secret", "")
}
