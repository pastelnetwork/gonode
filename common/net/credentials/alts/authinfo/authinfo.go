package authinfo

import "google.golang.org/grpc/credentials"

// altsAuthInfo exposes security information from the ALTS handshake to the
// application. altsAuthInfo is immutable and implements credentials.AuthInfo.
type altsAuthInfo struct {
	credentials.CommonAuthInfo
}

// New returns a new altsAuthInfo object given handshaker results.
func New() credentials.AuthInfo {
	return newAuthInfo()
}

func newAuthInfo() *altsAuthInfo {
	return &altsAuthInfo{
		CommonAuthInfo: credentials.CommonAuthInfo{SecurityLevel: credentials.PrivacyAndIntegrity},
	}
}

// AuthType identifies the context as providing ALTS authentication information.
func (s *altsAuthInfo) AuthType() string {
	return "alts"
}
