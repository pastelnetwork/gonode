module github.com/pastelnetwork/gonode/raptorq

go 1.16

require (
	github.com/pastelnetwork/gonode/common v0.0.0-20210624142025-67ad31676597
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

replace github.com/pastelnetwork/gonode/raptorq/grpc => ./node/grpc
