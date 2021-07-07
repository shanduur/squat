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
	./build/squat

.PHONY: all
all: build run

.PHONY: docker
docker:
	cd Docker && docker buildx build --tag squat:0.1 .
