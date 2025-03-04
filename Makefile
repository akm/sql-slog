.PHONY: default
default: build lint test

.PHONY: build
build:
	go build ./...

GOLANG_TOOL_PATH_TO_BIN=$(shell go env GOPATH)
GOLANGCI_LINT_CLI_VERSION?=latest
GOLANGCI_LINT_CLI_MODULE=github.com/golangci/golangci-lint/cmd/golangci-lint
GOLANGCI_LINT_CLI=$(GOLANG_TOOL_PATH_TO_BIN)/bin/golangci-lint
$(GOLANGCI_LINT_CLI):
	$(MAKE) golangci-lint-cli-install
golangci-lint-cli-install:
	go install $(GOLANGCI_LINT_CLI_MODULE)@$(GOLANGCI_LINT_CLI_VERSION)

.PHONY: lint
lint: $(GOLANGCI_LINT_CLI)
	golangci-lint run


GODOC_CLI_VERSION=latest
GODOC_CLI_MODULE=golang.org/x/tools/cmd/godoc
GODOC_CLI=$(GOLANG_TOOL_PATH_TO_BIN)/bin/godoc
$(GODOC_CLI):
	$(MAKE) godoc-cli-install
godoc-cli-install:
	go install $(GODOC_CLI_MODULE)@$(GODOC_CLI_VERSION)

.PHONY: godoc
godoc: $(GODOC_CLI)
	@echo "Open http://localhost:6060/pkg/github.com/akm/sql-slog"
	godoc -http=:6060

# examples-logs-gen
examples-%:
	$(MAKE) -C examples $*

tests-%:
	$(MAKE) -C tests $*

GO_TEST_OPTIONS?=

.PHONY: test
test: test-unit tests-run

.PHONY: test-unit
test-unit:
	go test $(GO_TEST_OPTIONS) ./...

GO_COVERAGE_DIR=coverage/unit
$(GO_COVERAGE_DIR):
	mkdir -p $(GO_COVERAGE_DIR)

GO_COVERAGE_MERGED_DIR=coverage/merged
$(GO_COVERAGE_MERGED_DIR):
	mkdir -p $(GO_COVERAGE_MERGED_DIR)

GO_COVERAGE_HTML?=coverage.html
GO_COVERAGE_PROFILE?=coverage.txt
$(GO_COVERAGE_PROFILE):
	$(MAKE) test-coverage-profile

test-with-coverage: test-with-coverage-unit tests-run-with-coverage 

# See https://app.codecov.io/github/akm/go-requestid/new
.PHONY: test-with-coverage-unit
test-with-coverage-unit: $(GO_COVERAGE_DIR)
	go test -cover ./... -args -test.gocoverdir="$(GO_COVERAGE_DIR)"

.PHONY: test-coverage-profile
test-coverage-profile: $(GO_COVERAGE_DIR) $(GO_COVERAGE_MERGED_DIR)
	go tool covdata merge \
		-i $(GO_COVERAGE_DIR),tests/mysql/coverage/unit,tests/postgres/coverage/unit,tests/sqlite3/coverage/unit \
		-o $(GO_COVERAGE_MERGED_DIR)
	go tool covdata percent -i=$(GO_COVERAGE_MERGED_DIR) -o $(GO_COVERAGE_PROFILE)

.PHONY: test-coverage
test-coverage: test-coverage-profile
	go tool cover -html=$(GO_COVERAGE_PROFILE) -o $(GO_COVERAGE_HTML)
	@command -v open && open $(GO_COVERAGE_HTML) || echo "open $(GO_COVERAGE_HTML)"

METADATA_YAML=.project.yaml
$(METADATA_YAML): metadata-gen

METADATA_LINTERS=$(shell cat .golangci.yml | yq '... comments="" | .linters.enable | length')
.PHONY: metadata-gen
metadata-gen: 
	@echo "linters: $(METADATA_LINTERS)" > $(METADATA_YAML)

.PHONY: clean
clean: tests-clean examples-clean
	rm -rf coverage
	rm -f $(GO_COVERAGE_HTML) $(GO_COVERAGE_PROFILE)

.PHONY: clobber
clobber: tests-clobber examples-clobber clean
	rm -f $(METADATA_YAML)
