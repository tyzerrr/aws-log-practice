# syntax=docker/dockerfile:1
# check=error=true
FROM golang:1.26.4 AS builder
WORKDIR /batch
RUN --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download -x

RUN --mount=type=bind,target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o /artifacts/batch ./

FROM gcr.io/distroless/static-debian13

COPY --from=builder /artifacts/batch /batch

CMD ["/batch"]
