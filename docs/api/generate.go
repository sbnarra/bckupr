package api

//go:generate docker run --rm -v $PWD/../..:/local openapitools/openapi-generator-cli generate -i /local/api/specification.yml -g markdown -o /local/docs/api
