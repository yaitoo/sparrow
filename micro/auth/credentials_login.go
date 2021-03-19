package auth

import (
	"context"
)

// LoginCredential  authorize via login and password
type LoginCredential struct {
	RequireTLS bool
	Login      string
	Passwd     string
}

//GetRequestMetadata build login metadata
func (c *LoginCredential) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {

	return map[string]string{
		"x-login":  c.Login,
		"x-passwd": c.Passwd,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (c *LoginCredential) RequireTransportSecurity() bool {
	return c.RequireTLS
}

//LoginFromContext get login from context
func LoginFromContext(ctx context.Context) (string, string) {
	md := FromIncoming(ctx)
	return md.LastValue("x-login", ""), md.LastValue("x-passwd", "")
}
