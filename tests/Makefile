.PHONY: run
run: sqlite3-test postgres-test mysql-test

.PHONY: run-with-coverage
run-with-coverage: sqlite3-test-with-coverage postgres-test-with-coverage mysql-test-with-coverage

.PHONY: clean
clean: sqlite3-clean postgres-clean mysql-clean

.PHONY: clobber
clobber: sqlite3-clobber postgres-clobber mysql-clobber

mysql-%:
	$(MAKE) -C mysql $*

postgres-%:
	$(MAKE) -C postgres $*

sqlite3-%:
	$(MAKE) -C sqlite3 $*
