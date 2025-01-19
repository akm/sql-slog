GO_TEST_OPTIONS?=

.PHONY: test
test:
	go test $(GO_TEST_OPTIONS) ./...

GO_COVERAGE_DIR=coverage/unit
$(GO_COVERAGE_DIR):
	mkdir -p $(GO_COVERAGE_DIR)
GO_COVERAGE_HTML?=coverage.html
GO_COVERAGE_PROFILE?=coverage.txt
$(GO_COVERAGE_PROFILE):
	$(MAKE) test-with-coverage

# See https://app.codecov.io/github/akm/go-requestid/new
.PHONY: test-with-coverage
test-with-coverage: $(GO_COVERAGE_DIR)
	go test -cover -coverpkg=github.com/akm/sql-slog ./... -args -test.gocoverdir="$(GO_COVERAGE_DIR)"

.PHONY: test-coverage
test-coverage: $(GO_COVERAGE_PROFILE)
	go tool covdata percent -i=$(GO_COVERAGE_DIR) -o $(GO_COVERAGE_PROFILE)
	go tool cover -html=$(GO_COVERAGE_PROFILE) -o $(GO_COVERAGE_HTML)
	@command -v open && open $(GO_COVERAGE_HTML) || echo "open $(GO_COVERAGE_HTML)"

.PHONY: demo
demo:
	go run ./demo.go ${LOG_LEVEL} ${LOG_FORMAT}

DEMO_RESULTS_DIR=demo-results
$(DEMO_RESULTS_DIR):
	mkdir -p $(DEMO_RESULTS_DIR)

demo-log-text-%: $(DEMO_RESULTS_DIR)
	LOG_LEVEL=$* LOG_FORMAT=text $(MAKE) demo > $(DEMO_RESULTS_DIR)/$*-log.txt

.PHONY: demo-logs-gen
demo-logs-gen: demo-log-text-info demo-log-text-debug demo-log-text-trace demo-log-text-verbose

.PHONY: clean
clean:
	rm -rf coverage
	rm -f $(GO_COVERAGE_HTML) $(GO_COVERAGE_PROFILE)

.PHONY: clobber
clobber: clean
	rm -rf $(DEMO_RESULTS_DIR)
