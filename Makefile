.SILENT:
.PHONY:
.DEFAULT_GOAL := run-dev


run-dev:
	./scripts/run-dev.sh
stop-dev:
	./scripts/stop-dev.sh

SERVICE_MOCKS_SRC=internal/services/interface.go
SERVICE_MOCKS_DST=internal/services/mocks/mock.go
STORAGE_MOCKS_SRC=internal/storages/interface.go
STORAGE_MOCKS_DST=internal/storages/mocks/mock.go

premock:
	rm -rf ${SERVICE_MOCKS_DST} ${STORAGE_MOCKS_DST}

mocks: premock
	mockgen -source ${SERVICE_MOCKS_SRC} -destination ${SERVICE_MOCKS_DST}
	mockgen -source ${STORAGE_MOCKS_SRC} -destination ${STORAGE_MOCKS_DST}

unit-test: mocks
	go test -short -race ./...