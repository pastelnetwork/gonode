package version

// This variable is set at build time using -ldflags parameters. For example, we typically set this flag in circle.yml
// to the latest Git tag when building our Go apps:
//
// go build -o my-app -ldflags "-X github.com/pastelnetwork/gonode/common/version.version=v0.0.1"
//
// For more info, see: http://stackoverflow.com/a/11355611/483528

var (
	version = "latest"
)

// Version composes a version of the package
func Version() string {
	return version
}
