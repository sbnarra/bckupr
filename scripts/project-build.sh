set -e

echo "Building/Testing App"
./scripts/app-build-image.sh
./scripts/app-test-end2end.sh

echo "Generating Doc Snippets"
./scripts/docs-snippets-generate-help.sh
echo "Building mkdocs Image"
./scripts/docs-build-mkdocs-image.sh

echo "Generating HTML Docs"
git submodule init
git submodule update
./scripts/docs-build-site.sh
