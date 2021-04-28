package gen

import "embed"

// OpenAPIContent is a filesystem that contains generated OpenAPI documents.
//go:embed http/openapi3.json
var OpenAPIContent embed.FS
