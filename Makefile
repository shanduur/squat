.PHONY: build
build:
	go build -o ./build/

.PHONY: test
test:
	go test -cover ./...

.PHONY: tools
tools:
	go build -o ./build/ ./tools/gob-generator/

include dev.env
export
.PHONY: run
run:
	./build/squat -p ":8081"

.PHONY: all
all: build run

.PHONY: data
data: tools
	./build/gob-generator -i bin/data/data.json -o bin/data/data.gob

.PHONY: docker
docker:
	cd Docker && docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--output=type=registry \
		--no-cache \
		--tag shanduur/squat:0.1 \
		--tag shanduur/squat:latest \
		.
