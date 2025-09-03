# Use the official Golang image for building
FROM golang:1.25 AS builder
# Set working directory

WORKDIR /app
# Copy Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy source code
COPY . .
# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server .
# Use a minimal base image for final deployment
FROM alpine:latest
# Set working directory in the container
WORKDIR /app
RUN adduser -D gouser && chown -R gouser /app

# Copy the built binary from the builder stage
COPY --from=builder /app/server /app/server
# Expose the application port
COPY .env /app
RUN chown gouser:gouser /app/server

EXPOSE 80

# Run the application
USER gouser
CMD ["/app/server"]