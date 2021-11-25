module github.com/pastelnetwork/gonode/messaging

go 1.17

require (
	github.com/AsynkronIT/protoactor-go v0.0.0-20211124041449-becb6dbfc022
	github.com/gogo/protobuf v1.3.2
	github.com/pastelnetwork/gonode/common v0.0.0
	google.golang.org/grpc v1.42.0
)

require (
	github.com/Workiva/go-datastructures v1.0.53 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20190107190726-7ed82d9cb717 // indirect
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023 // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)

replace github.com/pastelnetwork/gonode/common => ../common
