module github.com/pastelnetwork/gonode/walletnode

go 1.16

require (
	github.com/anacrolix/sync v0.3.0 // indirect
	github.com/bradfitz/slice v0.0.0-20180809154707-2b758aa73013 // indirect
	github.com/dimfeld/httptreemux/v5 v5.3.0 // indirect
	github.com/go-errors/errors v1.4.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6 // indirect
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/p2p v0.0.0-20210615205941-9aa3fc8e6fee // indirect
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/pastelnetwork/gonode/proto v0.0.0
	github.com/sahilm/fuzzy v0.1.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0
	go4.org v0.0.0-20201209231011-d4a079459e60 // indirect
	goa.design/goa/v3 v3.3.1
	goa.design/plugins/v3 v3.3.1
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/image v0.0.0-20210216034530-4410531fe030 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	golang.org/x/tools v0.1.1 // indirect
	google.golang.org/genproto v0.0.0-20210510173355-fb37daa5cd7a // indirect
	google.golang.org/grpc v1.37.0
)

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/pastel => ../pastel
	github.com/pastelnetwork/gonode/proto => ../proto
	github.com/pastelnetwork/gonode/p2p => ../p2p
)
