package static

import (
	"embed"
)

//go:embed index.html static
var content embed.FS

func StaticFS() embed.FS {
	return content
}
