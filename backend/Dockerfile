FROM golang:1.26-bookworm AS builder

WORKDIR /app

COPY ./ ./
RUN go build -o ./blog-back ./cmd/blog/

FROM  debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/blog-back /app/blog-back
COPY ./templates/ /app/templates/

CMD ["/app/blog-back"]
