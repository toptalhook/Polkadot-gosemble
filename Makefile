SHELL := /bin/bash
CURRENT_DIR = $(shell pwd)
SRC_DIR = /src/examples/wasm/gosemble
BUILD_PATH = build/runtime.wasm
IMAGE = polkawasm/tinygo
TAG = 0.25.0
BRANCH_CONSERVATIVE_GC = new-polkawasm-target-release-$(TAG)
BRANCH_EXTALLOC_GC = new-polkawasm-target-extallocleak-gc-release-$(TAG)

.PHONY: build
build:
	@if [[ "$(GC)" == "conservative" ]]; then \
		cd tinygo; \
		git checkout $(BRANCH_CONSERVATIVE_GC); \
		cd ..; \
		docker build --tag $(IMAGE):$(TAG) -f tinygo/Dockerfile.polkawasm tinygo; \
		docker run --rm -v $(CURRENT_DIR):$(SRC_DIR) -w $(SRC_DIR) $(IMAGE):$(TAG) /bin/bash -c "tinygo build -target=polkawasm -o=$(SRC_DIR)/$(BUILD_PATH) $(SRC_DIR)/runtime/"; \
		echo "compiled with conservative GC"; \
	else \
		cd tinygo; \
		git checkout $(BRANCH_EXTALLOC_GC); \
		cd ..; \
		docker build --tag $(IMAGE):$(TAG)-extallocleak -f tinygo/Dockerfile.polkawasm tinygo; \
		docker run --rm -v $(CURRENT_DIR):$(SRC_DIR) -w $(SRC_DIR) $(IMAGE):$(TAG)-extallocleak /bin/bash -c "tinygo build -target=polkawasm -o=$(SRC_DIR)/$(BUILD_PATH) $(SRC_DIR)/runtime/"; \
		echo "compiled with extallocleak GC"; \
	fi

start-network:
	cp build/runtime.wasm substrate/bin/node-template/runtime.wasm; \
	cd substrate/bin/node-template; \
	cargo build --release; \
	cd ../..; \
	./target/release/node-template --dev --execution Wasm

test: test_unit test_integration

# TODO: ignore the integration tests
test_unit:
	@go test --tags "nonwasmenv" -v `go list ./... | grep -v runtime`

test_integration:
	@go test --tags="nonwasmenv" -v ./runtime/... -timeout 1200s
