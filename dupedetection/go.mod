module github.com/pastelnetwork/gonode/dupedetection

go 1.16

require (
	github.com/btcsuite/btcutil v1.0.2
	github.com/google/gofuzz v1.2.0
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/metadb v0.0.0-20210803173228-bd776aecea6a
	github.com/pastelnetwork/gonode/p2p v0.0.0-20210803173228-bd776aecea6a
	github.com/pastelnetwork/gonode/pastel v0.0.0-20210803173228-bd776aecea6a
	github.com/sbinet/npyio v0.5.2
	github.com/stretchr/testify v1.7.0
)

replace (
	github.com/pastelnetwork/gonode/common => ../common
	github.com/pastelnetwork/gonode/p2p => ../p2p
	github.com/pastelnetwork/gonode/pastel => ../pastel
)
