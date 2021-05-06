package design

import (
	//revive:disable:dot-imports
	//lint:ignore ST1001 disable warning dot import
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports
	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = Service("swagger", func() {
	Description("The swagger service serves the API swagger definition.")
	Meta("swagger:generate", "false")

	cors.Origin("localhost")
	HTTP(func() {
		Path("/swagger")
	})

	Files("/swagger.json", "http/openapi3.json", func() {
		Description("JSON document containing the API swagger definition")
	})
})
