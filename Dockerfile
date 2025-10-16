# build stage
FROM golang:1.23.2-alpine3.19 AS builder

WORKDIR /app

ARG GIT_VERSION

COPY ./app .
RUN go mod download
RUN current_time=$(date +"%Y-%m-%dT%H:%M:%SZ") && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -X main.buildTime=$current_time -X main.version=$GIT_VERSION" -o /out/main ./

# image stage
FROM alpine:3.19

RUN apk --no-cache add curl

COPY --from=builder /out/main /app/main
COPY entry-point.sh /app/entry-point.sh
RUN chmod +x /app/entry-point.sh

WORKDIR /app
ENTRYPOINT ["/app/entry-point.sh"]