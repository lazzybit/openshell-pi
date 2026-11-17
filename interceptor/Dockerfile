FROM golang:1.24.4-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/openshell-pi-interceptor ./cmd/openshell-pi-interceptor

FROM debian:trixie-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --system --uid 10001 --home-dir /nonexistent --shell /usr/sbin/nologin interceptor
COPY --from=build /out/openshell-pi-interceptor /usr/local/bin/openshell-pi-interceptor
USER 10001:10001
EXPOSE 50051
ENTRYPOINT ["/usr/local/bin/openshell-pi-interceptor"]
