package assets

import (
	"embed"
	"net/http"
)

var (
	//go:embed *.html *.js
	assets embed.FS
	FS     = http.FS(assets)
)
