set -e

# https://openapi-generator.tech/docs/generators/go-gin-server/

rm -rf internal/openapi/spec
docker run --rm -u $(id -u):$(id -g)  \
  -v ./:/bckupr openapitools/openapi-generator-cli generate \
  -i /bckupr/api/openapi.yml \
  -g go-gin-server \
  --additional-properties=interfaceOnly=true,apiPath=internal/openapi/spec,packageName=spec \
  -o /bckupr


# https://openapi-generator.tech/docs/generators/go/

rm -rf pkg/client
mkdir -p pkg/client
cp .openapi-generator-ignore pkg/client/.openapi-generator-ignore
docker run --rm -u $(id -u):$(id -g)  \
  -v ./:/bckupr openapitools/openapi-generator-cli generate \
  -i /bckupr/api/openapi.yml \
  -g go \
  --additional-properties=isGoSubmodule=true,withGoMod=false,packageName=client,generateInterfaces=true \
  -o /bckupr/pkg/client