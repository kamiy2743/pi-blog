FROM golang:1.26-bookworm AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY ./ ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ./blog-mcp ./cmd/app

FROM golang:1.26-bookworm

COPY --from=docker:29.4.1-cli /usr/local/bin/docker /usr/local/bin/docker
COPY --from=docker:29.4.1-cli /usr/local/libexec/docker/cli-plugins/docker-buildx /usr/local/libexec/docker/cli-plugins/docker-buildx
COPY --from=docker:29.4.1-cli /usr/local/libexec/docker/cli-plugins/docker-compose /usr/local/libexec/docker/cli-plugins/docker-compose
COPY --from=builder /app/blog-mcp /usr/local/bin/blog-mcp

CMD ["blog-mcp"]
