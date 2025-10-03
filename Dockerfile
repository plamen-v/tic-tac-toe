FROM golang:1.25 AS builder
WORKDIR /server

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Download dependencies (cache this layer)
RUN go mod download

# Copy all source code
COPY src/ ./src

# Build the Go app (using Go modules)
RUN CGO_ENABLED=0 GOOS=linux go build -o tic-tac-toe ./src

FROM golang:1.25
ARG APP_PORT
WORKDIR /server

COPY --from=builder /server/tic-tac-toe .
COPY config.yaml /server/config.yaml

# Expose port (adjust if your app listens on different port)
EXPOSE $APP_PORT

# Run the binary
CMD ["./tic-tac-toe", "--config", "config.yaml"]

