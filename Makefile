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
	cd ./dupedetection/ddclient && go generate ./...
	cd ./p2p && go generate ./...
	cd ./common && go generate ./...
	#cd ./metadb && go generate ./...
	cd ./pastel && go generate ./...
	cd ./raptorq/node && go generate ./...
	cd ./supernode/node && go generate ./...
	cd ./walletnode/node && go generate ./...
gen-proto:
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative process_userdata_sn.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_sense_sn.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_nft_sn.proto
	cd ./proto/supernode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative common_sn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative process_userdata_wn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative download_nft.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_sense_wn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative register_nft_wn.proto
	cd ./proto/walletnode/protobuf && protoc --go_out=.. --go-grpc_out=.. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative common_wn.proto
integration-tests:
	cd ./integration && INTEGRATION_TEST_ENV=true go test -v --timeout=10m ./...
build:
	cd ./supernode && go build ./...
	cd ./walletnode && go build ./...