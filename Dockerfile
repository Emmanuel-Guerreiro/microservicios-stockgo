# Stage 1: Build the Go binary
FROM golang:1.22.6-bullseye AS builder

# Set the working directory
WORKDIR /go/src/stockgo

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stockgo .

# Stage 2: Run the binary
FROM alpine:latest

# Set up a non-root user (optional but recommended for security)
RUN adduser -D appuser

# ENV MONGO_URL=mongodb://host.docker.internal:27017
# ENV RABBIT_URL=amqp://host.docker.internal
# ENV FLUENT_URL=host.docker.internal:24224

# Copy the binary from the builder stage
COPY --from=builder /go/src/stockgo/stockgo /usr/local/bin/stockgo

# Run as the non-root user
USER appuser

# Expose the application's port (adjust according to your app)
EXPOSE 3000

# Set the default command to execute the binary
CMD ["stockgo"]
