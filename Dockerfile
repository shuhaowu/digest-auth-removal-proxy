# --- Build stage ---
FROM golang:1.25 AS builder

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o digest-auth-removal-proxy .

# --- Final image ---
FROM ubuntu:24.04

COPY --from=builder /src/digest-auth-removal-proxy /usr/local/bin/digest-auth-removal-proxy

# Entrypoint script
COPY docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD []

