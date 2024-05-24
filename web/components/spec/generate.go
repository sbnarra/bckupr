package spec

//go:generate rm -rf api model ApiClient.js index.js
//go:generate $PWD/../../../scripts/docker-run-for-user.sh -v $PWD/../../..:/local openapitools/openapi-generator-cli:v7.5.0 generate -i /local/api/specification.yml -g typescript-fetch -o /local/web/components/spec --additional-properties=usePromises=true,sourceFolder=
