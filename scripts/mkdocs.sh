docker run --rm \
    -v ./:/bckupr -w /bckupr \
    -p 8000:8000 scripts/mkdocs $@ --config-file docs/mkdocs.yml