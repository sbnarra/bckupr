package api

// go:generate rm -rf Apis Models
//go:generate $PWD/../../scripts/docker-run-for-user.sh -v $PWD/../..:/local openapitools/openapi-generator-cli generate -i /local/api/specification.yml -g markdown -o /local/docs/api
