# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Runtime stage
FROM alpine:3.22

RUN apk add --no-cache \
        ca-certificates \
        fontconfig \
        font-dejavu \
        texlive-most \
        texlive-xetex \
    && addgroup -S app \
    && adduser -S -G app app \
    && mkdir -p /tmp/go-latex \
    && chmod 1777 /tmp/go-latex

ENV TMPDIR=/tmp/go-latex

WORKDIR /app
COPY --from=builder /app/main ./main
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod 755 /usr/local/bin/docker-entrypoint.sh \
    && chown app:app /app/main

USER app

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]

EXPOSE 8080

CMD ["/app/main"]