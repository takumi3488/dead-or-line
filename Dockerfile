ARG GO_VERSION=1.25.1
FROM golang:${GO_VERSION}-alpine3.21 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd


FROM gcr.io/distroless/static-debian12@sha256:cd64bec9cec257044ce3a8dd3620cf83b387920100332f2b041f19c4d2febf93

COPY --from=builder /usr/local/bin/app /app

ENTRYPOINT ["/app"]
