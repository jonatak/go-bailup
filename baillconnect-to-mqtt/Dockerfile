ARG BUILD_FROM=ghcr.io/home-assistant/base:latest
FROM $BUILD_FROM

ARG BUILD_VERSION
ARG BUILD_ARCH

RUN apk add --no-cache bash ca-certificates curl \
  && mkdir -p /app

RUN set -eux; \
  test -n "${BUILD_VERSION}"; \
  test -n "${BUILD_ARCH}"; \
  script_ref="${BUILD_VERSION}"; \
  case "${script_ref}" in v*) ;; *) script_ref="v${script_ref}" ;; esac; \
  curl -fsSL \
    "https://raw.githubusercontent.com/jonatak/baillconnect-to-mqtt/${script_ref}/scripts/install-latest-release.sh" \
    -o /tmp/install-latest-release.sh; \
  XDG_BIN_HOME=/app BUILD_VERSION="${BUILD_VERSION}" BUILD_ARCH="${BUILD_ARCH}" \
    bash /tmp/install-latest-release.sh

CMD ["/app/baillconnect-to-mqtt", "--config", "/data/options.json"]
