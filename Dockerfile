# syntax=docker/dockerfile:1

ARG GO_VERSION=1.27
ARG NODE_VERSION=24

# ---- Frontend assets (Tailwind / daisyUI / htmx via Parcel) ----
FROM node:${NODE_VERSION}-slim AS assets
WORKDIR /src
COPY package.json package-lock.json .postcssrc ./
RUN --mount=type=cache,target=/root/.npm npm ci --no-audit --no-fund
COPY assets ./assets
# Tailwind scans the .templ files for class names
COPY internal/templates ./internal/templates
RUN npm run build

# ---- Go build ----
FROM golang:${GO_VERSION} AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
COPY --from=assets /src/internal/templates/static ./internal/templates/static
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go tool templ generate && \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/mid-blog ./cmd/web

# ---- Dev image: Go toolchain + air for live reload (source is bind-mounted) ----
FROM golang:${GO_VERSION} AS dev
RUN go install github.com/air-verse/air@latest
WORKDIR /src
EXPOSE 4231
CMD ["air", "-c", ".air.docker.toml"]

# ---- Production image ----
FROM gcr.io/distroless/static-debian12:nonroot AS prod
COPY --from=build /out/mid-blog /mid-blog
EXPOSE 4231
USER nonroot:nonroot
ENTRYPOINT ["/mid-blog"]
