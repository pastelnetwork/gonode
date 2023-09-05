VERSION = $(shell git describe --tag)
arch = "amd64"
ifeq ($(OS),Windows_NT)
	os = "windows"
	ext = ".exe"
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		os = "linux"
		ext = ""
	endif
	ifeq ($(UNAME_S),Darwin)
		os = "darwin"
		ext = ""
	endif
endif
LDFLAGS='-s -w -X github.com/pastelnetwork/gonode/common/version.version=$(VERSION)'
BINARY_SN ?= supernode
BINARY_WN ?= walletnode
BINARY_BRIDGE ?= bridge
BINARY_HERMES ?= hermes

wn-unit-tests:
	cd ./walletnode && go test -race ./...
	cd ./common && go test -race ./...
	cd ./pastel && go test -race ./...
	cd ./raptorq && go test -race ./...
sn-unit-tests:
	cd ./supernode && go test -race --timeout=5m ./...
	cd ./dupedetection && go test -race --timeout=5m ./...
	cd ./raptorq && go test -race --timeout=5m ./...
	cd ./p2p && go test -race --timeout=5m ./...
	cd ./pastel && go test -race --timeout=5m ./...
	cd ./common && go test -race --timeout=5m ./...
	cd ./bridge && go test -race --timeout=5m ./...
	cd ./hermes && go test -race --timeout=5m ./...
gen-mock:
	# rm -r -f ./dupedetection/ddclient/mocks
	# rm -r -f ./p2p/mocks
	# rm -r -f ./pastel/mocks
	# rm -r -f ./raptorq/node/mocks
	# rm -r -f ./bridge/node/mocks
	# rm -r -f ./supernode/node/mocks
	# rm -r -f ./walletnode/node/mocks
	# rm -r -f ./hermes/service/mocks
	# rm -r -f ./hermes/store/mocks
	cd ./dupedetection/ddclient && go generate ./...
	cd ./p2p && go generate ./...
	cd ./common && go generate ./...
	cd ./bridge && go generate ./...
	cd ./hermes/service && go generate ./...
	cd ./hermes/store && go generate ./...
	cd ./pastel && go generate ./...
	cd ./raptorq/node && go generate ./...
	cd ./supernode/node && go generate ./...
	cd ./walletnode/node && go generate ./...
gen-proto:
	cd ./dupedetection && protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative dd-server.proto	
	cd ./proto/hermes/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative hermes_p2p.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_cascade_sn.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_sense_sn.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_nft_sn.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative common_sn.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative storage_challenge.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative self_healing.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_collection_sn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative download_nft.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_sense_wn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_cascade_wn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_collection_wn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_nft_wn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative common_wn.proto
	cd ./proto/bridge/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative common_bridge.proto
	cd ./proto/bridge/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative download_data.proto
integration-tests:
	cd ./integration && INTEGRATION_TEST_ENV=true go test -v --timeout=20m ./...
build-linux:
	cd ./supernode && CC=musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' -o supernode
	cd ./walletnode && CC=musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' -o walletnode
	cd ./hermes && CC=musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' -o hermes
clean-proto:
	rm -f ./proto/supernode/*.go
	rm -f ./proto/walletnode/*.go

release:
	#cd ./walletnode && go build -ldflags=$(LDFLAGS) -o $(BINARY_WN)
	#strip -v walletnode/$(BINARY_WN) -o dist/$(BINARY_WN)-$(os)-$(arch)$(ext)

	#cd ./bridge && go build -ldflags=$(LDFLAGS) -o $(BINARY_BRIDGE)
	#strip -v bridge/$(BINARY_BRIDGE) -o dist/$(BINARY_BRIDGE)-$(os)-$(arch)$(ext)

	docker build -f ./scripts/release.Dockerfile --build-arg BUILD_VERSION=$(VERSION) -t gonode_release .
	docker create -ti --name gn_builder gonode_release bash
	docker cp gn_builder:/walletnode/walletnode-win64.exe ./dist/$(BINARY_WN)-win64.exe
	docker cp gn_builder:/walletnode/walletnode-linux-amd64 ./dist/$(BINARY_WN)-linux-amd64
	docker cp gn_builder:/supernode/supernode-linux-amd64 ./dist/$(BINARY_SN)-linux-amd64
	docker cp gn_builder:/bridge/bridge-win64.exe ./dist/$(BINARY_BRIDGE)-win64.exe
	docker cp gn_builder:/bridge/bridge-linux-amd64 ./dist/$(BINARY_BRIDGE)-linux-amd64
	docker cp gn_builder:/hermes/hermes-linux-amd64 ./dist/$(BINARY_HERMES)-linux-amd64
	docker rm -f gn_builder

goa-gen:
	cd ./walletnode/api && goa gen github.com/pastelnetwork/gonode/walletnode/api/design