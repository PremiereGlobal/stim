FROM golang:1.13 as builder

ARG VERSION=default
ARG GOOS=linux
ENV GO111MODULE=on

WORKDIR /go/src/github.com/PremiereGlobal/stim/
COPY ./ .

RUN CGO_ENABLED=0 GOOS=${GOOS} go build -mod vendor -ldflags "-s -w -X github.com/PremiereGlobal/stim/stim.version=${VERSION}" -v -a -o bin/stim .

# Stage 2

FROM alpine:3.11

ENV STIM_PATH=/stim
ENV STIM_CACHE_PATH=/cache

VOLUME /stim
VOLUME /cache

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/PremiereGlobal/stim/bin/stim /usr/bin

ENTRYPOINT ["stim"]
