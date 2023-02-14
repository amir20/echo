
FROM --platform=$BUILDPLATFORM golang:1.19.3-alpine AS builder

WORKDIR /build
COPY go.mod ./
COPY *.go ./

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -ldflags "-s -w"  -o echo

FROM scratch
ENV PATH /bin
COPY --from=builder /build/echo /echo

ENTRYPOINT ["/echo"]
