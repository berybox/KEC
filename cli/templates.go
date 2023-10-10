package main

import (
	_ "embed"
)

//go:embed appHelp.tmpl
var appHelpTemplate string

//go:embed cmdHelp.tmpl
var commandHelpTemplate string
