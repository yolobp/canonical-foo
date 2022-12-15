# Mostly copied from:
# https://github.com/GoogleContainerTools/skaffold/blob/main/examples/cross-platform-builds/Dockerfile

FROM golang:1.19 as builder

WORKDIR /code
COPY cmd/hello-foo/main.go .
COPY go.mod .

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -a -trimpath -o /app main.go

FROM gcr.io/distroless/static
# Define GOTRACEBACK to mark this container as using the Go language runtime
# for `skaffold debug` (https://skaffold.dev/docs/workflows/debug/).
ENV GOTRACEBACK=single
COPY --from=builder /app .

ENV PORT 8080
ENTRYPOINT ["/app"]
