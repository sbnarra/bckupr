# docker run --rm \
#   -v ./:/bckupr openapitools/openapi-generator-cli generate \
#   -i /bckupr/api/openapi.yml \
#   -g go \
#   -o /bckupr/out/client

docker run --rm -u $(id -u):$(id -g)  \
  -v ./:/bckupr openapitools/openapi-generator-cli generate \
  -i /bckupr/api/openapi.yml \
  -g go-gin-server \
  --additional-properties=interfaceOnly=true,apiPath=internal/openapi \
  -o /bckupr