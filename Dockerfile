ARG GO_VERSION=1.25.1
FROM golang:${GO_VERSION}-alpine3.21 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd


FROM gcr.io/distroless/static-debian12@sha256:4b2a093ef4649bccd586625090a3c668b254cfe180dee54f4c94f3e9bd7e381e

COPY --from=builder /usr/local/bin/app /app

ENTRYPOINT ["/app"]
