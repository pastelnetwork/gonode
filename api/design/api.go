package design

import . "goa.design/goa/v3/dsl"

// API describes the global properties of the API server.
var _ = API("walletnode", func() {
	Title("WalletNode REST API")
	Version("1.0")

	Server("walletnode", func() {
		Services("arts", "swagger")

		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})
