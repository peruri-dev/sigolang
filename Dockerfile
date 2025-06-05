# Use the official Golang image as the builder
FROM golang:1.22 AS builder

# Set the working directory inside the builder
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Ensure dependencies are correct by tidying up go modules
RUN go mod tidy
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application (replace ./cmd/main.go with your main application path if different)
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./main.go

# Use a minimal base image to keep the final image small
FROM alpine:latest

# Set environment variables (if needed)
ENV PORT=3000

# Copy the binary from the builder stage to the final stage
COPY --from=builder /api /api

# Expose the application port
EXPOSE 3000

# Command to run the application
CMD ["/api"]