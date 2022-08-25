
.PHONY: build

build:
	GOOS=linux CGO_ENABLED=0 go build -o bootstrap

clean:
	rm -f bootstrap

.DEFAULT_GOAL := build
