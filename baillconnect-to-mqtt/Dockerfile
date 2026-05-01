ARG BUILD_FROM=ghcr.io/home-assistant/base:latest
FROM $BUILD_FROM

ARG BUILD_VERSION
ARG BUILD_ARCH

RUN apk add --no-cache ca-certificates curl \
  && mkdir -p /app

RUN set -eux; \
  test -n "${BUILD_VERSION}"; \
  test -n "${BUILD_ARCH}"; \
  case "${BUILD_ARCH}" in \
    amd64) asset_arch="amd64" ;; \
    aarch64) asset_arch="arm64" ;; \
    armv7) asset_arch="armv7" ;; \
    *) echo "unsupported architecture: ${BUILD_ARCH}" >&2; exit 1 ;; \
  esac; \
  curl -fsSL \
    "https://github.com/jonatak/baillconnect-to-mqtt/releases/download/v${BUILD_VERSION}/checksums.txt" \
    -o /tmp/checksums.txt; \
  curl -fsSL \
    "https://github.com/jonatak/baillconnect-to-mqtt/releases/download/v${BUILD_VERSION}/baillconnect-to-mqtt-linux-${asset_arch}" \
    -o /tmp/baillconnect-to-mqtt; \
  expected="$(grep " baillconnect-to-mqtt-linux-${asset_arch}$" /tmp/checksums.txt | cut -d ' ' -f 1)"; \
  echo "${expected}  /tmp/baillconnect-to-mqtt" | sha256sum -c -; \
  mv /tmp/baillconnect-to-mqtt /app/baillconnect-to-mqtt; \
  chmod +x /app/baillconnect-to-mqtt

CMD ["/app/baillconnect-to-mqtt", "--config", "/data/options.json"]
