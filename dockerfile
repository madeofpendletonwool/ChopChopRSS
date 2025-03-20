FROM golang:1.18-alpine AS builder

WORKDIR /app

# Install git for dependency fetching
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies explicitly
RUN go mod download && go mod verify

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o chopchoprss

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/chopchoprss .

EXPOSE 8090

ENTRYPOINT ["./chopchoprss", "serve"]
