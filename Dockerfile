FROM golang:1.22.3-bookworm AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0
WORKDIR /build
COPY . .
RUN go mod tidy && go build -ldflags "-s -w" -o main

FROM quay.io/kubevirt/sidecar-shim:v1.2.0
COPY --from=builder /build/main /usr/bin/onDefineDomain

ENTRYPOINT [ "/sidecar-shim" ]