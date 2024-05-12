# docker build -t scripts/mkdocs -f scripts/mkdocs.Dockerfile .

docker run --rm \
    -v ./docs/gh-pages:/bckupr/docs -w /bckupr \
    -p 8000:8000 scripts/mkdocs $@ --config-file docs/mkdocs.yml