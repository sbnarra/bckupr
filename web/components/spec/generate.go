package spec

//go:generate rm -rf api model ApiClient.js index.js
//go:generate $PWD/../../../scripts/docker-run-for-user.sh -v $PWD/../../..:/local openapitools/openapi-generator-cli generate -i /local/api/specification.yml -g javascript -o /local/web/components/spec --additional-properties=sourceFolder=
