module github.com/pastelnetwork/gonode/walletnode

go 1.16

require (
	github.com/DataDog/zstd v1.4.8 // indirect
	github.com/anacrolix/sync v0.3.0 // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/go-errors/errors v1.4.0 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/p2p v0.0.0-20210615205941-9aa3fc8e6fee
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/pastelnetwork/gonode/probe v0.0.0-20210708185245-02e8c0ae2dd3 // indirect
	github.com/pastelnetwork/gonode/proto v0.0.0
	github.com/pastelnetwork/gonode/raptorq v0.0.0
	github.com/sahilm/fuzzy v0.1.0
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/vektra/mockery/v2 v2.9.0 // indirect
	go.uber.org/atomic v1.8.0 // indirect
	goa.design/goa/v3 v3.4.3
	goa.design/plugins/v3 v3.4.3
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	google.golang.org/grpc v1.38.0
)

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/p2p => ../p2p
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/proto => ../proto
	github.com/pastelnetwork/gonode/raptorq => ../raptorq
)
