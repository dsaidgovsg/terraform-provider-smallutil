FROM golang:1.11-alpine as builder
ARG UPX_VERSION=3.95

WORKDIR /app

# Cache the fetched Go packages
RUN apk add --no-cache gcc git musl-dev
COPY ./go.mod ./go.sum ./
RUN go mod download

# Then build the binary
COPY ./ ./

RUN set -euo pipefail && \
    go build -v -ldflags "-linkmode external -extldflags -static -s -w"; \
    wget https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-amd64_linux.tar.xz; \
    tar xvf upx-${UPX_VERSION}-amd64_linux.tar.xz; \
    mv upx-${UPX_VERSION}-amd64_linux/upx /usr/local/bin/; \
    rm -r upx-${UPX_VERSION}-amd64_linux upx-${UPX_VERSION}-amd64_linux.tar.xz; \
    upx --best terraform-provider-test; \
    :

FROM alpine:3.9 as release
WORKDIR /app
COPY --from=builder /app/terraform-provider-test ./

CMD ["/app/terraform-provider-test"]
