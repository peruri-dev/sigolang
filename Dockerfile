# Use the official Golang image as the builder
FROM golang:1.24-alpine AS builder

# Update Certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Set the working directory inside the builder
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

## Declare arguments for private dependency
#ARG GIT_USERNAME
## Set the Git configuration with environment variables
#RUN --mount=type=secret,id=GIT_TOKEN \
##    git config --global url."https://${GIT_USERNAME}:$(head -n 1 /run/secrets/GIT_TOKEN)@github.com".insteadOf "https://github.com"

# Copy the rest of the application code
COPY . .

# Declare arguments for build version
#ARG GIT_COMMIT
#RUN mkdir -p public
#RUN echo "$GIT_COMMIT" > public/version.txt

# Build the application (replace ./cmd/main.go with your main application path if different)
RUN --mount=type=cache,target=/go/pkg/mod CGO_ENABLED=0 GOOS=linux go build -o api ./main.go

# Use a minimal base image to keep the final image small
FROM alpine:latest

# Add dumb-init for entry point wrapper
RUN apk add dumb-init

# Set environment variables (if needed)
ENV PORT=3000

# Copy the binary from the builder stage to the final stage
COPY --from=builder /app/api /api
RUN chmod +x /api

COPY --from=builder /app/public /public

# Import system files from the builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# Expose the application port
EXPOSE 3000

# Command to run the application
CMD ["dumb-init", "/docker-entrypoint.sh"]