module github.com/pastelnetwork/gonode/walletnode

go 1.16

require (
	github.com/DataDog/zstd v1.4.8
	github.com/anacrolix/sync v0.3.0 // indirect
	github.com/btcsuite/btcutil v1.0.2
	github.com/go-errors/errors v1.4.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/metadb v0.0.0
	github.com/pastelnetwork/gonode/p2p v0.0.0
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/pastelnetwork/gonode/proto v0.0.0
	github.com/pastelnetwork/gonode/raptorq v0.0.0
	github.com/sahilm/fuzzy v0.1.0
	github.com/spf13/viper v1.7.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli/v2 v2.3.0 // indirect
	go.uber.org/atomic v1.8.0 // indirect
	goa.design/goa/v3 v3.4.3
	goa.design/plugins/v3 v3.4.3
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	google.golang.org/genproto v0.0.0-20210729151513-df9385d47c1b // indirect
	google.golang.org/grpc v1.39.0
)

replace (
	github.com/kolesa-team/go-webp => ../go-webp
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/dupedetection => ../dupedetection
	github.com/pastelnetwork/gonode/metadb => ../metadb
	github.com/pastelnetwork/gonode/p2p => ../p2p
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/probe => ../probe
	github.com/pastelnetwork/gonode/proto => ../proto
	github.com/pastelnetwork/gonode/raptorq => ../raptorq
	github.com/pastelnetwork/gonode/supernode => ../supernode
)
