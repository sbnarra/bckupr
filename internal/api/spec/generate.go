package spec

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest --config embedded.yml ../../../api/specification.yml
//go:generate echo generated oapi server embedded
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest --config models.yml ../../../api/specification.yml
//go:generate echo generated oapi server models
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest --config spec.yml ../../../api/specification.yml
//go:generate echo generated oapi server spec
