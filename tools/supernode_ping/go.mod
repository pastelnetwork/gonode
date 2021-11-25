module github.com/pastelnetwork/gonode/tools/supernode_ping

go 1.17

require (
	github.com/pastelnetwork/gonode/common v0.0.0-20210828164036-dff1dd45f2d0
	github.com/pastelnetwork/gonode/proto v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.42.0
)

require (
	github.com/go-errors/errors v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jpillora/longestcommon v0.0.0-20161227235612-adb9d91ee629 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023 // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace (
	github.com/pastelnetwork/gonode/common => ../../common
	github.com/pastelnetwork/gonode/proto => ../../proto
)
