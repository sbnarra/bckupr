set -e

if [ $# -eq 0 ]; then
    echo "Missing commit message"
    exit 1
fi

./scripts/docs-snippets-generate-help.sh
./scripts/docs-build-site.sh

cd docs/site
git add .
git commit -am "$*" && git push
cd -
git add docs
git commit -m "Docs Published: $*" && git push