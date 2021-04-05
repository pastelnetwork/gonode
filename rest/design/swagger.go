package design

import "goa.design/goa/v3/dsl"

var _ = dsl.Service("swagger", func() {
	dsl.Description("The swagger service serves the API swagger definition.")
	dsl.HTTP(func() {
		dsl.Path("/swagger")
	})

	dsl.Files("/swagger.json", "gen/http/openapi3.json", func() {
		dsl.Description("JSON document containing the API swagger definition")
	})
})
