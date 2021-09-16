
#/usr/bin
export APP_VERSION=$(git describe --tag)
go build -ldflags "-X github.com/pastelnetwork/gonode/common/version.version=$APP_VERSION" ./