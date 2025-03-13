FROM golang:latest AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
# Don't copy sensitive files during build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api_server ./cmd/server

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/api_server .
# Create a directory for configuration files
RUN mkdir -p /envs
# We'll mount the key file at runtime instead of baking it into the image
# Use an environment variable to specify the key file location
ENV GOOGLE_CONFIG_PATH=/envs/key.json
CMD ["/api_server"]