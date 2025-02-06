,PHONY: logs-gen
logs-gen: logs-sqlite3-gen logs-postgres-gen logs-mysql-gen

.PHONY: clean
clean: logs-sqlite3-clean logs-postgres-clean logs-mysql-clean

.PHONY: clobber
clobber: clean logs-sqlite3-clobber logs-postgres-clobber logs-mysql-clobber

logs-sqlite3-%:
	$(MAKE) -C logs-sqlite3 $*

logs-postgres-%:
	$(MAKE) -C logs-postgres $*

logs-mysql-%:
	$(MAKE) -C logs-mysql $*
