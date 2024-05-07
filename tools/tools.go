//go:build tools
// +build tools

package main

import (
	_ "github.com/contiamo/openapi-generator-go"

	_ "github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/go-swagger/go-swagger/cmd/swagger@latest"
)
