FROM golang:1.11 as builder
ADD . /go/src/github.com/eliotstocker/github-pr-close-resource
WORKDIR /go/src/github.com/eliotstocker/github-pr-close-resource
ENV TARGET=linux ARCH=amd64
RUN make build

FROM alpine:3.8 as resource
COPY --from=builder /go/src/github.com/eliotstocker/github-pr-close-resource/build /opt/resource
RUN apk add --update --no-cache \
    git \
    openssh \
    && chmod +x /opt/resource/*
ADD scripts/install_git_crypt.sh install_git_crypt.sh
RUN ./install_git_crypt.sh && rm ./install_git_crypt.sh

FROM resource
LABEL MAINTAINER=eliotstocker
