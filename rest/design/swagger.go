package design

import (
	//revive:disable:dot-imports
	. "goa.design/goa/v3/dsl"
	//revive:enable:dot-imports
)

var _ = Service("swagger", func() {
	Description("The swagger service serves the API swagger definition.")
	HTTP(func() {
		Path("/swagger")
	})

	Files("/swagger.json", "gen/http/openapi3.json", func() {
		Description("JSON document containing the API swagger definition")
	})
})
