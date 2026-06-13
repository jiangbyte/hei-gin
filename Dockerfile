FROM golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git tzdata

COPY go.work go.work.sum ./
COPY go.mod go.sum ./
COPY sdk/go.mod sdk/go.sum ./sdk/
COPY plugins/plugin-sys/go.mod plugins/plugin-sys/go.sum ./plugins/plugin-sys/
COPY plugins/plugin-client/go.mod plugins/plugin-client/go.sum ./plugins/plugin-client/
COPY plugins/plugin-im/go.mod plugins/plugin-im/go.sum ./plugins/plugin-im/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/hei-gin \
    ./main.go

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app \
    && mkdir -p /app/uploads \
    && chown -R app:app /app

COPY --from=builder /out/hei-gin /app/hei-gin
COPY config.example.yaml /app/config.example.yaml

ENV HEI_CONFIG=/app/config.yaml

EXPOSE 18885

USER app

ENTRYPOINT ["/app/hei-gin"]
