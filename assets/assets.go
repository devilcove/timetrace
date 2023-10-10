package assets

import _ "embed"

var (
	//go:embed logo.png
	Logo []byte
	//go:embed logo-small.png
	SmallLogo []byte
)
