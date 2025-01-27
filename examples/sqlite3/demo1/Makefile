.PHONY: run
run:
	go run . ${LOG_LEVEL} ${LOG_FORMAT}

RESULTS_DIR=results
$(RESULTS_DIR):
	mkdir -p $(RESULTS_DIR)

gen-text-logs-%: $(RESULTS_DIR)
	LOG_LEVEL=$* LOG_FORMAT=text $(MAKE) run > $(RESULTS_DIR)/$*-log.txt

.PHONY: gen
gen: gen-text-logs-info gen-text-logs-debug gen-text-logs-trace gen-text-logs-verbose

.PHONY: clean
clean:

.PHONY: clobber
clobber: clean
	rm -rf $(RESULTS_DIR)
