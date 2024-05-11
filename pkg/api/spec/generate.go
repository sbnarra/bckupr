package spec

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest --config models.yml ../../../api/specification.yml
//go:generate echo generated oapi client models
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest --config client.yml ../../../api/specification.yml
//go:generate echo generated oapi client client
