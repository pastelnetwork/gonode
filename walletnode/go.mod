module github.com/pastelnetwork/gonode/walletnode

go 1.16

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/go-errors/errors v1.4.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/pastelnetwork/gonode/pastel v0.0.0
	github.com/pastelnetwork/gonode/proto v0.0.0
	github.com/pelletier/go-toml v1.9.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0
	goa.design/goa/v3 v3.4.2
	goa.design/plugins/v3 v3.4.2
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/image v0.0.0-20210504121937-7319ad40d33e // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
	golang.org/x/term v0.0.0-20210503060354-a79de5458b56 // indirect
	golang.org/x/tools v0.1.2 // indirect
	google.golang.org/genproto v0.0.0-20210524171403-669157292da3 // indirect
	google.golang.org/grpc v1.38.0
	gopkg.in/ini.v1 v1.62.0 // indirect
)

replace github.com/pastelnetwork/gonode/common => ../common

replace github.com/pastelnetwork/gonode/proto => ../proto

replace github.com/pastelnetwork/gonode/pastel => ../pastel
