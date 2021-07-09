package design

import (
	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports
)

// API describes the global properties of the API server.
var _ = API("walletnode", func() {
	Title("WalletNode REST API")
	Version("1.0")

	Server("walletnode", func() {
		Services("artworks", "swagger")
		Services("userdatas", "swagger")

		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})
