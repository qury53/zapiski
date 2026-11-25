FROM golang:1.26.9-bookworm AS builder

WORKDIR /src
ENV CGO_ENABLED=0 GOOS=linux

COPY go.work go.work.sum ./
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod   shared/go.sum   ./shared/
COPY sso/go.mod      sso/go.sum      ./sso/
COPY agent/go.mod    agent/go.sum    ./agent/
COPY api/go.mod      api/go.sum      ./api/

RUN --mount=type=cache,target=/go/pkg/mod \
    for m in platform shared sso agent api; do (cd "$m" && go mod download) || exit 1; done

COPY platform ./platform
COPY shared ./shared
COPY agent ./agent

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/agent ./agent/cmd

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/agent /app/agent

EXPOSE 5678
USER nonroot:nonroot
ENTRYPOINT ["/app/agent"]