
generate:
	go generate ./...

doc:
	bash ./scripts/mkdocs-build.sh
	bash ./scripts/docs-build-site.sh

clean:
	go clean ./...
	sudo rm -rf bckupr docs/site

check: generate docs
	# go test ./...
	go build