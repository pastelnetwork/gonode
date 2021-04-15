package design

import (
	//revive:disable:dot-imports
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports
)

// API describes the global properties of the API server.
var _ = API("walletnode", func() {
	Title("WalletNode RESTFull API")
	Version("1.0")

	Server("walletnode", func() {
		Services("artworks", "swagger")

		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})
