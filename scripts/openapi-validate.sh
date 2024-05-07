docker run --rm -u $(id -u):$(id -g)  \
  -v ./:/bckupr openapitools/openapi-generator-cli validate \
  -i /bckupr/api/openapi-system.yml