FROM golang:1.26-bookworm AS builder

WORKDIR /app

COPY ./ ./
RUN go build -o ./blog-back ./cmd/

FROM  debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/blog-back /app/blog-back
COPY ./templates/ /app/templates/

RUN groupadd --gid 1000 go
RUN useradd --uid 1000 --gid 1000 --no-create-home --shell /usr/sbin/nologin go
USER go:go

CMD ["/app/blog-back"]
