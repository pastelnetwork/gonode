module github.com/pastelnetwork/gonode/walletnode

go 1.16

require (
	github.com/DataDog/zstd v1.4.8 // indirect
	github.com/dimfeld/httptreemux/v5 v5.3.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/pastelnetwork/gonode/proto v0.0.0
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0
	goa.design/goa/v3 v3.3.1
	goa.design/plugins/v3 v3.3.1
	golang.org/x/image v0.0.0-20210216034530-4410531fe030 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210511113859-b0526f3d8744 // indirect
	golang.org/x/tools v0.1.1 // indirect
	google.golang.org/genproto v0.0.0-20210510173355-fb37daa5cd7a // indirect
	google.golang.org/grpc v1.37.0
)

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/proto => ../proto
)
