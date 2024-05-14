package spec

//go:generate sudo rm -rf api model ApiClient.js index.js
//go:generate docker run --rm -v $PWD/../../..:/local openapitools/openapi-generator-cli generate -i /local/api/specification.yml -g javascript -o /local/web/app/spec --additional-properties=sourceFolder=
