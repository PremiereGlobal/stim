FROM golang:1.11-alpine as builder

ARG VERSION=default
ARG GOOS=linux

WORKDIR /go/src/github.com/PremiereGlobal/stim/
COPY ./ .

RUN CGO_ENABLED=0 GOOS=${GOOS} go build -ldflags "-X github.com/PremiereGlobal/stim/stim.version=${VERSION}" -v -a -o bin/stim .

# Stage 2

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/PremiereGlobal/stim/bin/stim /usr/bin

ENTRYPOINT ["stim"]
