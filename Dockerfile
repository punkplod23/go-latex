# Build stage
FROM ubuntu:24.04 AS builder

ARG GO_VERSION=1.24.1

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
    && curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tar.gz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/main .

# Runtime stage
FROM ubuntu:24.04

# Install only the runtime dependencies needed to compile TikZ/PGFPlots documents
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        fontconfig \
        texlive \
        texlive-latex-extra \
        texlive-pictures \
        fonts-dejavu \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main ./main

EXPOSE 8080

CMD ["/app/main"]