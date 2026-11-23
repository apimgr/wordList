FROM casjaysdev/go:latest AS builder

ARG TARGETARCH
ARG VERSION=devel
ARG COMMIT_ID=N/A
ARG BUILD_EPOCH=0
ARG OFFICIAL_SITE=

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p /out && CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} GOFLAGS=-buildvcs=false go build -trimpath -ldflags "-s -w -X main.Version=${VERSION} -X main.CommitID=${COMMIT_ID} -X main.BuildEpoch=${BUILD_EPOCH} -X main.OfficialSite=${OFFICIAL_SITE}" -o /out/wordList ./src

FROM alpine:3.22
RUN apk add --no-cache bash ca-certificates curl su-exec tini && addgroup -S wordlist && adduser -S -G wordlist wordlist && mkdir -p /config/wordList /data/wordList && chown -R wordlist:wordlist /config /data
COPY --from=builder /out/wordList /usr/local/bin/wordList
COPY docker/rootfs/ /
RUN chmod 755 /usr/local/bin/wordList /usr/local/bin/entrypoint.sh
ENV TZ=UTC
EXPOSE 80
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 CMD curl --fail --silent http://127.0.0.1:80/api/v1/server/healthz || exit 1
USER root
STOPSIGNAL SIGTERM
ENTRYPOINT ["tini", "--", "/usr/local/bin/entrypoint.sh"]
