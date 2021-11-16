module github.com/pastelnetwork/gonode/messaging

go 1.17

require (
	github.com/AsynkronIT/protoactor-go v0.0.0-20211115123807-aa1fee2f52fc
	github.com/gogo/protobuf v1.3.2
	github.com/pastelnetwork/gonode/common v0.0.0
	google.golang.org/grpc v1.42.0
)

require (
	github.com/Workiva/go-datastructures v1.0.53 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20190107190726-7ed82d9cb717 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)

replace github.com/pastelnetwork/gonode/common => ../common
