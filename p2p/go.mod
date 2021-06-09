module github.com/pastelnetwork/gonode/p2p

go 1.16

require (
	github.com/anacrolix/envpprof v1.1.1 // indirect
	github.com/anacrolix/sync v0.2.0 // indirect
	github.com/anacrolix/utp v0.1.0
	github.com/bradfitz/iter v0.0.0-20191230175014-e8f45d346db8 // indirect
	github.com/ccding/go-stun/stun v0.0.0-20200514191101-4dc67bcdb029
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mattn/go-sqlite3 v1.14.7
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1 // indirect
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect

)

replace github.com/pastelnetwork/gonode/common => ../common
