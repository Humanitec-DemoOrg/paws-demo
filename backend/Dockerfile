FROM golang:1.18-alpine AS builder

WORKDIR /app

# Copy just the files needed for download modules to take advantage of caching in Docker for local development
# go.sum is used for cache invalidation
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

COPY . .

# https://stackoverflow.com/questions/36279253/go-compiled-binary-wont-run-in-an-alpine-docker-container-on-ubuntu-host
ENV CGO_ENABLED=0

RUN GOOS=linux GOARCH=amd64 go build -o /opt/caws-demo/backend ./cmd/main.go

FROM alpine AS final

WORKDIR /bin
COPY --from=builder /opt/caws-demo /opt/caws-demo

WORKDIR /opt/caws-demo
ENTRYPOINT ["/opt/caws-demo/backend"]
