module github.com/pastelnetwork/walletnode

go 1.16

require (
	github.com/dimfeld/httptreemux/v5 v5.3.0 // indirect
	github.com/getkin/kin-openapi v0.53.0 // indirect
	github.com/pastelnetwork/go-commons v0.0.0-00010101000000-000000000000
	github.com/pastelnetwork/go-pastel v0.0.0-00010101000000-000000000000
	github.com/rs/cors v1.7.0
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	goa.design/goa/v3 v3.3.1
	golang.org/x/net v0.0.0-20210331212208-0fccb6fa2b5c // indirect
	golang.org/x/sys v0.0.0-20210403161142-5e06dd20ab57 // indirect
	golang.org/x/text v0.3.6 // indirect
)

replace (
	github.com/pastelnetwork/go-commons => ../go-commons
	github.com/pastelnetwork/go-pastel => ../go-pastel
)
