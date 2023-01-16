.SILENT:
.PHONY:
.DEFAULT_GOAL := run-dev


run-dev:
	./scripts/run-dev.sh
stop-dev:
	./scripts/stop-dev.sh

SERVICE_MOCKS_SRC=internal/services/interface.go
SERVICE_MOCKS_DST=internal/services/mocks/mock.go

premock:
	rm -rf ${SERVICE_MOCKS_DST}

mocks: premock
	mockgen -source ${SERVICE_MOCKS_SRC} -destination ${SERVICE_MOCKS_DST}

unit-test:
	go test -short -race ./...