.PHONY: build
build:
	go build -o ./build/

.PHONY: test
test:
	go test -cover ./...

.PHONY: run
run:
	./build/squat

.PHONY: all
all: build run
