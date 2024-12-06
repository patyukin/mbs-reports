FROM golang:1.23.2-alpine3.20 AS builder

ENV config=docker

WORKDIR /app

COPY . /app

ENV YAML_CONFIG_FILE_PATH=config.yaml

RUN go mod tidy && \
    go mod download && \
    go get github.com/githubnemo/CompileDaemon && \
    go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -o ./go_build_report cmd/report/main.go" --command=./go_build_report
