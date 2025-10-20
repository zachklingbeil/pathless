# Build stage
FROM --platform=$BUILDPLATFORM golang:latest AS go_builder
WORKDIR /app

# Cache dependencies
RUN --mount=type=cache,target=/go/pkg/mod/ \
    # --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

# Build static binary
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /app/pathless .

# Final stage
FROM scratch
COPY --from=go_builder /app/pathless /pathless
USER 10001
ENTRYPOINT ["/pathless"]