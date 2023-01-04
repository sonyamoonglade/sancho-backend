.SILENT:
.PHONY:
.DEFAULT_GOAL := run-dev


run-dev:
	./scripts/run-dev.sh
stop-dev:
	./scripts/stop-dev.sh