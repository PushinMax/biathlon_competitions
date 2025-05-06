BINARY_NAME=biathlon
GO=go
GOFLAGS=-mod=readonly
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html
TEST_FLAGS=-v -race


DEFAULT_EVENTS_PATH=configs/events
DEFAULT_CONFIG_PATH=configs/config.json

.PHONY: all test test-cover cover clean

all: test build


test:
	$(GO) test $(GOFLAGS) $(TEST_FLAGS) ./...

test-cover:
	$(GO) test $(GOFLAGS) $(TEST_FLAGS) -coverprofile=$(COVERAGE_FILE) ./...

cover: test-cover
	$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

cover-view:
	$(GO) tool cover -func=$(COVERAGE_FILE)

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) ./cmd/...


run:
	$(GO) run $(GOFLAGS) ./cmd/... -event=$(DEFAULT_EVENTS_PATH) -config=$(DEFAULT_CONFIG_PATH)

run-custom:
	$(GO) run $(GOFLAGS) ./cmd/... -event=$(EVENT) -config=$(CONFIG)	

clean:
	$(GO) clean
	rm -f $(BINARY_NAME) $(COVERAGE_FILE) $(COVERAGE_HTML)