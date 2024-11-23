FROM golang:1.22.5-alpine AS builder

# Install required system packages and update certificates
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates

# Add Maintainer Info to the Image
LABEL maintainer="Mehrdad Amini <pcmehrdad@gmail.com>"
LABEL description="solana rpc proxy"

# Set the Current Working Directory inside the container
# We'll use a path that maintains the same relative structure
WORKDIR /build/gosol

# Copy the gosol module
COPY ./gosol/ .

# Copy the handler_socket2 module to maintain the relative path structure
COPY ./handler_socket2/ ../handler_socket2/

# since ../handler_socket2 exists relative to the gosol directory
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/solproxy ./main/main.go

# Start a new stage from busybox
FROM scratch
#FROM alpine:3.15

WORKDIR /app

# Copy binary and certs
COPY --from=builder /app/solproxy .
COPY --from=builder /build/gosol/main/server-status.html .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose ports
EXPOSE 7777
EXPOSE 7778

# Command to run the executable
CMD ["./solproxy"]