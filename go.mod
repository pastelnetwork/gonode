module github.com/pastelnetwork/walletnode

go 1.16

require (
	github.com/dimfeld/httptreemux/v5 v5.3.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/pastelnetwork/go-commons v0.0.3
	github.com/pastelnetwork/go-pastel v0.0.1
	github.com/pastelnetwork/supernode-proto v0.0.0-20210422111752-b8a67357d93e // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	goa.design/goa/v3 v3.3.1
	goa.design/plugins/v3 v3.3.1
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210403161142-5e06dd20ab57 // indirect
	google.golang.org/grpc v1.37.0 // indirect
)

replace github.com/pastelnetwork/go-pastel => ../go-pastel

replace github.com/pastelnetwork/supernode-proto => ../supernode-proto

replace github.com/pastelnetwork/go-commons => ../go-commons
