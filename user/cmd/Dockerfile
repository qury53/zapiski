FROM golang:1.26.9-bookworm AS builder

WORKDIR /src
ENV CGO_ENABLED=0 GOOS=linux

COPY go.work go.work.sum ./
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod   shared/go.sum   ./shared/
COPY sso/go.mod      sso/go.sum      ./sso/
COPY agent/go.mod    agent/go.sum    ./agent/
COPY api/go.mod      api/go.sum      ./api/
COPY user/go.mod                     ./user/

RUN --mount=type=cache,target=/go/pkg/mod \
    for m in platform shared sso agent api user; do (cd "$m" && go mod download) || exit 1; done

COPY platform ./platform
COPY shared   ./shared
COPY user     ./user

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/user ./user/cmd

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/user /app/user

EXPOSE 50053 9093
USER nonroot:nonroot
ENTRYPOINT ["/app/user"]
