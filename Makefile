# Customize Make
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
.RECIPEPREFIX = >
# end

build:
> echo "Building Hastie..."
> go build ./cmd/hastie
> cp hastie $(HOME)/bin/
.PHONY: build

test:
> echo "Testing..."
> go test --cover -v ./cmd/hastie
.PHONY: test

multiarch: build
> VERSION=`hastie -version | sed -e 's/hastie v//'`
> echo "Building multiarch for $(VERSION)..."
> gox -osarch="linux/amd64 darwin/amd64 windows/amd64" -output "{{.Dir}}-$(VERSION)-{{.OS}}/{{.Dir}}"
> for arch in linux darwin windows; do
>   tar cf hastie-$(VERSION)-$(arch).tar hastie-$(VERSION)-$(arch)
>   gzip hastie-$(VERSION)-$(arch).tar
>   rm -rf hastie-$(VERSION)-$(arch)
> done
.PHONY: multiarch
