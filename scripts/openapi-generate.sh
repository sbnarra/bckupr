# https://openapi-generator.tech/docs/generators/go-gin-server/

rm -rf internal/openapi/generated
docker run --rm -u $(id -u):$(id -g)  \
  -v ./:/bckupr openapitools/openapi-generator-cli generate \
  -i /bckupr/api/openapi.yml \
  -g go-gin-server \
  --additional-properties=interfaceOnly=true,apiPath=internal/openapi/spec,packageName=spec \
  -o /bckupr


# https://openapi-generator.tech/docs/generators/go/

rm -rf pkg/api2
docker run --rm \
  -v ./:/bckupr openapitools/openapi-generator-cli generate \
  -i /bckupr/api/openapi.yml \
  -g go \
  --additional-properties=isGoSubmodule=true,withGoMod=false,packageName=client \
  -o /bckupr/pkg/api2