FROM golang:1.18-alpine3.14
WORKDIR /app
RUN go install github.com/githubnemo/CompileDaemon@latest && \
    mkdir -p /etc/backend && \
    touch /etc/backend/config.yml
COPY . .
ENTRYPOINT CompileDaemon -polling -build="go build -o /build/backend cmd/main.go" -command="/build/backend --config-path=/etc/backend/config.yml"
