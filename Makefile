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
LDFLAGS="-s -w -X github.com/pastelnetwork/gonode/common/version.version=$(VERSION)"
BINARY_SN ?= supernode
BINARY_WN ?= walletnode

wn-unit-tests:
	cd ./walletnode && go test -race ./...
	cd ./common && go test -race ./...
	cd ./pastel && go test -race ./...
	cd ./raptorq && go test -race ./...
sn-unit-tests:
	cd ./supernode && go test -race ./...
	cd ./dupedetection && go test -race ./...
	#cd ./metadb && go test -race ./...
	cd ./raptorq && go test -race ./...
	cd ./p2p && go test -race ./...
	cd ./pastel && go test -race ./...
	cd ./common && go test -race ./...
gen-mock:
	rm -r -f ./dupedetection/ddclient/mocks
	rm -r -f ./p2p/mocks
	rm -r -f ./pastel/mocks
	rm -r -f ./raptorq/node/mocks
	rm -r -f ./supernode/node/mocks
	rm -r -f ./walletnode/node/mocks
	cd ./dupedetection/ddclient && go generate ./...
	cd ./p2p && go generate ./...
	cd ./common && go generate ./...
	#cd ./metadb && go generate ./...
	cd ./pastel && go generate ./...
	cd ./raptorq/node && go generate ./...
	cd ./supernode/node && go generate ./...
	cd ./walletnode/node && go generate ./...
gen-proto:
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_cascade_sn.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_sense_sn.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_nft_sn.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative common_sn.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative storage_challenge.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative download_nft.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_sense_wn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_cascade_wn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_nft_wn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative common_wn.proto
integration-tests:
	cd ./integration && INTEGRATION_TEST_ENV=true go test -v --timeout=10m ./...
build:
	cd ./supernode && go build ./...
	cd ./walletnode && go build ./...
clean-proto:
	rm -f ./proto/supernode/*.go
	rm -f ./proto/walletnode/*.go

release:
	cd ./walletnode && go build -ldflags=$(LDFLAGS) -o $(BINARY_WN)
	strip -v walletnode/$(BINARY_WN) -o dist/$(BINARY_WN)-$(os)-$(arch)$(ext)

	docker build -f ./scripts/release.Dockerfile --build-arg LD_FLAGS=$(LDFLAGS) -t gonode_release .
	docker create -ti --name gn_builder gonode_release bash
	docker cp gn_builder:/walletnode/walletnode-win32.exe ./dist/$(BINARY_WN)-win32.exe
	docker cp gn_builder:/walletnode/walletnode-win64.exe ./dist/$(BINARY_WN)-win64.exe
	docker cp gn_builder:/walletnode/walletnode-linux-amd64 ./dist/$(BINARY_WN)-linux-amd64
	docker cp gn_builder:/supernode/supernode-linux-amd64 ./dist/$(BINARY_SN)-linux-amd64
	docker rm -f gn_builder

