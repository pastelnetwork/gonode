package docs

import "embed"

// SwaggerContent is a filesystem that contains swagger html content.
//go:embed swagger
var SwaggerContent embed.FS
