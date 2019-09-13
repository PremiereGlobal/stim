FROM golang:1.12.5-alpine as builder

ARG VERSION=default
ARG GOOS=linux

WORKDIR /go/src/github.com/PremiereGlobal/stim/
COPY ./ .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=${GOOS} go build -mod=vendor -ldflags "-X github.com/PremiereGlobal/stim/stim.version=${VERSION}" -v -a -o bin/stim .

# Stage 2

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/PremiereGlobal/stim/bin/stim /usr/bin

ENTRYPOINT ["stim"]
