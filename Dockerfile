
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /build
COPY go.mod ./
COPY *.go ./

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -ldflags "-s -w"  -o echo

FROM alpine:latest
ENV PATH=/bin
COPY --from=builder /build/echo /echo

ENTRYPOINT ["/echo"]
