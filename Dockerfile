# syntax=docker/dockerfile:1

# ── Stage 1: Build ───────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

ARG VERSION=dev

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /src
COPY go.mod ./
COPY go.sum* ./
RUN go mod download 2>/dev/null || true
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
      -ldflags="-s -w -X main.version=${VERSION}" \
      -trimpath \
      -o /coldkey \
      ./cmd/coldkey

# ── Stage 2: Runtime ────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12:nonroot

LABEL org.opencontainers.image.source="https://github.com/pike00/coldkey"

COPY --from=builder /coldkey /coldkey
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

USER nonroot:nonroot
WORKDIR /out

ENTRYPOINT ["/coldkey"]
