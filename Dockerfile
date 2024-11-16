ARG GO_VERSION=1.23.3
FROM golang:${GO_VERSION}-alpine3.20 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd


FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /usr/local/bin/app /usr/local/bin/app

ENTRYPOINT ["/usr/local/bin/app"]
