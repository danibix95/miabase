
DEFAULT_BUILD_FLAGS = -ldflags="-w -s"
DEFAULT_TEST_FLAGS = -timeout 30s -failfast
ADVANCED_BUILD_FLAGS = -trimpath -mod=readonly

.PHONY: test clean

all: clean build

build: compile test

compile:
	@mkdir -p bin
	@go build ${DEFAULT_BUILD_FLAGS} ${ADVANCED_BUILD_FLAGS} -o bin ./...

test:
	@go test ${DEFAULT_TEST_FLAGS} -race ./...

cover:
	@go test ${DEFAULT_TEST_FLAGS} -cover -coverprofile=coverage.out ./...

bench:
	@go test ${DEFAULT_TEST_FLAGS} -bench=. -benchmem ./...

show-coverage: cover
	@go tool cover -html=coverage.out

clean:
	@rm -rf bin/
