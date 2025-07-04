FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for dependency fetching
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./
COPY chopchop.png ./

# Download dependencies explicitly
RUN go mod download && go mod verify

# Copy the rest of the source code
COPY . .
COPY chopchop.png /chopchop.png

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o chopchoprss

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /app/
COPY --from=builder /app/chopchoprss .
COPY --from=builder /chopchop.png /app/chopchop.png

# Create a volume for persistent storage
VOLUME ["/data"]

# Create a wrapper script to run different commands
RUN echo '#!/bin/bash' > /app/entrypoint.sh && \
    echo 'if [ "$1" = "serve" ] || [ -z "$1" ]; then' >> /app/entrypoint.sh && \
    echo '  exec /app/chopchoprss serve "$@"' >> /app/entrypoint.sh && \
    echo 'else' >> /app/entrypoint.sh && \
    echo '  exec /app/chopchoprss "$@"' >> /app/entrypoint.sh && \
    echo 'fi' >> /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

EXPOSE 8090

ENV CHOPCHOP_CONFIG_DIR=/data

ENTRYPOINT ["/app/entrypoint.sh"]
