package assets

import (
	"embed"
	"net/http"
)

//go:embed *.html *.js
var assets embed.FS

var FS = http.FS(assets)
