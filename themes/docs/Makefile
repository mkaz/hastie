# Customize Make
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
.RECIPEPREFIX = >
# end

build:
> echo "Building scss..."
> sassc -t compact scss/main.scss static/style.css
.PHONY: build

