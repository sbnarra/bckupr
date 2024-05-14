package api

//go:generate sudo rm -rf Apis Models
//go:generate docker run --rm -v $PWD/../..:/local openapitools/openapi-generator-cli generate -i /local/api/specification.yml -g markdown -o /local/docs/api
