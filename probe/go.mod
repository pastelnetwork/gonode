module github.com/pastelnetwork/gonode/probe

go 1.16

require (
	github.com/aclements/go-moremath v0.0.0-20210112150236-f10218a38794
	github.com/c-bata/goptuna v0.8.1
	github.com/corona10/goimghdr v0.0.0-20190614101314-9af2afa93d77
	github.com/dgryski/go-onlinestats v0.0.0-20170612111826-1c7d19468768
	github.com/galeone/tensorflow v2.4.0-rc0.0.20210202175351-640a390c2283+incompatible
	github.com/gitchander/permutation v0.0.0-20210302120832-6ab79d7de174
	github.com/kofalt/go-memoize v0.0.0-20200917044458-9b55a8d73e1c
	github.com/kzahedi/goent v0.0.0-20190403094137-49773660fa36
	github.com/mattn/go-sqlite3 v1.14.6
	github.com/montanaflynn/stats v0.6.5
	github.com/mxschmitt/golang-combinations v1.1.0
	github.com/pa-m/sklearn v0.0.0-20200711083454-beb861ee48b1
	github.com/pastelnetwork/gonode/common v0.0.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/profile v1.5.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887 // indirect
	gonum.org/v1/gonum v0.9.1
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.12
)

replace github.com/pastelnetwork/gonode/common => ../common
